package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rack-ac/circadian/internal/config"
)

const (
	authPath  = "/oauth/oauth2/auth"
	tokenPath = "/oauth/oauth2/token"
)

var allScopes = []string{
	"read:recovery",
	"read:cycles",
	"read:workout",
	"read:sleep",
	"read:profile",
	"read:body_measurement",
}

type Manager struct {
	cfg    config.Config
	store  *TokenStore
	client *http.Client
	logger *log.Logger

	mu sync.Mutex
}

func NewManager(cfg config.Config, store *TokenStore, logger *log.Logger) *Manager {
	return &Manager{
		cfg:    cfg,
		store:  store,
		client: &http.Client{Timeout: 20 * time.Second},
		logger: logger,
	}
}

func (m *Manager) GetAccessToken(ctx context.Context) (string, error) {
	token, err := m.store.Load()
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			token, err = m.AuthorizeInteractive(ctx)
			if err != nil {
				return "", err
			}
			return token.AccessToken, nil
		}
		return "", err
	}

	if token.AccessToken != "" && time.Until(token.Expiry) > 30*time.Second {
		return token.AccessToken, nil
	}

	token, err = m.RefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrMissingRefreshToken) {
			token, err = m.AuthorizeInteractive(ctx)
			if err != nil {
				return "", err
			}
			return token.AccessToken, nil
		}
		return "", err
	}

	return token.AccessToken, nil
}

var ErrMissingRefreshToken = errors.New("refresh token missing")

func (m *Manager) RefreshToken(ctx context.Context, current *Token) (*Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	latest, err := m.store.Load()
	if err == nil && latest.AccessToken != "" && time.Until(latest.Expiry) > 30*time.Second {
		return latest, nil
	}
	if err == nil {
		current = latest
	}

	if current == nil || current.RefreshToken == "" {
		return nil, ErrMissingRefreshToken
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", current.RefreshToken)
	form.Set("client_id", m.cfg.ClientID)
	form.Set("client_secret", m.cfg.ClientSecret)

	token, err := m.exchange(ctx, form)
	if err != nil {
		return nil, fmt.Errorf("refresh WHOOP token: %w", err)
	}
	if token.RefreshToken == "" {
		token.RefreshToken = current.RefreshToken
	}
	if err := m.store.Save(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Manager) EnsureFreshToken(ctx context.Context) (*Token, error) {
	current, err := m.store.Load()
	if err != nil && !errors.Is(err, ErrTokenNotFound) {
		return nil, err
	}
	if current != nil {
		token, refreshErr := m.RefreshToken(ctx, current)
		if refreshErr == nil {
			return token, nil
		}
		if !errors.Is(refreshErr, ErrMissingRefreshToken) {
			m.logger.Printf("WHOOP token refresh failed, falling back to browser auth: %v", refreshErr)
		}
	}
	return m.AuthorizeInteractive(ctx)
}

func (m *Manager) AuthorizeInteractive(ctx context.Context) (*Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if token, err := m.store.Load(); err == nil && token.AccessToken != "" && time.Until(token.Expiry) > 30*time.Second {
		return token, nil
	}

	state, err := randomState()
	if err != nil {
		return nil, fmt.Errorf("generate oauth state: %w", err)
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	redirectURL, err := url.Parse(m.cfg.RedirectURI)
	if err != nil {
		return nil, fmt.Errorf("parse redirect URI: %w", err)
	}

	serverErrCh := make(chan error, 1)
	httpServer, listenErr := startCallbackServer(redirectURL, state, codeCh, errCh, serverErrCh)
	if listenErr != nil {
		return nil, listenErr
	}
	defer shutdownServer(httpServer)

	authURL := m.authorizationURL(state)
	m.logger.Printf("opening WHOOP authorization in browser")
	if err := openBrowser(authURL); err != nil {
		m.logger.Printf("browser open failed, open this URL manually: %s", authURL)
	} else {
		m.logger.Printf("if your browser did not open, authorize manually: %s", authURL)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("authorization canceled: %w", ctx.Err())
	case err := <-errCh:
		return nil, err
	case err := <-serverErrCh:
		return nil, err
	case code := <-codeCh:
		form := url.Values{}
		form.Set("grant_type", "authorization_code")
		form.Set("code", code)
		form.Set("redirect_uri", m.cfg.RedirectURI)
		form.Set("client_id", m.cfg.ClientID)
		form.Set("client_secret", m.cfg.ClientSecret)

		token, err := m.exchange(ctx, form)
		if err != nil {
			return nil, fmt.Errorf("exchange WHOOP authorization code: %w", err)
		}
		if err := m.store.Save(token); err != nil {
			return nil, err
		}
		return token, nil
	}
}

func (m *Manager) authorizationURL(state string) string {
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", m.cfg.ClientID)
	query.Set("redirect_uri", m.cfg.RedirectURI)
	query.Set("scope", strings.Join(allScopes, " "))
	query.Set("state", state)

	return strings.TrimRight(m.cfg.APIRoot, "/") + authPath + "?" + query.Encode()
}

func (m *Manager) exchange(ctx context.Context, form url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(m.cfg.APIRoot, "/")+tokenPath,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read token response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("WHOOP token endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		ExpiresIn    int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if payload.AccessToken == "" {
		return nil, fmt.Errorf("WHOOP token response did not include access_token")
	}

	expiry := time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	if payload.ExpiresIn == 0 {
		expiry = time.Now().Add(50 * time.Minute)
	}

	return &Token{
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
		TokenType:    payload.TokenType,
		Scope:        payload.Scope,
		Expiry:       expiry,
	}, nil
}

func startCallbackServer(redirectURL *url.URL, expectedState string, codeCh chan<- string, errCh chan<- error, serverErrCh chan<- error) (*http.Server, error) {
	host := redirectURL.Host
	if !strings.Contains(host, ":") {
		if redirectURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("listen for WHOOP callback on %s: %w", host, err)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	mux.HandleFunc(redirectURL.Path, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Query().Get("error") != "":
			http.Error(w, "WHOOP authorization failed. You can close this window.", http.StatusBadRequest)
			errCh <- fmt.Errorf("WHOOP authorization failed: %s", r.URL.Query().Get("error"))
		case r.URL.Query().Get("state") != expectedState:
			http.Error(w, "Invalid OAuth state. You can close this window.", http.StatusBadRequest)
			errCh <- fmt.Errorf("invalid OAuth state returned by WHOOP")
		case r.URL.Query().Get("code") == "":
			http.Error(w, "WHOOP did not return an authorization code. You can close this window.", http.StatusBadRequest)
			errCh <- fmt.Errorf("WHOOP did not return an authorization code")
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.WriteString(w, "<html><body><h1>WHOOP authorization complete.</h1><p>You can close this window and return to your MCP client.</p></body></html>")
			codeCh <- r.URL.Query().Get("code")
		}
	})

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- fmt.Errorf("WHOOP callback server failed: %w", err)
		}
	}()

	return server, nil
}

func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func randomState() (string, error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func openBrowser(target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	return cmd.Start()
}

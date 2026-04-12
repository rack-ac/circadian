package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rack-ac/circadian/internal/version"
)

const (
	defaultRedirectURI = "http://127.0.0.1:8976/callback"
	defaultAPIRoot     = "https://api.prod.whoop.com"
	defaultMinInterval = 350 * time.Millisecond
	defaultHTTPHost    = "127.0.0.1"
	defaultHTTPPort    = "8080"
	defaultHTTPPath    = "/mcp"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	APIRoot      string
	TokenPath    string
	Version      string
	MinInterval  time.Duration
	HTTPHost     string
	HTTPPort     string
	HTTPPath     string
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home directory: %w", err)
	}

	cfg := Config{
		ClientID:     os.Getenv("WHOOP_CLIENT_ID"),
		ClientSecret: os.Getenv("WHOOP_CLIENT_SECRET"),
		RedirectURI:  envOrDefault("WHOOP_REDIRECT_URI", defaultRedirectURI),
		APIRoot:      envOrDefault("WHOOP_API_ROOT", defaultAPIRoot),
		TokenPath:    filepath.Join(home, ".circadian", "tokens.json"),
		Version:      envOrDefault("CIRCADIAN_VERSION", version.Version),
		MinInterval:  defaultMinInterval,
		HTTPHost:     envOrDefault("CIRCADIAN_HOST", defaultHTTPHost),
		HTTPPort:     envOrDefault("CIRCADIAN_PORT", defaultHTTPPort),
		HTTPPath:     envOrDefault("CIRCADIAN_PATH", defaultHTTPPath),
	}

	if cfg.ClientID == "" {
		return Config{}, fmt.Errorf("WHOOP_CLIENT_ID is required")
	}
	if cfg.ClientSecret == "" {
		return Config{}, fmt.Errorf("WHOOP_CLIENT_SECRET is required")
	}
	if cfg.MinInterval, err = loadMinInterval(); err != nil {
		return Config{}, err
	}
	if cfg.HTTPPath == "" || cfg.HTTPPath[0] != '/' {
		return Config{}, fmt.Errorf("CIRCADIAN_PATH must start with '/'")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadMinInterval() (time.Duration, error) {
	raw := os.Getenv("WHOOP_MIN_INTERVAL_MS")
	if raw == "" {
		return defaultMinInterval, nil
	}

	ms, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("WHOOP_MIN_INTERVAL_MS must be an integer number of milliseconds")
	}
	if ms < 0 {
		return 0, fmt.Errorf("WHOOP_MIN_INTERVAL_MS must be >= 0")
	}

	return time.Duration(ms) * time.Millisecond, nil
}

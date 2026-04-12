package whoop

import (
	"net/http"
	"strings"
	"time"

	"github.com/rack-ac/circadian/internal/config"
	"github.com/rack-ac/circadian/internal/oauth"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *oauth.Manager
	limiter    *rateLimiter
}

func NewClient(cfg config.Config, auth *oauth.Manager) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.APIRoot, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		auth:    auth,
		limiter: newRateLimiter(cfg.MinInterval),
	}
}

type CollectionParams struct {
	Limit     int
	Start     string
	End       string
	NextToken string
}

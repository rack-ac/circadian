package server

import (
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rack-ac/circadian/internal/config"
	"github.com/rack-ac/circadian/internal/oauth"
	"github.com/rack-ac/circadian/internal/whoop"
)

type services struct {
	cfg         config.Config
	store       *oauth.TokenStore
	authManager *oauth.Manager
	client      *whoop.Client
}

func New(cfg config.Config) (*mcp.Server, error) {
	logger := log.New(os.Stderr, "circadian: ", log.LstdFlags)
	store := oauth.NewTokenStore(cfg.TokenPath)
	authManager := oauth.NewManager(cfg, store, logger)
	client := whoop.NewClient(cfg, authManager)

	svc := services{
		cfg:         cfg,
		store:       store,
		authManager: authManager,
		client:      client,
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "circadian",
		Version: cfg.Version,
	}, &mcp.ServerOptions{
		Instructions: "Use this server to access WHOOP user profile, body metrics, cycles, recoveries, sleeps, and workouts.",
	})
	server.AddReceivingMiddleware(newLoggingMiddleware())

	registerLifecycleTools(server, svc)
	registerDataTools(server, svc)

	return server, nil
}

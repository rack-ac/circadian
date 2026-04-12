package main

import (
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rack-ac/circadian/internal/config"
	whserver "github.com/rack-ac/circadian/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	server, err := whserver.New(cfg)
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	logger := log.New(os.Stderr, "circadian: ", log.LstdFlags)
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle(cfg.HTTPPath, handler)

	addr := cfg.HTTPHost + ":" + cfg.HTTPPort
	logger.Printf("starting circadian %s on http://%s%s (min interval %s)", cfg.Version, addr, cfg.HTTPPath, cfg.MinInterval)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}

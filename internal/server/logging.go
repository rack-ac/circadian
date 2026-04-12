package server

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newLoggingMiddleware() mcp.Middleware {
	logger := log.New(os.Stderr, "circadian: ", log.LstdFlags)

	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			start := time.Now()
			sessionID := "unknown"
			if session := req.GetSession(); session != nil {
				sessionID = session.ID()
			}

			logger.Printf("mcp request session=%s method=%s", sessionID, method)
			result, err := next(ctx, method, req)
			if err != nil {
				logger.Printf("mcp response session=%s method=%s status=error duration=%s error=%v", sessionID, method, time.Since(start), err)
				return result, err
			}

			logger.Printf("mcp response session=%s method=%s status=ok duration=%s", sessionID, method, time.Since(start))
			return result, nil
		}
	}
}

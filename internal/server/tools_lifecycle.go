package server

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rack-ac/circadian/internal/oauth"
)

func registerLifecycleTools(server *mcp.Server, svc services) {
	mcp.AddTool(server, &mcp.Tool{Name: "whoop_auth_status", Description: "Inspect local WHOOP auth state without triggering browser auth."}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		token, err := svc.store.Load()
		if err != nil && err != oauth.ErrTokenNotFound {
			return nil, nil, err
		}

		status := map[string]any{
			"configured":    true,
			"redirect_uri":  svc.cfg.RedirectURI,
			"has_token":     token != nil,
			"authenticated": token != nil && token.AccessToken != "",
		}
		if token != nil {
			status["expires_at"] = token.Expiry.UTC().Format(time.RFC3339)
			status["expired"] = time.Now().After(token.Expiry)
		}

		return shapedResult(status, nil)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_reauthorize", Description: "Force a fresh WHOOP browser authorization and replace the stored token."}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		token, err := svc.authManager.AuthorizeInteractive(ctx)
		if err != nil {
			return nil, nil, err
		}
		return shapedResult(map[string]any{
			"authenticated": true,
			"expires_at":    token.Expiry.UTC().Format(time.RFC3339),
		}, nil)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_revoke_access", Description: "Revoke WHOOP OAuth access and delete the local token file."}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		err := svc.client.RevokeAccess(ctx)
		if err != nil {
			return nil, nil, err
		}
		if err := svc.store.Delete(); err != nil {
			return nil, nil, err
		}
		return shapedResult(map[string]any{
			"revoked": true,
		}, nil)
	})
}

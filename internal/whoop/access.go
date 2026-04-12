package whoop

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) RevokeAccess(ctx context.Context) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	token, err := c.auth.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("authorize WHOOP revoke request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/developer/v2/user/access", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request WHOOP revoke endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read WHOOP revoke response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("WHOOP API returned 401 while revoking access")
	default:
		return fmt.Errorf("WHOOP revoke endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
}

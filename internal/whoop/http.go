package whoop

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (c *Client) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}

		token, err := c.auth.GetAccessToken(ctx)
		if err != nil {
			return fmt.Errorf("authorize WHOOP request: %w", err)
		}

		reqURL := c.baseURL + path
		if query != nil && len(query) > 0 {
			reqURL += "?" + query.Encode()
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request WHOOP API: %w", err)
			if wait(ctx, attempt, "") != nil {
				return lastErr
			}
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		closeErr := resp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("read WHOOP response: %w", readErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close WHOOP response body: %w", closeErr)
		}

		switch {
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			if err := json.Unmarshal(body, out); err != nil {
				return fmt.Errorf("decode WHOOP response: %w", err)
			}
			return nil
		case resp.StatusCode == http.StatusUnauthorized:
			if _, err := c.auth.EnsureFreshToken(ctx); err == nil {
				continue
			}
			return fmt.Errorf("WHOOP API returned 401 and token refresh failed; re-authentication is required")
		case resp.StatusCode == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("WHOOP API rate limited the request")
			c.limiter.Delay(parseRetryAfter(resp.Header.Get("Retry-After")))
			if wait(ctx, attempt, resp.Header.Get("Retry-After")) != nil {
				return lastErr
			}
			continue
		case resp.StatusCode >= 500 && resp.StatusCode <= 599:
			lastErr = fmt.Errorf("WHOOP API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
			if wait(ctx, attempt, "") != nil {
				return lastErr
			}
			continue
		default:
			return fmt.Errorf("WHOOP API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("WHOOP request failed after retries")
}

func collectionQuery(params CollectionParams) url.Values {
	values := url.Values{}
	if params.Limit > 0 {
		values.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Start != "" {
		values.Set("start", params.Start)
	}
	if params.End != "" {
		values.Set("end", params.End)
	}
	if params.NextToken != "" {
		values.Set("nextToken", params.NextToken)
	}
	return values
}

func wait(ctx context.Context, attempt int, retryAfter string) error {
	delay := time.Duration(1<<attempt) * time.Second
	if parsed := parseRetryAfter(retryAfter); parsed > 0 {
		delay = parsed
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(value); err == nil {
		return time.Until(when)
	}
	return 0
}

package whoop

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestParseRetryAfterSeconds(t *testing.T) {
	t.Parallel()

	got := parseRetryAfter("3")
	if got != 3*time.Second {
		t.Fatalf("parseRetryAfter() = %v, want %v", got, 3*time.Second)
	}
}

func TestParseRetryAfterHTTPDate(t *testing.T) {
	t.Parallel()

	retryAt := time.Now().Add(2 * time.Second).UTC()
	got := parseRetryAfter(retryAt.Format(http.TimeFormat))
	if got < time.Second || got > 3*time.Second {
		t.Fatalf("parseRetryAfter() = %v, want about 2s", got)
	}
}

func TestRateLimiterDelay(t *testing.T) {
	t.Parallel()

	limiter := newRateLimiter(0)
	limiter.Delay(40 * time.Millisecond)

	start := time.Now()
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if elapsed := time.Since(start); elapsed < 35*time.Millisecond {
		t.Fatalf("Wait() elapsed = %v, want at least 35ms", elapsed)
	}
}

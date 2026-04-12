package whoop

import (
	"context"
	"sync"
	"time"
)

type rateLimiter struct {
	mu          sync.Mutex
	minInterval time.Duration
	nextAllowed time.Time
}

func newRateLimiter(minInterval time.Duration) *rateLimiter {
	return &rateLimiter{minInterval: minInterval}
}

func (r *rateLimiter) Wait(ctx context.Context) error {
	if r == nil {
		return nil
	}

	for {
		r.mu.Lock()
		waitFor := time.Until(r.nextAllowed)
		if waitFor <= 0 {
			r.nextAllowed = time.Now().Add(r.minInterval)
			r.mu.Unlock()
			return nil
		}
		r.mu.Unlock()

		timer := time.NewTimer(waitFor)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func (r *rateLimiter) Delay(delay time.Duration) {
	if r == nil || delay <= 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	until := time.Now().Add(delay)
	if until.After(r.nextAllowed) {
		r.nextAllowed = until
	}
}

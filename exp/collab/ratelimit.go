package collab

import (
	"sync"
	"time"
)

// tokenBucket implements a simple token bucket rate limiter.
// NOT thread-safe — use directly only in single-goroutine contexts (e.g. ReadPump).
type tokenBucket struct {
	tokens    float64
	capacity  float64
	rate      float64 // tokens per second
	lastCheck time.Time
}

// allow checks if a token is available and consumes it. Not thread-safe.
// A zero-valued bucket (rate == 0) always allows — this is the pass-through default.
func (b *tokenBucket) allow() bool {
	if b.rate == 0 {
		return true
	}
	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.lastCheck = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// IPRateLimiter applies per-IP rate limiting using token buckets.
type IPRateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64
	capacity float64
}

// NewIPRateLimiter creates a rate limiter with the given rate (tokens/s) and burst capacity.
func NewIPRateLimiter(rate, capacity float64) *IPRateLimiter {
	return &IPRateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     rate,
		capacity: capacity,
	}
}

// Allow checks if the given IP is allowed to proceed.
// Returns whether the request is allowed and the duration to wait before retrying.
func (rl *IPRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[ip]
	if !ok {
		b = &tokenBucket{
			tokens:    rl.capacity,
			capacity:  rl.capacity,
			rate:      rl.rate,
			lastCheck: time.Now(),
		}
		rl.buckets[ip] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.lastCheck = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	if b.tokens >= 1 {
		b.tokens--
		return true, 0
	}

	retryAfter := time.Duration((1 - b.tokens) / b.rate * float64(time.Second))
	return false, retryAfter
}

// Cleanup removes stale buckets that have been idle for longer than staleThreshold.
func (rl *IPRateLimiter) Cleanup(staleThreshold time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, b := range rl.buckets {
		if now.Sub(b.lastCheck) > staleThreshold {
			delete(rl.buckets, ip)
		}
	}
}

// StartCleanup runs periodic cleanup of stale buckets. Stops when done is closed.
func (rl *IPRateLimiter) StartCleanup(interval, staleThreshold time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.Cleanup(staleThreshold)
			case <-done:
				return
			}
		}
	}()
}

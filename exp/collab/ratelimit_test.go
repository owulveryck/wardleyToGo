package collab

import (
	"testing"
	"time"
)

func TestTokenBucketAllow(t *testing.T) {
	b := tokenBucket{tokens: 3, capacity: 5, rate: 1, lastCheck: time.Now()}
	// Should allow 3 requests (3 tokens available)
	for i := 0; i < 3; i++ {
		if !b.allow() {
			t.Fatalf("expected allow on request %d", i)
		}
	}
	// 4th should fail
	if b.allow() {
		t.Fatal("expected deny after exhausting tokens")
	}
}

func TestTokenBucketRefill(t *testing.T) {
	b := tokenBucket{tokens: 0, capacity: 5, rate: 10, lastCheck: time.Now().Add(-time.Second)}
	// After 1 second at rate=10, should have 10 tokens (capped at 5)
	if !b.allow() {
		t.Fatal("expected allow after refill")
	}
	// Should have 4 tokens left (5 refilled - 1 consumed)
	if b.tokens < 3.9 || b.tokens > 4.1 {
		t.Fatalf("expected ~4 tokens, got %f", b.tokens)
	}
}

func TestIPRateLimiterAllow(t *testing.T) {
	rl := NewIPRateLimiter(10, 3) // 10/s, burst 3

	// First 3 should succeed
	for i := 0; i < 3; i++ {
		allowed, _ := rl.Allow("192.168.1.1")
		if !allowed {
			t.Fatalf("expected allow on request %d", i)
		}
	}

	// 4th should fail
	allowed, retryAfter := rl.Allow("192.168.1.1")
	if allowed {
		t.Fatal("expected deny after exhausting burst")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected positive retryAfter, got %v", retryAfter)
	}
}

func TestIPRateLimiterDifferentIPs(t *testing.T) {
	rl := NewIPRateLimiter(10, 2)

	// Exhaust IP1
	rl.Allow("ip1")
	rl.Allow("ip1")
	allowed, _ := rl.Allow("ip1")
	if allowed {
		t.Fatal("ip1 should be denied")
	}

	// IP2 should still work
	allowed, _ = rl.Allow("ip2")
	if !allowed {
		t.Fatal("ip2 should be allowed")
	}
}

func TestIPRateLimiterRetryAfter(t *testing.T) {
	rl := NewIPRateLimiter(1, 1) // 1 token/s, burst 1

	// Consume the token
	rl.Allow("test-ip")

	// Should return retryAfter ~1 second
	allowed, retryAfter := rl.Allow("test-ip")
	if allowed {
		t.Fatal("expected deny")
	}
	if retryAfter < 500*time.Millisecond || retryAfter > 1500*time.Millisecond {
		t.Fatalf("expected retryAfter ~1s, got %v", retryAfter)
	}
}

func TestIPRateLimiterCleanup(t *testing.T) {
	rl := NewIPRateLimiter(10, 5)

	rl.Allow("stale-ip")
	rl.Allow("fresh-ip")

	// Manually backdate the stale bucket
	rl.mu.Lock()
	rl.buckets["stale-ip"].lastCheck = time.Now().Add(-20 * time.Minute)
	rl.mu.Unlock()

	rl.Cleanup(10 * time.Minute)

	rl.mu.Lock()
	_, staleExists := rl.buckets["stale-ip"]
	_, freshExists := rl.buckets["fresh-ip"]
	rl.mu.Unlock()

	if staleExists {
		t.Fatal("stale bucket should have been cleaned up")
	}
	if !freshExists {
		t.Fatal("fresh bucket should not be cleaned up")
	}
}

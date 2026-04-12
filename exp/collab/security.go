package collab

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// clientIP extracts the client IP from the request,
// checking X-Forwarded-For, X-Real-IP, then falling back to RemoteAddr.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.IndexByte(xff, ','); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// adminAuth returns middleware that requires a bearer token or query param token.
// If token is empty, all requests pass through (dev mode).
func adminAuth(next http.HandlerFunc, token string) http.HandlerFunc {
	if token == "" {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "Bearer "+token {
			next(w, r)
			return
		}
		if r.URL.Query().Get("token") == token {
			next(w, r)
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

// rateLimitMiddleware wraps a handler with per-IP rate limiting.
// Returns HTTP 429 with Retry-After header when rate exceeded.
// If limiter is nil, requests pass through.
func rateLimitMiddleware(next http.HandlerFunc, limiter *IPRateLimiter) http.HandlerFunc {
	if limiter == nil {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		allowed, retryAfter := limiter.Allow(ip)
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// securityHeaders adds standard security headers to all responses.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self' ws: wss:")
		next.ServeHTTP(w, r)
	})
}

// ConnTracker tracks active WebSocket connections per IP and per session.
type ConnTracker struct {
	mu         sync.Mutex
	perIP      map[string]int
	perSession map[string]int
}

// NewConnTracker creates a new connection tracker.
func NewConnTracker() *ConnTracker {
	return &ConnTracker{
		perIP:      make(map[string]int),
		perSession: make(map[string]int),
	}
}

// TryAdd attempts to register a new connection. Returns false if either limit is exceeded.
func (ct *ConnTracker) TryAdd(ip, sessionID string, maxPerIP, maxPerSession int) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if ct.perIP[ip] >= maxPerIP {
		return false
	}
	if ct.perSession[sessionID] >= maxPerSession {
		return false
	}

	ct.perIP[ip]++
	ct.perSession[sessionID]++
	return true
}

// Remove decrements the counters for the given IP and session.
func (ct *ConnTracker) Remove(ip, sessionID string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if ct.perIP[ip] > 0 {
		ct.perIP[ip]--
		if ct.perIP[ip] == 0 {
			delete(ct.perIP, ip)
		}
	}
	if ct.perSession[sessionID] > 0 {
		ct.perSession[sessionID]--
		if ct.perSession[sessionID] == 0 {
			delete(ct.perSession, sessionID)
		}
	}
}

// ConnCounts returns the current connection counts (for testing).
func (ct *ConnTracker) ConnCounts(ip, sessionID string) (ipCount, sessionCount int) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.perIP[ip], ct.perSession[sessionID]
}

// maxConnPerIP is the default maximum WebSocket connections per IP address.
const maxConnPerIP = 50

// maxConnPerSession is the default maximum clients per session.
const maxConnPerSession = 50

// msgRateLimit is the maximum messages per second per WebSocket client.
const msgRateLimit = 30.0

// msgRateBurst is the burst capacity for per-client message rate limiting.
const msgRateBurst = 30.0

// newMsgBucket creates a token bucket for per-client message rate limiting.
func newMsgBucket() tokenBucket {
	return tokenBucket{
		tokens:    msgRateBurst,
		capacity:  msgRateBurst,
		rate:      msgRateLimit,
		lastCheck: time.Now(),
	}
}

package collab

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIP_XForwardedFor(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Forwarded-For", "203.0.113.50, 70.41.3.18, 150.172.238.178")
	if got := clientIP(r); got != "203.0.113.50" {
		t.Fatalf("expected 203.0.113.50, got %s", got)
	}
}

func TestClientIP_XForwardedForSingle(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Forwarded-For", "10.0.0.1")
	if got := clientIP(r); got != "10.0.0.1" {
		t.Fatalf("expected 10.0.0.1, got %s", got)
	}
}

func TestClientIP_XRealIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Real-IP", "10.0.0.2")
	if got := clientIP(r); got != "10.0.0.2" {
		t.Fatalf("expected 10.0.0.2, got %s", got)
	}
}

func TestClientIP_RemoteAddr(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "192.168.1.1:12345"
	if got := clientIP(r); got != "192.168.1.1" {
		t.Fatalf("expected 192.168.1.1, got %s", got)
	}
}

func TestAdminAuth_DevMode(t *testing.T) {
	called := false
	handler := adminAuth(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}, "")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)

	if !called {
		t.Fatal("handler should be called in dev mode (empty token)")
	}
}

func TestAdminAuth_BearerToken(t *testing.T) {
	called := false
	handler := adminAuth(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}, "secret123")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer secret123")
	handler(w, r)

	if !called {
		t.Fatal("handler should be called with valid bearer token")
	}
}

func TestAdminAuth_QueryToken(t *testing.T) {
	called := false
	handler := adminAuth(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}, "secret123")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?token=secret123", nil)
	handler(w, r)

	if !called {
		t.Fatal("handler should be called with valid query token")
	}
}

func TestAdminAuth_Unauthorized(t *testing.T) {
	handler := adminAuth(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called without auth")
	}, "secret123")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminAuth_WrongToken(t *testing.T) {
	handler := adminAuth(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called with wrong token")
	}, "secret123")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?token=wrongtoken", nil)
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	limiter := NewIPRateLimiter(10, 2) // burst 2
	callCount := 0
	handler := rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		callCount++
	}, limiter)

	// First 2 should pass
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = "1.2.3.4:1234"
		handler(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, w.Code)
		}
	}

	// 3rd should be rate limited
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "1.2.3.4:1234"
	handler(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}

	if callCount != 2 {
		t.Fatalf("expected handler called 2 times, got %d", callCount)
	}
}

func TestRateLimitMiddleware_NilLimiter(t *testing.T) {
	called := false
	handler := rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)

	if !called {
		t.Fatal("handler should pass through with nil limiter")
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":       "DENY",
		"Referrer-Policy":       "strict-origin-when-cross-origin",
	}
	for header, expected := range checks {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("header %s: expected %q, got %q", header, expected, got)
		}
	}
	if w.Header().Get("Content-Security-Policy") == "" {
		t.Error("expected Content-Security-Policy header")
	}
}

func TestConnTracker_TryAdd(t *testing.T) {
	ct := NewConnTracker()

	// Should succeed
	if !ct.TryAdd("ip1", "s1", 2, 3) {
		t.Fatal("expected first add to succeed")
	}
	if !ct.TryAdd("ip1", "s1", 2, 3) {
		t.Fatal("expected second add to succeed")
	}

	// Should fail: ip1 at max (2)
	if ct.TryAdd("ip1", "s1", 2, 3) {
		t.Fatal("expected add to fail (IP limit)")
	}

	// Different IP should work
	if !ct.TryAdd("ip2", "s1", 2, 3) {
		t.Fatal("expected add with different IP to succeed")
	}
}

func TestConnTracker_SessionLimit(t *testing.T) {
	ct := NewConnTracker()

	for i := 0; i < 3; i++ {
		ct.TryAdd("ip"+string(rune('0'+i)), "s1", 10, 3)
	}

	// Should fail: session at max (3)
	if ct.TryAdd("ip9", "s1", 10, 3) {
		t.Fatal("expected add to fail (session limit)")
	}

	// Different session should work
	if !ct.TryAdd("ip9", "s2", 10, 3) {
		t.Fatal("expected add with different session to succeed")
	}
}

func TestConnTracker_Remove(t *testing.T) {
	ct := NewConnTracker()

	ct.TryAdd("ip1", "s1", 10, 10)
	ct.TryAdd("ip1", "s1", 10, 10)

	ipCount, sessionCount := ct.ConnCounts("ip1", "s1")
	if ipCount != 2 || sessionCount != 2 {
		t.Fatalf("expected (2,2), got (%d,%d)", ipCount, sessionCount)
	}

	ct.Remove("ip1", "s1")

	ipCount, sessionCount = ct.ConnCounts("ip1", "s1")
	if ipCount != 1 || sessionCount != 1 {
		t.Fatalf("expected (1,1), got (%d,%d)", ipCount, sessionCount)
	}

	ct.Remove("ip1", "s1")

	ipCount, sessionCount = ct.ConnCounts("ip1", "s1")
	if ipCount != 0 || sessionCount != 0 {
		t.Fatalf("expected (0,0), got (%d,%d)", ipCount, sessionCount)
	}
}

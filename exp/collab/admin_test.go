package collab

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminIndex_Empty(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "No active sessions") {
		t.Error("expected empty state message")
	}
}

func TestAdminIndex_WithSessions(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, s.ID[:12]) {
		t.Error("expected session ID in page")
	}
	// Tokens are masked (value="*****") but present in data-url attributes
	if !strings.Contains(body, "*****") {
		t.Error("expected masked token placeholder in page")
	}
	if !strings.Contains(body, "data-url=\"/s/"+s.ID+"/"+s.WriteAccessID+"\"") {
		t.Error("expected write access URL in data-url attribute")
	}
	if !strings.Contains(body, "data-url=\"/s/"+s.ID+"/"+s.ReadAccessID+"\"") {
		t.Error("expected read access URL in data-url attribute")
	}
}

func TestAdminCreateSession(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("POST", "/session", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want 303", w.Code)
	}
	loc := w.Header().Get("Location")
	if loc != "/" {
		t.Errorf("redirect location = %q, want %q", loc, "/")
	}

	sessions := hub.Sessions()
	if len(sessions) != 1 {
		t.Errorf("sessions count = %d, want 1", len(sessions))
	}
}

func TestAdminDeleteSession(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("POST", "/session/"+s.ID+"/delete", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want 303", w.Code)
	}

	if hub.GetSession(s.ID) != nil {
		t.Error("session should be deleted")
	}
}

func TestAdminDeleteNonexistent(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("POST", "/session/nonexistent/delete", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should redirect without error
	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want 303", w.Code)
	}
}

func TestAdminMultipleSessions(t *testing.T) {
	hub := NewHub()
	hub.CreateSession() //nolint:errcheck
	hub.CreateSession() //nolint:errcheck
	hub.CreateSession() //nolint:errcheck
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "3 active sessions") {
		t.Error("expected '3 active sessions' in page")
	}
}

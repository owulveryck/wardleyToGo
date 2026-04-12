package collab

import (
	"testing"
	"time"
)

func TestHubCreateSession(t *testing.T) {
	h := NewHub()
	s, _ := h.CreateSession()

	if s.ID == "" {
		t.Fatal("session ID is empty")
	}
	if s.WriteAccessID == "" {
		t.Fatal("WriteAccessID is empty")
	}
	if s.ReadAccessID == "" {
		t.Fatal("ReadAccessID is empty")
	}
	if s.WriteAccessID == s.ReadAccessID {
		t.Fatal("WriteAccessID and ReadAccessID should be different")
	}

	// Should be findable
	found := h.GetSession(s.ID)
	if found == nil {
		t.Fatal("GetSession returned nil")
	}
	if found.ID != s.ID {
		t.Errorf("found.ID = %q, want %q", found.ID, s.ID)
	}
}

func TestHubGetSessionNotFound(t *testing.T) {
	h := NewHub()
	if h.GetSession("nonexistent") != nil {
		t.Fatal("expected nil for unknown session")
	}
}

func TestHubDeleteSession(t *testing.T) {
	h := NewHub()
	s, _ := h.CreateSession()
	id := s.ID

	h.DeleteSession(id)

	if h.GetSession(id) != nil {
		t.Fatal("session should have been deleted")
	}
}

func TestHubDeleteNonexistent(t *testing.T) {
	h := NewHub()
	// Should not panic
	h.DeleteSession("nonexistent")
}

func TestHubSessions(t *testing.T) {
	h := NewHub()
	h.CreateSession()
	h.CreateSession()
	h.CreateSession()

	sessions := h.Sessions()
	if len(sessions) != 3 {
		t.Errorf("expected 3 sessions, got %d", len(sessions))
	}
}

func TestHubMultipleSessions(t *testing.T) {
	h := NewHub()
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		s, _ := h.CreateSession()
		if ids[s.ID] {
			t.Fatalf("duplicate session ID: %s", s.ID)
		}
		ids[s.ID] = true
	}
}

func TestHubCleanupIdle(t *testing.T) {
	h := NewHub()
	s, _ := h.CreateSession()

	// Simulate idle session by backdating lastActive
	s.mu.Lock()
	s.lastActive = time.Now().Add(-2 * time.Hour)
	s.mu.Unlock()

	cleaned := h.CleanupIdle(1 * time.Hour)
	if cleaned != 1 {
		t.Errorf("cleaned = %d, want 1", cleaned)
	}
	if h.GetSession(s.ID) != nil {
		t.Fatal("idle session should have been removed")
	}
}

func TestHubCleanupIdleKeepsActive(t *testing.T) {
	h := NewHub()
	s, _ := h.CreateSession()

	// Session is freshly created, should not be cleaned
	cleaned := h.CleanupIdle(1 * time.Hour)
	if cleaned != 0 {
		t.Errorf("cleaned = %d, want 0", cleaned)
	}
	if h.GetSession(s.ID) == nil {
		t.Fatal("active session should not be removed")
	}
}

func TestHubCleanupIdleKeepsWithClients(t *testing.T) {
	h := NewHub()
	s, _ := h.CreateSession()

	// Add a fake client
	s.mu.Lock()
	s.clients["fake"] = &Client{id: "fake", send: make(chan []byte, 1)}
	s.lastActive = time.Now().Add(-2 * time.Hour) // old but has client
	s.mu.Unlock()

	cleaned := h.CleanupIdle(1 * time.Hour)
	if cleaned != 0 {
		t.Errorf("cleaned = %d, want 0 (session has clients)", cleaned)
	}
}

func TestGenerateIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID(8)
		if len(id) != 8 {
			t.Errorf("id length = %d, want 8", len(id))
		}
		if ids[id] {
			t.Fatalf("duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

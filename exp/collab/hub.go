package collab

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

const maxSessions = 100

// Hub manages all active collaboration sessions.
type Hub struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a new session with generated access IDs.
// Returns an error if the maximum number of sessions is reached.
func (h *Hub) CreateSession() (*Session, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.sessions) >= maxSessions {
		return nil, errors.New("maximum active sessions reached")
	}

	s := &Session{
		ID:            generateID(16),
		WriteAccessID: generateID(8),
		ReadAccessID:  generateID(8),
		doc:           NewDocument(),
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte, 256),
		stopCh:        make(chan struct{}),
		created:       time.Now(),
		lastActive:    time.Now(),
	}
	h.sessions[s.ID] = s
	go s.Run()
	return s, nil
}

// GetSession returns a session by ID, or nil if not found.
func (h *Hub) GetSession(id string) *Session {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sessions[id]
}

// DeleteSession stops and removes a session.
func (h *Hub) DeleteSession(id string) {
	h.mu.Lock()
	s := h.sessions[id]
	delete(h.sessions, id)
	h.mu.Unlock()

	if s != nil {
		s.Stop()
	}
}

// Sessions returns a snapshot of all active sessions.
func (h *Hub) Sessions() []*Session {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]*Session, 0, len(h.sessions))
	for _, s := range h.sessions {
		result = append(result, s)
	}
	return result
}

// CleanupIdle removes sessions that have had no connected clients
// for longer than maxIdle.
func (h *Hub) CleanupIdle(maxIdle time.Duration) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	var toDelete []string
	for id, s := range h.sessions {
		s.mu.RLock()
		clientCount := len(s.clients)
		lastActive := s.lastActive
		s.mu.RUnlock()

		if clientCount == 0 && now.Sub(lastActive) > maxIdle {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		s := h.sessions[id]
		delete(h.sessions, id)
		s.Stop()
	}
	return len(toDelete)
}

// StartCleanupTicker starts a background goroutine that periodically
// cleans up idle sessions. It stops when the done channel is closed.
func (h *Hub) StartCleanupTicker(interval, maxIdle time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				h.CleanupIdle(maxIdle)
			case <-done:
				return
			}
		}
	}()
}

func generateID(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

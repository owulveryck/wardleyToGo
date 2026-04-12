package collab

import (
	"fmt"
	"sync"
	"time"
)

// Session represents a collaborative editing session.
type Session struct {
	ID            string
	WriteAccessID string
	ReadAccessID  string

	doc     *Document
	history []*OpPayload // ring buffer for OT, last N ops

	mu         sync.RWMutex
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	stopCh     chan struct{}
	stopped    bool
	created    time.Time
	lastActive time.Time
}

const maxHistory = 100

// Input validation limits.
const (
	maxNameLen       = 50
	maxOpLines       = 500
	maxLineLen       = 1000
	maxDocumentLines = 5000
)

// User colors palette.
var userColors = []string{
	"#e74c3c", "#2ecc71", "#3498db", "#9b59b6", "#e67e22",
	"#1abc9c", "#f1c40f", "#e84393", "#00cec9", "#6c5ce7",
}

// ResolveAccess determines the access mode from an access ID.
// Returns "rw", "ro", or "" (unknown/rejected).
func (s *Session) ResolveAccess(accessID string) string {
	if accessID == s.WriteAccessID {
		return "rw"
	}
	if accessID == s.ReadAccessID {
		return "ro"
	}
	return ""
}

// ClientCount returns the number of connected clients.
func (s *Session) ClientCount() (rw, ro int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		if c.mode == "rw" {
			rw++
		} else {
			ro++
		}
	}
	return
}

// Created returns the session creation time.
func (s *Session) Created() time.Time {
	return s.created
}

// Run is the main loop for the session, handling register/unregister/broadcast.
func (s *Session) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.id] = client
			s.lastActive = time.Now()
			s.mu.Unlock()

			// Send welcome to the new client
			s.sendWelcome(client)

			// Notify others
			s.broadcastUserJoined(client)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.send)
				s.lastActive = time.Now()
			}
			s.mu.Unlock()

			// Notify others
			s.broadcastUserLeft(client)

		case message := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.send <- message:
				default:
					// Client buffer full, skip
				}
			}
			s.mu.RUnlock()

		case <-s.stopCh:
			s.mu.Lock()
			for _, client := range s.clients {
				close(client.send)
			}
			s.clients = make(map[string]*Client)
			s.mu.Unlock()
			return
		}
	}
}

// Stop signals the session to stop its Run loop.
func (s *Session) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.stopped {
		s.stopped = true
		close(s.stopCh)
	}
}

// validateOp checks that an operation's content is within limits.
func validateOp(op *OpPayload, currentDocLen int) error {
	if len(op.Lines) > maxOpLines {
		return fmt.Errorf("too many lines: %d (max %d)", len(op.Lines), maxOpLines)
	}
	for _, line := range op.Lines {
		if len(line) > maxLineLen {
			return fmt.Errorf("line too long: %d chars (max %d)", len(line), maxLineLen)
		}
	}
	resultingLines := currentDocLen
	switch op.Type {
	case "insert":
		resultingLines += len(op.Lines)
	case "replace":
		resultingLines = resultingLines - op.LineCount + len(op.Lines)
	}
	if resultingLines > maxDocumentLines {
		return fmt.Errorf("document would exceed %d lines", maxDocumentLines)
	}
	return nil
}

// HandleOp processes an incoming operation from a client.
func (s *Session) HandleOp(client *Client, op *OpPayload) {
	if client.mode != "rw" {
		return // read-only clients cannot edit
	}

	// Validate input before acquiring lock
	if err := validateOp(op, len(s.doc.Lines())); err != nil {
		errMsg, _ := MarshalMessage(MsgError, ErrorPayload{Message: err.Error()})
		select {
		case client.send <- errMsg:
		default:
		}
		return
	}

	s.doc.mu.Lock()

	// Transform against any ops the client hasn't seen
	if op.Version < s.doc.version {
		start := len(s.history) - int(s.doc.version-op.Version)
		if start < 0 {
			start = 0
		}
		for _, histOp := range s.history[start:] {
			TransformInPlace(op, histOp)
		}
	}

	err := s.doc.applyLocked(op)
	version := s.doc.version

	// Store in history ring buffer
	if err == nil {
		s.history = append(s.history, op)
		if len(s.history) > maxHistory {
			s.history = s.history[len(s.history)-maxHistory:]
		}
	}

	s.doc.mu.Unlock()

	if err != nil {
		errMsg, _ := MarshalMessage(MsgError, ErrorPayload{Message: err.Error()})
		select {
		case client.send <- errMsg:
		default:
		}
		return
	}

	// Send ack to the originating client
	ackMsg, _ := MarshalMessage(MsgAck, AckPayload{Version: version})
	select {
	case client.send <- ackMsg:
	default:
	}

	// Broadcast to all other clients
	serverOp := ServerOpPayload{
		ClientID:  client.id,
		Type:      op.Type,
		LineStart: op.LineStart,
		LineCount: op.LineCount,
		Lines:     op.Lines,
		Version:   version,
	}
	opMsg, _ := MarshalMessage(MsgOp, serverOp)

	s.mu.RLock()
	for _, c := range s.clients {
		if c.id != client.id {
			select {
			case c.send <- opMsg:
			default:
			}
		}
	}
	s.mu.RUnlock()
}

// HandleCursor broadcasts a cursor update from a client.
func (s *Session) HandleCursor(client *Client, cursor *CursorPayload) {
	msg, _ := MarshalMessage(MsgCursor, ServerCursorPayload{
		ClientID: client.id,
		Line:     cursor.Line,
		Ch:       cursor.Ch,
	})

	s.mu.RLock()
	for _, c := range s.clients {
		if c.id != client.id {
			select {
			case c.send <- msg:
			default:
			}
		}
	}
	s.mu.RUnlock()
}

func (s *Session) sendWelcome(client *Client) {
	s.mu.RLock()
	users := make([]UserInfo, 0, len(s.clients))
	for _, c := range s.clients {
		if c.id != client.id {
			users = append(users, UserInfo{
				ID:    c.id,
				Name:  c.name,
				Color: c.color,
				Mode:  c.mode,
			})
		}
	}
	s.mu.RUnlock()

	welcome := WelcomePayload{
		ClientID: client.id,
		Mode:     client.mode,
		Document: s.doc.Lines(),
		Version:  s.doc.Version(),
		Users:    users,
	}
	msg, _ := MarshalMessage(MsgWelcome, welcome)
	select {
	case client.send <- msg:
	default:
	}
}

func (s *Session) broadcastUserJoined(client *Client) {
	msg, _ := MarshalMessage(MsgUserJoined, UserJoinedPayload{
		ID:    client.id,
		Name:  client.name,
		Color: client.color,
		Mode:  client.mode,
	})
	s.mu.RLock()
	for _, c := range s.clients {
		if c.id != client.id {
			select {
			case c.send <- msg:
			default:
			}
		}
	}
	s.mu.RUnlock()
}

func (s *Session) broadcastUserLeft(client *Client) {
	msg, _ := MarshalMessage(MsgUserLeft, UserLeftPayload{ID: client.id})
	s.mu.RLock()
	for _, c := range s.clients {
		select {
		case c.send <- msg:
		default:
		}
	}
	s.mu.RUnlock()
}

// nextColor picks a color for a new client.
func (s *Session) nextColor() string {
	s.mu.RLock()
	idx := len(s.clients) % len(userColors)
	s.mu.RUnlock()
	return userColors[idx]
}

package collab

import (
	"encoding/json"
	"testing"
	"time"
)

func newTestSession() *Session {
	s := &Session{
		ID:            "test-session",
		WriteAccessID: "write123",
		ReadAccessID:  "read1234",
		doc:           NewDocumentFromText("anchor Customer\ncomponent App : III.5"),
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte, 256),
		stopCh:        make(chan struct{}),
		created:       time.Now(),
		lastActive:    time.Now(),
	}
	go s.Run()
	return s
}

func newTestClient(id, name, mode string) *Client {
	return &Client{
		id:    id,
		name:  name,
		color: "#e74c3c",
		mode:  mode,
		send:  make(chan []byte, 256),
	}
}

func TestResolveAccess(t *testing.T) {
	s := &Session{
		WriteAccessID: "wABC1234",
		ReadAccessID:  "rXYZ5678",
	}

	tests := []struct {
		accessID string
		want     string
	}{
		{"wABC1234", "rw"},
		{"rXYZ5678", "ro"},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		got := s.ResolveAccess(tt.accessID)
		if got != tt.want {
			t.Errorf("ResolveAccess(%q) = %q, want %q", tt.accessID, got, tt.want)
		}
	}
}

func TestSessionRegisterAndWelcome(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	client := newTestClient("u1", "Alice", "rw")
	client.session = s

	// Register client
	s.register <- client

	// Should receive welcome message
	select {
	case msg := <-client.send:
		var envelope Message
		if err := json.Unmarshal(msg, &envelope); err != nil {
			t.Fatal(err)
		}
		if envelope.Type != MsgWelcome {
			t.Errorf("expected welcome, got %q", envelope.Type)
		}

		var welcome WelcomePayload
		json.Unmarshal(envelope.Payload, &welcome)
		if welcome.ClientID != "u1" {
			t.Errorf("clientId = %q, want %q", welcome.ClientID, "u1")
		}
		if welcome.Mode != "rw" {
			t.Errorf("mode = %q, want %q", welcome.Mode, "rw")
		}
		if len(welcome.Document) != 2 {
			t.Errorf("document lines = %d, want 2", len(welcome.Document))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for welcome")
	}
}

func TestSessionBroadcastUserJoined(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	// Register first client
	c1 := newTestClient("u1", "Alice", "rw")
	c1.session = s
	s.register <- c1
	<-c1.send // consume welcome

	// Register second client
	c2 := newTestClient("u2", "Bob", "ro")
	c2.session = s
	s.register <- c2
	<-c2.send // consume welcome

	// c1 should receive user_joined for c2
	select {
	case msg := <-c1.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgUserJoined {
			t.Errorf("expected user_joined, got %q", envelope.Type)
		}
		var joined UserJoinedPayload
		json.Unmarshal(envelope.Payload, &joined)
		if joined.ID != "u2" {
			t.Errorf("joined.ID = %q, want %q", joined.ID, "u2")
		}
		if joined.Mode != "ro" {
			t.Errorf("joined.Mode = %q, want %q", joined.Mode, "ro")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for user_joined")
	}
}

func TestSessionUnregister(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	c1 := newTestClient("u1", "Alice", "rw")
	c1.session = s
	s.register <- c1
	<-c1.send // welcome

	c2 := newTestClient("u2", "Bob", "rw")
	c2.session = s
	s.register <- c2
	<-c2.send // welcome
	<-c1.send // user_joined for c2

	// Unregister c2
	s.unregister <- c2

	// c1 should receive user_left
	select {
	case msg := <-c1.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgUserLeft {
			t.Errorf("expected user_left, got %q", envelope.Type)
		}
		var left UserLeftPayload
		json.Unmarshal(envelope.Payload, &left)
		if left.ID != "u2" {
			t.Errorf("left.ID = %q, want %q", left.ID, "u2")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for user_left")
	}
}

func TestSessionHandleOpRW(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	c1 := newTestClient("u1", "Alice", "rw")
	c1.session = s
	s.register <- c1
	<-c1.send // welcome

	c2 := newTestClient("u2", "Bob", "rw")
	c2.session = s
	s.register <- c2
	<-c2.send // welcome
	<-c1.send // user_joined

	// c1 sends an insert operation
	op := &OpPayload{
		Type:      "insert",
		LineStart: 1,
		Lines:     []string{"component API : II.5"},
		Version:   0,
	}
	s.HandleOp(c1, op)

	// c1 should get ack
	select {
	case msg := <-c1.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgAck {
			t.Errorf("expected ack, got %q", envelope.Type)
		}
		var ack AckPayload
		json.Unmarshal(envelope.Payload, &ack)
		if ack.Version != 1 {
			t.Errorf("ack version = %d, want 1", ack.Version)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for ack")
	}

	// c2 should get the op broadcast
	select {
	case msg := <-c2.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgOp {
			t.Errorf("expected op, got %q", envelope.Type)
		}
		var serverOp ServerOpPayload
		json.Unmarshal(envelope.Payload, &serverOp)
		if serverOp.ClientID != "u1" {
			t.Errorf("op clientId = %q, want %q", serverOp.ClientID, "u1")
		}
		if serverOp.Type != "insert" {
			t.Errorf("op type = %q, want %q", serverOp.Type, "insert")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for op broadcast")
	}

	// Verify document was updated
	lines := s.doc.Lines()
	if len(lines) != 3 {
		t.Fatalf("doc lines = %d, want 3", len(lines))
	}
	if lines[1] != "component API : II.5" {
		t.Errorf("line[1] = %q, want %q", lines[1], "component API : II.5")
	}
}

func TestSessionHandleOpRO_Ignored(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	roClient := newTestClient("u1", "Alice", "ro")
	roClient.session = s
	s.register <- roClient
	<-roClient.send // welcome

	initialVersion := s.doc.Version()

	// Read-only client tries to send an op
	op := &OpPayload{
		Type:      "insert",
		LineStart: 0,
		Lines:     []string{"hacked"},
		Version:   0,
	}
	s.HandleOp(roClient, op)

	// Document should be unchanged
	if s.doc.Version() != initialVersion {
		t.Errorf("document version changed for ro client: got %d, want %d", s.doc.Version(), initialVersion)
	}

	// No ack should be sent
	select {
	case msg := <-roClient.send:
		t.Errorf("ro client should not receive messages after op, got: %s", string(msg))
	case <-time.After(100 * time.Millisecond):
		// Expected: no message
	}
}

func TestSessionHandleCursor(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	c1 := newTestClient("u1", "Alice", "rw")
	c1.session = s
	s.register <- c1
	<-c1.send // welcome

	c2 := newTestClient("u2", "Bob", "rw")
	c2.session = s
	s.register <- c2
	<-c2.send // welcome
	<-c1.send // user_joined

	// c1 sends cursor
	s.HandleCursor(c1, &CursorPayload{Line: 5, Ch: 10})

	// c2 should receive cursor
	select {
	case msg := <-c2.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgCursor {
			t.Errorf("expected cursor, got %q", envelope.Type)
		}
		var cursor ServerCursorPayload
		json.Unmarshal(envelope.Payload, &cursor)
		if cursor.ClientID != "u1" {
			t.Errorf("cursor clientId = %q, want %q", cursor.ClientID, "u1")
		}
		if cursor.Line != 5 || cursor.Ch != 10 {
			t.Errorf("cursor pos = (%d,%d), want (5,10)", cursor.Line, cursor.Ch)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for cursor")
	}

	// c1 should NOT receive its own cursor
	select {
	case msg := <-c1.send:
		t.Errorf("sender should not receive own cursor, got: %s", string(msg))
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}

func TestSessionClientCount(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	c1 := newTestClient("u1", "Alice", "rw")
	c1.session = s
	s.register <- c1
	<-c1.send

	c2 := newTestClient("u2", "Bob", "ro")
	c2.session = s
	s.register <- c2
	<-c2.send
	<-c1.send // user_joined

	rw, ro := s.ClientCount()
	if rw != 1 || ro != 1 {
		t.Errorf("ClientCount = (%d rw, %d ro), want (1, 1)", rw, ro)
	}
}

func TestSessionStop(t *testing.T) {
	s := newTestSession()

	c := newTestClient("u1", "Alice", "rw")
	c.session = s
	s.register <- c
	<-c.send // welcome

	s.Stop()

	// c.send should be closed after Stop
	time.Sleep(50 * time.Millisecond)
	_, ok := <-c.send
	if ok {
		t.Error("client send channel should be closed after Stop")
	}
}

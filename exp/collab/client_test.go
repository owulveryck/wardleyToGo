package collab

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// startTestWSServer creates a test HTTP server that upgrades to WebSocket,
// creates a Client, and runs the read/write pumps.
func startTestWSServer(t *testing.T, session *Session, mode string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			t.Logf("accept error: %v", err)
			return
		}

		client := &Client{
			id:      generateID(6),
			name:    "TestUser",
			color:   "#e74c3c",
			mode:    mode,
			session: session,
			conn:    conn,
			send:    make(chan []byte, sendBufLen),
		}
		session.register <- client

		ctx := context.Background()
		go client.WritePump(ctx)
		client.ReadPump(ctx)
	}))
}

func TestClientReadWritePump(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	server := startTestWSServer(t, s, "rw")
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:], nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Should receive welcome
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var msg Message
	json.Unmarshal(data, &msg)
	if msg.Type != MsgWelcome {
		t.Fatalf("expected welcome, got %q", msg.Type)
	}

	// Send a hello
	helloMsg, _ := MarshalMessage(MsgHello, HelloPayload{Name: "TestClient"})
	err = conn.Write(ctx, websocket.MessageText, helloMsg)
	if err != nil {
		t.Fatal(err)
	}

	// Send an op
	opMsg, _ := MarshalMessage(MsgOp, OpPayload{
		Type:      "insert",
		LineStart: 1,
		Lines:     []string{"component Test : III.5"},
		Version:   0,
	})
	err = conn.Write(ctx, websocket.MessageText, opMsg)
	if err != nil {
		t.Fatal(err)
	}

	// Should receive ack
	_, data, err = conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(data, &msg)
	if msg.Type != MsgAck {
		t.Fatalf("expected ack, got %q", msg.Type)
	}

	// Verify document was updated
	lines := s.doc.Lines()
	if len(lines) != 3 {
		t.Fatalf("doc lines = %d, want 3: %v", len(lines), lines)
	}
}

func TestClientDisconnection(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	server := startTestWSServer(t, s, "rw")
	defer server.Close()

	// Connect observer to see user_left
	observer := newTestClient("observer", "Observer", "rw")
	observer.session = s
	s.register <- observer
	<-observer.send // welcome

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:], nil)
	if err != nil {
		t.Fatal(err)
	}

	// Read welcome
	_, _, err = conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Observer gets user_joined
	<-observer.send

	// Close connection
	conn.Close(websocket.StatusNormalClosure, "bye")

	// Observer should receive user_left
	select {
	case msg := <-observer.send:
		var envelope Message
		json.Unmarshal(msg, &envelope)
		if envelope.Type != MsgUserLeft {
			t.Errorf("expected user_left, got %q", envelope.Type)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for user_left after disconnect")
	}
}

func TestClientROCannotSendOps(t *testing.T) {
	s := newTestSession()
	defer s.Stop()

	server := startTestWSServer(t, s, "ro")
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:], nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Read welcome
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var msg Message
	json.Unmarshal(data, &msg)
	if msg.Type != MsgWelcome {
		t.Fatalf("expected welcome, got %q", msg.Type)
	}

	// Verify mode is ro
	var welcome WelcomePayload
	json.Unmarshal(msg.Payload, &welcome)
	if welcome.Mode != "ro" {
		t.Errorf("mode = %q, want %q", welcome.Mode, "ro")
	}

	initialVersion := s.doc.Version()

	// Try to send an op (should be ignored by session)
	opMsg, _ := MarshalMessage(MsgOp, OpPayload{
		Type:      "insert",
		LineStart: 0,
		Lines:     []string{"hacked"},
		Version:   0,
	})
	conn.Write(ctx, websocket.MessageText, opMsg)

	// Wait a bit for the server to process
	time.Sleep(200 * time.Millisecond)

	// Document should be unchanged
	if s.doc.Version() != initialVersion {
		t.Errorf("document version changed for ro client")
	}
}

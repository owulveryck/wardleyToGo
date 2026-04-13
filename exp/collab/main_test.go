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

func TestCORSHeaders(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "http://localhost:8080", "")

	req := httptest.NewRequest("OPTIONS", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("OPTIONS status = %d, want 200", w.Code)
	}
	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:8080" {
		t.Errorf("CORS origin = %q, want %q", origin, "http://localhost:8080")
	}
}

func TestCORSWildcard(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "", "")

	req := httptest.NewRequest("OPTIONS", "/s/fake/fake", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("CORS origin = %q, want %q", origin, "*")
	}
}

func TestWebSocketInvalidSession(t *testing.T) {
	hub := NewHub()
	handler := NewServeMux(hub, "http://localhost:8080", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/nonexistent/fake", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
	if resp != nil && resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestWebSocketInvalidAccessID(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "http://localhost:8080", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/"+s.ID+"/wrongaccess", nil)
	if err == nil {
		t.Fatal("expected error for invalid access ID")
	}
	if resp != nil && resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", resp.StatusCode)
	}
}

func TestWebSocketWriteAccess(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/"+s.ID+"/"+s.WriteAccessID+"?name=Alice", &websocket.DialOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Should receive welcome with mode=rw
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var msg Message
	json.Unmarshal(data, &msg)
	if msg.Type != MsgWelcome {
		t.Fatalf("expected welcome, got %q", msg.Type)
	}
	var welcome WelcomePayload
	json.Unmarshal(msg.Payload, &welcome)
	if welcome.Mode != "rw" {
		t.Errorf("mode = %q, want rw", welcome.Mode)
	}
}

func TestWebSocketReadAccess(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/"+s.ID+"/"+s.ReadAccessID+"?name=Bob", &websocket.DialOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var msg Message
	json.Unmarshal(data, &msg)
	var welcome WelcomePayload
	json.Unmarshal(msg.Payload, &welcome)
	if welcome.Mode != "ro" {
		t.Errorf("mode = %q, want ro", welcome.Mode)
	}
}

func TestFullCollaborationScenario(t *testing.T) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect writer
	writerConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/"+s.ID+"/"+s.WriteAccessID+"?name=Alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer writerConn.Close(websocket.StatusNormalClosure, "")

	// Read welcome
	_, _, err = writerConn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Connect reader
	readerConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:]+"/s/"+s.ID+"/"+s.ReadAccessID+"?name=Bob", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer readerConn.Close(websocket.StatusNormalClosure, "")

	// Reader receives welcome
	_, _, err = readerConn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Writer receives user_joined for reader
	_, data, err := writerConn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var joinMsg Message
	json.Unmarshal(data, &joinMsg)
	if joinMsg.Type != MsgUserJoined {
		t.Fatalf("expected user_joined, got %q", joinMsg.Type)
	}

	// Writer sends an insert
	opMsg, _ := MarshalMessage(MsgOp, OpPayload{
		Type:      "insert",
		LineStart: 0,
		Lines:     []string{"anchor Customer", "component App : III.5"},
		Version:   0,
	})
	writerConn.Write(ctx, websocket.MessageText, opMsg)

	// Writer receives ack
	_, data, err = writerConn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var ackMsg Message
	json.Unmarshal(data, &ackMsg)
	if ackMsg.Type != MsgAck {
		t.Fatalf("expected ack, got %q", ackMsg.Type)
	}

	// Reader receives the op
	_, data, err = readerConn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var opBroadcast Message
	json.Unmarshal(data, &opBroadcast)
	if opBroadcast.Type != MsgOp {
		t.Fatalf("expected op, got %q", opBroadcast.Type)
	}

	var serverOp ServerOpPayload
	json.Unmarshal(opBroadcast.Payload, &serverOp)
	if serverOp.Type != "insert" {
		t.Errorf("op type = %q, want insert", serverOp.Type)
	}
	if len(serverOp.Lines) != 2 {
		t.Errorf("op lines = %d, want 2", len(serverOp.Lines))
	}

	// Verify document state
	// Document starts with [""] (1 empty line), insert at 0 adds 2 lines before it
	lines := s.doc.Lines()
	if len(lines) != 3 {
		t.Fatalf("doc lines = %d, want 3: %v", len(lines), lines)
	}
	if lines[0] != "anchor Customer" {
		t.Errorf("line[0] = %q, want %q", lines[0], "anchor Customer")
	}
	if lines[1] != "component App : III.5" {
		t.Errorf("line[1] = %q, want %q", lines[1], "component App : III.5")
	}
}

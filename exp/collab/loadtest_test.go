package collab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// TestLoadMultipleWriters simulates N concurrent writers over WebSocket
// and measures throughput and correctness.
func TestLoadMultipleWriters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test in short mode")
	}

	const (
		numWriters   = 5
		opsPerWriter = 100
	)

	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect all writers and start reader goroutines to drain messages
	conns := make([]*websocket.Conn, numWriters)
	var drainWg sync.WaitGroup
	for i := 0; i < numWriters; i++ {
		url := fmt.Sprintf("ws%s/s/%s/%s?name=Writer%d", server.URL[4:], s.ID, s.WriteAccessID, i)
		conn, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			t.Fatalf("writer %d: dial error: %v", i, err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		conn.SetReadLimit(1 << 20)
		conns[i] = conn

		// Drain all incoming messages (welcome, user_joined, acks, ops) in background
		drainWg.Add(1)
		go func(c *websocket.Conn) {
			defer drainWg.Done()
			for {
				_, _, err := c.Read(ctx)
				if err != nil {
					return
				}
			}
		}(conn)
	}

	// Give time for all connections to be established and welcomed
	time.Sleep(100 * time.Millisecond)

	// Start writing
	var totalOps atomic.Int64
	start := time.Now()
	var wg sync.WaitGroup

	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(writerIdx int) {
			defer wg.Done()
			conn := conns[writerIdx]
			for i := 0; i < opsPerWriter; i++ {
				op := OpPayload{
					Type:      "insert",
					LineStart: 0,
					Lines:     []string{fmt.Sprintf("writer%d-line%d", writerIdx, i)},
					Version:   0, // let server transform
				}
				msg, _ := MarshalMessage(MsgOp, op)
				err := conn.Write(ctx, websocket.MessageText, msg)
				if err != nil {
					t.Errorf("writer %d op %d: write error: %v", writerIdx, i, err)
					return
				}
				totalOps.Add(1)
			}
		}(w)
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Wait for server to process all ops
	time.Sleep(200 * time.Millisecond)

	total := totalOps.Load()
	opsPerSec := float64(total) / elapsed.Seconds()
	t.Logf("Load test: %d writers, %d total ops in %v (%.0f ops/sec)",
		numWriters, total, elapsed, opsPerSec)

	// Verify document has expected line count
	docLines := s.doc.Lines()
	expectedMinLines := int(total) // at least one line per op
	if len(docLines) < expectedMinLines {
		t.Errorf("document has %d lines, expected at least %d", len(docLines), expectedMinLines)
	}

	t.Logf("Final document: %d lines, version %d", len(docLines), s.doc.Version())

	// Cancel context to stop drain goroutines
	cancel()
	drainWg.Wait()
}

// TestLoadReadersAndWriters simulates a mix of readers and writers.
func TestLoadReadersAndWriters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test in short mode")
	}

	const (
		numWriters   = 3
		numReaders   = 5
		opsPerWriter = 50
	)

	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect writers
	writerConns := make([]*websocket.Conn, numWriters)
	var writerDrainWg sync.WaitGroup
	for i := 0; i < numWriters; i++ {
		url := fmt.Sprintf("ws%s/s/%s/%s?name=Writer%d", server.URL[4:], s.ID, s.WriteAccessID, i)
		conn, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		conn.SetReadLimit(1 << 20)
		writerConns[i] = conn

		// Drain writer messages (welcome, acks, etc.)
		writerDrainWg.Add(1)
		go func(c *websocket.Conn) {
			defer writerDrainWg.Done()
			for {
				_, _, err := c.Read(ctx)
				if err != nil {
					return
				}
			}
		}(conn)
	}

	// Connect readers
	readerConns := make([]*websocket.Conn, numReaders)
	var readerOps [5]atomic.Int64
	var readerWg sync.WaitGroup
	for i := 0; i < numReaders; i++ {
		url := fmt.Sprintf("ws%s/s/%s/%s?name=Reader%d", server.URL[4:], s.ID, s.ReadAccessID, i)
		conn, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		conn.SetReadLimit(1 << 20)
		readerConns[i] = conn

		// Readers count received ops
		readerWg.Add(1)
		go func(idx int, c *websocket.Conn) {
			defer readerWg.Done()
			for {
				_, data, err := c.Read(ctx)
				if err != nil {
					return
				}
				var msg Message
				json.Unmarshal(data, &msg)
				if msg.Type == MsgOp {
					readerOps[idx].Add(1)
				}
			}
		}(i, conn)
	}

	// Give time for all connections to be established
	time.Sleep(100 * time.Millisecond)

	// Writers send ops
	start := time.Now()
	var writeWg sync.WaitGroup
	for w := 0; w < numWriters; w++ {
		writeWg.Add(1)
		go func(idx int) {
			defer writeWg.Done()
			conn := writerConns[idx]
			for i := 0; i < opsPerWriter; i++ {
				op := OpPayload{
					Type:      "insert",
					LineStart: 0,
					Lines:     []string{fmt.Sprintf("w%d-op%d", idx, i)},
					Version:   0,
				}
				msg, _ := MarshalMessage(MsgOp, op)
				conn.Write(ctx, websocket.MessageText, msg)
			}
		}(w)
	}
	writeWg.Wait()
	elapsed := time.Since(start)

	// Wait for readers to receive everything
	time.Sleep(500 * time.Millisecond)

	// Cancel context to stop reader goroutines
	cancel()
	readerWg.Wait()
	writerDrainWg.Wait()

	totalWriterOps := numWriters * opsPerWriter
	t.Logf("Load test: %d writers (%d ops each), %d readers, completed in %v",
		numWriters, opsPerWriter, numReaders, elapsed)

	for i := 0; i < numReaders; i++ {
		received := readerOps[i].Load()
		t.Logf("  Reader %d received %d/%d ops (%.1f%%)",
			i, received, totalWriterOps, float64(received)/float64(totalWriterOps)*100)
		if received < int64(totalWriterOps*80/100) {
			t.Errorf("reader %d received only %d ops (expected >= 80%% of %d)", i, received, totalWriterOps)
		}
	}
}

// BenchmarkWebSocketRoundtrip measures end-to-end latency for a single op.
func BenchmarkWebSocketRoundtrip(b *testing.B) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	handler := NewServeMux(hub, "", "")
	server := httptest.NewServer(handler)
	defer server.Close()

	ctx := context.Background()
	url := fmt.Sprintf("ws%s/s/%s/%s?name=Bench", server.URL[4:], s.ID, s.WriteAccessID)
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Read welcome
	_, _, _ = conn.Read(ctx)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := OpPayload{
			Type:      "insert",
			LineStart: 0,
			Lines:     []string{fmt.Sprintf("bench-line-%d", i)},
			Version:   0,
		}
		msg, _ := MarshalMessage(MsgOp, op)
		conn.Write(ctx, websocket.MessageText, msg)

		// Read ack
		_, _, err := conn.Read(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWebSocketBroadcast measures broadcast throughput with N clients.
func BenchmarkWebSocketBroadcast_2Clients(b *testing.B)  { benchWSBroadcast(b, 2) }
func BenchmarkWebSocketBroadcast_5Clients(b *testing.B)  { benchWSBroadcast(b, 5) }
func BenchmarkWebSocketBroadcast_10Clients(b *testing.B) { benchWSBroadcast(b, 10) }

func benchWSBroadcast(b *testing.B, numClients int) {
	hub := NewHub()
	s, _ := hub.CreateSession()
	mux := http.NewServeMux()

	adminHandler := NewAdminHandler(hub, "", "", nil)
	mux.HandleFunc("GET /{$}", adminHandler.Index)
	mux.HandleFunc("GET /s/{sessionId}/{accessId}", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue("sessionId")
		accessID := r.PathValue("accessId")
		session := hub.GetSession(sessionID)
		if session == nil {
			http.Error(w, "not found", 404)
			return
		}
		mode := session.ResolveAccess(accessID)
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			return
		}
		client := &Client{
			id:      generateID(6),
			name:    r.URL.Query().Get("name"),
			color:   session.nextColor(),
			mode:    mode,
			session: session,
			conn:    conn,
			send:    make(chan []byte, sendBufLen),
		}
		session.register <- client
		ctx := r.Context()
		go client.WritePump(ctx)
		client.ReadPump(ctx)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	ctx := context.Background()

	// Connect all clients
	conns := make([]*websocket.Conn, numClients)
	for i := 0; i < numClients; i++ {
		url := fmt.Sprintf("ws%s/s/%s/%s?name=Client%d", server.URL[4:], s.ID, s.WriteAccessID, i)
		conn, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			b.Fatal(err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		conns[i] = conn
		_, _, _ = conn.Read(ctx) // welcome
	}

	// Wait for all clients to be registered and user_joined messages sent
	time.Sleep(100 * time.Millisecond)

	// Start reader goroutines for all clients except sender
	// Use a dedicated read context so we can stop readers cleanly
	readCtx, readCancel := context.WithCancel(ctx)
	var readWg sync.WaitGroup
	for i := 1; i < numClients; i++ {
		readWg.Add(1)
		go func(conn *websocket.Conn) {
			defer readWg.Done()
			for {
				_, _, err := conn.Read(readCtx)
				if err != nil {
					return
				}
			}
		}(conns[i])
	}

	// Also drain user_joined messages from sender
	// by reading with a short timeout just once
	sender := conns[0]

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := OpPayload{
			Type:      "insert",
			LineStart: 0,
			Lines:     []string{"bench"},
			Version:   0,
		}
		msg, _ := MarshalMessage(MsgOp, op)
		sender.Write(ctx, websocket.MessageText, msg)

		// Read ack
		_, _, err := sender.Read(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	readCancel()
	readWg.Wait()
}

package collab

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func drainChan(ch chan []byte) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

func newBenchSession(docLines int) *Session {
	lines := make([]string, docLines)
	for i := range lines {
		lines[i] = fmt.Sprintf("component Line%d : III.%d", i, i%10)
	}
	s := &Session{
		ID:            "bench-session",
		WriteAccessID: "write",
		ReadAccessID:  "read",
		doc:           &Document{lines: lines},
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

func addBenchClients(s *Session, n int) []*Client {
	clients := make([]*Client, n)
	for i := range clients {
		c := &Client{
			id:    fmt.Sprintf("u%d", i),
			name:  fmt.Sprintf("User%d", i),
			color: userColors[i%len(userColors)],
			mode:  "rw",
			send:  make(chan []byte, 1024),
		}
		c.session = s
		clients[i] = c
		s.register <- c
		<-c.send // consume welcome
	}
	// Drain user_joined messages
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			// Previous clients receive user_joined for the new one
		}
	}
	// Drain all pending messages
	time.Sleep(10 * time.Millisecond)
	for _, c := range clients {
		drainChan(c.send)
	}
	return clients
}

// --- HandleOp benchmarks ---

func BenchmarkHandleOp_1Client(b *testing.B)  { benchHandleOp(b, 1) }
func BenchmarkHandleOp_5Clients(b *testing.B) { benchHandleOp(b, 5) }
func BenchmarkHandleOp_10Clients(b *testing.B) { benchHandleOp(b, 10) }

func benchHandleOp(b *testing.B, numClients int) {
	s := newBenchSession(50)
	defer s.Stop()
	clients := addBenchClients(s, numClients)

	// Drain any remaining messages
	time.Sleep(10 * time.Millisecond)
	for _, c := range clients {
		for len(c.send) > 0 {
			<-c.send
		}
	}

	op := &OpPayload{Type: "replace", LineStart: 25, LineCount: 1, Lines: []string{"updated"}, Version: 0}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op.Version = s.doc.Version()
		s.HandleOp(clients[0], op)
		// Drain ack from sender
		<-clients[0].send
		// Drain op from other clients
		for j := 1; j < numClients; j++ {
			<-clients[j].send
		}
	}
}

// --- HandleCursor benchmarks ---

func BenchmarkHandleCursor_5Clients(b *testing.B)  { benchHandleCursor(b, 5) }
func BenchmarkHandleCursor_10Clients(b *testing.B) { benchHandleCursor(b, 10) }

func benchHandleCursor(b *testing.B, numClients int) {
	s := newBenchSession(50)
	defer s.Stop()
	clients := addBenchClients(s, numClients)

	time.Sleep(10 * time.Millisecond)
	for _, c := range clients {
		for len(c.send) > 0 {
			<-c.send
		}
	}

	cursor := &CursorPayload{Line: 10, Ch: 5}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.HandleCursor(clients[0], cursor)
		// Drain cursor messages from other clients
		for j := 1; j < numClients; j++ {
			<-clients[j].send
		}
	}
}

// --- Concurrent access benchmarks ---

func BenchmarkHandleOp_Concurrent_5Writers(b *testing.B) {
	benchHandleOpConcurrent(b, 5)
}

func BenchmarkHandleOp_Concurrent_10Writers(b *testing.B) {
	benchHandleOpConcurrent(b, 10)
}

func benchHandleOpConcurrent(b *testing.B, numWriters int) {
	s := newBenchSession(100)
	defer s.Stop()
	clients := addBenchClients(s, numWriters)

	time.Sleep(10 * time.Millisecond)
	for _, c := range clients {
		for len(c.send) > 0 {
			<-c.send
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	var wg sync.WaitGroup
	opsPerWriter := b.N / numWriters
	if opsPerWriter == 0 {
		opsPerWriter = 1
	}

	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(clientIdx int) {
			defer wg.Done()
			c := clients[clientIdx]
			for i := 0; i < opsPerWriter; i++ {
				op := &OpPayload{
					Type:      "replace",
					LineStart: (clientIdx*10 + i) % 50,
					LineCount: 1,
					Lines:     []string{fmt.Sprintf("writer%d-op%d", clientIdx, i)},
					Version:   s.doc.Version(),
				}
				s.HandleOp(c, op)
			}
		}(w)
	}
	wg.Wait()

	// Drain all messages
	for _, c := range clients {
		for len(c.send) > 0 {
			<-c.send
		}
	}
}

// --- OT transform under load ---

func BenchmarkHandleOp_WithVersionLag(b *testing.B) {
	s := newBenchSession(50)
	defer s.Stop()
	clients := addBenchClients(s, 2)

	time.Sleep(10 * time.Millisecond)
	for _, c := range clients {
		for len(c.send) > 0 {
			<-c.send
		}
	}

	// Pre-fill history with 50 ops so transforms are exercised
	for i := 0; i < 50; i++ {
		op := &OpPayload{
			Type:      "insert",
			LineStart: i % 30,
			Lines:     []string{fmt.Sprintf("history-%d", i)},
			Version:   s.doc.Version(),
		}
		s.HandleOp(clients[0], op)
		<-clients[0].send // ack
		<-clients[1].send // broadcast
	}

	// Now benchmark with version lag: client sends ops with old version.
	// Use replace to keep document size constant and avoid hitting maxDocumentLines.
	// After OT transforms, a replace may have its LineCount reduced to 0 if a
	// concurrent replace in history covers the same range, causing applyReplace to
	// reject the op. In that case only the sender gets an error, no broadcast is
	// sent. We use non-blocking reads for the broadcast to avoid deadlock.
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a version that is 20 behind — triggers 20 transforms
		op := &OpPayload{
			Type:      "replace",
			LineStart: 10,
			LineCount: 1,
			Lines:     []string{"lagged-edit"},
			Version:   s.doc.Version() - 20,
		}
		s.HandleOp(clients[1], op)
		<-clients[1].send // ack or error
		// Broadcast only arrives on success; op may fail after OT transform
		select {
		case <-clients[0].send:
		default:
		}
	}
}

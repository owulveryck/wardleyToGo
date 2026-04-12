package collab

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
)

// ServerConfig holds the configuration for the collaboration server.
type ServerConfig struct {
	Addr       string
	Frontend   string
	AdminToken string // if set, admin routes require this token; empty = dev mode
}

// RunServer starts the collaboration server. It blocks until the context is cancelled.
func RunServer(ctx context.Context, cfg ServerConfig) error {
	hub := NewHub()

	done := make(chan struct{})
	defer close(done)
	hub.StartCleanupTicker(5*time.Minute, 1*time.Hour, done)

	// Rate limiters
	adminLimiter := NewIPRateLimiter(10, 20)
	wsUpgradeLimiter := NewIPRateLimiter(5, 20)
	sessionCreateLimiter := NewIPRateLimiter(10.0/3600.0, 10)
	connTracker := NewConnTracker()

	adminLimiter.StartCleanup(5*time.Minute, 10*time.Minute, done)
	wsUpgradeLimiter.StartCleanup(5*time.Minute, 10*time.Minute, done)
	sessionCreateLimiter.StartCleanup(5*time.Minute, 10*time.Minute, done)

	mux := http.NewServeMux()

	adminHandler := NewAdminHandler(hub, cfg.Frontend, cfg.AdminToken, sessionCreateLimiter)
	mux.HandleFunc("GET /{$}", rateLimitMiddleware(adminAuth(adminHandler.Index, cfg.AdminToken), adminLimiter))
	mux.HandleFunc("POST /session", rateLimitMiddleware(adminAuth(adminHandler.CreateSession, cfg.AdminToken), adminLimiter))
	mux.HandleFunc("POST /session/{id}/delete", rateLimitMiddleware(adminAuth(adminHandler.DeleteSession, cfg.AdminToken), adminLimiter))

	mux.HandleFunc("GET /s/{sessionId}/{accessId}", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, hub, cfg.Frontend, wsUpgradeLimiter, connTracker)
	})

	handler := securityHeaders(corsMiddleware(mux, cfg.Frontend))

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	log.Printf("Collaboration server listening on %s (frontend: %s)", cfg.Addr, cfg.Frontend)
	return server.ListenAndServe()
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, hub *Hub, frontend string, wsLimiter *IPRateLimiter, connTracker *ConnTracker) {
	ip := clientIP(r)

	// Rate limit WebSocket upgrades
	if wsLimiter != nil {
		allowed, retryAfter := wsLimiter.Allow(ip)
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	}

	sessionID := r.PathValue("sessionId")
	accessID := r.PathValue("accessId")

	session := hub.GetSession(sessionID)
	if session == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	mode := session.ResolveAccess(accessID)
	if mode == "" {
		http.Error(w, "invalid access ID", http.StatusForbidden)
		return
	}

	// Check connection limits
	if connTracker != nil {
		if !connTracker.TryAdd(ip, sessionID, maxConnPerIP, maxConnPerSession) {
			http.Error(w, "too many connections", http.StatusTooManyRequests)
			return
		}
	}

	origins := []string{"*"}
	if frontend != "" {
		origins = []string{frontend}
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: origins,
	})
	if err != nil {
		// Release connection tracker slot if accept fails
		if connTracker != nil {
			connTracker.Remove(ip, sessionID)
		}
		log.Printf("websocket accept: %v", err)
		return
	}

	name := r.URL.Query().Get("name")
	if len(name) > maxNameLen {
		name = name[:maxNameLen]
	}

	client := &Client{
		id:      generateID(6),
		name:    name,
		color:   session.nextColor(),
		mode:    mode,
		session: session,
		conn:    conn,
		send:    make(chan []byte, sendBufLen),
		ip:      ip,
	}

	// Enable per-client message rate limiting when rate limiters are active
	if wsLimiter != nil {
		client.msgBucket = newMsgBucket()
	}

	// Set up connection tracking cleanup
	if connTracker != nil {
		client.onClose = func() {
			connTracker.Remove(ip, sessionID)
		}
	}

	session.register <- client

	ctx := r.Context()
	go client.WritePump(ctx)
	client.ReadPump(ctx)
}

func corsMiddleware(next http.Handler, frontend string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := frontend
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Main is the entry point for the collaboration server binary.
func Main() {
	addr := flag.String("addr", ":8081", "listen address")
	frontend := flag.String("frontend", "http://localhost:8080", "frontend URL for CORS and link generation")
	adminToken := flag.String("admin-token", "", "token for admin authentication (empty = dev mode, no auth)")
	flag.Parse()

	// Normalize frontend: remove trailing slash
	fe := strings.TrimRight(*frontend, "/")

	ctx := context.Background()
	if err := RunServer(ctx, ServerConfig{Addr: *addr, Frontend: fe, AdminToken: *adminToken}); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// NewServeMux creates the HTTP mux for testing purposes.
// Pass adminToken="" for dev mode (no auth). Limiters are nil (no rate limiting in tests).
func NewServeMux(hub *Hub, frontend string, adminToken string) http.Handler {
	mux := http.NewServeMux()

	adminHandler := NewAdminHandler(hub, frontend, adminToken, nil)
	mux.HandleFunc("GET /{$}", adminAuth(adminHandler.Index, adminToken))
	mux.HandleFunc("POST /session", adminAuth(adminHandler.CreateSession, adminToken))
	mux.HandleFunc("POST /session/{id}/delete", adminAuth(adminHandler.DeleteSession, adminToken))

	mux.HandleFunc("GET /s/{sessionId}/{accessId}", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, hub, frontend, nil, nil)
	})

	return securityHeaders(corsMiddleware(mux, frontend))
}

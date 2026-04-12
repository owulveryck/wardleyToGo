package collab

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

var adminTmpl = template.Must(template.New("admin.html").Funcs(template.FuncMap{
	"timeAgo": timeAgo,
	"truncate": func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "..."
	},
}).ParseFS(templateFS, "templates/admin.html"))

// AdminHandler handles the admin UI for session management.
type AdminHandler struct {
	hub                  *Hub
	frontend             string
	adminToken           string
	sessionCreateLimiter *IPRateLimiter
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(hub *Hub, frontend, adminToken string, sessionCreateLimiter *IPRateLimiter) *AdminHandler {
	return &AdminHandler{hub: hub, frontend: frontend, adminToken: adminToken, sessionCreateLimiter: sessionCreateLimiter}
}

// SessionView is the data passed to the admin template for each session.
type SessionView struct {
	ID            string
	WriteAccessID string
	ReadAccessID  string
	WriteURL      string
	ReadURL       string
	RWCount       int
	ROCount       int
	Created       time.Time
}

// Index renders the admin page listing all sessions.
func (h *AdminHandler) Index(w http.ResponseWriter, r *http.Request) {
	sessions := h.hub.Sessions()

	views := make([]SessionView, 0, len(sessions))
	for _, s := range sessions {
		rw, ro := s.ClientCount()
		views = append(views, SessionView{
			ID:            s.ID,
			WriteAccessID: s.WriteAccessID,
			ReadAccessID:  s.ReadAccessID,
			WriteURL:      h.buildWSURL(s.ID, s.WriteAccessID),
			ReadURL:       h.buildWSURL(s.ID, s.ReadAccessID),
			RWCount:       rw,
			ROCount:       ro,
			Created:       s.created,
		})
	}

	sort.Slice(views, func(i, j int) bool {
		return views[i].Created.After(views[j].Created)
	})

	data := struct {
		Sessions   []SessionView
		Frontend   string
		AdminToken string
	}{
		Sessions:   views,
		Frontend:   h.frontend,
		AdminToken: h.adminToken,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	adminTmpl.Execute(w, data)
}

// CreateSession creates a new session and redirects to the admin page.
func (h *AdminHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if h.sessionCreateLimiter != nil {
		ip := clientIP(r)
		if allowed, _ := h.sessionCreateLimiter.Allow(ip); !allowed {
			http.Error(w, "too many sessions created", http.StatusTooManyRequests)
			return
		}
	}

	_, err := h.hub.CreateSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	redirect := "/"
	if h.adminToken != "" {
		redirect = "/?token=" + h.adminToken
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// DeleteSession deletes a session and redirects to the admin page.
func (h *AdminHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.hub.DeleteSession(id)

	redirect := "/"
	if h.adminToken != "" {
		redirect = "/?token=" + h.adminToken
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (h *AdminHandler) buildWSURL(sessionID, accessID string) string {
	// Convert frontend URL scheme to ws/wss
	wsScheme := "ws"
	host := h.frontend
	if len(host) > 8 && host[:8] == "https://" {
		wsScheme = "wss"
		host = host[8:]
	} else if len(host) > 7 && host[:7] == "http://" {
		host = host[7:]
	}
	_ = wsScheme // The WS URL uses the collab server host, not frontend

	// For the WebSocket URL, we use the request host (collab server)
	// but since we don't have the request here, we derive from frontend
	// The admin page will show the full path
	return "/s/" + sessionID + "/" + accessID
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	default:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	}
}

package collab

import "encoding/json"

// Message is the envelope for all WebSocket messages.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// MarshalMessage creates a Message with the given type and payload.
// It builds the JSON in a single pass to avoid double-marshaling.
func MarshalMessage(typ string, payload any) ([]byte, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Build {"type":"...","payload":...} directly to avoid marshaling twice.
	// Pre-allocate: {"type":"" + typ + ","payload":} = 24 + len(typ) + len(p) + 1
	buf := make([]byte, 0, 24+len(typ)+len(p))
	buf = append(buf, `{"type":"`...)
	buf = append(buf, typ...)
	buf = append(buf, `","payload":`...)
	buf = append(buf, p...)
	buf = append(buf, '}')
	return buf, nil
}

// --- Client → Server messages ---

// HelloPayload is sent by a client upon connecting.
type HelloPayload struct {
	Name string `json:"name"`
}

// OpPayload is sent by a client (rw) to edit the document.
type OpPayload struct {
	Type      string   `json:"type"` // "insert", "delete", "replace"
	LineStart int      `json:"lineStart"`
	LineCount int      `json:"lineCount"` // lines affected (delete/replace)
	Lines     []string `json:"lines"`     // new content (insert/replace)
	Version   uint64   `json:"version"`   // client's last known server version
}

// CursorPayload is sent by a client to share cursor position.
type CursorPayload struct {
	Line int `json:"line"`
	Ch   int `json:"ch"`
}

// --- Server → Client messages ---

// UserInfo describes a connected user.
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Mode  string `json:"mode"` // "rw" or "ro"
}

// WelcomePayload is sent to a client upon joining a session.
type WelcomePayload struct {
	ClientID string     `json:"clientId"`
	Mode     string     `json:"mode"` // "rw" or "ro"
	Document []string   `json:"document"`
	Version  uint64     `json:"version"`
	Users    []UserInfo `json:"users"`
}

// ServerOpPayload is an operation broadcast to other clients.
type ServerOpPayload struct {
	ClientID  string   `json:"clientId"`
	Type      string   `json:"type"`
	LineStart int      `json:"lineStart"`
	LineCount int      `json:"lineCount"`
	Lines     []string `json:"lines"`
	Version   uint64   `json:"version"`
}

// AckPayload confirms a client's operation was applied.
type AckPayload struct {
	Version uint64 `json:"version"`
}

// UserJoinedPayload notifies that a user joined the session.
type UserJoinedPayload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Mode  string `json:"mode"`
}

// UserLeftPayload notifies that a user left the session.
type UserLeftPayload struct {
	ID string `json:"id"`
}

// ServerCursorPayload broadcasts a remote cursor position.
type ServerCursorPayload struct {
	ClientID string `json:"clientId"`
	Line     int    `json:"line"`
	Ch       int    `json:"ch"`
}

// FullSyncPayload sends the complete document for resync.
type FullSyncPayload struct {
	Document []string `json:"document"`
	Version  uint64   `json:"version"`
}

// ErrorPayload sends an error message to the client.
type ErrorPayload struct {
	Message string `json:"message"`
}

// Message type constants.
const (
	MsgHello      = "hello"
	MsgOp         = "op"
	MsgCursor     = "cursor"
	MsgWelcome    = "welcome"
	MsgAck        = "ack"
	MsgUserJoined = "user_joined"
	MsgUserLeft   = "user_left"
	MsgFullSync   = "full_sync"
	MsgError      = "error"
)

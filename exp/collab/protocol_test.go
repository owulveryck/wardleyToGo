package collab

import (
	"encoding/json"
	"testing"
)

func TestMarshalMessage(t *testing.T) {
	tests := []struct {
		name    string
		msgType string
		payload any
	}{
		{
			name:    "hello",
			msgType: MsgHello,
			payload: HelloPayload{Name: "Alice"},
		},
		{
			name:    "op insert",
			msgType: MsgOp,
			payload: OpPayload{
				Type:      "insert",
				LineStart: 5,
				LineCount: 0,
				Lines:     []string{"component Foo : III.5"},
				Version:   42,
			},
		},
		{
			name:    "cursor",
			msgType: MsgCursor,
			payload: CursorPayload{Line: 10, Ch: 5},
		},
		{
			name:    "welcome",
			msgType: MsgWelcome,
			payload: WelcomePayload{
				ClientID: "abc123",
				Mode:     "rw",
				Document: []string{"title: Test", "anchor User"},
				Version:  1,
				Users: []UserInfo{
					{ID: "u1", Name: "Bob", Color: "#e74c3c", Mode: "rw"},
				},
			},
		},
		{
			name:    "server op",
			msgType: MsgOp,
			payload: ServerOpPayload{
				ClientID:  "u1",
				Type:      "delete",
				LineStart: 3,
				LineCount: 2,
				Lines:     nil,
				Version:   43,
			},
		},
		{
			name:    "ack",
			msgType: MsgAck,
			payload: AckPayload{Version: 44},
		},
		{
			name:    "user joined",
			msgType: MsgUserJoined,
			payload: UserJoinedPayload{ID: "u2", Name: "Charlie", Color: "#2ecc71", Mode: "ro"},
		},
		{
			name:    "user left",
			msgType: MsgUserLeft,
			payload: UserLeftPayload{ID: "u2"},
		},
		{
			name:    "full sync",
			msgType: MsgFullSync,
			payload: FullSyncPayload{
				Document: []string{"line1", "line2"},
				Version:  10,
			},
		},
		{
			name:    "error",
			msgType: MsgError,
			payload: ErrorPayload{Message: "invalid access"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalMessage(tt.msgType, tt.payload)
			if err != nil {
				t.Fatalf("MarshalMessage: %v", err)
			}

			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				t.Fatalf("Unmarshal envelope: %v", err)
			}
			if msg.Type != tt.msgType {
				t.Errorf("type = %q, want %q", msg.Type, tt.msgType)
			}
			if len(msg.Payload) == 0 {
				t.Error("payload is empty")
			}
		})
	}
}

func TestRoundTripOpPayload(t *testing.T) {
	original := OpPayload{
		Type:      "replace",
		LineStart: 7,
		LineCount: 2,
		Lines:     []string{"new line 1", "new line 2"},
		Version:   99,
	}

	data, err := MarshalMessage(MsgOp, original)
	if err != nil {
		t.Fatal(err)
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatal(err)
	}

	var decoded OpPayload
	if err := json.Unmarshal(msg.Payload, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Type != original.Type {
		t.Errorf("Type = %q, want %q", decoded.Type, original.Type)
	}
	if decoded.LineStart != original.LineStart {
		t.Errorf("LineStart = %d, want %d", decoded.LineStart, original.LineStart)
	}
	if decoded.LineCount != original.LineCount {
		t.Errorf("LineCount = %d, want %d", decoded.LineCount, original.LineCount)
	}
	if decoded.Version != original.Version {
		t.Errorf("Version = %d, want %d", decoded.Version, original.Version)
	}
	if len(decoded.Lines) != len(original.Lines) {
		t.Fatalf("Lines len = %d, want %d", len(decoded.Lines), len(original.Lines))
	}
	for i, l := range decoded.Lines {
		if l != original.Lines[i] {
			t.Errorf("Lines[%d] = %q, want %q", i, l, original.Lines[i])
		}
	}
}

func TestRoundTripWelcomePayload(t *testing.T) {
	original := WelcomePayload{
		ClientID: "client-42",
		Mode:     "ro",
		Document: []string{"title: My Map", "anchor Customer", "component App : III.5"},
		Version:  5,
		Users: []UserInfo{
			{ID: "u1", Name: "Alice", Color: "#e74c3c", Mode: "rw"},
			{ID: "u2", Name: "Bob", Color: "#2ecc71", Mode: "ro"},
		},
	}

	data, err := MarshalMessage(MsgWelcome, original)
	if err != nil {
		t.Fatal(err)
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatal(err)
	}

	var decoded WelcomePayload
	if err := json.Unmarshal(msg.Payload, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ClientID != original.ClientID {
		t.Errorf("ClientID = %q, want %q", decoded.ClientID, original.ClientID)
	}
	if decoded.Mode != original.Mode {
		t.Errorf("Mode = %q, want %q", decoded.Mode, original.Mode)
	}
	if decoded.Version != original.Version {
		t.Errorf("Version = %d, want %d", decoded.Version, original.Version)
	}
	if len(decoded.Document) != len(original.Document) {
		t.Fatalf("Document len = %d, want %d", len(decoded.Document), len(original.Document))
	}
	if len(decoded.Users) != len(original.Users) {
		t.Fatalf("Users len = %d, want %d", len(decoded.Users), len(original.Users))
	}
	if decoded.Users[1].Name != "Bob" {
		t.Errorf("Users[1].Name = %q, want %q", decoded.Users[1].Name, "Bob")
	}
}

func TestRoundTripServerCursorPayload(t *testing.T) {
	original := ServerCursorPayload{ClientID: "u3", Line: 15, Ch: 23}

	data, err := MarshalMessage(MsgCursor, original)
	if err != nil {
		t.Fatal(err)
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatal(err)
	}

	var decoded ServerCursorPayload
	if err := json.Unmarshal(msg.Payload, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ClientID != original.ClientID {
		t.Errorf("ClientID = %q, want %q", decoded.ClientID, original.ClientID)
	}
	if decoded.Line != original.Line || decoded.Ch != original.Ch {
		t.Errorf("position = (%d,%d), want (%d,%d)", decoded.Line, decoded.Ch, original.Line, original.Ch)
	}
}

func TestMessageEnvelopeStructure(t *testing.T) {
	data, err := MarshalMessage(MsgError, ErrorPayload{Message: "test error"})
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's valid JSON with exactly "type" and "payload" keys
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["type"]; !ok {
		t.Error("missing 'type' key")
	}
	if _, ok := raw["payload"]; !ok {
		t.Error("missing 'payload' key")
	}
	if len(raw) != 2 {
		t.Errorf("expected 2 keys, got %d", len(raw))
	}
}

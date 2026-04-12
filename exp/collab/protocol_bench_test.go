package collab

import (
	"encoding/json"
	"testing"
)

func BenchmarkMarshalMessageOp(b *testing.B) {
	payload := OpPayload{
		Type:      "insert",
		LineStart: 42,
		LineCount: 0,
		Lines:     []string{"component App : III.5"},
		Version:   100,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgOp, payload)
	}
}

func BenchmarkMarshalMessageServerOp(b *testing.B) {
	payload := ServerOpPayload{
		ClientID:  "abc123",
		Type:      "replace",
		LineStart: 10,
		LineCount: 1,
		Lines:     []string{"component Updated : IV.2"},
		Version:   200,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgOp, payload)
	}
}

func BenchmarkMarshalMessageAck(b *testing.B) {
	payload := AckPayload{Version: 42}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgAck, payload)
	}
}

func BenchmarkMarshalMessageCursor(b *testing.B) {
	payload := ServerCursorPayload{ClientID: "u1", Line: 15, Ch: 23}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgCursor, payload)
	}
}

func BenchmarkUnmarshalOpMessage(b *testing.B) {
	data, _ := MarshalMessage(MsgOp, OpPayload{
		Type:      "insert",
		LineStart: 42,
		Lines:     []string{"component App : III.5"},
		Version:   100,
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg Message
		json.Unmarshal(data, &msg)
		var op OpPayload
		json.Unmarshal(msg.Payload, &op)
	}
}

func BenchmarkUnmarshalCursorMessage(b *testing.B) {
	data, _ := MarshalMessage(MsgCursor, CursorPayload{Line: 5, Ch: 10})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg Message
		json.Unmarshal(data, &msg)
		var cursor CursorPayload
		json.Unmarshal(msg.Payload, &cursor)
	}
}

func BenchmarkMarshalWelcome_SmallDoc(b *testing.B) {
	payload := WelcomePayload{
		ClientID: "abc",
		Mode:     "rw",
		Document: []string{"anchor Customer", "component App : III.5"},
		Version:  1,
		Users:    []UserInfo{{ID: "u1", Name: "Alice", Color: "#e74c3c", Mode: "rw"}},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgWelcome, payload)
	}
}

func BenchmarkMarshalWelcome_LargeDoc(b *testing.B) {
	doc := make([]string, 200)
	for i := range doc {
		doc[i] = "component Line : III.5"
	}
	users := make([]UserInfo, 10)
	for i := range users {
		users[i] = UserInfo{ID: "u", Name: "User", Color: "#e74c3c", Mode: "rw"}
	}
	payload := WelcomePayload{
		ClientID: "abc",
		Mode:     "rw",
		Document: doc,
		Version:  500,
		Users:    users,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalMessage(MsgWelcome, payload)
	}
}

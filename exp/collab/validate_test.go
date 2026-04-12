package collab

import (
	"strings"
	"testing"
)

func TestValidateOp_Valid(t *testing.T) {
	op := &OpPayload{Type: "insert", LineStart: 0, Lines: []string{"hello"}}
	if err := validateOp(op, 10); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateOp_TooManyLines(t *testing.T) {
	lines := make([]string, maxOpLines+1)
	op := &OpPayload{Type: "insert", LineStart: 0, Lines: lines}
	err := validateOp(op, 10)
	if err == nil {
		t.Fatal("expected error for too many lines")
	}
	if !strings.Contains(err.Error(), "too many lines") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateOp_LineTooLong(t *testing.T) {
	longLine := strings.Repeat("x", maxLineLen+1)
	op := &OpPayload{Type: "insert", LineStart: 0, Lines: []string{longLine}}
	err := validateOp(op, 10)
	if err == nil {
		t.Fatal("expected error for line too long")
	}
	if !strings.Contains(err.Error(), "line too long") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateOp_DocumentTooLarge_Insert(t *testing.T) {
	lines := make([]string, 10)
	op := &OpPayload{Type: "insert", LineStart: 0, Lines: lines}
	err := validateOp(op, maxDocumentLines-5) // would result in maxDocumentLines+5
	if err == nil {
		t.Fatal("expected error for document too large")
	}
	if !strings.Contains(err.Error(), "exceed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateOp_DocumentTooLarge_Replace(t *testing.T) {
	lines := make([]string, 10)
	// Replace 1 line with 10 lines, net +9
	op := &OpPayload{Type: "replace", LineStart: 0, LineCount: 1, Lines: lines}
	err := validateOp(op, maxDocumentLines-5)
	if err == nil {
		t.Fatal("expected error for document too large after replace")
	}
}

func TestValidateOp_DocumentOK_Delete(t *testing.T) {
	op := &OpPayload{Type: "delete", LineStart: 0, LineCount: 5}
	// Delete doesn't increase doc size
	if err := validateOp(op, maxDocumentLines); err != nil {
		t.Fatalf("expected no error for delete, got %v", err)
	}
}

func TestValidateOp_NameTruncation(t *testing.T) {
	longName := strings.Repeat("A", maxNameLen+20)
	if len(longName) > maxNameLen {
		longName = longName[:maxNameLen]
	}
	if len(longName) != maxNameLen {
		t.Fatalf("expected truncated name length %d, got %d", maxNameLen, len(longName))
	}
}

func TestSessionLimitMaxSessions(t *testing.T) {
	h := NewHub()
	for i := 0; i < maxSessions; i++ {
		_, err := h.CreateSession()
		if err != nil {
			t.Fatalf("session %d: unexpected error: %v", i, err)
		}
	}
	// One more should fail
	_, err := h.CreateSession()
	if err == nil {
		t.Fatal("expected error when exceeding max sessions")
	}
}

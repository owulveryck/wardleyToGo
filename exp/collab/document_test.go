package collab

import (
	"strings"
	"testing"
)

func docLines(d *Document) string {
	return strings.Join(d.Lines(), "|")
}

// --- Apply tests ---

func TestApplyInsertAtBeginning(t *testing.T) {
	d := NewDocumentFromText("line0\nline1\nline2")
	err := d.Apply(&OpPayload{Type: "insert", LineStart: 0, Lines: []string{"new"}})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "new|line0|line1|line2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if d.Version() != 1 {
		t.Errorf("version = %d, want 1", d.Version())
	}
}

func TestApplyInsertInMiddle(t *testing.T) {
	d := NewDocumentFromText("line0\nline1\nline2")
	err := d.Apply(&OpPayload{Type: "insert", LineStart: 2, Lines: []string{"newA", "newB"}})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "line0|line1|newA|newB|line2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyInsertAtEnd(t *testing.T) {
	d := NewDocumentFromText("line0\nline1")
	err := d.Apply(&OpPayload{Type: "insert", LineStart: 2, Lines: []string{"line2"}})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "line0|line1|line2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyInsertEmptyLines(t *testing.T) {
	d := NewDocumentFromText("line0")
	err := d.Apply(&OpPayload{Type: "insert", LineStart: 0, Lines: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	if d.Version() != 0 { // no-op, version should not change
		t.Errorf("version = %d, want 0", d.Version())
	}
}

func TestApplyInsertOutOfBounds(t *testing.T) {
	d := NewDocumentFromText("line0")
	err := d.Apply(&OpPayload{Type: "insert", LineStart: 5, Lines: []string{"x"}})
	if err == nil {
		t.Fatal("expected error for out of bounds insert")
	}
}

func TestApplyDeleteSingleLine(t *testing.T) {
	d := NewDocumentFromText("line0\nline1\nline2")
	err := d.Apply(&OpPayload{Type: "delete", LineStart: 1, LineCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "line0|line2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyDeleteMultipleLines(t *testing.T) {
	d := NewDocumentFromText("a\nb\nc\nd\ne")
	err := d.Apply(&OpPayload{Type: "delete", LineStart: 1, LineCount: 3})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "a|e"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyDeleteOutOfBounds(t *testing.T) {
	d := NewDocumentFromText("line0\nline1")
	err := d.Apply(&OpPayload{Type: "delete", LineStart: 1, LineCount: 5})
	if err == nil {
		t.Fatal("expected error for out of bounds delete")
	}
}

func TestApplyReplace(t *testing.T) {
	d := NewDocumentFromText("a\nb\nc\nd")
	err := d.Apply(&OpPayload{Type: "replace", LineStart: 1, LineCount: 2, Lines: []string{"B", "C", "C2"}})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "a|B|C|C2|d"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyReplaceWithFewer(t *testing.T) {
	d := NewDocumentFromText("a\nb\nc\nd")
	err := d.Apply(&OpPayload{Type: "replace", LineStart: 1, LineCount: 2, Lines: []string{"X"}})
	if err != nil {
		t.Fatal(err)
	}
	got := docLines(d)
	want := "a|X|d"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyReplaceOutOfBounds(t *testing.T) {
	d := NewDocumentFromText("a\nb")
	err := d.Apply(&OpPayload{Type: "replace", LineStart: 1, LineCount: 5, Lines: []string{"x"}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyUnknownType(t *testing.T) {
	d := NewDocument()
	err := d.Apply(&OpPayload{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestEmptyDocument(t *testing.T) {
	d := NewDocument()
	// An empty document has 1 empty line (matches CodeMirror model)
	if len(d.Lines()) != 1 {
		t.Errorf("expected 1 line, got %d lines", len(d.Lines()))
	}
	if d.Lines()[0] != "" {
		t.Errorf("expected empty first line, got %q", d.Lines()[0])
	}
	if d.Text() != "" {
		t.Errorf("expected empty text, got %q", d.Text())
	}
	// Replace the empty line (simulates typing in an empty editor)
	err := d.Apply(&OpPayload{Type: "replace", LineStart: 0, LineCount: 1, Lines: []string{"first"}})
	if err != nil {
		t.Fatal(err)
	}
	if d.Text() != "first" {
		t.Errorf("got %q, want %q", d.Text(), "first")
	}
}

func TestNewDocumentFromText(t *testing.T) {
	d := NewDocumentFromText("a\nb\nc")
	lines := d.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestNewDocumentFromTextEmpty(t *testing.T) {
	d := NewDocumentFromText("")
	// An empty text still has 1 empty line (matches editor model)
	if len(d.Lines()) != 1 {
		t.Errorf("expected 1 line for empty text, got %d", len(d.Lines()))
	}
	if d.Lines()[0] != "" {
		t.Errorf("expected empty first line, got %q", d.Lines()[0])
	}
}

func TestEmptyDocumentReplace(t *testing.T) {
	// Simulates the exact scenario: client connects to fresh session, types text
	d := NewDocument()
	err := d.Apply(&OpPayload{Type: "replace", LineStart: 0, LineCount: 1, Lines: []string{"typed text"}})
	if err != nil {
		t.Fatalf("replace on fresh document failed: %v", err)
	}
	if d.Text() != "typed text" {
		t.Errorf("got %q, want %q", d.Text(), "typed text")
	}
}

func TestVersionIncrements(t *testing.T) {
	d := NewDocumentFromText("a\nb\nc")
	if d.Version() != 0 {
		t.Fatalf("initial version = %d, want 0", d.Version())
	}
	d.Apply(&OpPayload{Type: "insert", LineStart: 0, Lines: []string{"x"}})
	if d.Version() != 1 {
		t.Errorf("after insert: version = %d, want 1", d.Version())
	}
	d.Apply(&OpPayload{Type: "delete", LineStart: 0, LineCount: 1})
	if d.Version() != 2 {
		t.Errorf("after delete: version = %d, want 2", d.Version())
	}
	d.Apply(&OpPayload{Type: "replace", LineStart: 0, LineCount: 1, Lines: []string{"y"}})
	if d.Version() != 3 {
		t.Errorf("after replace: version = %d, want 3", d.Version())
	}
}

// --- Transform tests ---

func TestTransformInsertVsInsert_Before(t *testing.T) {
	// op inserts at line 5, concurrent inserts 3 lines at line 2
	op := &OpPayload{Type: "insert", LineStart: 5, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "insert", LineStart: 2, Lines: []string{"a", "b", "c"}}
	result := Transform(op, concurrent)
	if result.LineStart != 8 { // 5 + 3
		t.Errorf("LineStart = %d, want 8", result.LineStart)
	}
}

func TestTransformInsertVsInsert_After(t *testing.T) {
	// op inserts at line 2, concurrent inserts at line 5
	op := &OpPayload{Type: "insert", LineStart: 2, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "insert", LineStart: 5, Lines: []string{"a", "b"}}
	result := Transform(op, concurrent)
	if result.LineStart != 2 { // unchanged
		t.Errorf("LineStart = %d, want 2", result.LineStart)
	}
}

func TestTransformInsertVsInsert_SamePosition(t *testing.T) {
	// Both insert at line 3 — the later one shifts
	op := &OpPayload{Type: "insert", LineStart: 3, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "insert", LineStart: 3, Lines: []string{"a"}}
	result := Transform(op, concurrent)
	if result.LineStart != 4 { // shifted by 1
		t.Errorf("LineStart = %d, want 4", result.LineStart)
	}
}

func TestTransformDeleteVsInsert_Before(t *testing.T) {
	// Delete lines 5-6, insert 2 lines at line 2
	op := &OpPayload{Type: "delete", LineStart: 5, LineCount: 2}
	concurrent := &OpPayload{Type: "insert", LineStart: 2, Lines: []string{"a", "b"}}
	result := Transform(op, concurrent)
	if result.LineStart != 7 { // 5 + 2
		t.Errorf("LineStart = %d, want 7", result.LineStart)
	}
}

func TestTransformDeleteVsInsert_After(t *testing.T) {
	// Delete lines 1-2, insert at line 5
	op := &OpPayload{Type: "delete", LineStart: 1, LineCount: 2}
	concurrent := &OpPayload{Type: "insert", LineStart: 5, Lines: []string{"a"}}
	result := Transform(op, concurrent)
	if result.LineStart != 1 { // unchanged
		t.Errorf("LineStart = %d, want 1", result.LineStart)
	}
}

func TestTransformInsertVsDelete_Before(t *testing.T) {
	// Insert at line 5, delete lines 1-2
	op := &OpPayload{Type: "insert", LineStart: 5, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "delete", LineStart: 1, LineCount: 2}
	result := Transform(op, concurrent)
	if result.LineStart != 3 { // 5 - 2
		t.Errorf("LineStart = %d, want 3", result.LineStart)
	}
}

func TestTransformInsertVsDelete_InsertInDeletedRange(t *testing.T) {
	// Insert at line 3, delete lines 2-5
	op := &OpPayload{Type: "insert", LineStart: 3, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "delete", LineStart: 2, LineCount: 4}
	result := Transform(op, concurrent)
	if result.LineStart != 2 { // moved to deletion point
		t.Errorf("LineStart = %d, want 2", result.LineStart)
	}
}

func TestTransformDeleteVsDelete_NonOverlapping(t *testing.T) {
	// Delete lines 5-6, concurrent deletes lines 1-2
	op := &OpPayload{Type: "delete", LineStart: 5, LineCount: 2}
	concurrent := &OpPayload{Type: "delete", LineStart: 1, LineCount: 2}
	result := Transform(op, concurrent)
	if result.LineStart != 3 { // 5 - 2
		t.Errorf("LineStart = %d, want 3", result.LineStart)
	}
	if result.LineCount != 2 { // unchanged
		t.Errorf("LineCount = %d, want 2", result.LineCount)
	}
}

func TestTransformDeleteVsDelete_Overlapping(t *testing.T) {
	// Delete lines 2-5, concurrent deletes lines 4-7
	op := &OpPayload{Type: "delete", LineStart: 2, LineCount: 4}
	concurrent := &OpPayload{Type: "delete", LineStart: 4, LineCount: 4}
	result := Transform(op, concurrent)
	// Overlap is lines 4-5 (2 lines). Our range shrinks by 2.
	if result.LineCount != 2 { // 4 - 2
		t.Errorf("LineCount = %d, want 2", result.LineCount)
	}
}

func TestTransformReplaceVsInsert(t *testing.T) {
	// Replace lines 5-6, insert 3 lines at line 2
	op := &OpPayload{Type: "replace", LineStart: 5, LineCount: 2, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "insert", LineStart: 2, Lines: []string{"a", "b", "c"}}
	result := Transform(op, concurrent)
	if result.LineStart != 8 { // 5 + 3
		t.Errorf("LineStart = %d, want 8", result.LineStart)
	}
}

func TestTransformInsertVsReplace(t *testing.T) {
	// Insert at line 10, replace lines 3-5 with 1 line (delta = -2)
	op := &OpPayload{Type: "insert", LineStart: 10, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "replace", LineStart: 3, LineCount: 3, Lines: []string{"merged"}}
	result := Transform(op, concurrent)
	if result.LineStart != 8 { // 10 + (1 - 3) = 8
		t.Errorf("LineStart = %d, want 8", result.LineStart)
	}
}

// --- Integration: Apply + Transform scenario ---

func TestConcurrentEditsScenario(t *testing.T) {
	// Two users start from the same document
	doc := NewDocumentFromText("anchor Customer\ncomponent App : III.5\ncomponent DB : IV.5")

	// User A inserts at line 1
	opA := &OpPayload{Type: "insert", LineStart: 1, Lines: []string{"component API : II.5"}, Version: 0}
	// User B inserts at line 2
	opB := &OpPayload{Type: "insert", LineStart: 2, Lines: []string{"component Cache : III.8"}, Version: 0}

	// Server receives A first, applies it
	if err := doc.Apply(opA); err != nil {
		t.Fatal(err)
	}

	// Server transforms B against A, then applies
	transformedB := Transform(opB, opA)
	if err := doc.Apply(transformedB); err != nil {
		t.Fatal(err)
	}

	lines := doc.Lines()
	expected := []string{
		"anchor Customer",
		"component API : II.5",
		"component App : III.5",
		"component Cache : III.8",
		"component DB : IV.5",
	}
	if len(lines) != len(expected) {
		t.Fatalf("got %d lines, want %d:\n%v", len(lines), len(expected), lines)
	}
	for i, l := range lines {
		if l != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, l, expected[i])
		}
	}
}

func TestConcurrentDeleteAndInsert(t *testing.T) {
	doc := NewDocumentFromText("a\nb\nc\nd\ne")

	// User A deletes line 1 ("b")
	opA := &OpPayload{Type: "delete", LineStart: 1, LineCount: 1, Version: 0}
	// User B inserts at line 3
	opB := &OpPayload{Type: "insert", LineStart: 3, Lines: []string{"x"}, Version: 0}

	doc.Apply(opA)
	transformedB := Transform(opB, opA)
	doc.Apply(transformedB)

	got := docLines(doc)
	// After A: a|c|d|e. B wanted to insert at 3 (before "d"), shifted to 2.
	// After B: a|c|x|d|e
	want := "a|c|x|d|e"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLinesReturnsCopy(t *testing.T) {
	d := NewDocumentFromText("a\nb")
	lines := d.Lines()
	lines[0] = "modified"
	if d.Lines()[0] != "a" {
		t.Error("Lines() should return a copy, but modification leaked through")
	}
}

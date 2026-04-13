package collab

import (
	"errors"
	"slices"
	"strings"
	"sync"
)

// Document represents a collaborative text document as a line array.
type Document struct {
	mu      sync.RWMutex
	lines   []string
	version uint64
}

// NewDocument creates an empty document.
// Initializes with one empty line to match editor models (e.g. CodeMirror)
// which always have at least one line even when empty.
func NewDocument() *Document {
	return &Document{lines: []string{""}}
}

// NewDocumentFromText creates a document from source text.
func NewDocumentFromText(text string) *Document {
	var lines []string
	if text == "" {
		lines = []string{""}
	} else {
		lines = strings.Split(text, "\n")
	}
	return &Document{lines: lines}
}

// Lines returns a copy of the document's line array.
func (d *Document) Lines() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cp := make([]string, len(d.lines))
	copy(cp, d.lines)
	return cp
}

// Text returns the document as a single string.
func (d *Document) Text() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return strings.Join(d.lines, "\n")
}

// Version returns the current version.
func (d *Document) Version() uint64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.version
}

// Apply applies an operation to the document and increments the version.
// The caller must hold no lock; Apply acquires its own.
func (d *Document) Apply(op *OpPayload) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.applyLocked(op)
}

func (d *Document) applyLocked(op *OpPayload) error {
	switch op.Type {
	case "insert":
		return d.applyInsert(op)
	case "delete":
		return d.applyDelete(op)
	case "replace":
		return d.applyReplace(op)
	default:
		return errors.New("unknown operation type: " + op.Type)
	}
}

func (d *Document) applyInsert(op *OpPayload) error {
	if op.LineStart < 0 || op.LineStart > len(d.lines) {
		return errors.New("insert: lineStart out of bounds")
	}
	if len(op.Lines) == 0 {
		return nil // no-op
	}

	d.lines = slices.Insert(d.lines, op.LineStart, op.Lines...)
	d.version++
	return nil
}

func (d *Document) applyDelete(op *OpPayload) error {
	if op.LineStart < 0 || op.LineCount <= 0 {
		return errors.New("delete: invalid lineStart or lineCount")
	}
	end := op.LineStart + op.LineCount
	if end > len(d.lines) {
		return errors.New("delete: range exceeds document length")
	}

	d.lines = slices.Delete(d.lines, op.LineStart, end)
	d.version++
	return nil
}

func (d *Document) applyReplace(op *OpPayload) error {
	if op.LineStart < 0 || op.LineCount <= 0 {
		return errors.New("replace: invalid lineStart or lineCount")
	}
	end := op.LineStart + op.LineCount
	if end > len(d.lines) {
		return errors.New("replace: range exceeds document length")
	}

	newCount := len(op.Lines)
	if newCount == op.LineCount {
		// Same number of lines: replace in-place, no allocation
		copy(d.lines[op.LineStart:], op.Lines)
	} else {
		// Different count: delete then insert
		d.lines = slices.Replace(d.lines, op.LineStart, end, op.Lines...)
	}
	d.version++
	return nil
}

// Transform adjusts op's line indices to account for a concurrent operation
// that has already been applied. Returns a new transformed operation.
func Transform(op *OpPayload, concurrent *OpPayload) *OpPayload {
	result := &OpPayload{
		Type:      op.Type,
		LineStart: op.LineStart,
		LineCount: op.LineCount,
		Lines:     op.Lines,
		Version:   op.Version,
	}
	TransformInPlace(result, concurrent)
	return result
}

// TransformInPlace adjusts op's line indices in place without allocating.
func TransformInPlace(op *OpPayload, concurrent *OpPayload) {
	switch concurrent.Type {
	case "insert":
		transformAgainstInsert(op, concurrent)
	case "delete":
		transformAgainstDelete(op, concurrent)
	case "replace":
		transformAgainstReplace(op, concurrent)
	}
}

func transformAgainstInsert(op *OpPayload, ins *OpPayload) {
	insertedCount := len(ins.Lines)
	if insertedCount == 0 {
		return
	}

	switch op.Type {
	case "insert":
		// If our insert is at or after the concurrent insert, shift forward
		if op.LineStart >= ins.LineStart {
			op.LineStart += insertedCount
		}
	case "delete", "replace":
		if op.LineStart >= ins.LineStart {
			// Our range is entirely after the insert
			op.LineStart += insertedCount
		} else if op.LineStart+op.LineCount > ins.LineStart {
			// The insert falls within our range — expand the range
			op.LineCount += insertedCount
		}
	}
}

func transformAgainstDelete(op *OpPayload, del *OpPayload) {
	delEnd := del.LineStart + del.LineCount

	switch op.Type {
	case "insert":
		if op.LineStart >= delEnd {
			op.LineStart -= del.LineCount
		} else if op.LineStart > del.LineStart {
			// Insert point was within the deleted range — move to deletion point
			op.LineStart = del.LineStart
		}
	case "delete", "replace":
		opEnd := op.LineStart + op.LineCount

		if op.LineStart >= delEnd {
			// Our range is entirely after the deletion
			op.LineStart -= del.LineCount
		} else if opEnd <= del.LineStart {
			// Our range is entirely before the deletion — no change
		} else {
			// Ranges overlap — adjust
			newStart := op.LineStart
			if newStart < del.LineStart {
				newStart = op.LineStart
			} else {
				newStart = del.LineStart
			}
			// Reduce lineCount by the overlap
			overlapStart := max(op.LineStart, del.LineStart)
			overlapEnd := min(opEnd, delEnd)
			overlap := overlapEnd - overlapStart
			if overlap < 0 {
				overlap = 0
			}
			op.LineStart = newStart
			op.LineCount -= overlap
			if op.LineCount < 0 {
				op.LineCount = 0
			}
		}
	}
}

func transformAgainstReplace(op *OpPayload, repl *OpPayload) {
	// A replace is conceptually a delete+insert at the same position.
	// The net effect is: lines shift by (len(repl.Lines) - repl.LineCount).
	delta := len(repl.Lines) - repl.LineCount
	replEnd := repl.LineStart + repl.LineCount

	switch op.Type {
	case "insert":
		if op.LineStart >= replEnd {
			op.LineStart += delta
		} else if op.LineStart > repl.LineStart {
			// Insert was within the replaced range — move to end of replacement
			op.LineStart = repl.LineStart + len(repl.Lines)
		}
	case "delete", "replace":
		if op.LineStart >= replEnd {
			op.LineStart += delta
		} else if op.LineStart+op.LineCount <= repl.LineStart {
			// Entirely before — no change
		} else {
			// Overlap: adjust similarly to delete, then apply delta
			opEnd := op.LineStart + op.LineCount
			overlapStart := max(op.LineStart, repl.LineStart)
			overlapEnd := min(opEnd, replEnd)
			overlap := overlapEnd - overlapStart
			if overlap < 0 {
				overlap = 0
			}

			if op.LineStart >= repl.LineStart {
				op.LineStart += delta
				if op.LineStart < repl.LineStart {
					op.LineStart = repl.LineStart
				}
			}
			op.LineCount -= overlap
			if op.LineCount < 0 {
				op.LineCount = 0
			}
		}
	}
}

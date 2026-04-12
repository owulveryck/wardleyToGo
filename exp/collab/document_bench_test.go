package collab

import (
	"fmt"
	"testing"
)

// --- Helpers ---

// makeDoc creates a document with n lines.
func makeDoc(n int) *Document {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("component Line%d : III.%d", i, i%10)
	}
	return &Document{lines: lines}
}

// --- Apply benchmarks ---

func BenchmarkApplyInsert_10Lines(b *testing.B)   { benchApplyInsert(b, 10) }
func BenchmarkApplyInsert_100Lines(b *testing.B)  { benchApplyInsert(b, 100) }
func BenchmarkApplyInsert_1000Lines(b *testing.B) { benchApplyInsert(b, 1000) }

func benchApplyInsert(b *testing.B, docSize int) {
	op := &OpPayload{Type: "insert", LineStart: docSize / 2, Lines: []string{"new line"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(docSize)
		d.Apply(op)
	}
}

func BenchmarkApplyInsertMulti_100Lines(b *testing.B) {
	op := &OpPayload{Type: "insert", LineStart: 50, Lines: []string{"a", "b", "c", "d", "e"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(100)
		d.Apply(op)
	}
}

func BenchmarkApplyDelete_10Lines(b *testing.B)   { benchApplyDelete(b, 10) }
func BenchmarkApplyDelete_100Lines(b *testing.B)  { benchApplyDelete(b, 100) }
func BenchmarkApplyDelete_1000Lines(b *testing.B) { benchApplyDelete(b, 1000) }

func benchApplyDelete(b *testing.B, docSize int) {
	op := &OpPayload{Type: "delete", LineStart: docSize / 2, LineCount: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(docSize)
		d.Apply(op)
	}
}

func BenchmarkApplyReplace_10Lines(b *testing.B)   { benchApplyReplace(b, 10) }
func BenchmarkApplyReplace_100Lines(b *testing.B)  { benchApplyReplace(b, 100) }
func BenchmarkApplyReplace_1000Lines(b *testing.B) { benchApplyReplace(b, 1000) }

func benchApplyReplace(b *testing.B, docSize int) {
	op := &OpPayload{Type: "replace", LineStart: docSize / 2, LineCount: 1, Lines: []string{"replaced line"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(docSize)
		d.Apply(op)
	}
}

// --- Sustained editing: many ops on same document ---

func BenchmarkSustainedEditing_100Ops(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		d := makeDoc(50)
		for j := 0; j < 100; j++ {
			lineCount := len(d.lines)
			pos := j % lineCount
			switch j % 3 {
			case 0:
				d.Apply(&OpPayload{Type: "insert", LineStart: pos, Lines: []string{fmt.Sprintf("inserted %d", j)}})
			case 1:
				if lineCount > 1 {
					d.Apply(&OpPayload{Type: "delete", LineStart: pos, LineCount: 1})
				}
			case 2:
				d.Apply(&OpPayload{Type: "replace", LineStart: pos, LineCount: 1, Lines: []string{fmt.Sprintf("replaced %d", j)}})
			}
		}
	}
}

// --- Transform benchmarks ---

func BenchmarkTransformInsertVsInsert(b *testing.B) {
	op := &OpPayload{Type: "insert", LineStart: 50, Lines: []string{"x"}}
	concurrent := &OpPayload{Type: "insert", LineStart: 20, Lines: []string{"a", "b", "c"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(op, concurrent)
	}
}

func BenchmarkTransformDeleteVsInsert(b *testing.B) {
	op := &OpPayload{Type: "delete", LineStart: 50, LineCount: 3}
	concurrent := &OpPayload{Type: "insert", LineStart: 20, Lines: []string{"a", "b"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(op, concurrent)
	}
}

func BenchmarkTransformDeleteVsDelete(b *testing.B) {
	op := &OpPayload{Type: "delete", LineStart: 10, LineCount: 5}
	concurrent := &OpPayload{Type: "delete", LineStart: 8, LineCount: 4}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(op, concurrent)
	}
}

func BenchmarkTransformReplaceVsReplace(b *testing.B) {
	op := &OpPayload{Type: "replace", LineStart: 20, LineCount: 3, Lines: []string{"x", "y"}}
	concurrent := &OpPayload{Type: "replace", LineStart: 10, LineCount: 2, Lines: []string{"a", "b", "c"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(op, concurrent)
	}
}

// Transform chain: simulate version lag with N transforms
func BenchmarkTransformChain_10(b *testing.B)  { benchTransformChain(b, 10) }
func BenchmarkTransformChain_50(b *testing.B)  { benchTransformChain(b, 50) }
func BenchmarkTransformChain_100(b *testing.B) { benchTransformChain(b, 100) }

func benchTransformChain(b *testing.B, chainLen int) {
	// Pre-build history of concurrent ops
	history := make([]*OpPayload, chainLen)
	for i := range history {
		history[i] = &OpPayload{Type: "insert", LineStart: i * 2, Lines: []string{"line"}}
	}

	op := &OpPayload{Type: "insert", LineStart: 50, Lines: []string{"my line"}, Version: 0}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		current := op
		for _, h := range history {
			current = Transform(current, h)
		}
	}
}

// --- Lines() / Text() copy benchmarks ---

func BenchmarkDocLines_100(b *testing.B)  { benchDocLines(b, 100) }
func BenchmarkDocLines_1000(b *testing.B) { benchDocLines(b, 1000) }

func benchDocLines(b *testing.B, n int) {
	d := makeDoc(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Lines()
	}
}

func BenchmarkDocText_100(b *testing.B)  { benchDocText(b, 100) }
func BenchmarkDocText_1000(b *testing.B) { benchDocText(b, 1000) }

func benchDocText(b *testing.B, n int) {
	d := makeDoc(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Text()
	}
}

// --- In-place insert benchmark (to compare with allocating version) ---

func BenchmarkApplyInsertAtStart_1000Lines(b *testing.B) {
	op := &OpPayload{Type: "insert", LineStart: 0, Lines: []string{"new line"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(1000)
		d.Apply(op)
	}
}

func BenchmarkApplyInsertAtEnd_1000Lines(b *testing.B) {
	op := &OpPayload{Type: "insert", LineStart: 1000, Lines: []string{"new line"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(1000)
		d.Apply(op)
	}
}

func BenchmarkApplyDeleteMiddle_1000Lines(b *testing.B) {
	op := &OpPayload{Type: "delete", LineStart: 500, LineCount: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := makeDoc(1000)
		d.Apply(op)
	}
}

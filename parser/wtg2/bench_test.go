package wtg2

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"
	"testing"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/layout"
)

var testData []byte

func loadTestData(b *testing.B) string {
	b.Helper()
	if testData == nil {
		var err error
		testData, err = os.ReadFile("testdata/example.wtg2")
		if err != nil {
			b.Fatal(err)
		}
	}
	return string(testData)
}

// BenchmarkEndToEnd measures the full WTG2→SVG pipeline matching the WASM hot path.
func BenchmarkEndToEnd(b *testing.B) {
	raw := loadTestData(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, err := NewParser(bytes.NewBufferString(raw))
		if err != nil {
			b.Fatal(err)
		}
		doc, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
		result, err := BuildMap(doc)
		if err != nil {
			b.Fatal(err)
		}
		var buf bytes.Buffer
		imgSize := image.Rect(0, 0, 1300, 900)
		mapSize := image.Rect(30, 50, 1270, 850)
		enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
		if err != nil {
			b.Fatal(err)
		}
		style := svgmap.NewOctoStyle(result.Stages)
		enc.Init(style)
		if err := enc.Encode(result.Map); err != nil {
			b.Fatal(err)
		}
		enc.Close()
		_ = buf.String()
	}
}

// BenchmarkEndToEnd_WithAnnotations measures the full pipeline with annotations enabled.
func BenchmarkEndToEnd_WithAnnotations(b *testing.B) {
	raw := loadTestData(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, err := NewParser(bytes.NewBufferString(raw))
		if err != nil {
			b.Fatal(err)
		}
		doc, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
		result, err := BuildMap(doc)
		if err != nil {
			b.Fatal(err)
		}
		var buf bytes.Buffer
		legendWidth := 0
		if result.Legend && len(result.LegendItems) > 0 {
			legendWidth = svgmap.LegendWidth
		}
		imgSize := image.Rect(0, 0, 1300+legendWidth, 900)
		mapSize := image.Rect(30, 50, 1270, 850)
		enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
		if err != nil {
			b.Fatal(err)
		}
		indicators := svgmap.AllEvolutionIndications()
		if result.Legend && len(result.LegendItems) > 0 {
			indicators = append(indicators, &svgmap.Legend{Items: result.LegendItems})
		}
		style := svgmap.NewOctoStyle(result.Stages, indicators...)
		enc.Init(style)
		if err := enc.Encode(result.Map); err != nil {
			b.Fatal(err)
		}
		enc.Close()
		_ = buf.String()
	}
}

// BenchmarkParse measures only the lexer+parser stage.
func BenchmarkParse(b *testing.B) {
	raw := loadTestData(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, err := NewParser(bytes.NewBufferString(raw))
		if err != nil {
			b.Fatal(err)
		}
		_, err = p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildMap measures builder+layout from a pre-parsed Document.
func BenchmarkBuildMap(b *testing.B) {
	raw := loadTestData(b)
	// Parse once to get a reference Document, then re-parse each iteration
	// because BuildMap may mutate the Document.
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, err := NewParser(bytes.NewBufferString(raw))
		if err != nil {
			b.Fatal(err)
		}
		doc, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
		_, err = BuildMap(doc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEncode measures only the SVG encoding stage from a pre-built Map.
func BenchmarkEncode(b *testing.B) {
	raw := loadTestData(b)
	p, err := NewParser(bytes.NewBufferString(raw))
	if err != nil {
		b.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		b.Fatal(err)
	}
	result, err := BuildMap(doc)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		imgSize := image.Rect(0, 0, 1300, 900)
		mapSize := image.Rect(30, 50, 1270, 850)
		enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
		if err != nil {
			b.Fatal(err)
		}
		style := svgmap.NewOctoStyle(result.Stages)
		enc.Init(style)
		if err := enc.Encode(result.Map); err != nil {
			b.Fatal(err)
		}
		enc.Close()
		_ = buf.String()
	}
}

// BenchmarkLexer measures only the tokenization stage.
func BenchmarkLexer(b *testing.B) {
	raw := loadTestData(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := NewLexer(raw)
		for {
			tok := l.Next()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

// BenchmarkParsePosition measures evolution position parsing (called per component).
func BenchmarkParsePosition(b *testing.B) {
	positions := []string{"III.5", "I.0", "IV.9", "II", "II.7", "I.2", "III.8", "IV.3"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range positions {
			_, _ = ParsePosition(p)
		}
	}
}

// BenchmarkDocumentToLayoutGraph measures the AST-to-layout-graph conversion.
func BenchmarkDocumentToLayoutGraph(b *testing.B) {
	raw := loadTestData(b)
	p, err := NewParser(bytes.NewBufferString(raw))
	if err != nil {
		b.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = documentToLayoutGraph(doc)
	}
}

// BenchmarkSpreadOverlappingEdges measures the O(n^2) edge overlap detection.
func BenchmarkSpreadOverlappingEdges(b *testing.B) {
	raw := loadTestData(b)
	p, err := NewParser(bytes.NewBufferString(raw))
	if err != nil {
		b.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		b.Fatal(err)
	}
	result, err := BuildMap(doc)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spreadOverlappingEdges(result.Map)
	}
}

// BenchmarkComputeAnimationData measures animation metadata computation.
func BenchmarkComputeAnimationData(b *testing.B) {
	raw := loadTestData(b)
	p, err := NewParser(bytes.NewBufferString(raw))
	if err != nil {
		b.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		b.Fatal(err)
	}

	lg := documentToLayoutGraph(doc)
	layouter := layout.New(layout.DefaultOptions())
	yPositions, err := layouter.Layout(lg)
	if err != nil {
		b.Fatal(err)
	}

	// Build node dict and map the same way BuildMap does, but reuse for benchmarking
	result, err := BuildMap(doc)
	if err != nil {
		b.Fatal(err)
	}

	// Re-parse to get a fresh doc (BuildMap may mutate)
	p2, _ := NewParser(bytes.NewBufferString(raw))
	doc2, _ := p2.Parse()
	lg2 := documentToLayoutGraph(doc2)
	_ = yPositions

	// Reconstruct nodeDict
	nodeDict := make(map[string]*nodeEntry)
	evolvedMap := make(map[int64]int64)
	for _, c := range result.Map.Components() {
		nodeDict[fmt.Sprintf("%v", c)] = &nodeEntry{node: c}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = computeAnimationData(lg2, nodeDict, doc2, result.Map, evolvedMap)
	}
}

// BenchmarkEndToEnd_BufferReuse compares buffer reuse (WASM pattern) vs fresh allocation.
func BenchmarkEndToEnd_BufferReuse(b *testing.B) {
	raw := loadTestData(b)

	b.Run("FreshBuffer", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p, err := NewParser(bytes.NewBufferString(raw))
			if err != nil {
				b.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				b.Fatal(err)
			}
			result, err := BuildMap(doc)
			if err != nil {
				b.Fatal(err)
			}
			var buf bytes.Buffer
			imgSize := image.Rect(0, 0, 1300, 900)
			mapSize := image.Rect(30, 50, 1270, 850)
			enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
			if err != nil {
				b.Fatal(err)
			}
			style := svgmap.NewOctoStyle(result.Stages)
			enc.Init(style)
			if err := enc.Encode(result.Map); err != nil {
				b.Fatal(err)
			}
			enc.Close()
			_ = buf.String()
		}
	})

	b.Run("ReusedBuffer", func(b *testing.B) {
		var buf bytes.Buffer
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			p, err := NewParser(bytes.NewBufferString(raw))
			if err != nil {
				b.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				b.Fatal(err)
			}
			result, err := BuildMap(doc)
			if err != nil {
				b.Fatal(err)
			}
			imgSize := image.Rect(0, 0, 1300, 900)
			mapSize := image.Rect(30, 50, 1270, 850)
			enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
			if err != nil {
				b.Fatal(err)
			}
			style := svgmap.NewOctoStyle(result.Stages)
			enc.Init(style)
			if err := enc.Encode(result.Map); err != nil {
				b.Fatal(err)
			}
			enc.Close()
			_ = buf.String()
		}
	})
}

// generateLargeWTG2 creates a synthetic WTG2 document with n components.
func generateLargeWTG2(n int) string {
	var sb strings.Builder
	sb.WriteString("title: Large Benchmark Map\n")
	sb.WriteString("stages: Genesis, Custom, Product, Commodity\n\n")

	// Create 2 anchors
	sb.WriteString("anchor User : III.5\n")
	sb.WriteString("anchor Business : II.5\n\n")

	// Create n components spread across evolution stages
	evolutions := []string{"I.2", "I.5", "I.8", "II.2", "II.5", "II.8", "III.2", "III.5", "III.8", "IV.2", "IV.5"}
	for i := 0; i < n; i++ {
		evo := evolutions[i%len(evolutions)]
		fmt.Fprintf(&sb, "component C%d : %s\n", i, evo)
	}
	sb.WriteString("\n")

	// Create edges: anchors -> first components, then a chain
	sb.WriteString("User -> C0\n")
	sb.WriteString("Business -> C1\n")
	for i := 0; i < n-1; i++ {
		fmt.Fprintf(&sb, "C%d -> C%d\n", i, i+1)
	}
	// Add cross-links for realism
	for i := 0; i < n-3; i += 3 {
		fmt.Fprintf(&sb, "C%d -> C%d\n", i, i+2)
	}

	return sb.String()
}

// BenchmarkEndToEnd_LargeMap measures the full pipeline with a large synthetic map.
func BenchmarkEndToEnd_LargeMap(b *testing.B) {
	for _, size := range []int{20, 50} {
		raw := generateLargeWTG2(size)
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p, err := NewParser(bytes.NewBufferString(raw))
				if err != nil {
					b.Fatal(err)
				}
				doc, err := p.Parse()
				if err != nil {
					b.Fatal(err)
				}
				result, err := BuildMap(doc)
				if err != nil {
					b.Fatal(err)
				}
				var buf bytes.Buffer
				imgSize := image.Rect(0, 0, 1300, 900)
				mapSize := image.Rect(30, 50, 1270, 850)
				enc, err := svgmap.NewEncoder(&buf, imgSize, mapSize)
				if err != nil {
					b.Fatal(err)
				}
				style := svgmap.NewOctoStyle(result.Stages)
				enc.Init(style)
				if err := enc.Encode(result.Map); err != nil {
					b.Fatal(err)
				}
				enc.Close()
				_ = buf.String()
			}
		})
	}
}

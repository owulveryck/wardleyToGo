package wtg2

import (
	"bytes"
	"image"
	"os"
	"testing"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
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

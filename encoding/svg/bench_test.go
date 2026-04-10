package svgmap

import (
	"bytes"
	"image"
	"testing"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

// buildBenchMap creates a realistic Wardley Map for benchmarking without
// importing the parser (which would cause an import cycle).
func buildBenchMap() (*wardleyToGo.Map, []Evolution) {
	m := wardleyToGo.NewMap(0)
	m.Title = "Benchmark Map"

	// Create anchors
	a1 := wardley.NewAnchor(1)
	a1.Label = "User"
	a1.Placement = image.Pt(50, 3)
	_ = m.AddComponent(a1)

	a2 := wardley.NewAnchor(2)
	a2.Label = "Business"
	a2.Placement = image.Pt(70, 3)
	_ = m.AddComponent(a2)

	// Create components
	names := []string{"Web App", "API", "Auth", "Database", "Cache", "CDN",
		"Cloud", "DNS", "Logging", "Monitoring", "Queue", "Worker",
		"Storage", "Config", "Gateway"}
	comps := make([]*wardley.Component, len(names))
	for i, name := range names {
		c := wardley.NewComponent(int64(10 + i))
		c.Label = name
		c.Placement = image.Pt(20+i*5, 10+i*6)
		if c.Placement.X > 99 {
			c.Placement.X = 99
		}
		if c.Placement.Y > 99 {
			c.Placement.Y = 99
		}
		c.Configured = true
		comps[i] = c
		_ = m.AddComponent(c)
	}

	// Create edges: anchor -> first few components, then chain
	edges := [][2]int64{
		{1, 10}, {1, 11}, {2, 12},
		{10, 11}, {10, 12}, {11, 13}, {11, 14},
		{12, 13}, {13, 15}, {13, 16},
		{14, 17}, {14, 18}, {15, 19},
		{16, 20}, {17, 21}, {18, 22},
		{19, 23}, {20, 24}, {21, 22},
	}
	for _, e := range edges {
		from := m.Node(e[0])
		to := m.Node(e[1])
		if from == nil || to == nil {
			continue
		}
		collab := &wardley.Collaboration{
			F:              from,
			T:              to,
			RenderingLayer: 1,
		}
		_ = m.SetCollaboration(collab)
	}

	stages := DefaultEvolution
	return m, stages
}

func BenchmarkGenerateJsData(b *testing.B) {
	m, _ := buildBenchMap()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateJsData(m)
	}
}

func BenchmarkGenerateCSSData(b *testing.B) {
	m, _ := buildBenchMap()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateCSSData(m)
	}
}

func BenchmarkEncoderEncode(b *testing.B) {
	m, stages := buildBenchMap()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		imgSize := image.Rect(0, 0, 1300, 900)
		mapSize := image.Rect(30, 50, 1270, 850)
		enc, err := NewEncoder(&buf, imgSize, mapSize)
		if err != nil {
			b.Fatal(err)
		}
		style := NewOctoStyle(stages)
		enc.Init(style)
		if err := enc.Encode(m); err != nil {
			b.Fatal(err)
		}
		enc.Close()
	}
}

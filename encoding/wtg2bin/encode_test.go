package wtg2bin

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func TestRoundTripEmpty(t *testing.T) {
	doc := &wtg2.Document{}
	doc.Stages = [4]string{"", "", "", ""}
	assertRoundTrip(t, doc)
}

func TestRoundTripMetadataOnly(t *testing.T) {
	doc := &wtg2.Document{
		Title:    "Test Map",
		Date:     "2026-01-15",
		Author:   "Test Author",
		Scope:    "Test Scope",
		Question: "What is the question?",
		Doctrine: "context",
		Stages:   [4]string{"Genesis", "Custom", "Product", "Commodity"},
		Legend:   true,
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripNodeAllFields(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{
				Name:         "Full Component",
				Kind:         wtg2.KindComponent,
				Evolution:    "III.5",
				EvolvedTo:    "IV.2",
				Inertia:      2,
				InertiaKinds: []string{"tech", "human"},
				Type:         "build",
				Asset:        "tech",
				Color:        "#3498DB",
				Visibility:   0.75,
				Cost:         "1.2M€/an",
				Note:         "Key differentiator",
			},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripNodeMinimal(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{
				Name:       "Simple",
				Kind:       wtg2.KindComponent,
				Visibility: -1,
			},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripVisibilityEdgeCases(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "Unset", Kind: wtg2.KindComponent, Visibility: -1},
			{Name: "Zero", Kind: wtg2.KindComponent, Visibility: 0.0},
			{Name: "One", Kind: wtg2.KindComponent, Visibility: 1.0},
			{Name: "Half", Kind: wtg2.KindComponent, Visibility: 0.5},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllEvolutionCompact(t *testing.T) {
	var nodes []*wtg2.NodeDecl
	romans := []string{"I", "II", "III", "IV"}
	for _, r := range romans {
		// Bare roman numeral
		nodes = append(nodes, &wtg2.NodeDecl{
			Name: "Bare " + r, Kind: wtg2.KindComponent, Evolution: r, Visibility: -1,
		})
		// With digits 0-9
		for d := 0; d <= 9; d++ {
			evo := r + "." + string(rune('0'+d))
			nodes = append(nodes, &wtg2.NodeDecl{
				Name: "Evo " + evo, Kind: wtg2.KindComponent, Evolution: evo, Visibility: -1,
			})
		}
	}
	doc := &wtg2.Document{Stages: [4]string{"", "", "", ""}, Nodes: nodes}
	assertRoundTrip(t, doc)
}

func TestRoundTripEvolutionFallback(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "MultiDigit", Kind: wtg2.KindComponent, Evolution: "III.15", Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllNodeKinds(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "Comp", Kind: wtg2.KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "Anch", Kind: wtg2.KindAnchor, Evolution: "II.3", Visibility: -1},
			{Name: "Sub", Kind: wtg2.KindSubmap, Evolution: "IV.1", Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllInertia(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "None", Kind: wtg2.KindComponent, Inertia: 0, Visibility: -1},
			{Name: "Low", Kind: wtg2.KindComponent, Inertia: 1, Visibility: -1},
			{Name: "Med", Kind: wtg2.KindComponent, Inertia: 2, InertiaKinds: []string{"tech"}, Visibility: -1},
			{Name: "High", Kind: wtg2.KindComponent, Inertia: 3, InertiaKinds: []string{"tech", "financial", "human", "relational", "social"}, Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllTypes(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "Regular", Kind: wtg2.KindComponent, Visibility: -1},
			{Name: "Build", Kind: wtg2.KindComponent, Type: "build", Visibility: -1},
			{Name: "Buy", Kind: wtg2.KindComponent, Type: "buy", Visibility: -1},
			{Name: "Outsource", Kind: wtg2.KindComponent, Type: "outsource", Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllAssets(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes: []*wtg2.NodeDecl{
			{Name: "None", Kind: wtg2.KindComponent, Visibility: -1},
			{Name: "Tech", Kind: wtg2.KindComponent, Asset: "tech", Visibility: -1},
			{Name: "Fin", Kind: wtg2.KindComponent, Asset: "financial", Visibility: -1},
			{Name: "Human", Kind: wtg2.KindComponent, Asset: "human", Visibility: -1},
			{Name: "Rel", Kind: wtg2.KindComponent, Asset: "relational", Visibility: -1},
			{Name: "Social", Kind: wtg2.KindComponent, Asset: "social", Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripPipeline(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Pipelines: []*wtg2.PipelineDecl{
			{
				Name: "Engine",
				Members: []*wtg2.PipelineMemberDecl{
					{Name: "Classic", Position: "III.5"},
					{Name: "AI", Position: "II.3"},
					{Name: "Quantum", Position: "I.2"},
				},
			},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripEdges(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Edges: []*wtg2.EdgeDecl{
			{From: "A", To: "B"},
			{From: "C", To: "D", Label: "some label"},
			{From: "E", To: "F", Bidirectional: true},
			{From: "G", To: "H", Label: "bidir label", Bidirectional: true},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripGroups(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Groups: []*wtg2.GroupDecl{
			{Name: "Plain", Members: []string{"A", "B"}},
			{Name: "Colored", Color: "#FF0000", Members: []string{"C"}},
			{Name: "Team", Team: "explorer", Members: []string{"D", "E"}},
			{Name: "Full", Color: "#00FF00", Team: "town-planner", Members: []string{"F"}},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAnnotations(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Annotations: []*wtg2.AnnotationDecl{
			{Kind: "note", Text: "Some note", Target: "CompA"},
			{Kind: "warning", Text: "Danger!", Target: "CompB"},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllSignalTypes(t *testing.T) {
	types := []string{"accelerating", "stagnating", "declining", "co-evolution",
		"red-queen", "commoditization", "network-effects", "economies-of-scale"}
	var signals []*wtg2.SignalDecl
	for _, st := range types {
		signals = append(signals, &wtg2.SignalDecl{Type: st, Target: "Comp"})
	}
	doc := &wtg2.Document{Stages: [4]string{"", "", "", ""}, Signals: signals}
	assertRoundTrip(t, doc)
}

func TestRoundTripAllGameplayTypes(t *testing.T) {
	types := []string{"ILC", "open-source", "land-grab", "embrace-extend",
		"tower-moat", "FUD", "strangler-fig", "signal-distortion"}
	var gps []*wtg2.GameplayDecl
	for _, gt := range types {
		gps = append(gps, &wtg2.GameplayDecl{Type: gt, Target: "Comp"})
	}
	gps = append(gps, &wtg2.GameplayDecl{Type: "ILC", Text: "with text", Target: "Comp2"})
	doc := &wtg2.Document{Stages: [4]string{"", "", "", ""}, Gameplays: gps}
	assertRoundTrip(t, doc)
}

func TestRoundTripFocuses(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Focuses: []*wtg2.FocusDecl{
			{Target: "CompA"},
			{Target: "CompB"},
		},
	}
	assertRoundTrip(t, doc)
}

func TestRoundTripCompressedVsUncompressed(t *testing.T) {
	doc := &wtg2.Document{
		Title:  "Test",
		Stages: [4]string{"", "", "", ""},
		Nodes:  []*wtg2.NodeDecl{{Name: "A", Visibility: -1}},
	}

	compressed, err := EncodeWithOptions(doc, Options{Compress: true})
	if err != nil {
		t.Fatalf("compress encode: %v", err)
	}
	uncompressed, err := EncodeWithOptions(doc, Options{Compress: false})
	if err != nil {
		t.Fatalf("uncompress encode: %v", err)
	}

	doc1, err := Decode(compressed)
	if err != nil {
		t.Fatalf("decode compressed: %v", err)
	}
	doc2, err := Decode(uncompressed)
	if err != nil {
		t.Fatalf("decode uncompressed: %v", err)
	}

	if !reflect.DeepEqual(doc1, doc2) {
		t.Error("compressed and uncompressed decode to different documents")
	}
}

func TestRoundTripExampleFile(t *testing.T) {
	testFile(t, "../../wtg2/example.wtg2")
}

func TestRoundTripFullExampleFile(t *testing.T) {
	testFile(t, "../../wtg2/full_example.wtg2")
}

func testFile(t *testing.T, path string) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("cannot open %s: %v", path, err)
	}
	defer func() { _ = f.Close() }()

	p, err := wtg2.NewParser(f)
	if err != nil {
		t.Fatalf("creating parser: %v", err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	assertRoundTrip(t, doc)
}

func TestSizeComparison(t *testing.T) {
	paths := []string{"../../wtg2/example.wtg2", "../../wtg2/full_example.wtg2"}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			t.Skipf("cannot open %s: %v", path, err)
		}
		source, _ := os.ReadFile(path)
		_ = f.Close()

		f, _ = os.Open(path)
		p, err := wtg2.NewParser(f)
		if err != nil {
			_ = f.Close()
			t.Fatalf("parser: %v", err)
		}
		doc, err := p.Parse()
		f.Close()
		if err != nil {
			t.Fatalf("parse: %v", err)
		}

		uncompressed, err := EncodeWithOptions(doc, Options{Compress: false})
		if err != nil {
			t.Fatalf("encode uncompressed: %v", err)
		}
		compressed, err := Encode(doc)
		if err != nil {
			t.Fatalf("encode compressed: %v", err)
		}

		t.Logf("%s:", path)
		t.Logf("  source:       %5d bytes", len(source))
		t.Logf("  binary:       %5d bytes (%d%%)", len(uncompressed), len(uncompressed)*100/len(source))
		t.Logf("  binary+flate: %5d bytes (%d%%)", len(compressed), len(compressed)*100/len(source))
	}
}

func TestAllDoctrine(t *testing.T) {
	for _, d := range []string{"", "hygiene", "context", "excellence", "evolution"} {
		doc := &wtg2.Document{Doctrine: d, Stages: [4]string{"", "", "", ""}}
		assertRoundTrip(t, doc)
	}
}

func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"too short", []byte{0x57}},
		{"bad magic", []byte{0x00, 0x00, 0x01, 0x00}},
		{"bad version", []byte{0x57, 0x42, 0x99, 0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.data)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestDecodeCorruptCompressed(t *testing.T) {
	data := []byte{magic0, magic1, version, flagCompressed, 0xFF, 0xFE}
	_, err := Decode(data)
	if err == nil {
		t.Error("expected error for corrupt compressed data")
	}
}

func TestRoundTripUnknownDoctrine(t *testing.T) {
	doc := &wtg2.Document{Doctrine: "custom-doctrine", Stages: [4]string{"", "", "", ""}}
	assertRoundTrip(t, doc)
}

func TestRoundTripUnknownEnums(t *testing.T) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
		Nodes:  []*wtg2.NodeDecl{{Name: "X", Type: "custom-type", Asset: "custom-asset", Visibility: -1}},
		Groups: []*wtg2.GroupDecl{{Name: "G", Team: "custom-team", Members: []string{"X"}}},
		Signals: []*wtg2.SignalDecl{{Type: "land-grab", Target: "X"}},
		Gameplays: []*wtg2.GameplayDecl{{Type: "custom-play", Target: "X"}},
		Annotations: []*wtg2.AnnotationDecl{{Kind: "custom-kind", Text: "test", Target: "X"}},
	}
	assertRoundTrip(t, doc)
}

func TestEvolutionEncoding(t *testing.T) {
	tests := []struct {
		input   string
		compact bool
	}{
		{"", true},
		{"I", true}, {"II", true}, {"III", true}, {"IV", true},
		{"I.0", true}, {"I.9", true}, {"IV.0", true}, {"IV.9", true},
		{"III.5", true},
		{"III.15", false},
		{"V.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			b, ok := encodeEvolution(tt.input)
			if ok != tt.compact {
				t.Errorf("encodeEvolution(%q) compact=%v, want %v", tt.input, ok, tt.compact)
			}
			if ok {
				s, valid := decodeEvolution(b)
				if !valid {
					t.Errorf("decodeEvolution(0x%02X) returned invalid", b)
				}
				if s != tt.input {
					t.Errorf("round-trip: %q -> 0x%02X -> %q", tt.input, b, s)
				}
			}
		})
	}
}

func TestLegendFlag(t *testing.T) {
	doc := &wtg2.Document{Legend: true, Stages: [4]string{"", "", "", ""}}
	data, err := Encode(doc)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Legend {
		t.Error("Legend flag not preserved")
	}
}

func TestUnicodeStrings(t *testing.T) {
	doc := &wtg2.Document{
		Title:  "Stratégie 2026 — Données",
		Stages: [4]string{"Genèse", "Sur-mesure", "Produit", "Commodité"},
		Nodes: []*wtg2.NodeDecl{
			{Name: "Moteur de Calcul d'Itinéraire", Visibility: -1},
		},
	}
	assertRoundTrip(t, doc)
}

func assertRoundTrip(t *testing.T, doc *wtg2.Document) {
	t.Helper()

	for _, compress := range []bool{true, false} {
		label := "compressed"
		if !compress {
			label = "uncompressed"
		}
		t.Run(label, func(t *testing.T) {
			data, err := EncodeWithOptions(doc, Options{Compress: compress})
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			got, err := Decode(data)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if !reflect.DeepEqual(doc, got) {
				t.Errorf("round-trip mismatch")
				diff := documentDiff(doc, got)
				if diff != "" {
					t.Log(diff)
				}
			}
		})
	}
}

func documentDiff(a, b *wtg2.Document) string {
	var sb strings.Builder
	if a.Title != b.Title {
		sb.WriteString("Title: " + a.Title + " != " + b.Title + "\n")
	}
	if a.Date != b.Date {
		sb.WriteString("Date: " + a.Date + " != " + b.Date + "\n")
	}
	if a.Author != b.Author {
		sb.WriteString("Author: " + a.Author + " != " + b.Author + "\n")
	}
	if a.Scope != b.Scope {
		sb.WriteString("Scope: " + a.Scope + " != " + b.Scope + "\n")
	}
	if a.Question != b.Question {
		sb.WriteString("Question: " + a.Question + " != " + b.Question + "\n")
	}
	if a.Doctrine != b.Doctrine {
		sb.WriteString("Doctrine: " + a.Doctrine + " != " + b.Doctrine + "\n")
	}
	if a.Stages != b.Stages {
		sb.WriteString("Stages differ\n")
	}
	if a.Legend != b.Legend {
		sb.WriteString("Legend differs\n")
	}
	if len(a.Nodes) != len(b.Nodes) {
		sb.WriteString("Nodes count differs\n")
	}
	if len(a.Pipelines) != len(b.Pipelines) {
		sb.WriteString("Pipelines count differs\n")
	}
	if len(a.Edges) != len(b.Edges) {
		sb.WriteString("Edges count differs\n")
	}
	if len(a.Groups) != len(b.Groups) {
		sb.WriteString("Groups count differs\n")
	}
	if len(a.Annotations) != len(b.Annotations) {
		sb.WriteString("Annotations count differs\n")
	}
	if len(a.Signals) != len(b.Signals) {
		sb.WriteString("Signals count differs\n")
	}
	if len(a.Gameplays) != len(b.Gameplays) {
		sb.WriteString("Gameplays count differs\n")
	}
	if len(a.Focuses) != len(b.Focuses) {
		sb.WriteString("Focuses count differs\n")
	}
	return sb.String()
}

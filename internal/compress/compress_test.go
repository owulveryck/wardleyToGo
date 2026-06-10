package compress

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func parseWTG2(t *testing.T, input string) *wtg2.Document {
	t.Helper()
	p, err := wtg2.NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return doc
}

func roundTrip(t *testing.T, doc *wtg2.Document) *wtg2.Document {
	t.Helper()
	var buf bytes.Buffer
	if err := Compress(doc, &buf); err != nil {
		t.Fatalf("Compress: %v", err)
	}
	encodedSize := buf.Len()

	got, err := Decompress(bufio.NewReader(&buf))
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}
	t.Logf("compressed size: %d bytes", encodedSize)
	return got
}

func TestRoundTripEmpty(t *testing.T) {
	doc := &wtg2.Document{Stages: [4]string{}}
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal: %+v\ngot:      %+v", doc, got)
	}
}

func TestRoundTripMetaOnly(t *testing.T) {
	doc := parseWTG2(t, "title: My Map\nauthor: Alice\ndate: 2026-01-01\n")
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal: %+v\ngot:      %+v", doc, got)
	}
}

func TestRoundTripNodes(t *testing.T) {
	input := `
component CDN : IV.5 (buy)
component Compute : III.2 (outsource)
anchor User : IV.5
component Engine : II.7 !! >> III.5 {
  type: build
  color: #3498DB
  note: Our key differentiator
}
submap Payment : III.6
anchor Admin
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
		for i := range doc.Nodes {
			if i < len(got.Nodes) && !reflect.DeepEqual(doc.Nodes[i], got.Nodes[i]) {
				t.Errorf("Node[%d]:\n  want: %+v\n  got:  %+v", i, *doc.Nodes[i], *got.Nodes[i])
			}
		}
	}
}

func TestRoundTripEdges(t *testing.T) {
	input := `
component A : I.5
component B : II.5
component C : III.5

A -> B
B <-> C
A -[some label]-> C
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
		for i := range doc.Edges {
			if i < len(got.Edges) && !reflect.DeepEqual(doc.Edges[i], got.Edges[i]) {
				t.Errorf("Edge[%d]:\n  want: %+v\n  got:  %+v", i, *doc.Edges[i], *got.Edges[i])
			}
		}
	}
}

func TestRoundTripPipeline(t *testing.T) {
	input := `
component Engine : II.7

pipeline Engine {
  Classic : III.5
  AI : II.3
  Quantum : I.2
}
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal: %+v\ngot:      %+v", doc.Pipelines, got.Pipelines)
	}
}

func TestRoundTripGroups(t *testing.T) {
	input := `
component A : I.5
component B : II.5

group Team Alpha {
  color: #E74C3C
  team: explorer
  A
  B
}
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal groups: %+v\ngot groups:      %+v", doc.Groups, got.Groups)
	}
}

func TestRoundTripAnnotations(t *testing.T) {
	input := `
component A : I.5

note "Important thing" on A
warning "Risk here" on A
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
	}
}

func TestRoundTripSignals(t *testing.T) {
	input := `
component A : I.5
component B : II.5

signal accelerating on A
signal commoditization on B
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
	}
}

func TestRoundTripGameplays(t *testing.T) {
	input := `
component A : I.5

gameplay open-source on A
gameplay ILC "Invest heavily" on A
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
	}
}

func TestRoundTripFocusAndLegend(t *testing.T) {
	input := `
component A : I.5

focus A
legend
`
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch")
	}
}

func TestRoundTripStages(t *testing.T) {
	input := "stages: Genesis, Custom, Product, Commodity\n"
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal: %v\ngot:      %v", doc.Stages, got.Stages)
	}
}

func TestRoundTripDoctrine(t *testing.T) {
	input := "doctrine: hygiene\n"
	doc := parseWTG2(t, input)
	got := roundTrip(t, doc)
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("mismatch:\noriginal doctrine: %q\ngot doctrine:      %q", doc.Doctrine, got.Doctrine)
	}
}

func TestRoundTripExampleFiles(t *testing.T) {
	files := []string{
		"../../wtg2/example.wtg2",
		"../../wtg2/full_example.wtg2",
	}

	for _, path := range files {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Skipf("cannot read %s: %v", path, err)
			}

			doc := parseWTG2(t, string(data))

			var buf bytes.Buffer
			if err := Compress(doc, &buf); err != nil {
				t.Fatalf("Compress: %v", err)
			}
			compressedSize := buf.Len()
			originalSize := len(data)

			got, err := Decompress(bufio.NewReader(&buf))
			if err != nil {
				t.Fatalf("Decompress: %v", err)
			}

			if !reflect.DeepEqual(doc, got) {
				t.Errorf("roundtrip mismatch for %s", path)
				compareNodes(t, doc.Nodes, got.Nodes)
				compareEdges(t, doc.Edges, got.Edges)
			}

			ratio := float64(originalSize) / float64(compressedSize)
			t.Logf("%s: %d bytes → %d bytes (%.1fx compression)", path, originalSize, compressedSize, ratio)
		})
	}
}

func compareNodes(t *testing.T, a, b []*wtg2.NodeDecl) {
	t.Helper()
	if len(a) != len(b) {
		t.Errorf("node count: %d vs %d", len(a), len(b))
		return
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			t.Errorf("Node[%d] %q:\n  want: %+v\n  got:  %+v", i, a[i].Name, *a[i], *b[i])
		}
	}
}

func compareEdges(t *testing.T, a, b []*wtg2.EdgeDecl) {
	t.Helper()
	if len(a) != len(b) {
		t.Errorf("edge count: %d vs %d", len(a), len(b))
		return
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			t.Errorf("Edge[%d]:\n  want: %+v\n  got:  %+v", i, *a[i], *b[i])
		}
	}
}

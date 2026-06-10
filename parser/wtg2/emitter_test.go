package wtg2

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestEmitRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "meta only",
			input: "title: My Map\nauthor: Alice\ndate: 2026-01-01\nscope: test\nquestion: Why?\ndoctrine: hygiene\n",
		},
		{
			name:  "stages",
			input: "stages: Genesis, Custom, Product, Commodity\n",
		},
		{
			name:  "bare anchor",
			input: "anchor User\n",
		},
		{
			name:  "component with evolution",
			input: "component CDN : II.7\n",
		},
		{
			name:  "component with evolution and type",
			input: "component CDN : II.7 (buy)\n",
		},
		{
			name:  "component with evolution type and visibility",
			input: "component CDN : II.7 (buy) @0.9\n",
		},
		{
			name:  "component with move",
			input: "component Engine : II.7 >> III.5\n",
		},
		{
			name:  "component with inertia and move",
			input: "component Engine : II.7 !! >> III.5\n",
		},
		{
			name:  "component with qualified inertia",
			input: "component Engine : II.7 !!(tech,human) >> III.5\n",
		},
		{
			name:  "component with block",
			input: "component Engine : II.7 !! >> III.5 {\n  type: build\n  color: #3498DB\n  note: Our key differentiator\n}\n",
		},
		{
			name:  "submap",
			input: "submap Payment : III.6\n",
		},
		{
			name:  "pipeline",
			input: "pipeline Engine {\n  Classic : III.5\n  AI : II.3\n}\n",
		},
		{
			name:  "simple edge",
			input: "component A : I.5\ncomponent B : II.5\n\nA -> B\n",
		},
		{
			name:  "bidirectional edge",
			input: "component A : I.5\ncomponent B : II.5\n\nA <-> B\n",
		},
		{
			name:  "labeled edge",
			input: "component A : I.5\ncomponent B : II.5\n\nA -[some label]-> B\n",
		},
		{
			name:  "group",
			input: "component A : I.5\n\ngroup Team Alpha {\n  A\n}\n",
		},
		{
			name:  "group with color and team",
			input: "component A : I.5\n\ngroup Team Alpha {\n  color: #E74C3C\n  team: explorer\n  A\n}\n",
		},
		{
			name:  "annotation note",
			input: "component A : I.5\n\nnote \"Important thing\" on A\n",
		},
		{
			name:  "annotation warning",
			input: "component A : I.5\n\nwarning \"Risk here\" on A\n",
		},
		{
			name:  "signal",
			input: "component A : I.5\n\nsignal accelerating on A\n",
		},
		{
			name:  "gameplay",
			input: "component A : I.5\n\ngameplay open-source on A\n",
		},
		{
			name:  "gameplay with text",
			input: "component A : I.5\n\ngameplay ILC \"Invest heavily\" on A\n",
		},
		{
			name:  "focus",
			input: "component A : I.5\n\nfocus A\n",
		},
		{
			name:  "legend",
			input: "component A : I.5\n\nlegend\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("NewParser: %v", err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}

			var buf bytes.Buffer
			if err := Emit(&buf, doc); err != nil {
				t.Fatalf("Emit: %v", err)
			}

			p2, err := NewParser(strings.NewReader(buf.String()))
			if err != nil {
				t.Fatalf("NewParser (re-parse): %v", err)
			}
			doc2, err := p2.Parse()
			if err != nil {
				t.Fatalf("Parse (re-parse): %v\nemitted:\n%s", err, buf.String())
			}

			if !reflect.DeepEqual(doc, doc2) {
				t.Errorf("roundtrip mismatch\noriginal: %+v\nre-parsed: %+v\nemitted:\n%s", doc, doc2, buf.String())
			}
		})
	}
}

func TestEmitRoundTripFiles(t *testing.T) {
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

			p, err := NewParser(strings.NewReader(string(data)))
			if err != nil {
				t.Fatalf("NewParser: %v", err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}

			var buf bytes.Buffer
			if err := Emit(&buf, doc); err != nil {
				t.Fatalf("Emit: %v", err)
			}

			p2, err := NewParser(strings.NewReader(buf.String()))
			if err != nil {
				t.Fatalf("NewParser (re-parse): %v", err)
			}
			doc2, err := p2.Parse()
			if err != nil {
				t.Fatalf("Parse (re-parse): %v\nemitted:\n%s", err, buf.String())
			}

			if !reflect.DeepEqual(doc, doc2) {
				t.Errorf("roundtrip mismatch for %s", path)
				compareDocuments(t, doc, doc2)
			}
		})
	}
}

func compareDocuments(t *testing.T, a, b *Document) {
	t.Helper()
	if a.Title != b.Title {
		t.Errorf("Title: %q vs %q", a.Title, b.Title)
	}
	if a.Date != b.Date {
		t.Errorf("Date: %q vs %q", a.Date, b.Date)
	}
	if a.Author != b.Author {
		t.Errorf("Author: %q vs %q", a.Author, b.Author)
	}
	if a.Scope != b.Scope {
		t.Errorf("Scope: %q vs %q", a.Scope, b.Scope)
	}
	if a.Question != b.Question {
		t.Errorf("Question: %q vs %q", a.Question, b.Question)
	}
	if a.Doctrine != b.Doctrine {
		t.Errorf("Doctrine: %q vs %q", a.Doctrine, b.Doctrine)
	}
	if a.Stages != b.Stages {
		t.Errorf("Stages: %v vs %v", a.Stages, b.Stages)
	}
	if len(a.Nodes) != len(b.Nodes) {
		t.Errorf("Nodes: %d vs %d", len(a.Nodes), len(b.Nodes))
	} else {
		for i := range a.Nodes {
			if !reflect.DeepEqual(a.Nodes[i], b.Nodes[i]) {
				t.Errorf("Node[%d]: %+v vs %+v", i, *a.Nodes[i], *b.Nodes[i])
			}
		}
	}
	if len(a.Edges) != len(b.Edges) {
		t.Errorf("Edges: %d vs %d", len(a.Edges), len(b.Edges))
	} else {
		for i := range a.Edges {
			if !reflect.DeepEqual(a.Edges[i], b.Edges[i]) {
				t.Errorf("Edge[%d]: %+v vs %+v", i, *a.Edges[i], *b.Edges[i])
			}
		}
	}
	if a.Legend != b.Legend {
		t.Errorf("Legend: %v vs %v", a.Legend, b.Legend)
	}
}

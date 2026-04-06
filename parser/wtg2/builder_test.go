package wtg2

import (
	"os"
	"testing"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestBuildMap_SimpleChain(t *testing.T) {
	doc := &Document{
		Title:  "Test Map",
		Stages: [4]string{"Genesis", "Custom", "Product", "Commodity"},
		Nodes: []*NodeDecl{
			{Name: "User", Kind: KindAnchor, Visibility: -1},
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "DB", Kind: KindComponent, Evolution: "IV.3", Type: "buy", Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "User", To: "App"},
			{From: "App", To: "DB"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Map.Title != "Test Map" {
		t.Errorf("title = %q, want %q", result.Map.Title, "Test Map")
	}

	// Count nodes
	nodes := result.Map.Nodes()
	count := 0
	for nodes.Next() {
		count++
	}
	if count != 3 {
		t.Errorf("node count = %d, want 3", count)
	}

	// Count edges
	edges := result.Map.Edges()
	edgeCount := 0
	for edges.Next() {
		edgeCount++
	}
	if edgeCount != 2 {
		t.Errorf("edge count = %d, want 2", edgeCount)
	}

	// Check stages
	if len(result.Stages) != 4 {
		t.Fatalf("stages count = %d, want 4", len(result.Stages))
	}
	if result.Stages[0].Label != "Genesis" {
		t.Errorf("stage[0] = %q, want %q", result.Stages[0].Label, "Genesis")
	}
}

func TestBuildMap_WithEvolution(t *testing.T) {
	doc := &Document{
		Title:  "Evolution Test",
		Stages: [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", EvolvedTo: "III.5", Inertia: 2, Visibility: -1},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	// Should have 2 nodes: original + evolved
	nodes := result.Map.Nodes()
	count := 0
	for nodes.Next() {
		count++
	}
	if count != 2 {
		t.Errorf("node count = %d, want 2 (original + evolved)", count)
	}

	// Should have 1 evolution edge
	edges := result.Map.Edges()
	edgeCount := 0
	for edges.Next() {
		e := edges.Edge()
		collab, ok := e.(*wardley.Collaboration)
		if ok && collab.Type == wardley.EvolvedComponentEdge {
			edgeCount++
		}
	}
	if edgeCount != 1 {
		t.Errorf("evolution edge count = %d, want 1", edgeCount)
	}
}

func TestBuildMap_WithPipeline(t *testing.T) {
	doc := &Document{
		Title:  "Pipeline Test",
		Stages: [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", Visibility: -1},
		},
		Pipelines: []*PipelineDecl{
			{
				Name: "Engine",
				Members: []*PipelineMemberDecl{
					{Name: "Algo A", Position: "III.5"},
					{Name: "Algo B", Position: "II.3"},
				},
			},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	// Should have 3 nodes: pipeline parent + 2 members
	nodes := result.Map.Nodes()
	count := 0
	for nodes.Next() {
		count++
	}
	if count != 3 {
		t.Errorf("node count = %d, want 3", count)
	}
}

func TestBuildMap_ExampleFile(t *testing.T) {
	f, err := os.Open("testdata/example.wtg2")
	if err != nil {
		f, err = os.Open("testdata/example.wtg2")
		if err != nil {
			t.Skip("example.wtg2 not found")
		}
	}
	defer f.Close()

	p, err := NewParser(f)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Map.Title != "Plateforme de Navigation — Stratégie 2026" {
		t.Errorf("title = %q", result.Map.Title)
	}

	// Count nodes
	nodes := result.Map.Nodes()
	count := 0
	for nodes.Next() {
		count++
	}
	if count < 15 {
		t.Errorf("node count = %d, want >= 15", count)
	}

	// Count edges
	edges := result.Map.Edges()
	edgeCount := 0
	for edges.Next() {
		edgeCount++
	}
	if edgeCount < 10 {
		t.Errorf("edge count = %d, want >= 10", edgeCount)
	}
}

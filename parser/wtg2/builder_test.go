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
	count := len(result.Map.Components())
	if count != 3 {
		t.Errorf("node count = %d, want 3", count)
	}

	// Count edges
	edgeCount := len(result.Map.Collaborations())
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
	count := len(result.Map.Components())
	if count != 2 {
		t.Errorf("node count = %d, want 2 (original + evolved)", count)
	}

	// Should have 1 evolution edge
	edgeCount := 0
	for _, c := range result.Map.Collaborations() {
		collab, ok := c.(*wardley.Collaboration)
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
	count := len(result.Map.Components())
	if count != 3 {
		t.Errorf("node count = %d, want 3", count)
	}
}

func TestBuildMap_SignalsAnnotationsGameplays(t *testing.T) {
	doc := &Document{
		Title:    "Enriched Test",
		Doctrine: "context",
		Stages:   [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1, Asset: "tech", Cost: "500k/year"},
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", EvolvedTo: "III.5", Inertia: 2, InertiaKinds: []string{"tech", "human"}, Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "App", To: "Engine"},
		},
		Signals: []*SignalDecl{
			{Type: "accelerating", Target: "App"},
			{Type: "co-evolution", Target: "Engine"},
		},
		Annotations: []*AnnotationDecl{
			{Kind: "warning", Text: "SPOF", Target: "Engine"},
			{Kind: "note", Text: "Migration planned", Target: "App"},
		},
		Gameplays: []*GameplayDecl{
			{Type: "ILC", Target: "App"},
			{Type: "strangler-fig", Text: "Replace legacy", Target: "Engine"},
		},
		Groups: []*GroupDecl{
			{Name: "Dev Team", Team: "explorer", Members: []string{"App", "Engine"}},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	// Check signals are attached to components
	for _, n := range result.Map.Components() {
		comp, ok := n.(*wardley.Component)
		if !ok {
			continue
		}
		switch comp.Label {
		case "App":
			if len(comp.Signals) != 1 || comp.Signals[0].Type != "accelerating" {
				t.Errorf("App signals = %v, want [accelerating]", comp.Signals)
			}
			if len(comp.Annotations) != 1 || comp.Annotations[0].Kind != "note" {
				t.Errorf("App annotations = %v, want [note]", comp.Annotations)
			}
			if len(comp.Gameplays) != 1 || comp.Gameplays[0].Type != "ILC" {
				t.Errorf("App gameplays = %v, want [ILC]", comp.Gameplays)
			}
			if comp.Asset != "tech" {
				t.Errorf("App asset = %q, want %q", comp.Asset, "tech")
			}
			if comp.Cost != "500k/year" {
				t.Errorf("App cost = %q, want %q", comp.Cost, "500k/year")
			}
		case "Engine":
			if len(comp.Signals) != 1 || comp.Signals[0].Type != "co-evolution" {
				t.Errorf("Engine signals = %v, want [co-evolution]", comp.Signals)
			}
			if len(comp.Annotations) != 1 || comp.Annotations[0].Kind != "warning" {
				t.Errorf("Engine annotations = %v, want [warning]", comp.Annotations)
			}
			if len(comp.Gameplays) != 1 || comp.Gameplays[0].Type != "strangler-fig" {
				t.Errorf("Engine gameplays = %v, want [strangler-fig]", comp.Gameplays)
			}
			if comp.InertiaKinds[0] != "tech" || comp.InertiaKinds[1] != "human" {
				t.Errorf("Engine inertia kinds = %v, want [tech human]", comp.InertiaKinds)
			}
		}
	}

	// Check evolution edge has InertiaKinds
	for _, c := range result.Map.Collaborations() {
		collab, ok := c.(*wardley.Collaboration)
		if !ok {
			continue
		}
		if collab.Type == wardley.EvolvedComponentEdge {
			if len(collab.InertiaKinds) != 2 {
				t.Errorf("evolution edge InertiaKinds = %v, want [tech human]", collab.InertiaKinds)
			}
		}
	}

	// Check group has team type
	for _, n := range result.Map.Components() {
		group, ok := n.(*wardley.Group)
		if ok && group.Label == "Dev Team" {
			if group.TeamType != "explorer" {
				t.Errorf("group TeamType = %q, want %q", group.TeamType, "explorer")
			}
		}
	}
}

func TestBuildMap_LegendItems(t *testing.T) {
	doc := &Document{
		Title:  "Legend Test",
		Stages: [4]string{"I", "II", "III", "IV"},
		Legend: true,
		Nodes: []*NodeDecl{
			{Name: "User", Kind: KindAnchor, Visibility: -1},
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "DB", Kind: KindComponent, Evolution: "IV.3", Type: "buy", Visibility: -1},
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", EvolvedTo: "III.5", Inertia: 2, Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "User", To: "App"},
		},
		Groups: []*GroupDecl{
			{Name: "Team", Members: []string{"App"}},
		},
		Signals: []*SignalDecl{
			{Type: "accelerating", Target: "App"},
		},
		Gameplays: []*GameplayDecl{
			{Type: "ILC", Target: "App"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Legend {
		t.Error("expected Legend=true")
	}
	if len(result.LegendItems) == 0 {
		t.Fatal("expected LegendItems to be populated")
	}

	// Check expected categories are present
	categories := make(map[string]bool)
	types := make(map[string]bool)
	for _, item := range result.LegendItems {
		categories[item.Category] = true
		types[item.Type] = true
	}

	for _, want := range []string{"Components", "Edges", "Signals", "Gameplays", "Other"} {
		if !categories[want] {
			t.Errorf("missing category %q", want)
		}
	}
	for _, want := range []string{"component", "buy", "evolved", "edge", "evolved_edge", "inertia", "group", "signal", "gameplay"} {
		if !types[want] {
			t.Errorf("missing legend type %q", want)
		}
	}
}

func TestBuildMap_LegendDisabled(t *testing.T) {
	doc := &Document{
		Title:  "No Legend",
		Stages: [4]string{"I", "II", "III", "IV"},
		Legend: false,
		Nodes: []*NodeDecl{
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Legend {
		t.Error("expected Legend=false")
	}
	if len(result.LegendItems) != 0 {
		t.Errorf("expected no LegendItems, got %d", len(result.LegendItems))
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
	defer func() { _ = f.Close() }()

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
	count := len(result.Map.Components())
	if count < 15 {
		t.Errorf("node count = %d, want >= 15", count)
	}

	// Count edges
	edgeCount := len(result.Map.Collaborations())
	if edgeCount < 10 {
		t.Errorf("edge count = %d, want >= 10", edgeCount)
	}
}

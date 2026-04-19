package wtg2

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	wardleyToGo "github.com/owulveryck/wardleyToGo"
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
	if result.Stages[0].Label != "I" {
		t.Errorf("stage[0].Label = %q, want %q", result.Stages[0].Label, "I")
	}
	if result.Stages[0].ZoneLabel != "Genesis" {
		t.Errorf("stage[0].ZoneLabel = %q, want %q", result.Stages[0].ZoneLabel, "Genesis")
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

	for _, want := range []string{"Components", "Edges", "Groups", "Signals", "Gameplays"} {
		if !categories[want] {
			t.Errorf("missing category %q", want)
		}
	}
	for _, want := range []string{"component", "buy", "evolved", "edge", "evolved_edge", "inertia", "group", "signal_accelerating", "gameplay"} {
		if !types[want] {
			t.Errorf("missing legend type %q", want)
		}
	}
}

func TestBuildMap_LegendQualifiedInertia(t *testing.T) {
	doc := &Document{
		Title:  "Qualified Inertia Legend",
		Stages: [4]string{"I", "II", "III", "IV"},
		Legend: true,
		Nodes: []*NodeDecl{
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", EvolvedTo: "III.5", Inertia: 2, InertiaKinds: []string{"tech", "human"}, Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "App", To: "Engine"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	types := make(map[string]bool)
	var techColor, humanColor bool
	for _, item := range result.LegendItems {
		types[item.Type] = true
		if item.Type == "inertia_tech" && item.Color != nil {
			techColor = true
		}
		if item.Type == "inertia_human" && item.Color != nil {
			humanColor = true
		}
	}

	if types["inertia"] {
		t.Error("unqualified 'inertia' should not appear when only qualified kinds are used")
	}
	for _, want := range []string{"inertia_tech", "inertia_human"} {
		if !types[want] {
			t.Errorf("missing legend type %q", want)
		}
	}
	if !techColor {
		t.Error("inertia_tech should have a non-nil Color")
	}
	if !humanColor {
		t.Error("inertia_human should have a non-nil Color")
	}
}

func TestBuildMap_LegendMixedInertia(t *testing.T) {
	doc := &Document{
		Title:  "Mixed Inertia Legend",
		Stages: [4]string{"I", "II", "III", "IV"},
		Legend: true,
		Nodes: []*NodeDecl{
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Inertia: 1, Visibility: -1},
			{Name: "Engine", Kind: KindComponent, Evolution: "II.7", EvolvedTo: "III.5", Inertia: 2, InertiaKinds: []string{"financial"}, Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "App", To: "Engine"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	types := make(map[string]bool)
	for _, item := range result.LegendItems {
		types[item.Type] = true
	}

	if !types["inertia"] {
		t.Error("unqualified 'inertia' should appear for App")
	}
	if !types["inertia_financial"] {
		t.Error("'inertia_financial' should appear for Engine")
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
	// LegendItems are always pre-computed (so callers can enable legend
	// via CLI flags), but Legend=false means they should not be displayed.
	if len(result.LegendItems) == 0 {
		t.Error("expected LegendItems to be pre-computed even when Legend=false")
	}
}

func TestBuildMap_WithFocus(t *testing.T) {
	doc := &Document{
		Title:  "Focus Test",
		Stages: [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "User", Kind: KindAnchor, Visibility: -1},
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "API", Kind: KindComponent, Evolution: "II.5", Visibility: -1},
			{Name: "DB", Kind: KindComponent, Evolution: "IV.2", Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "User", To: "App"},
			{From: "App", To: "API"},
			{From: "API", To: "DB"},
		},
		Focuses: []*FocusDecl{
			{Target: "App"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Focus == nil {
		t.Fatal("expected Focus to be set")
	}

	// App, API, DB should be focused (App + descendants)
	// User should NOT be focused
	focusedLabels := make(map[string]bool)
	for _, c := range result.Map.Components() {
		if result.Focus.ComponentIDs[c.ID()] {
			if comp, ok := c.(*wardley.Component); ok {
				focusedLabels[comp.Label] = true
			}
			if anchor, ok := c.(*wardley.Anchor); ok {
				focusedLabels[anchor.Label] = true
			}
		}
	}

	if !focusedLabels["App"] {
		t.Error("App should be focused")
	}
	if !focusedLabels["API"] {
		t.Error("API should be focused (descendant of App)")
	}
	if !focusedLabels["DB"] {
		t.Error("DB should be focused (descendant of App via API)")
	}
	if focusedLabels["User"] {
		t.Error("User should NOT be focused")
	}

	// Edges App->API and API->DB should be focused, User->App should NOT
	if len(result.Focus.EdgeKeys) != 2 {
		t.Errorf("focused edge count = %d, want 2", len(result.Focus.EdgeKeys))
	}
}

func TestBuildMap_WithFocusUnknownTarget(t *testing.T) {
	doc := &Document{
		Title:  "Focus Unknown",
		Stages: [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
		},
		Focuses: []*FocusDecl{
			{Target: "DoesNotExist"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Focus != nil {
		t.Error("expected Focus to be nil for unknown target")
	}
}

func TestBuildMap_WithFocusGroup(t *testing.T) {
	doc := &Document{
		Title:  "Focus Group Test",
		Stages: [4]string{"I", "II", "III", "IV"},
		Nodes: []*NodeDecl{
			{Name: "User", Kind: KindAnchor, Visibility: -1},
			{Name: "App", Kind: KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "DB", Kind: KindComponent, Evolution: "IV.2", Visibility: -1},
		},
		Edges: []*EdgeDecl{
			{From: "User", To: "App"},
			{From: "App", To: "DB"},
		},
		Groups: []*GroupDecl{
			{Name: "Backend", Members: []string{"App", "DB"}},
		},
		Focuses: []*FocusDecl{
			{Target: "App"},
		},
	}

	result, err := BuildMap(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result.Focus == nil {
		t.Fatal("expected Focus to be set")
	}

	// Group "Backend" contains App which is focused
	if len(result.Focus.GroupIDs) != 1 {
		t.Errorf("focused group count = %d, want 1", len(result.Focus.GroupIDs))
	}
}

func TestBuildMap_FocusDoesNotChangePositions(t *testing.T) {
	makeDoc := func(withFocus bool) *Document {
		doc := &Document{
			Title:  "Focus Position Test",
			Stages: [4]string{"Genesis", "Custom", "Product", "Commodity"},
			Nodes: []*NodeDecl{
				{Name: "Automobiliste", Kind: KindAnchor, Visibility: -1},
				{Name: "Collectivite", Kind: KindAnchor, Visibility: -1},
				{Name: "AppMobile", Kind: KindComponent, Evolution: "III.7", Visibility: -1},
				{Name: "API", Kind: KindComponent, Evolution: "II.5", Visibility: -1},
				{Name: "Itineraire", Kind: KindComponent, Evolution: "II.8", Visibility: -1},
				{Name: "Alertes", Kind: KindComponent, Evolution: "III.2", Visibility: -1},
				{Name: "CDN", Kind: KindComponent, Evolution: "IV.5", Visibility: -1},
				{Name: "Moteur", Kind: KindComponent, Evolution: "II.3", Visibility: -1},
				{Name: "FluxTrafic", Kind: KindComponent, Evolution: "II.5", Visibility: -1},
				{Name: "Donnees", Kind: KindComponent, Evolution: "III.0", Visibility: -1},
				{Name: "Cloud", Kind: KindComponent, Evolution: "IV.3", Visibility: -1},
				{Name: "OSM", Kind: KindComponent, Evolution: "IV.0", Visibility: -1},
			},
			Edges: []*EdgeDecl{
				{From: "Automobiliste", To: "AppMobile"},
				{From: "Collectivite", To: "API"},
				{From: "Collectivite", To: "Alertes"},
				{From: "AppMobile", To: "Itineraire"},
				{From: "AppMobile", To: "Alertes"},
				{From: "AppMobile", To: "CDN"},
				{From: "Itineraire", To: "Moteur"},
				{From: "Alertes", To: "FluxTrafic"},
				{From: "API", To: "Donnees"},
				{From: "Moteur", To: "Donnees"},
				{From: "Moteur", To: "Cloud"},
				{From: "FluxTrafic", To: "Cloud"},
				{From: "CDN", To: "Cloud"},
				{From: "Donnees", To: "OSM"},
			},
			Pipelines: []*PipelineDecl{
				{
					Name: "Moteur",
					Members: []*PipelineMemberDecl{
						{Name: "AlgoClassique", Position: "I.5"},
						{Name: "AlgoPredictif", Position: "II.5"},
					},
				},
			},
		}
		if withFocus {
			doc.Focuses = []*FocusDecl{{Target: "AppMobile"}}
		}
		return doc
	}

	resultNoFocus, err := BuildMap(makeDoc(false))
	if err != nil {
		t.Fatal(err)
	}
	resultWithFocus, err := BuildMap(makeDoc(true))
	if err != nil {
		t.Fatal(err)
	}

	// Build label -> position maps for comparison
	posNoFocus := make(map[string][2]int)
	for _, c := range resultNoFocus.Map.Components() {
		posNoFocus[componentLabel(c)] = [2]int{c.GetPosition().X, c.GetPosition().Y}
	}
	posWithFocus := make(map[string][2]int)
	for _, c := range resultWithFocus.Map.Components() {
		posWithFocus[componentLabel(c)] = [2]int{c.GetPosition().X, c.GetPosition().Y}
	}

	for label, pos := range posNoFocus {
		focusPos, ok := posWithFocus[label]
		if !ok {
			t.Errorf("component %q missing in focus result", label)
			continue
		}
		if pos != focusPos {
			t.Errorf("component %q: position without focus (%d,%d) != with focus (%d,%d)",
				label, pos[0], pos[1], focusPos[0], focusPos[1])
		}
	}
}

func TestBuildMap_FocusDoesNotChangePositions_FullParse(t *testing.T) {
	base := `title Wardley Map
stages Genesis, Custom, Product, Commodity
anchor Automobiliste
anchor Collectivite
component AppMobile III.7
component API II.5
component Itineraire II.8
component Alertes III.2
component CDN IV.5
component Moteur II.3
component FluxTrafic II.5
component Donnees III.0
component Cloud IV.3
component OSM IV.0
pipeline Moteur {
  AlgoClassique I.5
  AlgoPredictif II.5
}
Automobiliste -> AppMobile
Collectivite -> API
Collectivite -> Alertes
AppMobile -> Itineraire
AppMobile -> Alertes
AppMobile -> CDN
Itineraire -> Moteur
Alertes -> FluxTrafic
API -> Donnees
Moteur -> Donnees
Moteur -> Cloud
FluxTrafic -> Cloud
CDN -> Cloud
Donnees -> OSM
`
	withFocus := base + "focus AppMobile\n"

	parseAndBuild := func(input string) *BuildResult {
		t.Helper()
		p, err := NewParser(bytes.NewBufferString(input))
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
		return result
	}

	resultNoFocus := parseAndBuild(base)
	resultWithFocus := parseAndBuild(withFocus)

	posNoFocus := make(map[string][2]int)
	for _, c := range resultNoFocus.Map.Components() {
		posNoFocus[componentLabel(c)] = [2]int{c.GetPosition().X, c.GetPosition().Y}
	}
	posWithFocus := make(map[string][2]int)
	for _, c := range resultWithFocus.Map.Components() {
		posWithFocus[componentLabel(c)] = [2]int{c.GetPosition().X, c.GetPosition().Y}
	}

	for label, pos := range posNoFocus {
		focusPos, ok := posWithFocus[label]
		if !ok {
			t.Errorf("component %q missing in focus result", label)
			continue
		}
		if pos != focusPos {
			t.Errorf("component %q: position without focus (%d,%d) != with focus (%d,%d)",
				label, pos[0], pos[1], focusPos[0], focusPos[1])
		}
	}
}

func componentLabel(c wardleyToGo.Component) string {
	switch v := c.(type) {
	case *wardley.Component:
		return v.Label
	case *wardley.Anchor:
		return "anchor:" + v.Label
	case *wardley.EvolvedComponent:
		return "evolved:" + v.Label
	case *wardley.Group:
		return "group:" + v.Label
	default:
		return fmt.Sprintf("unknown:%d", c.ID())
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

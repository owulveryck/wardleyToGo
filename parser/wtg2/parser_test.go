package wtg2

import (
	"os"
	"strings"
	"testing"
)

func TestParseMeta(t *testing.T) {
	input := `title: My Wardley Map
date: 2026-01-15
author: Test Author
scope: B2C mobile
question: "Where to invest?"
`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if doc.Title != "My Wardley Map" {
		t.Errorf("title = %q, want %q", doc.Title, "My Wardley Map")
	}
	if doc.Date != "2026-01-15" {
		t.Errorf("date = %q, want %q", doc.Date, "2026-01-15")
	}
	if doc.Author != "Test Author" {
		t.Errorf("author = %q, want %q", doc.Author, "Test Author")
	}
	if doc.Question != "Where to invest?" {
		t.Errorf("question = %q, want %q", doc.Question, "Where to invest?")
	}
}

func TestParseDefaultStages(t *testing.T) {
	input := `title: Test`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	want := [4]string{"", "", "", ""}
	if doc.Stages != want {
		t.Errorf("default stages = %v, want %v", doc.Stages, want)
	}
}

func TestParseStages(t *testing.T) {
	input := `stages: Genèse, Sur-mesure, Produit, Commodité`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	want := [4]string{"Genèse", "Sur-mesure", "Produit", "Commodité"}
	if doc.Stages != want {
		t.Errorf("stages = %v, want %v", doc.Stages, want)
	}
}

func TestParseComponentShorthand(t *testing.T) {
	input := `Application Mobile : III.5`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
	}
	n := doc.Nodes[0]
	if n.Name != "Application Mobile" {
		t.Errorf("name = %q, want %q", n.Name, "Application Mobile")
	}
	if n.Evolution != "III.5" {
		t.Errorf("evolution = %q, want %q", n.Evolution, "III.5")
	}
}

func TestParseComponentWithType(t *testing.T) {
	input := `Données OSM : III.8 (buy)`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
	}
	n := doc.Nodes[0]
	if n.Name != "Données OSM" {
		t.Errorf("name = %q, want %q", n.Name, "Données OSM")
	}
	if n.Type != "buy" {
		t.Errorf("type = %q, want %q", n.Type, "buy")
	}
	if n.Evolution != "III.8" {
		t.Errorf("evolution = %q, want %q", n.Evolution, "III.8")
	}
}

func TestParseComponentWithEvolution(t *testing.T) {
	input := `Moteur de Calcul d'Itinéraire : II.7 !! >> III.5 {
  type: build
  color: #3498DB
  note: "Notre différenciant clé"
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
	}
	n := doc.Nodes[0]
	if n.Name != "Moteur de Calcul d'Itinéraire" {
		t.Errorf("name = %q, want %q", n.Name, "Moteur de Calcul d'Itinéraire")
	}
	if n.Evolution != "II.7" {
		t.Errorf("evolution = %q, want %q", n.Evolution, "II.7")
	}
	if n.EvolvedTo != "III.5" {
		t.Errorf("evolvedTo = %q, want %q", n.EvolvedTo, "III.5")
	}
	if n.Inertia != 2 {
		t.Errorf("inertia = %d, want 2", n.Inertia)
	}
	if n.Type != "build" {
		t.Errorf("type = %q, want %q", n.Type, "build")
	}
	if n.Color != "#3498DB" {
		t.Errorf("color = %q, want %q", n.Color, "#3498DB")
	}
	if n.Note != "Notre différenciant clé" {
		t.Errorf("note = %q, want %q", n.Note, "Notre différenciant clé")
	}
}

func TestParseEdgeChain(t *testing.T) {
	input := `A -> B -> C`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Edges) != 2 {
		t.Fatalf("got %d edges, want 2", len(doc.Edges))
	}
	if doc.Edges[0].From != "A" || doc.Edges[0].To != "B" {
		t.Errorf("edge[0]: %q -> %q, want A -> B", doc.Edges[0].From, doc.Edges[0].To)
	}
	if doc.Edges[1].From != "B" || doc.Edges[1].To != "C" {
		t.Errorf("edge[1]: %q -> %q, want B -> C", doc.Edges[1].From, doc.Edges[1].To)
	}
}

func TestParseAnnotatedEdge(t *testing.T) {
	input := `A -[Open Data, licence annuelle]-> B`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Edges) != 1 {
		t.Fatalf("got %d edges, want 1", len(doc.Edges))
	}
	e := doc.Edges[0]
	if e.Label != "Open Data, licence annuelle" {
		t.Errorf("label = %q, want %q", e.Label, "Open Data, licence annuelle")
	}
}

func TestParsePipeline(t *testing.T) {
	input := `pipeline Engine {
  Algo A : III.5
  Algo B : II.3
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Pipelines) != 1 {
		t.Fatalf("got %d pipelines, want 1", len(doc.Pipelines))
	}
	pl := doc.Pipelines[0]
	if pl.Name != "Engine" {
		t.Errorf("name = %q, want %q", pl.Name, "Engine")
	}
	if len(pl.Members) != 2 {
		t.Fatalf("got %d members, want 2", len(pl.Members))
	}
	if pl.Members[0].Name != "Algo A" {
		t.Errorf("member[0] = %q, want %q", pl.Members[0].Name, "Algo A")
	}
}

func TestParseAnnotation(t *testing.T) {
	input := `note "Candidat à l'externalisation" on Système de Paiement
warning "SPOF" on Moteur`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Annotations) != 2 {
		t.Fatalf("got %d annotations, want 2", len(doc.Annotations))
	}
	if doc.Annotations[0].Kind != "note" {
		t.Errorf("kind = %q, want %q", doc.Annotations[0].Kind, "note")
	}
	if doc.Annotations[0].Text != "Candidat à l'externalisation" {
		t.Errorf("text = %q", doc.Annotations[0].Text)
	}
	if doc.Annotations[0].Target != "Système de Paiement" {
		t.Errorf("target = %q, want %q", doc.Annotations[0].Target, "Système de Paiement")
	}
	if doc.Annotations[1].Kind != "warning" {
		t.Errorf("kind = %q, want %q", doc.Annotations[1].Kind, "warning")
	}
}

func TestParseSignal(t *testing.T) {
	input := `signal accelerating on Algo Prédictif IA`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Signals) != 1 {
		t.Fatalf("got %d signals, want 1", len(doc.Signals))
	}
	if doc.Signals[0].Type != "accelerating" {
		t.Errorf("type = %q, want %q", doc.Signals[0].Type, "accelerating")
	}
	if doc.Signals[0].Target != "Algo Prédictif IA" {
		t.Errorf("target = %q, want %q", doc.Signals[0].Target, "Algo Prédictif IA")
	}
}

func TestParseGroupWithColor(t *testing.T) {
	input := `group Backend {
  color: #E74C3C
  API
  DB
  Cache
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(doc.Groups))
	}
	g := doc.Groups[0]
	if g.Name != "Backend" {
		t.Errorf("name = %q, want %q", g.Name, "Backend")
	}
	if g.Color != "#E74C3C" {
		t.Errorf("color = %q, want %q", g.Color, "#E74C3C")
	}
	if len(g.Members) != 3 {
		t.Fatalf("got %d members, want 3", len(g.Members))
	}
	if g.Members[0] != "API" || g.Members[1] != "DB" || g.Members[2] != "Cache" {
		t.Errorf("members = %v, want [API DB Cache]", g.Members)
	}
}

func TestParseGroupWithoutColor(t *testing.T) {
	input := `group Frontend {
  App
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(doc.Groups))
	}
	g := doc.Groups[0]
	if g.Name != "Frontend" {
		t.Errorf("name = %q, want %q", g.Name, "Frontend")
	}
	if g.Color != "" {
		t.Errorf("color = %q, want empty", g.Color)
	}
	if len(g.Members) != 1 {
		t.Fatalf("got %d members, want 1", len(g.Members))
	}
}

func TestParseDoctrine(t *testing.T) {
	input := `title: My Map
doctrine: context
`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if doc.Doctrine != "context" {
		t.Errorf("doctrine = %q, want %q", doc.Doctrine, "context")
	}
}

func TestParseComponentWithAssetAndCost(t *testing.T) {
	input := `AI Team : II.3 {
  type: build
  asset: human
  cost: "1.2M/year, 12 FTEs"
  note: "Key differentiator"
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
	}
	n := doc.Nodes[0]
	if n.Asset != "human" {
		t.Errorf("asset = %q, want %q", n.Asset, "human")
	}
	if n.Cost != "1.2M/year, 12 FTEs" {
		t.Errorf("cost = %q, want %q", n.Cost, "1.2M/year, 12 FTEs")
	}
	if n.Type != "build" {
		t.Errorf("type = %q, want %q", n.Type, "build")
	}
}

func TestParseQualifiedInertia(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantInertia   int
		wantKinds     []string
		wantEvolution string
		wantEvolvedTo string
		wantType      string
	}{
		{
			name:          "unqualified inertia",
			input:         `Comp : II.7 !! >> III.5`,
			wantInertia:   2,
			wantKinds:     nil,
			wantEvolution: "II.7",
			wantEvolvedTo: "III.5",
		},
		{
			name:          "single qualified kind",
			input:         `Comp : II.7 !(financial) >> III.5`,
			wantInertia:   1,
			wantKinds:     []string{"financial"},
			wantEvolution: "II.7",
			wantEvolvedTo: "III.5",
		},
		{
			name:          "multiple qualified kinds",
			input:         `Comp : II.7 !!(tech,human) >> III.5`,
			wantInertia:   2,
			wantKinds:     []string{"tech", "human"},
			wantEvolution: "II.7",
			wantEvolvedTo: "III.5",
		},
		{
			name:          "triple inertia with kinds",
			input:         `Comp : II.7 !!!(tech,financial,social) >> III.5`,
			wantInertia:   3,
			wantKinds:     []string{"tech", "financial", "social"},
			wantEvolution: "II.7",
			wantEvolvedTo: "III.5",
		},
		{
			name:          "qualified inertia with type",
			input:         `Comp : II.7 !!(tech,human) >> III.5 (buy)`,
			wantInertia:   2,
			wantKinds:     []string{"tech", "human"},
			wantEvolution: "II.7",
			wantEvolvedTo: "III.5",
			wantType:      "buy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}

			if len(doc.Nodes) != 1 {
				t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
			}
			n := doc.Nodes[0]
			if n.Inertia != tt.wantInertia {
				t.Errorf("inertia = %d, want %d", n.Inertia, tt.wantInertia)
			}
			if n.Evolution != tt.wantEvolution {
				t.Errorf("evolution = %q, want %q", n.Evolution, tt.wantEvolution)
			}
			if n.EvolvedTo != tt.wantEvolvedTo {
				t.Errorf("evolvedTo = %q, want %q", n.EvolvedTo, tt.wantEvolvedTo)
			}
			if len(n.InertiaKinds) != len(tt.wantKinds) {
				t.Fatalf("inertiaKinds = %v, want %v", n.InertiaKinds, tt.wantKinds)
			}
			for i, k := range n.InertiaKinds {
				if k != tt.wantKinds[i] {
					t.Errorf("inertiaKinds[%d] = %q, want %q", i, k, tt.wantKinds[i])
				}
			}
			if n.Type != tt.wantType {
				t.Errorf("type = %q, want %q", n.Type, tt.wantType)
			}
		})
	}
}

func TestParseExtendedSignals(t *testing.T) {
	input := `signal co-evolution on DevOps
signal red-queen on Search Engine
signal commoditization on Cloud
signal network-effects on Platform
signal economies-of-scale on DC`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Signals) != 5 {
		t.Fatalf("got %d signals, want 5", len(doc.Signals))
	}

	expected := []struct{ typ, target string }{
		{"co-evolution", "DevOps"},
		{"red-queen", "Search Engine"},
		{"commoditization", "Cloud"},
		{"network-effects", "Platform"},
		{"economies-of-scale", "DC"},
	}
	for i, e := range expected {
		if doc.Signals[i].Type != e.typ {
			t.Errorf("signal[%d].type = %q, want %q", i, doc.Signals[i].Type, e.typ)
		}
		if doc.Signals[i].Target != e.target {
			t.Errorf("signal[%d].target = %q, want %q", i, doc.Signals[i].Target, e.target)
		}
	}
}

func TestParseGameplay(t *testing.T) {
	input := `gameplay ILC on Platform API
gameplay open-source "Commoditize compute" on Cloud Infra
gameplay strangler-fig on Legacy CRM`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Gameplays) != 3 {
		t.Fatalf("got %d gameplays, want 3", len(doc.Gameplays))
	}

	// Test simple gameplay (no description)
	if doc.Gameplays[0].Type != "ILC" {
		t.Errorf("gameplay[0].type = %q, want %q", doc.Gameplays[0].Type, "ILC")
	}
	if doc.Gameplays[0].Target != "Platform API" {
		t.Errorf("gameplay[0].target = %q, want %q", doc.Gameplays[0].Target, "Platform API")
	}
	if doc.Gameplays[0].Text != "" {
		t.Errorf("gameplay[0].text = %q, want empty", doc.Gameplays[0].Text)
	}

	// Test gameplay with description
	if doc.Gameplays[1].Type != "open-source" {
		t.Errorf("gameplay[1].type = %q, want %q", doc.Gameplays[1].Type, "open-source")
	}
	if doc.Gameplays[1].Text != "Commoditize compute" {
		t.Errorf("gameplay[1].text = %q, want %q", doc.Gameplays[1].Text, "Commoditize compute")
	}
	if doc.Gameplays[1].Target != "Cloud Infra" {
		t.Errorf("gameplay[1].target = %q, want %q", doc.Gameplays[1].Target, "Cloud Infra")
	}

	// Test another simple gameplay
	if doc.Gameplays[2].Type != "strangler-fig" {
		t.Errorf("gameplay[2].type = %q, want %q", doc.Gameplays[2].Type, "strangler-fig")
	}
	if doc.Gameplays[2].Target != "Legacy CRM" {
		t.Errorf("gameplay[2].target = %q, want %q", doc.Gameplays[2].Target, "Legacy CRM")
	}
}

func TestParseGroupWithTeam(t *testing.T) {
	input := `group R&D Team {
  team: explorer
  color: #E74C3C
  Algo
  Cache
}`
	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(doc.Groups))
	}
	g := doc.Groups[0]
	if g.Team != "explorer" {
		t.Errorf("team = %q, want %q", g.Team, "explorer")
	}
	if g.Color != "#E74C3C" {
		t.Errorf("color = %q, want %q", g.Color, "#E74C3C")
	}
	if len(g.Members) != 2 {
		t.Fatalf("got %d members, want 2", len(g.Members))
	}
}

func TestParseLegend(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"bare keyword", "legend\n", true},
		{"with colon", "legend: true\n", true},
		{"absent", "title: My Map\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if doc.Legend != tt.want {
				t.Errorf("legend = %v, want %v", doc.Legend, tt.want)
			}
		})
	}
}

func TestParseExampleFile(t *testing.T) {
	f, err := os.Open("testdata/example.wtg2")
	if err != nil {
		// Try from wtg2 directory
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

	if doc.Title != "Plateforme de Navigation — Stratégie 2026" {
		t.Errorf("title = %q", doc.Title)
	}
	if doc.Stages[0] != "Genèse" {
		t.Errorf("stage[0] = %q, want Genèse", doc.Stages[0])
	}

	// Should have anchors + components
	anchorCount := 0
	compCount := 0
	for _, n := range doc.Nodes {
		switch n.Kind {
		case KindAnchor:
			anchorCount++
		case KindComponent:
			compCount++
		case KindSubmap:
			compCount++
		}
	}
	if anchorCount != 2 {
		t.Errorf("anchor count = %d, want 2", anchorCount)
	}
	if compCount < 10 {
		t.Errorf("component count = %d, want >= 10", compCount)
	}

	// First anchor should have evolution
	for _, n := range doc.Nodes {
		if n.Kind == KindAnchor && n.Name == "Automobiliste" {
			if n.Evolution != "III.5" {
				t.Errorf("anchor Automobiliste evolution = %q, want %q", n.Evolution, "III.5")
			}
		}
		if n.Kind == KindAnchor && n.Name == "Collectivité Locale" {
			if n.Evolution != "" {
				t.Errorf("anchor Collectivité Locale evolution = %q, want empty", n.Evolution)
			}
		}
	}

	// Should have edges
	if len(doc.Edges) < 10 {
		t.Errorf("edge count = %d, want >= 10", len(doc.Edges))
	}

	// Should have a pipeline
	if len(doc.Pipelines) != 1 {
		t.Errorf("pipeline count = %d, want 1", len(doc.Pipelines))
	}

	// Should have annotations
	if len(doc.Annotations) < 5 {
		t.Errorf("annotation count = %d, want >= 5", len(doc.Annotations))
	}

	// Should have signals
	if len(doc.Signals) != 4 {
		t.Errorf("signal count = %d, want 4", len(doc.Signals))
	}
}

func TestParseFocus(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCount  int
		wantTarget string
	}{
		{
			name:       "simple name",
			input:      "focus App",
			wantCount:  1,
			wantTarget: "App",
		},
		{
			name:       "multi-word name",
			input:      "focus Moteur de Calcul d'Itinéraire",
			wantCount:  1,
			wantTarget: "Moteur de Calcul d'Itinéraire",
		},
		{
			name:      "multiple focuses",
			input:     "focus App\nfocus DB",
			wantCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if len(doc.Focuses) != tt.wantCount {
				t.Fatalf("got %d focuses, want %d", len(doc.Focuses), tt.wantCount)
			}
			if tt.wantTarget != "" && doc.Focuses[0].Target != tt.wantTarget {
				t.Errorf("target = %q, want %q", doc.Focuses[0].Target, tt.wantTarget)
			}
		})
	}
}

func TestParseAnchorWithEvolution(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantName      string
		wantEvolution string
		wantEvolvedTo string
	}{
		{
			name:          "anchor without evolution",
			input:         "anchor Customer",
			wantName:      "Customer",
			wantEvolution: "",
		},
		{
			name:          "anchor with evolution",
			input:         "anchor Customer : III.5",
			wantName:      "Customer",
			wantEvolution: "III.5",
		},
		{
			name:          "anchor with multi-word name and evolution",
			input:         "anchor Collectivité Locale : II.3",
			wantName:      "Collectivité Locale",
			wantEvolution: "II.3",
		},
		{
			name:          "anchor with evolution and evolvedTo",
			input:         "anchor User : II.3 >> III.5",
			wantName:      "User",
			wantEvolution: "II.3",
			wantEvolvedTo: "III.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if len(doc.Nodes) != 1 {
				t.Fatalf("got %d nodes, want 1", len(doc.Nodes))
			}
			n := doc.Nodes[0]
			if n.Kind != KindAnchor {
				t.Errorf("kind = %v, want KindAnchor", n.Kind)
			}
			if n.Name != tt.wantName {
				t.Errorf("name = %q, want %q", n.Name, tt.wantName)
			}
			if n.Evolution != tt.wantEvolution {
				t.Errorf("evolution = %q, want %q", n.Evolution, tt.wantEvolution)
			}
			if n.EvolvedTo != tt.wantEvolvedTo {
				t.Errorf("evolvedTo = %q, want %q", n.EvolvedTo, tt.wantEvolvedTo)
			}
		})
	}
}

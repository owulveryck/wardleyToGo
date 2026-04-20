package layout

import (
	"testing"
	"time"
)

func TestLayout_LinearChain(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	if pos["A"] >= pos["B"] {
		t.Errorf("A.Y=%d should be < B.Y=%d", pos["A"], pos["B"])
	}
	if pos["B"] >= pos["C"] {
		t.Errorf("B.Y=%d should be < C.Y=%d", pos["B"], pos["C"])
	}
	assertInRange(t, pos, "A")
	assertInRange(t, pos, "B")
	assertInRange(t, pos, "C")
}

func TestLayout_Diamond(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "C"},
			{Name: "D"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "B", To: "D"},
			{From: "C", To: "D"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	if pos["A"] >= pos["B"] {
		t.Errorf("A.Y=%d should be < B.Y=%d", pos["A"], pos["B"])
	}
	if pos["A"] >= pos["C"] {
		t.Errorf("A.Y=%d should be < C.Y=%d", pos["A"], pos["C"])
	}
	if pos["D"] <= pos["B"] {
		t.Errorf("D.Y=%d should be > B.Y=%d", pos["D"], pos["B"])
	}
	if pos["D"] <= pos["C"] {
		t.Errorf("D.Y=%d should be > C.Y=%d", pos["D"], pos["C"])
	}
}

func TestLayout_PipelineMembers(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "User", Kind: KindAnchor},
			{Name: "Engine"},
		},
		Edges: []Edge{
			{From: "User", To: "Engine"},
		},
		Pipelines: []Pipeline{
			{Parent: "Engine", Members: []string{"Algo A", "Algo B"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	parentY := pos["Engine"]
	if pos["Algo A"] != parentY {
		t.Errorf("Algo A.Y=%d should equal Engine.Y=%d", pos["Algo A"], parentY)
	}
	if pos["Algo B"] != parentY {
		t.Errorf("Algo B.Y=%d should equal Engine.Y=%d", pos["Algo B"], parentY)
	}
}

func TestLayout_SingleNode(t *testing.T) {
	g := &Graph{
		Nodes: []Node{{Name: "A"}},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	if _, ok := pos["A"]; !ok {
		t.Fatal("expected position for A")
	}
	assertInRange(t, pos, "A")
}

func TestLayout_Disconnected(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	for _, name := range []string{"A", "B", "C"} {
		assertInRange(t, pos, name)
	}
}

func TestLayout_AnchorsAtTop(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "Anchor1", Kind: KindAnchor},
			{Name: "Anchor2", Kind: KindAnchor},
			{Name: "Comp"},
		},
		Edges: []Edge{
			{From: "Anchor1", To: "Comp"},
			{From: "Anchor2", To: "Comp"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	opts := DefaultOptions()
	if pos["Anchor1"] > opts.MinY+opts.MinSpacing {
		t.Errorf("Anchor1.Y=%d should be near top", pos["Anchor1"])
	}
	if pos["Anchor2"] > opts.MinY+opts.MinSpacing {
		t.Errorf("Anchor2.Y=%d should be near top", pos["Anchor2"])
	}
	if pos["Comp"] <= pos["Anchor1"] {
		t.Errorf("Comp.Y=%d should be > Anchor1.Y=%d", pos["Comp"], pos["Anchor1"])
	}
}

func TestLayout_MinSpacing(t *testing.T) {
	// Build a deep chain to test spacing
	nodes := []Node{{Name: "N0", Kind: KindAnchor}}
	edges := []Edge{}
	for i := 1; i <= 15; i++ {
		name := "N" + string(rune('0'+i/10)) + string(rune('0'+i%10))
		if i < 10 {
			name = "N" + string(rune('0'+i))
		}
		nodes = append(nodes, Node{Name: name})
		prev := "N" + string(rune('0'+(i-1)/10)) + string(rune('0'+(i-1)%10))
		if i-1 < 10 {
			prev = "N" + string(rune('0'+i-1))
		}
		edges = append(edges, Edge{From: prev, To: name})
	}

	g := &Graph{Nodes: nodes, Edges: edges}
	opts := DefaultOptions()
	opts.MinSpacing = 3
	l := New(opts)
	pos := mustLayout(t, l, g)

	for _, name := range []string{"N0", "N1", "N2", "N3", "N4", "N5"} {
		assertInRange(t, pos, name)
	}
}

func TestLayout_LargeGraph(t *testing.T) {
	nodes := make([]Node, 50)
	edges := make([]Edge, 0, 70)
	nodes[0] = Node{Name: "root", Kind: KindAnchor}
	for i := 1; i < 50; i++ {
		name := nodeName(i)
		nodes[i] = Node{Name: name}
		// Connect to a parent in the previous "layer"
		parentIdx := (i - 1) / 3
		edges = append(edges, Edge{From: nodeName(parentIdx), To: name})
	}

	g := &Graph{Nodes: nodes, Edges: edges}
	l := New(DefaultOptions())

	start := time.Now()
	pos := mustLayout(t, l, g)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("layout took %v, expected < 100ms", elapsed)
	}

	for i := 0; i < 50; i++ {
		assertInRange(t, pos, nodeName(i))
	}
}

func TestLayout_ColonMemberEdges(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "User", Kind: KindAnchor},
			{Name: "Engine"},
		},
		Edges: []Edge{
			{From: "User", To: "Engine:AlgoA"},
		},
		Pipelines: []Pipeline{
			{Parent: "Engine", Members: []string{"AlgoA"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	if _, ok := pos["AlgoA"]; !ok {
		t.Fatal("expected position for AlgoA")
	}
}

func TestResolveMember(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Engine:AlgoA", "AlgoA"},
		{"simple name", "simple name"},
		{"A : B", "A : B"}, // space around colon — not a member ref
		{"Pipeline:Member", "Member"},
	}
	for _, tt := range tests {
		got := resolveMember(tt.input)
		if got != tt.want {
			t.Errorf("resolveMember(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func nodeName(i int) string {
	if i == 0 {
		return "root"
	}
	// Simple deterministic names
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i <= 26 {
		return string(letters[i-1])
	}
	return string(letters[(i-1)/26-1]) + string(letters[(i-1)%26])
}

func TestLayout_NegativeDyForce(t *testing.T) {
	// Two nodes placed so that dy < 0 and |dy| < 0.1 — exercises the
	// negative-dy clamping branch in forceSpread.
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			// B→C forces them into adjacent ranks even though they share a parent
			{From: "B", To: "C"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	if pos["A"] >= pos["B"] {
		t.Errorf("A.Y=%d should be < B.Y=%d", pos["A"], pos["B"])
	}
	if pos["B"] >= pos["C"] {
		t.Errorf("B.Y=%d should be < C.Y=%d", pos["B"], pos["C"])
	}
}

func TestLayout_MissingEdgeEndpoints(t *testing.T) {
	// Edges referencing non-existent nodes should be silently ignored
	// during force simulation.
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "A", To: "Ghost"},   // "Ghost" not in Nodes
			{From: "Missing", To: "B"}, // "Missing" not in Nodes
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	assertInRange(t, pos, "A")
	assertInRange(t, pos, "B")
	if _, ok := pos["Ghost"]; ok {
		t.Error("should not have position for Ghost")
	}
}

func TestLayout_LowerBoundClamp(t *testing.T) {
	// Use extreme options where MinY is high so positions might clamp.
	opts := DefaultOptions()
	opts.MinY = 40
	opts.MaxY = 60

	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "C"},
			{Name: "D"},
			{Name: "E"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "D"},
			{From: "D", To: "E"},
		},
	}

	l := New(opts)
	pos := mustLayout(t, l, g)

	for _, name := range []string{"A", "B", "C", "D", "E"} {
		y := pos[name]
		if y < opts.MinY || y > opts.MaxY {
			t.Errorf("%s.Y=%d out of range [%d, %d]", name, y, opts.MinY, opts.MaxY)
		}
	}
}

func TestLayout_DisconnectedNodeRank(t *testing.T) {
	// A disconnected node with an anchor exercises the maxRank++ fallback
	// in topoRanks for unvisited nodes.
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "Island"}, // not connected to anything
		},
		Edges: []Edge{
			{From: "A", To: "B"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	assertInRange(t, pos, "A")
	assertInRange(t, pos, "B")
	assertInRange(t, pos, "Island")
	// Island should be at or below B since it gets a higher rank
	if pos["Island"] < pos["A"] {
		t.Errorf("Island.Y=%d should be >= A.Y=%d", pos["Island"], pos["A"])
	}
}

func TestLayout_PipelineMemberEdgeInForce(t *testing.T) {
	// Edges involving pipeline members should be skipped during force simulation.
	g := &Graph{
		Nodes: []Node{
			{Name: "User", Kind: KindAnchor},
			{Name: "Engine"},
			{Name: "DB"},
		},
		Edges: []Edge{
			{From: "User", To: "Engine"},
			{From: "Engine:AlgoA", To: "DB"},
		},
		Pipelines: []Pipeline{
			{Parent: "Engine", Members: []string{"AlgoA", "AlgoB"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	// Pipeline members share parent Y
	if pos["AlgoA"] != pos["Engine"] {
		t.Errorf("AlgoA.Y=%d should equal Engine.Y=%d", pos["AlgoA"], pos["Engine"])
	}
	assertInRange(t, pos, "DB")
}

func TestLayout_VerticalSpreadWithPipeline(t *testing.T) {
	// Mimics the GPS navigation example map structure.
	// Verifies nodes are spread evenly, not bunched at bottom.
	g := &Graph{
		Nodes: []Node{
			{Name: "Anchor1", Kind: KindAnchor},
			{Name: "Anchor2", Kind: KindAnchor},
			{Name: "App"},
			{Name: "API"},
			{Name: "Display"},
			{Name: "Alerts"},
			{Name: "CDN"},
			{Name: "Engine"},
			{Name: "Feed"},
			{Name: "DataModel"},
			{Name: "Cloud"},
			{Name: "OSMData"},
			{Name: "Mobile"},
		},
		Edges: []Edge{
			{From: "Anchor1", To: "App"},
			{From: "Anchor2", To: "API"},
			{From: "Anchor2", To: "Alerts"},
			{From: "App", To: "Display"},
			{From: "App", To: "Alerts"},
			{From: "App", To: "CDN"},
			{From: "Display", To: "Engine"},
			{From: "Alerts", To: "Feed"},
			{From: "Alerts", To: "Engine:AlgoPredictive"},
			{From: "API", To: "DataModel"},
			{From: "Engine", To: "DataModel"},
			{From: "Engine", To: "Cloud"},
			{From: "Feed", To: "Cloud"},
			{From: "CDN", To: "Cloud"},
			{From: "DataModel", To: "OSMData"},
			{From: "Cloud", To: "Mobile"},
		},
		Pipelines: []Pipeline{
			{Parent: "Engine", Members: []string{"AlgoClassic", "AlgoPredictive", "AlgoQuantum"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	opts := DefaultOptions()
	for _, name := range []string{"Anchor1", "Anchor2", "App", "API", "Display", "Alerts", "CDN", "Engine", "Feed", "DataModel", "Cloud", "OSMData", "Mobile"} {
		assertInRange(t, pos, name)
	}

	// Pipeline members inherit parent Y
	if pos["AlgoClassic"] != pos["Engine"] {
		t.Errorf("AlgoClassic.Y=%d should equal Engine.Y=%d", pos["AlgoClassic"], pos["Engine"])
	}
	if pos["AlgoPredictive"] != pos["Engine"] {
		t.Errorf("AlgoPredictive.Y=%d should equal Engine.Y=%d", pos["AlgoPredictive"], pos["Engine"])
	}

	// Verify vertical spread: mid-depth nodes should be near the
	// middle of the range, not crammed near maxY. The threshold
	// accounts for same-rank nodes being slightly spread apart.
	midY := (opts.MinY + opts.MaxY) / 2
	if pos["Engine"] > midY+25 {
		t.Errorf("Engine.Y=%d is too far below midpoint %d, nodes are bunched at bottom", pos["Engine"], midY)
	}
	if pos["App"] > midY {
		t.Errorf("App.Y=%d should be in the upper half (midY=%d)", pos["App"], midY)
	}
}

func TestLayout_SameRankNodesCloseY(t *testing.T) {
	// B, C, D are all at rank 1. They should have similar Y values.
	g := &Graph{
		Nodes: []Node{
			{Name: "A", Kind: KindAnchor},
			{Name: "B"},
			{Name: "C"},
			{Name: "D"},
			{Name: "E"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "A", To: "D"},
			{From: "B", To: "E"},
			{From: "C", To: "E"},
			{From: "D", To: "E"},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	maxDiff := 0
	for _, n1 := range []string{"B", "C", "D"} {
		for _, n2 := range []string{"B", "C", "D"} {
			diff := pos[n1] - pos[n2]
			if diff < 0 {
				diff = -diff
			}
			if diff > maxDiff {
				maxDiff = diff
			}
		}
	}
	opts := DefaultOptions()
	if maxDiff > 2*opts.MinSpacing {
		t.Errorf("same-rank nodes B,C,D have Y spread of %d (B=%d, C=%d, D=%d), want <= %d",
			maxDiff, pos["B"], pos["C"], pos["D"], 2*opts.MinSpacing)
	}
}

func TestLayout_GPSNavigationMap(t *testing.T) {
	// Mimics the real GPS navigation Wardley Map that triggered the
	// cascading compression bug: 2 anchors, 6 rank-deep DAG, a pipeline,
	// and many same-rank nodes. Before the displacement cap fix, 11 of
	// 17 components were clamped at Y=97.
	g := &Graph{
		Nodes: []Node{
			{Name: "Automobiliste", Kind: KindAnchor},
			{Name: "Collectivite", Kind: KindAnchor},
			{Name: "AppMobile"},
			{Name: "API"},
			{Name: "Itineraire"},
			{Name: "Alertes"},
			{Name: "CDN"},
			{Name: "Paiement"},
			{Name: "Moteur"}, // pipeline parent
			{Name: "FluxTrafic"},
			{Name: "Donnees"},
			{Name: "Cloud"},
			{Name: "OSM"},
			{Name: "Hebergement"},
		},
		Edges: []Edge{
			{From: "Automobiliste", To: "AppMobile"},
			{From: "Collectivite", To: "API"},
			{From: "Collectivite", To: "Alertes"},
			{From: "AppMobile", To: "Itineraire"},
			{From: "AppMobile", To: "Alertes"},
			{From: "AppMobile", To: "CDN"},
			{From: "AppMobile", To: "Paiement"},
			{From: "Itineraire", To: "Moteur"},
			{From: "Alertes", To: "FluxTrafic"},
			{From: "Alertes", To: "Moteur:AlgoPredictif"},
			{From: "API", To: "Donnees"},
			{From: "Moteur", To: "Donnees"},
			{From: "Moteur", To: "Cloud"},
			{From: "FluxTrafic", To: "Cloud"},
			{From: "CDN", To: "Cloud"},
			{From: "Donnees", To: "OSM"},
			{From: "Cloud", To: "Hebergement"},
		},
		Pipelines: []Pipeline{
			{Parent: "Moteur", Members: []string{"AlgoClassique", "AlgoPredictif", "AlgoQuantique"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	opts := DefaultOptions()
	nonPipeline := []string{
		"Automobiliste", "Collectivite", "AppMobile", "API",
		"Itineraire", "Alertes", "CDN", "Paiement",
		"Moteur", "FluxTrafic", "Donnees", "Cloud", "OSM", "Hebergement",
	}
	for _, name := range nonPipeline {
		assertInRange(t, pos, name)
	}

	// Pipeline members inherit parent Y
	if pos["AlgoClassique"] != pos["Moteur"] {
		t.Errorf("AlgoClassique.Y=%d should equal Moteur.Y=%d", pos["AlgoClassique"], pos["Moteur"])
	}

	// Key assertion: no more than 3 non-pipeline nodes at the exact same Y.
	// Before the fix, 11 nodes were all at Y=97.
	yCounts := make(map[int]int)
	for _, name := range nonPipeline {
		yCounts[pos[name]]++
	}
	for y, count := range yCounts {
		if count > 3 {
			t.Errorf("%d non-pipeline nodes at Y=%d, want <= 3 (nodes bunched)", count, y)
		}
	}

	// Mid-depth nodes should be in the middle third, not at the bottom
	midY := (opts.MinY + opts.MaxY) / 2
	if pos["Moteur"] > midY+25 {
		t.Errorf("Moteur.Y=%d should be near middle (midY=%d), not at bottom", pos["Moteur"], midY)
	}

	// Anchors at top, deepest nodes at bottom, with spread in between
	if pos["OSM"] <= pos["Moteur"] {
		t.Errorf("OSM.Y=%d should be below Moteur.Y=%d", pos["OSM"], pos["Moteur"])
	}
	if pos["Hebergement"] <= pos["Cloud"] {
		t.Errorf("Hebergement.Y=%d should be below Cloud.Y=%d", pos["Hebergement"], pos["Cloud"])
	}
}

func TestLayout_PipelineMemberOutgoingEdge(t *testing.T) {
	// Regression test: when pipeline members have outgoing edges to
	// non-pipeline nodes, those downstream nodes must get proper ranks
	// via BFS (not inflated disconnected ranks). Before the fix,
	// pipeline members were unreachable from anchors, so downstream
	// nodes like C and D were pushed to the very bottom.
	g := &Graph{
		Nodes: []Node{
			{Name: "Anchor", Kind: KindAnchor},
			{Name: "Mid"},
			{Name: "Parent"}, // pipeline parent
			{Name: "C"},      // depends on pipeline member M1
			{Name: "D"},      // depends on C
		},
		Edges: []Edge{
			{From: "Anchor", To: "Mid"},
			{From: "Mid", To: "Parent"},
			{From: "Parent:M1", To: "C"},
			{From: "C", To: "D"},
		},
		Pipelines: []Pipeline{
			{Parent: "Parent", Members: []string{"M1", "M2"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	opts := DefaultOptions()

	// Pipeline members inherit parent Y
	if pos["M1"] != pos["Parent"] {
		t.Errorf("M1.Y=%d should equal Parent.Y=%d", pos["M1"], pos["Parent"])
	}
	if pos["M2"] != pos["Parent"] {
		t.Errorf("M2.Y=%d should equal Parent.Y=%d", pos["M2"], pos["Parent"])
	}

	// C must be below Parent (it depends on a member of Parent)
	if pos["C"] <= pos["Parent"] {
		t.Errorf("C.Y=%d should be > Parent.Y=%d", pos["C"], pos["Parent"])
	}

	// D must be below C
	if pos["D"] <= pos["C"] {
		t.Errorf("D.Y=%d should be > C.Y=%d", pos["D"], pos["C"])
	}

	// Key assertion: C should NOT be pushed to the very bottom.
	// With 5 effective ranks (0-4), C at rank 3 should be well below
	// midpoint but not crammed at maxY.
	midY := (opts.MinY + opts.MaxY) / 2
	if pos["C"] > midY+30 {
		t.Errorf("C.Y=%d is too far below midpoint %d — rank inflation bug", pos["C"], midY)
	}

	for _, name := range []string{"Anchor", "Mid", "Parent", "C", "D"} {
		assertInRange(t, pos, name)
	}
}

func TestLayout_MultiplePipelinesSpacing(t *testing.T) {
	// 4 pipeline parents (3 at the same rank) must not overlap visually.
	// This reproduces the bug where multiple pipelines were crammed together.
	g := &Graph{
		Nodes: []Node{
			{Name: "Anchor", Kind: KindAnchor},
			{Name: "Mid"},
			{Name: "P1"}, // pipeline parent
			{Name: "P2"}, // pipeline parent, same rank as P1
			{Name: "P3"}, // pipeline parent, same rank as P1
			{Name: "P4"}, // pipeline parent, one rank deeper
			{Name: "Leaf"},
		},
		Edges: []Edge{
			{From: "Anchor", To: "Mid"},
			{From: "Mid", To: "P1"},
			{From: "Mid", To: "P2"},
			{From: "Mid", To: "P3"},
			{From: "P1", To: "P4"},
			{From: "P4", To: "Leaf"},
		},
		Pipelines: []Pipeline{
			{Parent: "P1", Members: []string{"P1M1", "P1M2"}},
			{Parent: "P2", Members: []string{"P2M1", "P2M2"}},
			{Parent: "P3", Members: []string{"P3M1", "P3M2"}},
			{Parent: "P4", Members: []string{"P4M1", "P4M2"}},
		},
	}

	l := New(DefaultOptions())
	pos := mustLayout(t, l, g)

	opts := DefaultOptions()

	// Pipeline parents at the same rank should be spaced at least
	// MinSpacing * 2 apart (conservative check for 3× repulsion).
	minGap := opts.MinSpacing * 2
	for _, pair := range [][2]string{{"P1", "P2"}, {"P2", "P3"}, {"P1", "P3"}} {
		diff := pos[pair[0]] - pos[pair[1]]
		if diff < 0 {
			diff = -diff
		}
		if diff < minGap {
			t.Errorf("%s.Y=%d and %s.Y=%d are only %d apart, want >= %d",
				pair[0], pos[pair[0]], pair[1], pos[pair[1]], diff, minGap)
		}
	}

	// All nodes in range
	for _, name := range []string{"Anchor", "Mid", "P1", "P2", "P3", "P4", "Leaf"} {
		assertInRange(t, pos, name)
	}

	// Pipeline members inherit parent Y
	if pos["P1M1"] != pos["P1"] {
		t.Errorf("P1M1.Y=%d should equal P1.Y=%d", pos["P1M1"], pos["P1"])
	}
}

func TestLayout_Deterministic(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "Automobiliste", Kind: KindAnchor},
			{Name: "Collectivite", Kind: KindAnchor},
			{Name: "AppMobile"},
			{Name: "API"},
			{Name: "Itineraire"},
			{Name: "Alertes"},
			{Name: "CDN"},
			{Name: "Paiement"},
			{Name: "Moteur"},
			{Name: "FluxTrafic"},
			{Name: "Donnees"},
			{Name: "Cloud"},
			{Name: "OSM"},
			{Name: "Hebergement"},
		},
		Edges: []Edge{
			{From: "Automobiliste", To: "AppMobile"},
			{From: "Collectivite", To: "API"},
			{From: "Collectivite", To: "Alertes"},
			{From: "AppMobile", To: "Itineraire"},
			{From: "AppMobile", To: "Alertes"},
			{From: "AppMobile", To: "CDN"},
			{From: "AppMobile", To: "Paiement"},
			{From: "Itineraire", To: "Moteur"},
			{From: "Alertes", To: "FluxTrafic"},
			{From: "Alertes", To: "Moteur:AlgoPredictif"},
			{From: "API", To: "Donnees"},
			{From: "Moteur", To: "Donnees"},
			{From: "Moteur", To: "Cloud"},
			{From: "FluxTrafic", To: "Cloud"},
			{From: "CDN", To: "Cloud"},
			{From: "Donnees", To: "OSM"},
			{From: "Cloud", To: "Hebergement"},
		},
		Pipelines: []Pipeline{
			{Parent: "Moteur", Members: []string{"AlgoClassique", "AlgoPredictif", "AlgoQuantique"}},
		},
	}

	l := New(DefaultOptions())
	reference := mustLayout(t, l, g)

	for run := 0; run < 50; run++ {
		pos := mustLayout(t, l, g)
		for name, y := range reference {
			if pos[name] != y {
				t.Fatalf("run %d: %s.Y=%d, want %d (non-deterministic layout)",
					run, name, pos[name], y)
			}
		}
	}
}

func TestLayout_DeterministicDisconnected(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"},
		},
	}

	l := New(DefaultOptions())
	reference := mustLayout(t, l, g)

	for run := 0; run < 50; run++ {
		pos := mustLayout(t, l, g)
		for name, y := range reference {
			if pos[name] != y {
				t.Fatalf("run %d: %s.Y=%d, want %d (non-deterministic layout)",
					run, name, pos[name], y)
			}
		}
	}
}

func TestLayout_CycleReturnsError(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "A"},
		},
	}

	l := New(DefaultOptions())
	done := make(chan error, 1)
	go func() {
		_, err := l.Layout(g)
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected error for cyclic graph, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Layout with cyclic graph did not terminate")
	}
}

func TestLayout_SelfLoopReturnsError(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{Name: "A"},
			{Name: "B"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "B"},
		},
	}

	l := New(DefaultOptions())
	_, err := l.Layout(g)
	if err == nil {
		t.Fatal("expected error for self-loop, got nil")
	}
}

func mustLayout(t *testing.T, l Layouter, g *Graph) map[string]int {
	t.Helper()
	pos, err := l.Layout(g)
	if err != nil {
		t.Fatal(err)
	}
	return pos
}

func assertInRange(t *testing.T, pos map[string]int, name string) {
	t.Helper()
	opts := DefaultOptions()
	y, ok := pos[name]
	if !ok {
		t.Errorf("no position for %q", name)
		return
	}
	if y < opts.MinY || y > opts.MaxY {
		t.Errorf("%s.Y=%d out of range [%d, %d]", name, y, opts.MinY, opts.MaxY)
	}
}

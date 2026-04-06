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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)
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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
			{From: "A", To: "Ghost"},     // "Ghost" not in Nodes
			{From: "Missing", To: "B"},    // "Missing" not in Nodes
		},
	}

	l := New(DefaultOptions())
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

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
	pos := l.Layout(g)

	// Pipeline members share parent Y
	if pos["AlgoA"] != pos["Engine"] {
		t.Errorf("AlgoA.Y=%d should equal Engine.Y=%d", pos["AlgoA"], pos["Engine"])
	}
	assertInRange(t, pos, "DB")
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

package graph

import "testing"

func TestDirectedGraph_AddNodeAndEdge(t *testing.T) {
	g := NewDirectedGraph()
	n1 := g.NewNode()
	n2 := g.NewNode()
	if err := g.AddNode(n1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatal(err)
	}

	if g.Node(n1.ID()) == nil {
		t.Fatal("expected node n1")
	}
	if g.Nodes().Len() != 2 {
		t.Fatalf("expected 2 nodes, got %d", g.Nodes().Len())
	}

	e := SimpleEdge{F: n1, T: n2}
	if err := g.SetEdge(e); err != nil {
		t.Fatal(err)
	}

	if !g.HasEdgeFromTo(n1.ID(), n2.ID()) {
		t.Error("expected edge from n1 to n2")
	}
	if g.HasEdgeFromTo(n2.ID(), n1.ID()) {
		t.Error("unexpected reverse edge")
	}
	if g.Edges().Len() != 1 {
		t.Fatalf("expected 1 edge, got %d", g.Edges().Len())
	}
}

func TestDirectedGraph_RemoveEdge(t *testing.T) {
	g := NewDirectedGraph()
	n1 := g.NewNode()
	n2 := g.NewNode()
	if err := g.AddNode(n1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n2}); err != nil {
		t.Fatal(err)
	}
	g.RemoveEdge(n1.ID(), n2.ID())

	if g.HasEdgeFromTo(n1.ID(), n2.ID()) {
		t.Error("edge should have been removed")
	}
}

func TestDirectedGraph_FromTo(t *testing.T) {
	g := NewDirectedGraph()
	n1 := g.NewNode()
	n2 := g.NewNode()
	n3 := g.NewNode()
	if err := g.AddNode(n1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n3); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n2}); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n3}); err != nil {
		t.Fatal(err)
	}

	from := g.From(n1.ID())
	if from.Len() != 2 {
		t.Fatalf("expected 2 successors, got %d", from.Len())
	}

	to := g.To(n2.ID())
	if to.Len() != 1 {
		t.Fatalf("expected 1 predecessor, got %d", to.Len())
	}
}

func TestDepthFirst_Walk(t *testing.T) {
	g := NewDirectedGraph()
	n1 := g.NewNode()
	n2 := g.NewNode()
	n3 := g.NewNode()
	if err := g.AddNode(n1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n3); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n2}); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n2, T: n3}); err != nil {
		t.Fatal(err)
	}

	var visited []int64
	df := &DepthFirst{
		Visit: func(n Node) {
			visited = append(visited, n.ID())
		},
	}
	df.Walk(g, n1, nil)

	if len(visited) != 3 {
		t.Fatalf("expected 3 visited nodes, got %d", len(visited))
	}
}

func TestDirectedGraph_AddNode_DuplicateError(t *testing.T) {
	g := NewDirectedGraph()
	n := g.NewNode()
	if err := g.AddNode(n); err != nil {
		t.Fatal(err)
	}
	err := g.AddNode(n)
	if err == nil {
		t.Fatal("expected error when adding duplicate node")
	}
	if g.Nodes().Len() != 1 {
		t.Fatalf("expected 1 node after duplicate add, got %d", g.Nodes().Len())
	}
}

func TestDirectedGraph_SetEdge_SelfLoopError(t *testing.T) {
	g := NewDirectedGraph()
	n := g.NewNode()
	if err := g.AddNode(n); err != nil {
		t.Fatal(err)
	}
	err := g.SetEdge(SimpleEdge{F: n, T: n})
	if err == nil {
		t.Fatal("expected error for self-loop edge")
	}
	if g.Edges().Len() != 0 {
		t.Fatalf("expected 0 edges after self-loop attempt, got %d", g.Edges().Len())
	}
}

func TestDirectedGraph_SetEdge_ReplaceExisting(t *testing.T) {
	g := NewDirectedGraph()
	n1 := g.NewNode()
	n2 := g.NewNode()
	if err := g.AddNode(n1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n2}); err != nil {
		t.Fatal(err)
	}
	if err := g.SetEdge(SimpleEdge{F: n1, T: n2}); err != nil {
		t.Fatal("replacing existing edge should not error:", err)
	}
	if g.Edges().Len() != 1 {
		t.Fatalf("expected 1 edge after replace, got %d", g.Edges().Len())
	}
}

func TestIterators(t *testing.T) {
	nodes := NewNodes([]Node{simpleNode{0}, simpleNode{1}})
	if nodes.Len() != 2 {
		t.Fatalf("expected len 2")
	}
	count := 0
	for nodes.Next() {
		_ = nodes.Node()
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 iterations, got %d", count)
	}
	nodes.Reset()
	if !nodes.Next() {
		t.Fatal("expected Next to succeed after Reset")
	}

	edges := NewEdges([]Edge{SimpleEdge{F: simpleNode{0}, T: simpleNode{1}}})
	if edges.Len() != 1 {
		t.Fatal("expected 1 edge")
	}
	if !edges.Next() {
		t.Fatal("expected Next to succeed")
	}
	_ = edges.Edge()
	edges.Reset()
	if !edges.Next() {
		t.Fatal("expected Next to succeed after Reset")
	}
}

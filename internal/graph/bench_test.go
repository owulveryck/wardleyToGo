package graph

import "testing"

// buildTestGraph creates a directed graph with n nodes and roughly 2*n edges.
func buildTestGraph(n int) *DirectedGraph {
	g := NewDirectedGraph()
	nodes := make([]Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = g.NewNode()
		_ = g.AddNode(nodes[i])
	}
	// Create edges: each node connects to the next two (wrapping).
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if i != j {
			_ = g.SetEdge(simpleEdge{from: nodes[i], to: nodes[j]})
		}
		k := (i + 2) % n
		if i != k && j != k {
			_ = g.SetEdge(simpleEdge{from: nodes[i], to: nodes[k]})
		}
	}
	return g
}

type simpleEdge struct {
	from, to Node
}

func (e simpleEdge) From() Node { return e.from }
func (e simpleEdge) To() Node   { return e.to }

func BenchmarkFrom(b *testing.B) {
	g := buildTestGraph(20)
	ids := make([]int64, 0, 20)
	nodes := g.Nodes()
	for nodes.Next() {
		ids = append(ids, nodes.Node().ID())
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range ids {
			_ = g.From(id)
		}
	}
}

func BenchmarkNodes(b *testing.B) {
	g := buildTestGraph(20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Nodes()
	}
}

func BenchmarkEdges(b *testing.B) {
	g := buildTestGraph(20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Edges()
	}
}

func BenchmarkFrom_LargeGraph(b *testing.B) {
	g := buildTestGraph(100)
	ids := make([]int64, 0, 100)
	nodes := g.Nodes()
	for nodes.Next() {
		ids = append(ids, nodes.Node().ID())
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range ids {
			_ = g.From(id)
		}
	}
}

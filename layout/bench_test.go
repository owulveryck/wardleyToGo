package layout

import (
	"fmt"
	"testing"
)

// buildDenseGraph creates a graph with n nodes where each node connects to its
// next 3 nodes, producing O(3n) edges with multiple paths.
func buildDenseGraph(n int) *Graph {
	g := &Graph{
		Nodes: make([]Node, 0, n+1),
		Edges: make([]Edge, 0, 3*n),
	}
	g.Nodes = append(g.Nodes, Node{Name: "anchor", Kind: KindAnchor})
	for i := 1; i <= n; i++ {
		g.Nodes = append(g.Nodes, Node{Name: fmt.Sprintf("n%d", i), Kind: KindRegular})
	}
	g.Edges = append(g.Edges, Edge{From: "anchor", To: "n1"})
	for i := 1; i <= n; i++ {
		for j := 1; j <= 3; j++ {
			target := i + j
			if target <= n {
				g.Edges = append(g.Edges, Edge{
					From: fmt.Sprintf("n%d", i),
					To:   fmt.Sprintf("n%d", target),
				})
			}
		}
	}
	return g
}

func BenchmarkLayout(b *testing.B) {
	for _, size := range []int{10, 20, 50, 100} {
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			g := buildDenseGraph(size)
			l := New(DefaultOptions())
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = l.Layout(g)
			}
		})
	}
}

func BenchmarkForceSpread(b *testing.B) {
	for _, size := range []int{10, 20, 50, 100} {
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			g := buildDenseGraph(size)
			ranks, maxRank := topoRanks(g)
			pipelineMembers := make(map[string]bool)

			positions := make(map[string]float64, len(ranks))
			opts := DefaultOptions()
			minY := float64(opts.MinY)
			maxY := float64(opts.MaxY)
			for name, r := range ranks {
				if maxRank == 0 {
					positions[name] = minY
				} else {
					positions[name] = minY + float64(r)*(maxY-minY)/float64(maxRank)
				}
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Copy positions to avoid mutation across iterations
				pos := make(map[string]float64, len(positions))
				for k, v := range positions {
					pos[k] = v
				}
				pipelineParents := make(map[string]bool)
				_ = forceSpread(pos, g.Edges, pipelineMembers, pipelineParents, ranks, maxRank, opts)
			}
		})
	}
}

func BenchmarkEnforceRankOrder(b *testing.B) {
	g := buildDenseGraph(50)
	ranks, _ := topoRanks(g)
	opts := DefaultOptions()

	names := make([]string, 0, len(ranks))
	positions := make(map[string]float64, len(ranks))
	minY := float64(opts.MinY)
	maxY := float64(opts.MaxY)
	for name, r := range ranks {
		names = append(names, name)
		positions[name] = minY + float64(r)*(maxY-minY)/50
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipelineParents := make(map[string]bool)
		enforceRankOrder(positions, names, ranks, minY, maxY, float64(opts.MinSpacing), pipelineParents)
	}
}

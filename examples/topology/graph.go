package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

type node struct {
	id         int64
	visibility int
}

func (node *node) ID() int64 {
	return node.id
}

type edge struct {
	f, t       graph.Node
	visibility int
}

// From returns the from node of the edge.
func (e *edge) From() graph.Node {
	return e.f
}

// To returns the to node of the edge.
func (e *edge) To() graph.Node {
	return e.t
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (e *edge) ReversedEdge() graph.Edge {
	return &edge{
		f: e.t,
		t: e.f,
	}
}

type wmap struct {
	simple.DirectedGraph
}

type visibilityVisiter struct {
	g             graph.Directed
	maxVisibility int
}

func (v *visibilityVisiter) visit(srcNode graph.Node) {
	n := srcNode.(*node)
	// set the visibility of node n
	// given t_0, ..., t_n the nodes that can rean directly n (result of a call to g.To(n)) through edges e_0, ..., e_n
	// visibility is max((e_0.visibility + t_0.visibility), ..., (e_n.visibility + t_n.visibility))
	srcNodeVisibility := 0
	ts := v.g.To(n.ID())
	for ts.Next() {
		tX := ts.Node().(*node)
		eX := v.g.Edge(ts.Node().ID(), n.ID()).(*edge)
		eXVisibility := eX.visibility
		txToNVisibility := eXVisibility + tX.visibility
		if txToNVisibility > srcNodeVisibility {
			srcNodeVisibility = txToNVisibility
		}
	}
	// the node may have already been visited in some circumstances
	// in that case, we take the breatest visibility
	if srcNodeVisibility > n.visibility {
		n.visibility = srcNodeVisibility
	}
	if srcNodeVisibility > v.maxVisibility {
		v.maxVisibility = srcNodeVisibility
	}
}

// compute the visibility for each node and return the max visibility found
func setNodesVisibility(g graph.Directed) int {
	roots := findRoot(g)
	v := &visibilityVisiter{
		g: g,
	}
	bf := &traverse.BreadthFirst{
		Visit: v.visit,
	}
	for _, root := range roots {
		bf.Walk(g, root, nil)
	}
	return v.maxVisibility
}

func findRoot(g graph.Directed) []graph.Node {
	return nil

}

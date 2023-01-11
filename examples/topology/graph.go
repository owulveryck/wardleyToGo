package main

import (
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

type node struct {
	visibility    int
	evolutionStep int
	*wardley.Component
}

func (node *node) ID() int64 {
	return node.Component.ID()
}

type evolutionSetter struct {
	g           graph.Directed
	currentStep int
}

func (e *evolutionSetter) visit(srcNode graph.Node) {
	n := srcNode.(*node)
	n.evolutionStep = e.currentStep
	// if the node is a leaf (meaning the from is empty), move the cursor
	fs := e.g.From(n.ID())
	if fs.Len() == 0 {
		e.currentStep++

	}
}

type visibilityVisiter struct {
	g             graph.Directed
	maxVisibility int
}

func (v *visibilityVisiter) visit(srcNode graph.Node) {
	n := srcNode.(*node)
	// set the visibility of node n
	// given tX := t_0, ..., t_n the nodes that can rean directly n (result of a call to g.To(n)) through edges eX := e_0, ..., e_n
	// visibility is max((e_0.visibility + t_0.visibility), ..., (e_n.visibility + t_n.visibility))
	nVisibility := 0
	ts := v.g.To(n.ID())
	for ts.Next() {
		tX := ts.Node().(*node)
		eX := v.g.Edge(tX.ID(), n.ID()).(*wardley.Collaboration)
		eXVisibility := eX.Visibility
		rootToNVisibility := eXVisibility + tX.visibility
		if rootToNVisibility > nVisibility {
			nVisibility = rootToNVisibility
		}
	}
	// the node may have already been visited in some circumstances
	// in that case, we take the breatest visibility
	if nVisibility > n.visibility {
		n.visibility = nVisibility
	}
	if nVisibility > v.maxVisibility {
		v.maxVisibility = nVisibility
	}
}

// returns the max evolution
func setNodesEvolutionStep(g graph.Directed) int {
	roots := findRoot(g)
	e := &evolutionSetter{
		g: g,
	}
	df := &traverse.DepthFirst{
		Visit: e.visit,
	}
	for _, root := range roots {
		df.Walk(g, root, nil)
	}
	return e.currentStep
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

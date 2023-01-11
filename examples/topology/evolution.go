package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

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

// returns the max evolution
func setNodesEvolutionStep(g *scratchMapchMap) int {
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

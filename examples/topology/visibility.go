package main

import (
	"log"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

func setNodesvisibility(tempMap *mymap) {
	vis := setNodesVisibility(tempMap)
	nodes := tempMap.Nodes()
	for nodes.Next() {
		n := nodes.Node().(*node)
		log.Printf("%v: %v", n.c, n.visibility)
	}
	log.Printf("max visibility is %v", vis)
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
		eX := v.g.Edge(tX.ID(), n.ID()).(*edge)
		eXVisibility := eX.visibility
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

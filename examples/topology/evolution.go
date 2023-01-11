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

/*
func setNodesEvolution(tempMap *mymap) {
	vis := setNodesEvolutionStep(tempMap)
	nodes := tempMap.Nodes()
	for nodes.Next() {
		n := nodes.Node().(*node)
		for i := 0; i < n.evolutionStep; i++ {
			fmt.Printf("\t")
		}
		fmt.Println(n.c)
		//log.Printf("%v: %v", n.c, n.evolutionStep)
	}
	log.Printf("max evolution is %v", vis)
}
*/

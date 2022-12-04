package main

import (
	"fmt"
	"log"
	"os"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
)

func main() {

	g, err := parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(g, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	allShortestPaths := path.DijkstraAllPaths(g)
	roots := findRoot(g)
	leafs := findLeafs(g)
	var maxDepth int
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := g.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*edge).visibility
				}
				if currentVisibility > maxDepth {
					maxDepth = currentVisibility
				}
			}
		}
	}

	step := 100 / (maxDepth + 1)
	cs := &coordSetter{
		verticalStep: step,
	}
	for _, n := range roots {
		cs.walk(g, n.(*node), 0)
	}
	_ = step
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(g graph.Directed, n *node, visibility int) {
	n.visibility = visibility
	n.point.X = visibility * c.verticalStep
	from := g.From(n.ID())
	for from.Next() {
		c.walk(g, from.Node().(*node), g.Edge(n.ID(), from.Node().ID()).(*edge).visibility+visibility)
	}
}

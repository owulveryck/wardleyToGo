package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {

	inventory := make(map[string]*node, 0)
	edgeInventory := make([]*edge, 0)
	var link = regexp.MustCompile(`^\s*(.*\S)\s+(-+)\s+(.*)$`)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			log.Fatal("bad entry", scanner.Text())
		}
		if _, ok := inventory[elements[1]]; !ok {
			inventory[elements[1]] = &node{
				id:    int64(len(inventory)),
				label: elements[1],
			}
		}
		if _, ok := inventory[elements[3]]; !ok {
			inventory[elements[3]] = &node{
				id:    int64(len(inventory)),
				label: elements[3],
			}
		}
		edgeInventory = append(edgeInventory, &edge{
			from:       inventory[elements[1]],
			to:         inventory[elements[3]],
			visibility: len(elements[2]),
		})
	}

	g := simple.NewDirectedGraph()
	for _, n := range inventory {
		g.AddNode(n)
	}
	for _, e := range edgeInventory {
		g.SetEdge(e)
	}
	b, err := dot.Marshal(g, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	allShortestPaths := path.DijkstraAllPaths(g)
	roots := findRoot(g)
	leafs := findLeafs(g)
	var maxVis int
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := g.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*edge).visibility
				}
				if currentVisibility > maxVis {
					maxVis = currentVisibility
				}
			}
		}
	}

	fmt.Println(string(b))
	step := 100 / (maxVis + 1)
	cs := &coordSetter{
		verticalStep: step,
	}
	for _, n := range roots {
		cs.walk(g, n.(*node), 0)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
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

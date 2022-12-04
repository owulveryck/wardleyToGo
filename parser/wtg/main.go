package main

import (
	"image"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"

	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
)

func main() {

	m, err := initialize(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(m, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))

	allShortestPaths := path.DijkstraAllPaths(m)
	roots := findRoot(m)
	leafs := findLeafs(m)
	var maxDepth int
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := m.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*edge).visibility
				}
				if currentVisibility > maxDepth {
					maxDepth = currentVisibility
				}
			}
		}
	}

	step := 100 / maxDepth
	cs := &coordSetter{
		verticalStep: step,
	}
	for _, n := range roots {
		cs.walk(m, n, 0)
	}

	e, err := svgmap.NewEncoder(os.Stdout, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewWardleyStyle(svgmap.DefaultEvolution)
	e.Init(style)
	err = e.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(m *wardleyToGo.Map, n *wardley.Component, visibility int) {
	n.Placement.X = 50
	n.Placement.Y = visibility * c.verticalStep
	from := m.From(n.ID())
	for from.Next() {
		c.walk(m, from.Node().(*wardley.Component), m.Edge(n.ID(), from.Node().ID()).(*edge).visibility+visibility)
	}
}

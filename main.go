package main

import (
	"fmt"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/internal/encoding/svg"
	"github.com/owulveryck/wardleyToGo/internal/parser"
	"gonum.org/v1/gonum/graph/path"
)

func main() {
	width := 1000
	height := 700
	padLeft := 25
	padBottom := 30

	p := parser.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	svgmap.Encode(m, os.Stdout, width, height, padLeft, padBottom)
	allShortest := path.DijkstraAllPaths(m)
	it1 := m.Nodes()
	it2 := m.Nodes()
	for it1.Next() {
		from := it1.Node()
		for it2.Next() {
			to := it2.Node()
			if to == from {
				continue
			}
			p, _, _ := allShortest.Between(from.ID(), to.ID())
			if len(p) == 0 {
				continue
			}
			fmt.Print(p[0])
			for i := 1; i < len(p); i++ {
				fmt.Printf(" -- %v", p[i])
			}
			fmt.Println("")
		}
		it2.Reset()
	}
}

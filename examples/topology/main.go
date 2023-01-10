package main

import (
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/wtg"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"
)

func main() {
	p := wtg.NewParser()

	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			log.Println(err)
		}
	}
	sorted, err := topo.SortStabilized(p.WMap, func(ns []graph.Node) {
		log.Println(ns)
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range sorted {
		log.Println(n)
	}
}

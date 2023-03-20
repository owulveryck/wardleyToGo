package main

import (
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {
	p := wtg.NewParser()
	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	dict := make(map[int64]string, 0)
	fmt.Printf("title %v\n\n", p.WMap.Title)
	allComponents := p.WMap.Nodes()
	for allComponents.Next() {
		if n, ok := allComponents.Node().(*wardley.Component); ok {
			fmt.Printf("component %v [%v, %v]\n", n.Label, float64(100-n.Placement.Y)/100, float64(n.Placement.X)/100)
			switch n.Type {
			case wardley.BuildComponent:
				fmt.Printf("build %v\n", n.Label)
			case wardley.BuyComponent:
				fmt.Printf("buy %v\n", n.Label)
			case wardley.OutsourceComponent:
				fmt.Printf("outsource %v\n", n.Label)
			}
			dict[n.ID()] = n.Label
		}
		if n, ok := allComponents.Node().(*wardley.EvolvedComponent); ok {
			fmt.Printf("evolve %v %v\n", n.Label, float64(n.Placement.X)/100)
			dict[n.ID()] = n.Label
		}
	}
	fmt.Println("")
	allEdges := p.WMap.Edges()
	for allEdges.Next() {
		if e, ok := allEdges.Edge().(*wardley.Collaboration); ok {
			if dict[e.F.ID()] == dict[e.T.ID()] {
				continue
			}
			fmt.Printf("%v -> %v\n", dict[e.F.ID()], dict[e.T.ID()])
		}
	}
	fmt.Println("\nstyle wardley")
}

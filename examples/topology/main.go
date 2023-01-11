package main

import (
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
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
	setCoords(*p.WMap)
	//setNodesvisibility(*p.WMap)
	//setNodesEvolution(*p.WMap)
	/*
		nodes := p.WMap.Nodes()
		nodess := make([]*wardley.Component, nodes.Len())
		for i := 0; nodes.Next(); i++ {
			nodess[i] = nodes.Node().(*wardley.Component)
		}
		n := nodeSort{
			g:     p.WMap,
			nodes: nodess,
		}
		sort.Sort(n)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range n.nodes {
			log.Println(n)
		}
	*/
}

type nodeSort struct {
	g     *wardleyToGo.Map
	nodes []*wardley.Component
}

// Len is the number of elements in the collection.
func (n nodeSort) Len() int {
	return len(n.nodes)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (ns nodeSort) Less(i int, j int) bool {
	// find the top common nodes
	// get both path from the common node and the nodes, compute the visibility and then return
	for _, n := range ns.nodes {
		if topo.PathExistsIn(ns.g, n, ns.nodes[i]) && topo.PathExistsIn(ns.g, n, ns.nodes[j]) {
			log.Printf("common node between %v and %v is %v", ns.nodes[i], ns.nodes[j], n)
		}
	}
	return topo.PathExistsIn(ns.g, ns.nodes[i], ns.nodes[j])
}

// Swap swaps the elements with indexes i and j.
func (n nodeSort) Swap(i int, j int) {
	n.nodes[i], n.nodes[j] = n.nodes[j], n.nodes[i]
}

package wtg

import (
	"testing"
)

func TestParser_consolidateMap(t *testing.T) {
	t.Run("simple map", consolidateMapOK)
	t.Run("simple map with evolution", consolidateMapWithEvolution)
}
func consolidateMapOK(t *testing.T) {
	nodes := `
	node1 - node2
	`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Fatal(err)
	}
	err = p.consolidateMap()
	if err != nil {
		t.Fatal(err)
	}
	nodeIT := p.WMap.Nodes()
	if nodeIT.Len() != 2 {
		t.Fatal("expected two nodes")
	}
	edgeIT := p.WMap.Edges()
	if edgeIT.Len() != 1 {
		t.Fatal("expected one edge")
	}

}

func consolidateMapWithEvolution(t *testing.T) {
	nodes := `
	node1 - node2
	node1: |.x.|.>.|...|...|
	`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Fatal(err)
	}
	err = p.consolidateMap()
	if err != nil {
		t.Fatal(err)
	}
	nodeIT := p.WMap.Nodes()
	if nodeIT.Len() != 3 {
		t.Fatal("expected three nodes")
	}
	edgeIT := p.WMap.Edges()
	if edgeIT.Len() != 2 {
		t.Fatal("expected two edges")
	}

}

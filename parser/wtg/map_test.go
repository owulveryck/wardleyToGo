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
	inv := &inventorier{}
	err := inv.init(nodes)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.start()
	if err != nil {
		t.Fatal(err)
	}
	m, err := consolidateMap(inv.nodeInventory, inv.edgeInventory)
	if err != nil {
		t.Fatal(err)
	}
	nodeIT := m.Nodes()
	if nodeIT.Len() != 2 {
		t.Fatal("expected two nodes")
	}
	edgeIT := m.Edges()
	if edgeIT.Len() != 1 {
		t.Fatal("expected one edge")
	}

}

func consolidateMapWithEvolution(t *testing.T) {
	nodes := `
	node1 - node2
	node1: |.x.|.>.|...|...|
	`
	inv := &inventorier{}
	err := inv.init(nodes)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.start()
	m, err := consolidateMap(inv.nodeInventory, inv.edgeInventory)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	nodeIT := m.Nodes()
	if nodeIT.Len() != 3 {
		t.Fatal("expected three nodes")
	}
	edgeIT := m.Edges()
	if edgeIT.Len() != 2 {
		t.Fatal("expected two edges")
	}

}

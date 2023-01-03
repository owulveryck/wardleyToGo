package wtg

import (
	"testing"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestInventory(t *testing.T) {
	t.Run("empty", empty)
	t.Run("all commented", allCommented)
	t.Run("one node", oneNode)
	t.Run("two nodes", twoNodes)
	t.Run("one edge", oneEdge)
	t.Run("simple evolution", simpleEvolution)
	t.Run("simple evolution with comment", simpleEvolutionWithComment)
	t.Run("visibility on nil node", visibilityOnNilNode)
	t.Run("evolution on nil node", evolutionOnNilNode)
	t.Run("type", types)
}

func allCommented(t *testing.T) {
	nodes := `
	/*
	dsadsadsa
	dsadsadsadsa
	*/`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory != nil && len(p.edgeInventory) != 0 {
		t.Log(p.edgeInventory)
		t.Error("edge inventory should be empty")
	}
	if len(p.nodeInventory) != 0 {
		t.Error("node inventory should be empty")

	}
}

func empty(t *testing.T) {
	nodes := ` `
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory != nil && len(p.edgeInventory) != 0 {
		t.Log(p.edgeInventory)
		t.Error("edge inventory should be empty")
	}
	if len(p.nodeInventory) != 0 {
		t.Error("node inventory should be empty")

	}
}

func oneNode(t *testing.T) {
	nodes := `
		node1
		`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory != nil && len(p.edgeInventory) != 0 {
		t.Log(p.edgeInventory)
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.nodeInventory["node1"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node1")
	}
	//t.Error("test")
}
func twoNodes(t *testing.T) {
	nodes := `
		node1
		node2
		`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory != nil && len(p.edgeInventory) != 0 {
		t.Log(p.edgeInventory)
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.nodeInventory["node1"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node1")
	}
	if _, ok := p.nodeInventory["node2"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node2")
	}
	//t.Error("test")
}
func oneEdge(t *testing.T) {
	nodes := `
		node1 - node2
		`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory == nil || len(p.edgeInventory) != 1 {
		t.Log(p.edgeInventory)
		t.Error("edge inventory should not be empty")
		if p.edgeInventory[0].F != p.nodeInventory["node1"] {
			t.Error("bad from node")
		}
		if p.edgeInventory[0].T != p.nodeInventory["node2"] {
			t.Error("bad to node")
		}
		if p.edgeInventory[0].Visibility != 1 {
			t.Error("bad visibility")
		}
	}
	if _, ok := p.nodeInventory["node1"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node1")
	}
	if _, ok := p.nodeInventory["node2"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node2")
	}
	//t.Error("test")
}
func simpleEvolution(t *testing.T) {
	nodes := `
	node1: |.x.|...|...|...|
		`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory == nil || len(p.edgeInventory) != 0 {
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.nodeInventory["node1"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node1")
		return
	}
	if p.nodeInventory["node1"].Placement.X != 9 {
		t.Errorf("expected plactement to be 9, but is %v", p.nodeInventory["node1"].Placement.X)
	}
	//t.Error("test")
}
func simpleEvolutionWithComment(t *testing.T) {
	nodes := `
	node1: |.x.|...|...|...| // comment
		`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.edgeInventory == nil || len(p.edgeInventory) != 0 {
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.nodeInventory["node1"]; !ok {
		t.Log(p.nodeInventory)
		t.Error("should have node1")
		return
	}
	if p.nodeInventory["node1"].Placement.X != 9 {
		t.Errorf("expected plactement to be 9, but is %v", p.nodeInventory["node1"].Placement.X)
	}
	//t.Error("test")
}

func visibilityOnNilNode(t *testing.T) {
	nodes := `-- bla`
	p := NewParser()
	err := p.inventory(nodes)
	if err == nil {
		t.Error("expected error")
	}
}
func evolutionOnNilNode(t *testing.T) {
	nodes := `
	|...|...|...|...|`
	p := NewParser()
	err := p.inventory(nodes)
	if err == nil {
		t.Error("expected error")
	}
}
func types(t *testing.T) {
	nodes := `
	build: {
		type: build
	}
	buy: {
		type: buy
	}
	outsource: {
		type: outsource
	}
	`
	p := NewParser()
	err := p.inventory(nodes)
	if err != nil {
		t.Error("expected error")
	}
	for k, v := range p.nodeInventory {
		switch k {
		case "build":
			if v.Type != wardley.BuildComponent {
				t.Error("exected build type")
			}
		case "buy":
			if v.Type != wardley.BuyComponent {
				t.Error("exected buy type")
			}
		case "outsource":
			if v.Type != wardley.OutsourceComponent {
				t.Error("exected outsource type")
			}
		}
	}
}

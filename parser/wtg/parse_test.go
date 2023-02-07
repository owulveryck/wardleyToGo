package wtg

import (
	"testing"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func FuzzParse(f *testing.F) {
	full := `
	title: sample map // title is optional
/***************
  value chain 
****************/

business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

/***************
  definitions 
****************/

// you can inline the evolution
business: |....|....|...x.|.........|

public: |....|....|....|.x...|

// or create blocks
cup of tea: {
  evolution: |....|....|..x..|........|
  color: Green // you can set colors
}
cup: {
  type: buy
  evolution: |....|....|....|......x....|
}
tea: {
  type: buy
  evolution: |....|....|....|.....x....|
}
hot water: {
  evolution: |....|....|....|....x....|
  color: Blue
}
water: {
  type: outsource
  evolution: |....|....|....|.....x....|
}

// you can set the evolution with a >
kettle: {
  type: build
  evolution: |...|...x.|..>.|.......|
}
power: {
  type: outsource
  evolution: |...|...|....x|.....>..|
}

stage1: genesis / concept
stage2: custom / emerging
stage3: product / converging
stage4: commodity / accepted
	`
	//	f.SkipNow()
	testcases := []string{" ", "a - b", full}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		p := NewParser()
		err := p.parse(orig)
		t.Log(err)
		for _, err := range p.InvalidEntries {
			t.Log(err)
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("empty map", testParseEmpty)
	t.Run("bad content", testParseBadContent)
	t.Run("complete test ok", testParseCompleteOk)
}
func testParseBadContent(t *testing.T) {
	p := NewParser()
	err := p.parse("--")
	if err == nil {
		t.Errorf("expected an error")
	}
}
func testParseEmpty(t *testing.T) {
	p := NewParser()
	err := p.parse("")
	if err == nil {
		t.Errorf("expected an error")
	}
}
func testParseCompleteOk(t *testing.T) {
	sampleTeahop := `
	title: sample map // title is optional
/***************
  value chain 
****************/

business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

/***************
  definitions 
****************/

// you can inline the evolution
business: |....|....|...x.|.........|

public: |....|....|....|.x...|

// or create blocks
cup of tea: {
  evolution: |....|....|..x..|........|
  color: Green // you can set colors
}
cup: {
  type: buy
  evolution: |....|....|....|......x....|
}
tea: {
  type: buy
  evolution: |....|....|....|.....x....|
}
hot water: {
  evolution: |....|....|....|....x....|
  color: Blue
}
water: {
  type: outsource
  evolution: |....|....|....|.....x....|
}

// you can set the evolution with a >
kettle: {
  type: build
  evolution: |...|...x.|..>.|.......|
}
power: {
  type: outsource
  evolution: |...|...|....x|.....>..|
}

stage1: genesis / concept
stage2: custom / emerging
stage3: product / converging
stage4: commodity / accepted
	`
	p := NewParser()
	err := p.parse(sampleTeahop)
	if err != nil {
		t.Fatal(err)
	}
	edgesIT := p.WMap.Edges()
	if edgesIT.Len() != 10 {
		t.Fatalf("expected 10 links got %v", edgesIT.Len())
	}
	nodesIT := p.WMap.Nodes()
	if nodesIT.Len() != 11 {
		t.Fatalf("expected 11 nodes got %v", nodesIT.Len())
	}
}

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
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory != nil && len(p.EdgeInventory) != 0 {
		t.Log(p.EdgeInventory)
		t.Error("edge inventory should be empty")
	}
	if len(p.NodeInventory) != 0 {
		t.Error("node inventory should be empty")

	}
}

func empty(t *testing.T) {
	nodes := ` `
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}

	if p.EdgeInventory != nil && len(p.EdgeInventory) != 0 {
		t.Log(p.EdgeInventory)
		t.Error("edge inventory should be empty")
	}
	if len(p.NodeInventory) != 0 {
		t.Error("node inventory should be empty")

	}
}

func oneNode(t *testing.T) {
	nodes := `
		node1
		`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory != nil && len(p.EdgeInventory) != 0 {
		t.Log(p.EdgeInventory)
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.NodeInventory["node1"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node1")
	}
	//t.Error("test")
}
func twoNodes(t *testing.T) {
	nodes := `
		node1
		node2
		`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory != nil && len(p.EdgeInventory) != 0 {
		t.Log(p.EdgeInventory)
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.NodeInventory["node1"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node1")
	}
	if _, ok := p.NodeInventory["node2"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node2")
	}
	//t.Error("test")
}
func oneEdge(t *testing.T) {
	nodes := `
		node1 - node2
		`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory == nil || len(p.EdgeInventory) != 1 {
		t.Log(p.EdgeInventory)
		t.Error("edge inventory should not be empty")
		if p.EdgeInventory[0].F != p.NodeInventory["node1"] {
			t.Error("bad from node")
		}
		if p.EdgeInventory[0].T != p.NodeInventory["node2"] {
			t.Error("bad to node")
		}
		if p.EdgeInventory[0].Visibility != 1 {
			t.Error("bad visibility")
		}
	}
	if _, ok := p.NodeInventory["node1"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node1")
	}
	if _, ok := p.NodeInventory["node2"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node2")
	}
	//t.Error("test")
}
func simpleEvolution(t *testing.T) {
	nodes := `
	node1: |.x.|...|...|...|
		`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory == nil || len(p.EdgeInventory) != 0 {
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.NodeInventory["node1"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node1")
		return
	}
	if p.NodeInventory["node1"].Placement.X != 9 {
		t.Errorf("expected plactement to be 9, but is %v", p.NodeInventory["node1"].Placement.X)
	}
	//t.Error("test")
}
func simpleEvolutionWithComment(t *testing.T) {
	nodes := `
	node1: |.x.|...|...|...| // comment
		`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err != nil {
		t.Error(err)
	}
	if p.EdgeInventory == nil || len(p.EdgeInventory) != 0 {
		t.Error("edge inventory should be empty")
	}
	if _, ok := p.NodeInventory["node1"]; !ok {
		t.Log(p.NodeInventory)
		t.Error("should have node1")
		return
	}
	if p.NodeInventory["node1"].Placement.X != 9 {
		t.Errorf("expected plactement to be 9, but is %v", p.NodeInventory["node1"].Placement.X)
	}
	//t.Error("test")
}

func visibilityOnNilNode(t *testing.T) {
	nodes := `-- bla`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
	if err == nil {
		t.Error("expected error")
	}
}
func evolutionOnNilNode(t *testing.T) {
	nodes := `
	|...|.x.|...|...|`
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error(err)
	}
	err = p.Run()
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
	p := NewInventory()
	err := p.init(nodes)
	if err != nil {
		t.Error("expected error")
	}
	err = p.Run()
	if err != nil {
		t.Error("expected error")
	}
	for k, v := range p.NodeInventory {
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

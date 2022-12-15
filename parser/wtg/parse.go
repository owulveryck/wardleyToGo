package wtg

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
	"regexp"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/path"
)

const (
	nodergxp = `[\p{L}|\s]+`
)

var (
	node      = regexp.MustCompile(`^\s*(` + nodergxp + `):\s{\s*$`)
	endnode   = regexp.MustCompile(`^\s*}\s*$`)
	evolution = regexp.MustCompile(`^\s*evolution:\s*(|.*x?.*|.*x?.*|.*x?.*|.*x?.*|)\s*$`)
	nodeType  = regexp.MustCompile(`^\s*type:\s*(.*)\s*$`)
	link      = regexp.MustCompile(`^\s*(.*\S)\s+(-+)\s+(.*)$`)
)

type Parser struct {
	nodeInventory map[string]*wardley.Component
	edgeInventory []*wardley.Collaboration
	currentNode   *wardley.Component
	WMap          *wardleyToGo.Map
}

func NewParser() *Parser {
	return &Parser{
		nodeInventory: make(map[string]*wardley.Component, 0),
		edgeInventory: make([]*wardley.Collaboration, 0),
		WMap:          wardleyToGo.NewMap(0),
	}
}

func (p *Parser) DumpComponents(w io.Writer) {
	for n := range p.nodeInventory {
		fmt.Fprintf(w, "%v\n", n)
	}
}

func (p *Parser) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		p.parseComponents(scanner.Text())
		p.parseValueChain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err := p.consolidateMap()
	if err != nil {
		return err
	}
	err = setYCoords(p.WMap)
	if err != nil {
		return nil
	}
	return nil
}

func (p *Parser) parseComponents(s string) error {

	elements := node.FindStringSubmatch(s)
	if len(elements) == 2 {
		if p.currentNode != nil {
			return fmt.Errorf("parser error, nested nodes unsupported (%v is in %v)", elements[1], p.currentNode.Label)
		}
		if _, ok := p.nodeInventory[elements[1]]; !ok {
			c := wardley.NewComponent(int64(len(p.nodeInventory)))
			c.Label = elements[1]
			c.Placement = image.Pt(0, 50)
			p.nodeInventory[elements[1]] = c
		}
		p.currentNode = p.nodeInventory[elements[1]]
	}
	elements = endnode.FindStringSubmatch(s)
	if len(elements) == 1 {
		p.currentNode = nil
	}
	elements = evolution.FindStringSubmatch(s)
	if len(elements) == 2 && p.currentNode != nil {
		placement, err := computePlacement(elements[1])
		if err != nil {
			return err
		}
		p.currentNode.Placement.X = placement
	}
	elements = nodeType.FindStringSubmatch(s)
	if len(elements) == 2 && p.currentNode != nil {
		switch elements[1] {
		case "build":
			p.currentNode.Type = wardley.BuildComponent
		case "buy":
			p.currentNode.Type = wardley.BuyComponent
		case "outsource":
			p.currentNode.Type = wardley.OutsourceComponent
		}
	}
	return nil
}

func (p *Parser) parseValueChain(s string) error {
	// do not parse value chain in the node definition you fool
	if p.currentNode != nil {
		return nil
	}

	elements := link.FindStringSubmatch(s)
	if len(elements) != 4 {
		// log.Fatal("bad entry", scanner.Text())
		return nil
	}
	if _, ok := p.nodeInventory[elements[1]]; !ok {
		c := wardley.NewComponent(int64(len(p.nodeInventory)))
		c.Label = elements[1]
		c.Placement = image.Pt(0, 50)
		p.nodeInventory[elements[1]] = c
	}
	if _, ok := p.nodeInventory[elements[3]]; !ok {
		c := wardley.NewComponent(int64(len(p.nodeInventory)))
		c.Label = elements[3]
		c.Placement = image.Pt(0, 50)
		p.nodeInventory[elements[3]] = c
	}
	p.edgeInventory = append(p.edgeInventory, &wardley.Collaboration{
		F:          p.nodeInventory[elements[1]],
		T:          p.nodeInventory[elements[3]],
		Type:       wardley.RegularEdge,
		Visibility: len(elements[2]),
	})
	return nil
}

func (p *Parser) consolidateMap() error {
	for _, c := range p.nodeInventory {
		err := p.WMap.AddComponent(c)
		if err != nil {
			return err
		}
	}
	for _, e := range p.edgeInventory {
		err := p.WMap.SetCollaboration(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func setYCoords(m *wardleyToGo.Map) error {

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
					currentVisibility += e.(*wardley.Collaboration).Visibility
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
	nroots := len(roots)
	hsteps := 100 / (nroots + 1)
	for i, n := range roots {
		if n.Placement.X == 0 {
			n.Placement.X = hsteps * (i + 1)
		}
		cs.walk(m, n, 0)
	}

	return nil
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(m *wardleyToGo.Map, n *wardley.Component, visibility int) {
	n.Placement.Y = visibility * c.verticalStep
	from := m.From(n.ID())
	hsteps := 100 / (from.Len() + 1)
	i := 1
	for from.Next() {
		if from.Node().(*wardley.Component).Placement.X == 0 {
			from.Node().(*wardley.Component).Placement.X = hsteps * i
		}
		c.walk(m, from.Node().(*wardley.Component), m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
		i++
	}
}

func computePlacement(s string) (int, error) {
	currentStage := -1
	currentCursor := 0
	stages := make([]int, 5)
	cursor := 0
	stage := 0
	for _, c := range s {
		if c == '|' {
			currentCursor = 0
			currentStage++
			continue
		}
		if c != 'x' {
			currentCursor++
		}
		if c == 'x' {
			cursor = currentCursor
			stage = currentStage
		}
		stages[currentStage]++
	}
	stagePositions := []float64{0, 17.4, 40, 70, 100}
	position := stagePositions[stage] + (stagePositions[stage+1]-stagePositions[stage])*float64(cursor)/float64(stages[stage])
	return int(math.Round(position)), nil
}

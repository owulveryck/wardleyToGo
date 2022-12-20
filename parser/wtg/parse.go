package wtg

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"regexp"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
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
		placement, evolvedPosition, err := computePlacement(elements[1])
		if err != nil {
			return err
		}
		p.currentNode.Placement.X = placement
		if evolvedPosition != 0 {
			evolvedC := wardley.NewEvolvedComponent(p.currentNode.ID() + 1000) // FIXME
			evolvedC.Placement.X = evolvedPosition
			evolvedC.Placement.Y = p.currentNode.Placement.Y
			evolvedC.Label = p.currentNode.Label
			p.WMap.AddComponent(evolvedC)
			p.edgeInventory = append(p.edgeInventory, &wardley.Collaboration{
				F:    p.currentNode,
				T:    evolvedC,
				Type: wardley.EvolvedComponentEdge,
			})
		}
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

package wtg

import (
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

type Parser struct {
	nodeInventory  map[string]*wardley.Component
	edgeInventory  []*wardley.Collaboration
	currentNode    *wardley.Component
	currentEdge    *wardley.Collaboration
	visibilityOnly bool
	WMap           *wardleyToGo.Map
}

func NewParser() *Parser {
	return &Parser{
		nodeInventory:  make(map[string]*wardley.Component, 0),
		edgeInventory:  make([]*wardley.Collaboration, 0),
		visibilityOnly: true,
		WMap:           wardleyToGo.NewMap(0),
	}
}

func (p *Parser) Parse(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.parse(string(b))
}
func (p *Parser) parse(s string) error {

	// TODO consolidate map
	// TODO compute Y
	err := p.inventory(s)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) inventory(s string) error {
	l := newLexer(s, startState)
	l.Start()
	inComment := false
	for tok := range l.tokens {
		if inComment {
			if tok.Type == endBlockCommentToken {
				inComment = false
			}
			continue
		}
		switch tok.Type {
		case identifierToken:
			if p.currentEdge != nil {
				p.currentEdge.F = p.currentNode
				p.currentEdge.T = p.upsertNode(tok.Value)
				p.currentNode = nil
				p.currentEdge = nil
				continue
			}
			p.currentNode = p.upsertNode(tok.Value)
		case visibilityToken:
			if p.currentNode == nil {
				return errors.New("cannot set visibility on a nil source node")
			}
			p.currentEdge = p.insertEdge(tok.Value)
		case evolutionItem:
			if p.currentNode == nil {
				return errors.New("cannot set evolution on a nil node")
			}
			pos, evolutionPos, err := computeEvolutionPosition(tok.Value)
			if err != nil {
				return fmt.Errorf("cannot compute evolution for %v (%w)", tok.Value, err)
			}
			p.currentNode.Placement.X = pos
			p.currentNode.Configured = true
			p.visibilityOnly = false
			_ = evolutionPos
		case eofToken:
		case colonToken:
		case commentToken:
		case startBlockCommentToken:
			inComment = true
		case typeItem:
			if p.currentNode == nil {
				return errors.New("cannot set type on a nil node")
			}
			switch tok.Value {
			case "build":
				p.currentNode.Type = wardley.BuildComponent
			case "buy":
				p.currentNode.Type = wardley.BuyComponent
			case "outsource":
				p.currentNode.Type = wardley.OutsourceComponent
			default:
				log.Printf("unhandled component type: %v", tok.Value)
			}
		case typeToken:
		case startBlockToken:
			p.visibilityOnly = false
		case endBlockToken:
		default:
			log.Printf("unhandled element: %v", tok.Value)
		}
	}
	return nil
}

func (p *Parser) upsertNode(s string) *wardley.Component {
	if _, ok := p.nodeInventory[s]; !ok {
		c := wardley.NewComponent(int64(len(p.nodeInventory)))
		c.Label = s
		c.Placement = image.Pt(0, 50)
		p.nodeInventory[s] = c
	}
	return p.nodeInventory[s]
}
func (p *Parser) insertEdge(s string) *wardley.Collaboration {
	p.edgeInventory = append(p.edgeInventory, &wardley.Collaboration{
		Type:       wardley.RegularEdge,
		Visibility: len(s),
	})
	return p.edgeInventory[len(p.edgeInventory)-1]
}

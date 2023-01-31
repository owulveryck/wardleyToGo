package wtg

import (
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"strings"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

type Parser struct {
	nodeInventory   map[string]*wardley.Component
	edgeInventory   []*wardley.Collaboration
	currentNode     *wardley.Component
	currentEdge     *wardley.Collaboration
	visibilityOnly  bool
	WMap            *wardleyToGo.Map
	EvolutionStages []svgmap.Evolution
	ImageSize       image.Rectangle
	MapSize         image.Rectangle
	// InvalidEntries reports any invalid of unkonwn token
	InvalidEntries []error
}

func NewParser() *Parser {
	return &Parser{
		nodeInventory:   make(map[string]*wardley.Component, 0),
		edgeInventory:   make([]*wardley.Collaboration, 0),
		visibilityOnly:  true,
		WMap:            wardleyToGo.NewMap(0),
		EvolutionStages: svgmap.DefaultEvolution,
		InvalidEntries:  make([]error, 0),
	}
}

func (p *Parser) DumpComponents(w io.Writer) {
	for n := range p.nodeInventory {
		fmt.Fprintf(w, "%v\n", n)
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

	err := p.inventory(s)
	if err != nil {
		return fmt.Errorf("error in parsing: %w", err)
	}
	if len(p.nodeInventory) == 0 {
		return fmt.Errorf("no map")
	}
	err = p.consolidateMap()
	if err != nil {
		return fmt.Errorf("cannot consolidate map: %w", err)
	}
	SetCoords(*p.WMap, true)
	SetLabelAnchor(*p.WMap)
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
			p.currentNode.EvolutionPos = evolutionPos
		case eofToken:
		case colonToken:
		case evolutionToken:
		case commentToken:
		case startBlockCommentToken:
			inComment = true
		case stage1Token:
		case stage1Item:
			p.EvolutionStages[0].Label = tok.Value
		case stage2Token:
		case stage2Item:
			p.EvolutionStages[1].Label = tok.Value
		case stage3Token:
		case stage3Item:
			p.EvolutionStages[2].Label = tok.Value
		case stage4Token:
		case stage4Item:
			p.EvolutionStages[3].Label = tok.Value
		case titleToken:
		case titleItem:
			p.WMap.Title = strings.TrimSpace(tok.Value)
		case colorToken:
		case labelItem:
			if p.currentNode == nil {
				return errors.New("cannot set type on a nil node")
			}
			err := setLabelPlacement(p.currentNode, tok.Value)
			if err != nil {
				return err
			}
		case colorItem:
			if p.currentNode == nil {
				return errors.New("cannot set type on a nil node")
			}
			if col, ok := Colors[tok.Value]; ok {
				p.currentNode.Color = col
				continue
			}
			p.InvalidEntries = append(p.InvalidEntries, fmt.Errorf("unknown color %v", tok.Value))
		case typeItem:
			if p.currentNode == nil {
				return errors.New("cannot set type on a nil node")
			}
			switch tok.Value {
			case "pipeline":
				p.currentNode.Type = wardley.PipelineComponent
				p.currentNode.LabelPlacement.Y = -15
				p.currentNode.LabelPlacement.X = 0
				p.currentNode.Anchor = wardley.AdjustMiddle
			case "build":
				p.currentNode.Type = wardley.BuildComponent
			case "buy":
				p.currentNode.Type = wardley.BuyComponent
			case "outsource":
				p.currentNode.Type = wardley.OutsourceComponent
			default:
				p.InvalidEntries = append(p.InvalidEntries, fmt.Errorf("unhandled component type: %v", tok.Value))
			}
		case typeToken:
		case startBlockToken:
			p.visibilityOnly = false
		case endBlockToken:
		case singleLineCommentSeparator:
		default:
			p.InvalidEntries = append(p.InvalidEntries, fmt.Errorf("unhandled element: %v (%x)", tok.Value, tok.Value))
		}
	}
	return l.Err
}

func (p *Parser) upsertNode(s string) *wardley.Component {
	if _, ok := p.nodeInventory[s]; !ok {
		c := wardley.NewComponent(int64(len(p.nodeInventory)))
		c.Label = s
		c.LabelPlacement.X = 10
		//c.LabelPlacement.Y = 6
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

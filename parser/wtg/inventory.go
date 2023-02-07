package wtg

import (
	"image"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

// Inventory is the result of the parsing of the elements
// it does not create a consolidated map
// this can be used in an LSP
type Inventory struct {
	NodeInventory   map[string]*wardley.Component
	EdgeInventory   []*wardley.Collaboration
	tokens          []token
	offset          int
	EvolutionStages []svgmap.Evolution
	Title           string
}

func NewInventory() *Inventory {
	return &Inventory{
		NodeInventory:   make(map[string]*wardley.Component),
		EdgeInventory:   make([]*wardley.Collaboration, 0),
		tokens:          make([]token, 0),
		EvolutionStages: svgmap.DefaultEvolution,
	}
}

func (inv *Inventory) init(src string) error {
	l := newLexer(src, startState)
	l.Start()
	for tok := range l.tokens {
		switch {
		case tok.Type == eofToken:
			return nil
		case l.Err != nil:
			return l.Err
		default:
			inv.tokens = append(inv.tokens, tok)
		}
	}
	return nil
}

// peek the next token+delta
func (inv *Inventory) peek(delta int) token {
	if inv.offset+delta < len(inv.tokens) {
		return inv.tokens[inv.offset+delta]
	}
	return token{Type: eofToken}
}

func (inv *Inventory) currentToken() token {
	if inv.offset < len(inv.tokens) {
		return inv.tokens[inv.offset]
	}
	return token{Type: eofToken}
}
func upsertAnchor(src, anchor *wardley.Component) {
	src.Type = wardley.PipelineComponent
	anchor.PipelineReference = src
	for i := range src.PipelinedComponents {
		if src.PipelinedComponents[i] == anchor {
			return
		}
	}
	src.PipelinedComponents = append(src.PipelinedComponents, anchor)
}

func (inv *Inventory) upsertNode(s string) *wardley.Component {
	if _, ok := inv.NodeInventory[s]; !ok {
		c := wardley.NewComponent(int64(len(inv.NodeInventory)))
		c.Label = s
		c.LabelPlacement.X = 10
		//c.LabelPlacement.Y = 6
		c.Placement = image.Pt(0, 50)
		inv.NodeInventory[s] = c
	}
	return inv.NodeInventory[s]
}
func (inv *Inventory) insertEdge(s string) *wardley.Collaboration {
	inv.EdgeInventory = append(inv.EdgeInventory, &wardley.Collaboration{
		Type:       wardley.RegularEdge,
		Visibility: len(s),
	})
	return inv.EdgeInventory[len(inv.EdgeInventory)-1]
}

package wtg

import (
	"image"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

type inventorier struct {
	l               *lexer
	nodeInventory   map[string]*wardley.Component
	edgeInventory   []*wardley.Collaboration
	tokens          []token
	offset          int
	evolutionStages []svgmap.Evolution
	title           string
}

func (inv *inventorier) init(src string) error {
	inv.nodeInventory = make(map[string]*wardley.Component)
	inv.edgeInventory = make([]*wardley.Collaboration, 0)
	inv.tokens = make([]token, 0)
	inv.evolutionStages = svgmap.DefaultEvolution
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
func (inv *inventorier) peek(delta int) token {
	if inv.offset+delta < len(inv.tokens) {
		return inv.tokens[inv.offset+delta]
	}
	return token{Type: eofToken}
}

func (inv *inventorier) currentToken() token {
	if inv.offset < len(inv.tokens) {
		return inv.tokens[inv.offset]
	}
	return token{Type: eofToken}
}
func upsertAnchor(src, anchor *wardley.Component) {
	anchor.PipelineReference = src
	for i := range src.PipelinedComponents {
		if src.PipelinedComponents[i] == anchor {
			return
		}
	}
	src.PipelinedComponents = append(src.PipelinedComponents, anchor)
}

func (inv *inventorier) upsertNode(s string) *wardley.Component {
	if _, ok := inv.nodeInventory[s]; !ok {
		c := wardley.NewComponent(int64(len(inv.nodeInventory)))
		c.Label = s
		c.LabelPlacement.X = 10
		//c.LabelPlacement.Y = 6
		c.Placement = image.Pt(0, 50)
		inv.nodeInventory[s] = c
	}
	return inv.nodeInventory[s]
}
func (inv *inventorier) insertEdge(s string) *wardley.Collaboration {
	inv.edgeInventory = append(inv.edgeInventory, &wardley.Collaboration{
		Type:       wardley.RegularEdge,
		Visibility: len(s),
	})
	return inv.edgeInventory[len(inv.edgeInventory)-1]
}

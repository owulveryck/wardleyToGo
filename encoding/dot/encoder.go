package dot

import (
	"fmt"
	"image"
	"io"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

type Encoder struct {
	box    image.Rectangle
	canvas image.Rectangle
	w      io.Writer
}

func NewEncoder(w io.Writer, box image.Rectangle, canvas image.Rectangle) *Encoder {
	return &Encoder{
		box:    box,
		canvas: canvas,
		w:      w,
	}
}

func (e *Encoder) Encode(m *wardleyToGo.Map) error {
	fmt.Fprint(e.w, template)
	n := m.Nodes()
	for n.Next() {
		current := n.Node().(wardleyToGo.Component)
		pos := current.GetPosition()
		x := pos.X * 11
		y := 800 - pos.Y*8
		if component, ok := current.(*wardley.Anchor); ok {
			fmt.Fprintf(e.w, "\t\tnode_%v [label=\"%v\", shape=none, fontsize=14, pos=\"%v,%v\"];\n", current.ID(), component.Label, x, y)
			continue
		}
		if component, ok := current.(*wardley.Component); ok {
			fmt.Fprintf(e.w, "\t\tnode_%v [xlabel=\"%v\", fontsize=9, pos=\"%v,%v\"];\n", current.ID(), component.Label, x, y)
			continue
		}
		if component, ok := current.(*wardley.EvolvedComponent); ok {
			fmt.Fprintf(e.w, "\t\tnode_%v [xlabel=\"%v\", fontsize=9, color=red, pos=\"%v,%v\"];\n", current.ID(), component.Label, x, y)
			continue
		}
	}
	edges := m.Edges()
	for edges.Next() {
		ed := edges.Edge()
		if collaboration, ok := ed.(*wardley.Collaboration); ok {
			switch collaboration.GetType() {
			case wardley.RegularEdge:
				fmt.Fprintf(e.w, "\t\tnode_%v -> node_%v;\n", ed.From().ID(), ed.To().ID())
				continue
			case wardley.EvolvedComponentEdge:
				fmt.Fprintf(e.w, "\t\tnode_%v -> node_%v [style=dotted, color=red, arrowhead=normal];\n", ed.From().ID(), ed.To().ID())
				continue
			case wardley.EvolvedEdge:
				continue
			}
		}
		fmt.Fprintf(e.w, "\t\tnode_%v -> node_%v;\n", ed.From().ID(), ed.To().ID())
	}

	fmt.Fprint(e.w, "\t}\n}\n")
	return nil

}

const template = `
// Generate with neato -n -Tsvg test.gv
digraph {
        graph [
                rankdir = "LR"
                bgcolor = "white:lightgrey"
                style="filled"
                gradientangle = 180];
        node [shape=point, label="\n", margin="0.001"];
        edge [arrowhead=none]
        x0 [label="", fontsize=7, pos="0,0"];
        x1 [label="", fontsize=7, pos="1100,0"];
        y2 [label="", fontsize=7, pos="0,800"];
        vc [shape=none, label="V\na\nl\nu\ne\n \nC\nh\na\ni\nn", pos="-10,400"]
        genesis [shape=none, label="genesis", pos="0,-10"]
        custom [shape=none, label="custom", pos="187,-10"]
        customX [shape=none, label="", pos="187,800"]
        custom -> customX [style=dotted, arrowhead=none, color=grey]
        product [shape=none, label="product", pos="440,-10"]
        productX [shape=none, label="", pos="440,800"]
        product -> productX [style=dotted, arrowhead=none, color=grey]
        commodity [shape=none, label="commodity", pos="770,-10"]
        commodityX [shape=none, label="", pos="770,800"]
        commodity -> commodityX [style=dotted, arrowhead=none, color=grey]
        x0 -> x1 [arrowhead=normal]
        x0 -> y2 [arrowhead=normal]

        subgraph cluster_0 {
`

package parser

import (
	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/internal/wardley"
	"gonum.org/v1/gonum/graph"
)

type edge struct {
	toLabel   string
	fromLabel string
	T         graph.Node
	F         graph.Node
	edgeLabel string
}

func (e edge) From() graph.Node {
	return e.F
}

func (e edge) ReversedEdge() graph.Edge {
	return edge{
		F:         e.T,
		T:         e.F,
		toLabel:   e.fromLabel,
		fromLabel: e.toLabel,
		edgeLabel: e.edgeLabel,
	}
}

func (e edge) To() graph.Node {
	return e.T
}

func (e edge) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	fromCoord := e.F.(wardley.Element).GetCoordinates()
	toCoord := e.T.(wardley.Element).GetCoordinates()
	s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
		(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
		toCoord[1]*(width-padLeft)/100+padLeft,
		(height-padLeft)-toCoord[0]*(height-padLeft)/100,
		`stroke="grey"`, `stroke-width="1"`)
}

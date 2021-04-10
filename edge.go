package wardleyToGo

import (
	"github.com/owulveryck/wardleyToGo/components"
	"gonum.org/v1/gonum/graph"
)

type Edge struct {
	T         graph.Node
	F         graph.Node
	EdgeLabel string
	EdgeType  components.EdgeType
}

func (e Edge) From() graph.Node {
	return e.F
}

func (e Edge) ReversedEdge() graph.Edge {
	return Edge{
		F:         e.T,
		T:         e.F,
		EdgeLabel: e.EdgeLabel,
		EdgeType:  e.EdgeType,
	}
}

func (e Edge) To() graph.Node {
	return e.T
}

/*
func (e Edge) SVGTT(s *svg.SVG, width, height, padLeft, padBottom int) {
	fromCoord := e.F.(Element).GetCoordinates()
	toCoord := e.T.(Element).GetCoordinates()
	switch e.EdgeType {
	case CollaborationEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="grey"`, `stroke-width="1"`)
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke-dasharray="15 10"`, `stroke="rgb(250,216,120)"`, `stroke-width="10"`)
	case FacilitatingEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="grey"`, `stroke-width="1"`)
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke-dasharray="8 8"`, `stroke="rgb(200,159,182)"`, `stroke-width="10"`)
	case XAsAServiceEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="grey"`, `stroke-width="1"`)
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke-dasharray="2 14"`, `stroke="grey"`, `stroke-width="20"`)
	}
}

func (e Edge) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	fromCoord := e.F.(Element).GetCoordinates()
	toCoord := e.T.(Element).GetCoordinates()
	switch e.EdgeType {
	case RegularEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="grey"`, `stroke-width="1"`)
	case EvolvedComponentEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke-dasharray="5 5"`, `stroke="red"`, `stroke-width="1"`, `marker-end="url(#arrow)"`)
	case EvolvedEdge:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="red"`, `stroke-width="1"`)
	default:
		s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
			toCoord[1]*(width-padLeft)/100+padLeft,
			(height-padLeft)-toCoord[0]*(height-padLeft)/100,
			`stroke="grey"`, `stroke-width="1"`)
	}
}

*/

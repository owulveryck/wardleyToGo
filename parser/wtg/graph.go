package main

import (
	"image"
	"strconv"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

type node struct {
	id         int64
	label      string
	visibility int
	point      image.Point
}

// GetPosition of the element wrt a 100x100 map
func (n *node) GetPosition() image.Point {
	return n.point
}

func (n *node) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{
			Key:   "label",
			Value: n.label,
		},
	}

}

func (n *node) ID() int64 {
	return n.id
}

type edge struct {
	from       *wardley.Component
	to         *wardley.Component
	visibility int
}

func (e *edge) GetType() wardleyToGo.EdgeType {
	return 0
}

// From returns the from node of the edge.
func (e *edge) From() graph.Node {
	return e.from
}

// To returns the to node of the edge.
func (e *edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (e *edge) ReversedEdge() graph.Edge {
	return &edge{
		from:       e.to,
		to:         e.from,
		visibility: e.visibility,
	}
}

func (e *edge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{
			Key:   "arrowhead",
			Value: "none",
		},
		{
			Key:   "minlen",
			Value: strconv.Itoa(e.visibility),
		},
	}
}

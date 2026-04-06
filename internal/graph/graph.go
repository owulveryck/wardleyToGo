// Package graph provides a minimal directed graph implementation.
//
// It replaces the subset of gonum.org/v1/gonum/graph that is used by
// the wardleyToGo project: node/edge interfaces, slice-backed iterators,
// and a concrete DirectedGraph.
package graph

// Node is a graph node identified by a unique int64 ID.
type Node interface {
	ID() int64
}

// Edge is a directed edge between two nodes.
type Edge interface {
	From() Node
	To() Node
}

// Nodes is a slice-backed iterator over graph nodes.
type Nodes struct {
	nodes []Node
	pos   int
}

// NewNodes creates a Nodes iterator from a slice.
func NewNodes(nodes []Node) *Nodes {
	return &Nodes{nodes: nodes, pos: -1}
}

// Next advances the iterator and reports whether there is another node.
func (n *Nodes) Next() bool {
	n.pos++
	return n.pos < len(n.nodes)
}

// Node returns the current node. Must be called after Next returns true.
func (n *Nodes) Node() Node {
	return n.nodes[n.pos]
}

// Len returns the total number of nodes.
func (n *Nodes) Len() int {
	return len(n.nodes)
}

// Reset moves the iterator back to before the first node.
func (n *Nodes) Reset() {
	n.pos = -1
}

// Edges is a slice-backed iterator over graph edges.
type Edges struct {
	edges []Edge
	pos   int
}

// NewEdges creates an Edges iterator from a slice.
func NewEdges(edges []Edge) *Edges {
	return &Edges{edges: edges, pos: -1}
}

// Next advances the iterator and reports whether there is another edge.
func (e *Edges) Next() bool {
	e.pos++
	return e.pos < len(e.edges)
}

// Edge returns the current edge. Must be called after Next returns true.
func (e *Edges) Edge() Edge {
	return e.edges[e.pos]
}

// Len returns the total number of edges.
func (e *Edges) Len() int {
	return len(e.edges)
}

// Reset moves the iterator back to before the first edge.
func (e *Edges) Reset() {
	e.pos = -1
}

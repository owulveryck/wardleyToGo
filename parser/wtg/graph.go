package wtg

import (
	"sort"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type scratchMap struct {
	backend *simple.DirectedGraph
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (m *scratchMap) Node(id int64) graph.Node {
	return m.backend.Node(id)
}

// Nodes returns all the nodes in the graph.
//
// Nodes must not return nil.
func (m *scratchMap) Nodes() graph.Nodes {
	return m.backend.Nodes()
}

type mynodes struct {
	cursor int
	nodes  []*node
}

// Next advances the iterator and returns whether
// the next call to the item method will return a
// non-nil item.
//
// Next should be called prior to any call to the
// iterator's item retrieval method after the
// iterator has been obtained or reset.
//
// The order of iteration is implementation
// dependent.
func (m *mynodes) Next() bool {
	if m.cursor < len(m.nodes)-1 {
		m.cursor++
		return true
	}
	return false
}

// Len returns the number of items remaining in the
// iterator.
//
// If the number of items in the iterator is unknown,
// too large to materialize or too costly to calculate
// then Len may return a negative value.
// In this case the consuming function must be able
// to operate on the items of the iterator directly
// without materializing the items into a slice.
// The magnitude of a negative length has
// implementation-dependent semantics.
func (m *mynodes) Len() int {
	return len(m.nodes)
}

// Reset returns the iterator to its start position.
func (m *mynodes) Reset() {
	m.cursor = 0
}

func (m *mynodes) Node() graph.Node {
	return m.nodes[m.cursor]
}

// From returns all nodes that can be reached directly
// from the node with the given ID.
//
// From must not return nil.
func (m *scratchMap) From(id int64) graph.Nodes {
	nodes := m.backend.From(id)
	myn := &mynodes{
		nodes:  make([]*node, nodes.Len()),
		cursor: -1,
	}
	for i := 0; nodes.Next(); i++ {
		myn.nodes[i] = nodes.Node().(*node)
	}
	// TODO order the nodes by reverse visibility
	sort.Sort(nodeSorter(myn.nodes))
	return myn
}

// HasEdgeBetween returns whether an edge exists between
// nodes with IDs xid and yid without considering direction.
func (m *scratchMap) HasEdgeBetween(xid int64, yid int64) bool {
	return m.backend.HasEdgeBetween(xid, yid)
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise. The node v
// must be directly reachable from u as defined by the
// From method.
func (m *scratchMap) Edge(uid int64, vid int64) graph.Edge {
	return m.backend.Edge(uid, vid)
}

// HasEdgeFromTo returns whether an edge exists
// in the graph from u to v with IDs uid and vid.
func (m *scratchMap) HasEdgeFromTo(uid int64, vid int64) bool {
	return m.backend.HasEdgeFromTo(uid, vid)
}

// To returns all nodes that can reach directly
// to the node with the given ID.
//
// To must not return nil.
func (m *scratchMap) To(id int64) graph.Nodes {
	return m.backend.To(id)

}

func (m *scratchMap) AddNode(n graph.Node) {
	m.backend.AddNode(n)
}

func (m *scratchMap) SetEdge(e graph.Edge) {
	m.backend.SetEdge(e)
}

type node struct {
	visibility    int
	evolutionStep int
	c             *wardley.Component
}

func (node *node) ID() int64 {
	return node.c.ID()
}

type edge struct {
	f          *node
	t          *node
	visibility int
}

// From returns the from node of the edge.
func (e *edge) From() graph.Node {
	return e.f
}

// To returns the to node of the edge.
func (e *edge) To() graph.Node {
	return e.t
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (e *edge) ReversedEdge() graph.Edge {
	return &edge{
		f: e.t,
		t: e.f,
	}
}

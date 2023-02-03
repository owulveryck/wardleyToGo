package wardley

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/set/uid"
)

var (
	_ wardleyToGo.Backend = &Carte{}
)

type Carte struct {
	nodes map[int64]Element
	from  map[int64]map[int64]Collaboration
	to    map[int64]map[int64]Collaboration

	nodeIDs *uid.Set
}

func NewCarte() *Carte {
	return &Carte{
		nodes: make(map[int64]Element),
		from:  make(map[int64]map[int64]Collaboration),
		to:    make(map[int64]map[int64]Collaboration),

		nodeIDs: uid.NewSet(),
	}
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (c *Carte) Node(id int64) graph.Node {
	return c.nodes[id]
}

// Nodes returns all the nodes in the graph.
//
// Nodes must not return nil.
func (carte *Carte) Nodes() graph.Nodes {
	panic("not implemented") // TODO: Implement
}

// From returns all nodes that can be reached directly
// from the node with the given ID.
//
// From must not return nil.
func (carte *Carte) From(id int64) graph.Nodes {
	panic("not implemented") // TODO: Implement
}

// HasEdgeBetween returns whether an edge exists between
// nodes with IDs xid and yid without considering direction.
func (carte *Carte) HasEdgeBetween(xid int64, yid int64) bool {
	panic("not implemented") // TODO: Implement
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise. The node v
// must be directly reachable from u as defined by the
// From method.
func (carte *Carte) Edge(uid int64, vid int64) graph.Edge {
	panic("not implemented") // TODO: Implement
}

// HasEdgeFromTo returns whether an edge exists
// in the graph from u to v with IDs uid and vid.
func (carte *Carte) HasEdgeFromTo(uid int64, vid int64) bool {
	panic("not implemented") // TODO: Implement
}

// To returns all nodes that can reach directly
// to the node with the given ID.
//
// To must not return nil.
func (carte *Carte) To(id int64) graph.Nodes {
	panic("not implemented") // TODO: Implement
}
func (carte *Carte) Edges() graph.Edges {
	panic("not implemented") // TODO: Implement
}

// NewNode returns a new Node with a unique
// arbitrary ID.
func (carte *Carte) NewNode() graph.Node {
	panic("not implemented") // TODO: Implement
}

// AddNode adds a node to the graph. AddNode panics if
// the added node ID matches an existing node ID.
func (c *Carte) AddNode(n graph.Node) {
	if _, exists := c.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	c.nodes[n.ID()] = n
	c.nodeIDs.Use(n.ID())
	panic("not implemented") // TODO: Implement
}

// NewEdge returns a new Edge from the source to the destination node.
func (carte *Carte) NewEdge(from graph.Node, to graph.Node) graph.Edge {
	panic("not implemented") // TODO: Implement
}

// SetEdge adds an edge from one node to another.
// If the graph supports node addition the nodes
// will be added if they do not exist, otherwise
// SetEdge will panic.
// The behavior of an EdgeAdder when the IDs
// returned by e.From() and e.To() are equal is
// implementation-dependent.
// Whether e, e.From() and e.To() are stored
// within the graph is implementation dependent.
func (carte *Carte) SetEdge(e graph.Edge) {
	panic("not implemented") // TODO: Implement
}

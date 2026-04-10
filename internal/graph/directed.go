package graph

import "sort"

// simpleNode is a minimal Node used by DirectedGraph.NewNode.
type simpleNode struct {
	id int64
}

func (n simpleNode) ID() int64 { return n.id }

// DirectedGraph is a directed graph backed by adjacency maps.
type DirectedGraph struct {
	nodes  map[int64]Node
	from   map[int64]map[int64]Edge // from[uid][vid] = edge
	to     map[int64]map[int64]Edge // to[vid][uid] = edge
	nextID int64

	// cachedNodes holds the sorted node list. It is invalidated (set to nil)
	// whenever nodes are added or removed.
	cachedNodes []Node
}

// NewDirectedGraph returns an empty directed graph.
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		nodes: make(map[int64]Node),
		from:  make(map[int64]map[int64]Edge),
		to:    make(map[int64]map[int64]Edge),
	}
}

// NewNode returns a new node with a unique ID.
func (g *DirectedGraph) NewNode() Node {
	// Find an unused ID.
	for {
		if _, exists := g.nodes[g.nextID]; !exists {
			n := simpleNode{id: g.nextID}
			g.nextID++
			return n
		}
		g.nextID++
	}
}

// AddNode adds n to the graph. It panics if a node with the same ID
// already exists.
func (g *DirectedGraph) AddNode(n Node) {
	id := n.ID()
	if _, exists := g.nodes[id]; exists {
		panic("graph: node already exists")
	}
	g.nodes[id] = n
	g.cachedNodes = nil // invalidate cache
	if g.from[id] == nil {
		g.from[id] = make(map[int64]Edge)
	}
	if g.to[id] == nil {
		g.to[id] = make(map[int64]Edge)
	}
	// Keep nextID ahead of any manually-added node.
	if id >= g.nextID {
		g.nextID = id + 1
	}
}

// Node returns the node with the given ID, or nil.
func (g *DirectedGraph) Node(id int64) Node {
	return g.nodes[id]
}

// Nodes returns an iterator over all nodes, sorted by ID for determinism.
// The result is cached and reused until the graph is mutated.
func (g *DirectedGraph) Nodes() *Nodes {
	if g.cachedNodes == nil {
		ids := make([]int64, 0, len(g.nodes))
		for id := range g.nodes {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		nodes := make([]Node, len(ids))
		for i, id := range ids {
			nodes[i] = g.nodes[id]
		}
		g.cachedNodes = nodes
	}
	// Return a fresh iterator over the shared slice (iterator state is per-call).
	return NewNodes(g.cachedNodes)
}

// SetEdge adds or replaces an edge. Both endpoints must already exist.
func (g *DirectedGraph) SetEdge(e Edge) {
	from := e.From().ID()
	to := e.To().ID()
	if from == to {
		panic("graph: self-loop")
	}
	// Auto-init adjacency maps if nodes were added without them.
	if g.from[from] == nil {
		g.from[from] = make(map[int64]Edge)
	}
	if g.to[to] == nil {
		g.to[to] = make(map[int64]Edge)
	}
	g.from[from][to] = e
	g.to[to][from] = e
}

// RemoveEdge removes the edge from uid to vid, if it exists.
func (g *DirectedGraph) RemoveEdge(uid, vid int64) {
	delete(g.from[uid], vid)
	delete(g.to[vid], uid)
}

// HasEdgeFromTo reports whether an edge exists from uid to vid.
func (g *DirectedGraph) HasEdgeFromTo(uid, vid int64) bool {
	if m, ok := g.from[uid]; ok {
		_, exists := m[vid]
		return exists
	}
	return false
}

// Edge returns the edge from uid to vid, or nil.
func (g *DirectedGraph) Edge(uid, vid int64) Edge {
	if m, ok := g.from[uid]; ok {
		return m[vid]
	}
	return nil
}

// Edges returns an iterator over all edges, sorted for determinism.
func (g *DirectedGraph) Edges() *Edges {
	var edges []Edge
	// Collect all from IDs sorted.
	fromIDs := make([]int64, 0, len(g.from))
	for id := range g.from {
		fromIDs = append(fromIDs, id)
	}
	sort.Slice(fromIDs, func(i, j int) bool { return fromIDs[i] < fromIDs[j] })
	for _, fid := range fromIDs {
		toIDs := make([]int64, 0, len(g.from[fid]))
		for tid := range g.from[fid] {
			toIDs = append(toIDs, tid)
		}
		sort.Slice(toIDs, func(i, j int) bool { return toIDs[i] < toIDs[j] })
		for _, tid := range toIDs {
			edges = append(edges, g.from[fid][tid])
		}
	}
	return NewEdges(edges)
}

// From returns an iterator over nodes reachable from the node with the
// given ID via outgoing edges.
func (g *DirectedGraph) From(id int64) *Nodes {
	m := g.from[id]
	ids := make([]int64, 0, len(m))
	for tid := range m {
		ids = append(ids, tid)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	nodes := make([]Node, len(ids))
	for i, tid := range ids {
		nodes[i] = g.nodes[tid]
	}
	return NewNodes(nodes)
}

// To returns an iterator over nodes that have an edge to the node with
// the given ID.
func (g *DirectedGraph) To(id int64) *Nodes {
	m := g.to[id]
	ids := make([]int64, 0, len(m))
	for uid := range m {
		ids = append(ids, uid)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	nodes := make([]Node, len(ids))
	for i, uid := range ids {
		nodes[i] = g.nodes[uid]
	}
	return NewNodes(nodes)
}

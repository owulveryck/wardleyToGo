package graph

// DepthFirst performs a depth-first traversal of a graph.
type DepthFirst struct {
	// Visit is called once for each visited node.
	Visit func(Node)

	visited map[int64]bool
}

// Walker is the subset of a graph needed for DFS traversal.
type Walker interface {
	From(id int64) *Nodes
}

// Walk performs a depth-first walk starting at root. If until is non-nil
// and returns true for a node, Walk returns that node immediately.
// Otherwise it returns nil after visiting all reachable nodes.
func (d *DepthFirst) Walk(g Walker, root Node, until func(Node) bool) Node {
	if d.visited == nil {
		d.visited = make(map[int64]bool)
	}
	return d.walk(g, root, until)
}

func (d *DepthFirst) walk(g Walker, n Node, until func(Node) bool) Node {
	id := n.ID()
	if d.visited[id] {
		return nil
	}
	d.visited[id] = true

	if d.Visit != nil {
		d.Visit(n)
	}
	if until != nil && until(n) {
		return n
	}

	neighbors := g.From(id)
	for neighbors.Next() {
		if result := d.walk(g, neighbors.Node(), until); result != nil {
			return result
		}
	}
	return nil
}

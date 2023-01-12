package wtg

import "gonum.org/v1/gonum/graph"

func findRoot(g graph.Directed) []graph.Node {
	ret := make([]graph.Node, 0)
	nodes := g.Nodes()
	for nodes.Next() {
		if g.To(nodes.Node().ID()).Len() == 0 {
			ret = append(ret, nodes.Node())
		}
	}
	return ret
}

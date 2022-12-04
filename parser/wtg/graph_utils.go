package main

import "gonum.org/v1/gonum/graph"

func findLeafs(g graph.Directed) []graph.Node {
	ret := make([]graph.Node, 0)
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if g.From(n.ID()).Len() == 0 {
			ret = append(ret, n)
		}
	}
	return ret
}
func findRoot(g graph.Directed) []graph.Node {
	ret := make([]graph.Node, 0)
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if g.To(n.ID()).Len() == 0 {
			ret = append(ret, n)
		}
	}
	return ret
}

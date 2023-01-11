package main

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/simple"
)

func setCoords(m wardleyToGo.Map, withEvolution bool) {
	tempMap := &scratchMapchMap{backend: simple.NewDirectedGraph()}
	ns := m.Nodes()
	inventory := make(map[int64]*node)
	for ns.Next() {
		if c, ok := ns.Node().(*wardley.Component); ok {
			n := &node{
				c: c,
			}
			inventory[c.ID()] = n
			tempMap.AddNode(n)
		}
	}
	es := m.Edges()
	for es.Next() {
		tempMap.SetEdge(&edge{
			f:          inventory[es.Edge().From().ID()],
			t:          inventory[es.Edge().To().ID()],
			visibility: es.Edge().(*wardley.Collaboration).Visibility,
		})
	}
	maxVisibility := setNodesVisibility(tempMap)
	setY(tempMap, m, maxVisibility)
	if withEvolution {
		maxEvolution := setNodesEvolutionStep(tempMap)
		setX(tempMap, m, maxEvolution)
	}
	// TODO set the graph nodes
}

func setY(buf *scratchMapchMap, m wardleyToGo.Map, maxVisibility int) {
	vStep := 95 / maxVisibility
	_ = vStep

}
func setX(buf *scratchMapchMap, m wardleyToGo.Map, maxEvolution int) {
	hStep := 95 / maxEvolution
	_ = hStep

}

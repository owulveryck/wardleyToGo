package main

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
	"gonum.org/v1/gonum/graph/simple"
)

func setCoords(m wardleyToGo.Map, withEvolution bool) {
	tempMap := &scratchMapchMap{backend: simple.NewDirectedGraph()}
	ns := m.Nodes()
	inventory := make(map[int64]*node)
	for ns.Next() {
		if c, ok := ns.Node().(*wardley.EvolvedComponent); ok {
			n := &node{
				c: c.Component,
			}
			inventory[c.ID()] = n
			tempMap.AddNode(n)
		}
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
	vStep := 96 / maxVisibility
	allNodes := buf.Nodes()
	for allNodes.Next() {
		n := allNodes.Node().(*node)
		if c, ok := m.Node(n.ID()).(*wardley.Component); ok {
			c.Placement.Y = n.visibility*vStep + 2
		}
		if c, ok := m.Node(n.ID()).(*wardley.EvolvedComponent); ok {
			c.Placement.Y = n.visibility*vStep + 2
		}
	}

}
func setX(buf *scratchMapchMap, m wardleyToGo.Map, maxEvolution int) {
	hStep := 20 / maxEvolution
	allNodes := buf.Nodes()
	for allNodes.Next() {
		n := allNodes.Node().(*node)
		if nn, ok := m.Node(n.ID()).(*wardley.Component); ok {
			if !nn.Configured {
				nn.Placement.X = n.evolutionStep*hStep + 30
				nn.Color = wtg.Colors["Grey"]
			}
		}
	}

}

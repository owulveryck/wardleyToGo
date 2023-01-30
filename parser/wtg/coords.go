package wtg

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/simple"
)

// SetCoords sets the y anx x axis of the components
func SetCoords(m wardleyToGo.Map, withEvolution bool) {
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
	setEdgeAbsoluteVisibility(m)
	// TODO set the graph nodes
}

func setEdgeAbsoluteVisibility(m wardleyToGo.Map) {
	allNodes := m.Nodes()
	for allNodes.Next() {
		currNode := allNodes.Node()
		t := m.To(currNode.ID())
		for t.Next() {
			if e := m.Edge(t.Node().ID(), currNode.ID()); e != nil {
				if e, ok := e.(*wardley.Collaboration); ok {
					e.AbsoluteVisibility = currNode.(wardleyToGo.Chainer).GetAbsoluteVisibility()
				}
			}
		}
	}
}

func setY(buf *scratchMapchMap, m wardleyToGo.Map, maxVisibility int) {
	vStep := 50
	if maxVisibility != 0 {
		vStep = 94 / maxVisibility
	}
	allNodes := buf.Nodes()
	for allNodes.Next() {
		n := allNodes.Node().(*node)
		if c, ok := m.Node(n.ID()).(*wardley.Component); ok {
			c.Placement.Y = n.visibility*vStep + 2
			c.AbsoluteVisibility = n.visibility
		}
		if c, ok := m.Node(n.ID()).(*wardley.EvolvedComponent); ok {
			c.Placement.Y = n.visibility*vStep + 2
			c.AbsoluteVisibility = n.visibility
		}
	}

}
func setX(buf *scratchMapchMap, m wardleyToGo.Map, maxEvolution int) {
	hStep := 50
	if maxEvolution != 0 {
		hStep = 90 / maxEvolution
	}
	allNodes := buf.Nodes()
	for allNodes.Next() {
		n := allNodes.Node().(*node)
		if nn, ok := m.Node(n.ID()).(*wardley.Component); ok {
			if !nn.Configured {
				nn.Placement.X = n.evolutionStep*hStep + 10
				nn.Color = Colors["Grey"]
			}
		}
	}

}

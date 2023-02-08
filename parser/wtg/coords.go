package wtg

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/simple"
)

// SetLabelCoords sets the label coordinate with simple rules:
// set anchor to left unless all the the From and To nodes are more advanced on evolution axis
func SetLabelAnchor(m wardleyToGo.Map) {
	ns := m.Nodes()
	for ns.Next() {
		if n, ok := ns.Node().(*wardley.EvolvedComponent); ok {
			n.Anchor = wardley.AdjustStart
		}
		if n, ok := ns.Node().(*wardley.Component); ok {
			if n.Anchor != wardley.AdjustUndefined {
				continue
			}
			n.Anchor = wardley.AdjustEnd
			n.LabelPlacement.X = -n.LabelPlacement.X
			if n.Type == wardley.PipelineComponent {
				n.Anchor = wardley.AdjustMiddle
				continue
			}
			start := false
			it := m.From(n.ID())
			for it.Next() {
				if f, ok := it.Node().(wardleyToGo.Component); ok {
					if f.GetPosition().X < n.GetPosition().X {
						start = true
						continue
					}
				}
			}
			it = m.To(n.ID())
			if it.Len() == 0 {
				// this is a top level node
				n.Anchor = wardley.AdjustMiddle
				n.LabelPlacement.Y = -10
				n.LabelPlacement.X = 0
				continue

			}
			for it.Next() {
				if f, ok := it.Node().(wardleyToGo.Component); ok {
					if f.GetPosition().X < n.GetPosition().X {
						start = true
						continue
					}
				}
			}
			if start {
				n.LabelPlacement.X = -n.LabelPlacement.X
				n.Anchor = wardley.AdjustStart

			}
		}
	}
}

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

// setY sets the placement of the node according to the visibility of the component carried by the scratchmap
func setY(buf *scratchMapchMap, m wardleyToGo.Map, maxVisibility int) {
	vStep := 50
	if maxVisibility != 0 {
		vStep = 94 / maxVisibility
	}
	allNodes := buf.Nodes()
	for allNodes.Next() {
		n := allNodes.Node().(*node)
		if c, ok := m.Node(n.ID()).(*wardley.Component); ok {
			c.Placement.Y = n.visibility*vStep + 3
			c.AbsoluteVisibility = n.visibility
		}
		if c, ok := m.Node(n.ID()).(*wardley.EvolvedComponent); ok {
			c.Placement.Y = n.visibility*vStep + 3
			c.AbsoluteVisibility = n.visibility
		}
	}

}
func setX(buf *scratchMapchMap, m wardleyToGo.Map, maxEvolution int) {
	hStep := 50
	if maxEvolution != 0 {
		hStep = 80 / maxEvolution
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

func setLabelPlacement(n *wardley.Component, placement string) error {
	switch placement {
	case "N":
		n.Anchor = wardley.AdjustMiddle
		n.LabelPlacement.X = 0
		n.LabelPlacement.Y = -15
	case "NE":
		n.Anchor = wardley.AdjustStart
		n.LabelPlacement.Y = -15
		n.LabelPlacement.X = 11
	case "NW":
		n.Anchor = wardley.AdjustEnd
		n.LabelPlacement.Y = -15
		n.LabelPlacement.X = -11
	case "W":
		n.Anchor = wardley.AdjustEnd
		n.LabelPlacement.Y = 0
		n.LabelPlacement.X = -11
	case "S":
		n.Anchor = wardley.AdjustMiddle
		n.LabelPlacement.X = 0
		n.LabelPlacement.Y = 19
	case "SW":
		n.Anchor = wardley.AdjustEnd
		n.LabelPlacement.Y = 19
		n.LabelPlacement.X = -11
	case "SE":
		n.Anchor = wardley.AdjustStart
		n.LabelPlacement.Y = 19
		n.LabelPlacement.X = 11
	case "E":
		n.Anchor = wardley.AdjustStart
		n.LabelPlacement.Y = 0
		n.LabelPlacement.X = 11
	default:
		return fmt.Errorf("unknown placement: %v", placement)
	}
	return nil
}

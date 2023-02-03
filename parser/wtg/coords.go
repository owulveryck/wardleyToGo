package wtg

import (
	"fmt"
	"image"

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
		curr := ns.Node().(wardley.Element)
		n := &node{
			c: curr,
		}
		inventory[curr.ID()] = n
		tempMap.AddNode(n)
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
		c := m.Node(n.ID()).(wardley.Element)
		p := c.GetPosition()
		p.Y = n.visibility*vStep + 3
		c.SetPosition(p)
		c.SetAbsoluteVisibility(n.visibility)
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
		nn := m.Node(n.ID()).(wardley.Element)
		p := nn.GetPosition()
		p.X = n.evolutionStep*hStep + 10
		nn.SetPosition(p)
		nn.SetColor(Colors["Grey"])
	}
}

func setLabelPlacement(n wardley.Labeler, placement string) error {
	switch placement {
	case "N":
		n.SetLabelAnchor(wardley.AdjustMiddle)
		n.SetLabelPlacement(image.Point{0, -15})
	case "NE":
		n.SetLabelAnchor(wardley.AdjustStart)
		n.SetLabelPlacement(image.Point{11, -15})
	case "NW":
		n.SetLabelAnchor(wardley.AdjustEnd)
		n.SetLabelPlacement(image.Point{-11, -15})
	case "W":
		n.SetLabelAnchor(wardley.AdjustEnd)
		n.SetLabelPlacement(image.Point{-11, 0})
	case "S":
		n.SetLabelAnchor(wardley.AdjustMiddle)
		n.SetLabelPlacement(image.Point{0, 19})
	case "SW":
		n.SetLabelAnchor(wardley.AdjustEnd)
		n.SetLabelPlacement(image.Point{-11, 19})
	case "SE":
		n.SetLabelAnchor(wardley.AdjustStart)
		n.SetLabelPlacement(image.Point{11, 19})
	case "E":
		n.SetLabelAnchor(wardley.AdjustStart)
		n.SetLabelPlacement(image.Point{11, 0})
	default:
		return fmt.Errorf("unknown placement: %v", placement)
	}
	return nil
}

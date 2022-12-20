package wtg

import (
	"math"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/path"
)

func setYCoords(m *wardleyToGo.Map) error {

	allShortestPaths := path.DijkstraAllPaths(m)
	roots := findRoot(m)
	leafs := findLeafs(m)
	var maxDepth int
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := m.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*wardley.Collaboration).Visibility
				}
				if currentVisibility > maxDepth {
					maxDepth = currentVisibility
				}
			}
		}
	}

	step := 100 / maxDepth
	cs := &coordSetter{
		verticalStep: step,
	}
	nroots := len(roots)
	hsteps := 100 / (nroots + 1)
	for i, n := range roots {
		if n.Placement.X == 0 {
			n.Placement.X = hsteps * (i + 1)
		}
		cs.walk(m, n, 0)
	}

	return nil
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(m *wardleyToGo.Map, n *wardley.Component, visibility int) {
	n.Placement.Y = visibility * c.verticalStep
	from := m.From(n.ID())
	hsteps := 100 / (from.Len() + 1)
	i := 1
	for from.Next() {
		switch from.Node().(type) {
		case *wardley.Component:
			if m.Edge(n.ID(), from.Node().ID()) != nil {
				c.walk(m, from.Node().(*wardley.Component), m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
			}
		case *wardley.EvolvedComponent:
			if m.Edge(n.ID(), from.Node().ID()) != nil {
				c.walk(m, from.Node().(*wardley.EvolvedComponent).Component, m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
			}
		}
		switch n := from.Node().(type) {
		case *wardley.Component:
			if n.Placement.X == 0 {
				n.Placement.X = hsteps * i
			}
			if m.Edge(n.ID(), from.Node().ID()) != nil {
				c.walk(m, n, m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
			}
		case *wardley.EvolvedComponent:
			if n.Placement.X == 0 {
				n.Placement.X = hsteps * i
			}
			if m.Edge(n.ID(), from.Node().ID()) != nil {
				c.walk(m, n.Component, m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
			}
		}
		i++
	}
}

func computePlacement(s string) (int, int, error) {
	currentStage := -1
	currentCursor := 0
	stages := make([]int, 5)
	evolvedCursor := 0
	evolvedStage := 0
	cursor := 0
	stage := 0
	for _, c := range s {
		switch c {
		case '|':
			currentCursor = 0
			currentStage++
			continue
		case 'x':
			cursor = currentCursor
			stage = currentStage
		case '>':
			evolvedCursor = currentCursor
			evolvedStage = currentStage
		default:
			currentCursor++
			stages[currentStage]++
		}
	}
	stagePositions := []float64{0, 17.4, 40, 70, 100}
	position := stagePositions[stage] + (stagePositions[stage+1]-stagePositions[stage])*float64(cursor)/float64(stages[stage])
	evolvedPosition := stagePositions[evolvedStage] + (stagePositions[evolvedStage+1]-stagePositions[evolvedStage])*float64(evolvedCursor)/float64(stages[evolvedStage])
	return int(math.Round(position)), int(math.Round(evolvedPosition)), nil
}

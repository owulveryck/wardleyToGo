package wtg

import (
	"fmt"
	"math"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/path"
)

func computeEvolutionPosition(s string) (int, int, error) {
	// stages is an array containing the size of each stage
	// for example if the string is |...|....|.....|......|
	// then stages is [3,4,5,6]
	stages := make([]int, 4)
	// cursor is the place of the cursor in its stage (ex: |..x..|, x=3; |x|, x=0)
	cursor := 0
	// stage is the actual stage where the cursor is (x)
	stage := 0
	evolvedCursor := 0
	evolvedStage := 0
	// Position of each stage
	stagePositions := []float64{0, 17.4, 40, 70, 100}
	iteratorStage := -1
	iteratorCursor := 0
	x := 0
	sup := 0
	for _, c := range s {
		switch c {
		case '|':
			iteratorCursor = 0
			iteratorStage++
			continue
		case 'x':
			x++
			cursor = iteratorCursor
			stage = iteratorStage
			stages[iteratorStage]++
		case '>':
			sup++
			evolvedCursor = iteratorCursor
			evolvedStage = iteratorStage
			stages[iteratorStage]++
		default:
			iteratorCursor++
			if iteratorStage < 0 {
				return 0, 0, fmt.Errorf("expected | as a first element")
			}
			if iteratorStage >= len(stages) {
				return 0, 0, fmt.Errorf("too many |")
			}
			stages[iteratorStage]++
		}
	}
	if iteratorStage != 4 {
		return 0, 0, fmt.Errorf("expected 5x|")
	}
	if x != 1 {
		return 0, 0, fmt.Errorf("expeted one and only one x")
	}
	if sup > 1 {
		return 0, 0, fmt.Errorf("expeted one or less >")
	}
	position := 50.0
	if stages[stage] != 0 {
		percentageInCurrentStage := float64(cursor+1) / float64(stages[stage]+1)
		//log.Printf("\n\t%v\n\tstages: %v\n\tcurrent stage: %v\n\tcursor: %v\n\tpercentage: %v", s, stages, stage, cursor, percentageInCurrentStage)
		position = stagePositions[stage] + (stagePositions[stage+1]-stagePositions[stage])*percentageInCurrentStage
	}
	evolvedPosition := 0.0
	if stages[evolvedStage] != 0 && sup != 0 {
		percentageInCurrentStage := float64(evolvedCursor+1) / float64(stages[evolvedStage]+1)
		evolvedPosition = stagePositions[evolvedStage] + (stagePositions[evolvedStage+1]-stagePositions[evolvedStage])*percentageInCurrentStage
	}
	if position >= evolvedPosition && evolvedPosition != 0 {
		return 0, 0, fmt.Errorf("cannot have an evolution before the cursor")
	}
	return int(math.Round(position)), int(math.Round(evolvedPosition)), nil
}

type nodeSorter []wardleyToGo.Component

// Len is part of sort.Interface.
func (s nodeSorter) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s nodeSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s nodeSorter) Less(i, j int) bool {
	var labelI string
	var labelJ string
	if current, ok := s[i].(*wardley.Component); ok {
		labelI = current.Label
	}
	if current, ok := s[i].(*wardley.EvolvedComponent); ok {
		labelI = current.Label
	}
	if current, ok := s[j].(*wardley.Component); ok {
		labelJ = current.Label
	}
	if current, ok := s[j].(*wardley.EvolvedComponent); ok {
		labelJ = current.Label
	}
	return labelI < labelJ
}

func (p *Parser) computeX() {
	var walk func(m *wardleyToGo.Map, n *wardley.Component)
	walk = func(m *wardleyToGo.Map, n *wardley.Component) {
		fromIT := m.From(n.ID())
		hsteps := 100 / (fromIT.Len() + 1)
		i := 1
		for fromIT.Next() {
			switch fromNode := fromIT.Node().(type) {
			case *wardley.Component:
				walk(m, fromNode)
				if fromNode.Placement.X == 0 {
					fromNode.Placement.X = hsteps * i
				}
			case *wardley.EvolvedComponent:
				walk(m, fromNode.Component)
				if fromNode.Placement.X == 0 {
					fromNode.Placement.X = hsteps * i
				}
			}
			i++
		}
	}
	roots := findRoot(p.WMap)
	nroots := len(roots)
	hsteps := 100 / (nroots + 1)

	for i, n := range roots {
		if n.Placement.X == 0 {
			n.Placement.X = hsteps * (i + 1)
		}
		walk(p.WMap, n)
	}
}

func (p *Parser) computeY() {
	allShortestPaths := path.DijkstraAllPaths(p.WMap)
	roots := findRoot(p.WMap)
	leafs := findLeafs(p.WMap)
	maxDepth := 1
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := p.WMap.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*wardley.Collaboration).Visibility
				}
				if currentVisibility > maxDepth {
					maxDepth = currentVisibility
				}
			}
		}
	}

	step := 96 / maxDepth
	cs := &coordSetter{
		verticalStep: step,
	}
	for _, n := range roots {
		cs.walk(p.WMap, n, 0)
	}
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(m *wardleyToGo.Map, n *wardley.Component, visibility int) {
	n.Placement.Y = 2 + visibility*c.verticalStep
	fromIT := m.From(n.ID())
	for fromIT.Next() {
		switch fromNode := fromIT.Node().(type) {
		case *wardley.Component:
			c.walk(m, fromNode, m.Edge(n.ID(), fromNode.ID()).(*wardley.Collaboration).Visibility+visibility)
		case *wardley.EvolvedComponent:
			c.walk(m, fromNode.Component, m.Edge(n.ID(), fromNode.ID()).(*wardley.Collaboration).Visibility+visibility)
		}
	}
}

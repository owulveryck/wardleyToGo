package wtg

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

type evolutionSetter struct {
	g           graph.Directed
	currentStep int
}

func (e *evolutionSetter) visit(srcNode graph.Node) {
	n := srcNode.(*node)
	n.evolutionStep = e.currentStep
	// if the node is a leaf (meaning the from is empty), move the cursor
	fs := e.g.From(n.ID())
	if fs.Len() == 0 {
		e.currentStep++
	}
}

// returns the max evolution
func setNodesEvolutionStep(g *scratchMapchMap) int {
	roots := findRoot(g)
	e := &evolutionSetter{
		g: g,
	}
	df := &traverse.DepthFirst{
		Visit: e.visit,
	}
	for _, root := range roots {
		df.Walk(g, root, nil)
	}
	return e.currentStep
}

func computeEvolutionPosition(s string) (int, int, int, error) {
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
	inertia := 0
	x := 0
	sup := 0
	for _, c := range s {
		switch c {
		case ']':
			inertia = int(stagePositions[iteratorStage+1])
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
				return 0, 0, 0, fmt.Errorf("expected | as a first element")
			}
			if iteratorStage >= len(stages) {
				return 0, 0, 0, fmt.Errorf("too many |")
			}
			stages[iteratorStage]++
		}
	}
	if iteratorStage != 4 {
		return 0, 0, 0, fmt.Errorf("expected 5x|")
	}
	if x != 1 {
		return 0, 0, 0, fmt.Errorf("expeted one and only one x")
	}
	if sup > 1 {
		return 0, 0, 0, fmt.Errorf("expeted one or less >")
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
		return 0, 0, 0, fmt.Errorf("cannot have an evolution before the cursor")
	}
	return int(math.Round(position)), int(math.Round(evolvedPosition)), inertia, nil
}

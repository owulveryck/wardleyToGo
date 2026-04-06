package wtg2

import "testing"

func TestComputeYPositions_LinearChain(t *testing.T) {
	doc := &Document{
		Nodes: []*NodeDecl{
			{Name: "A", Kind: KindAnchor},
			{Name: "B", Kind: KindComponent},
			{Name: "C", Kind: KindComponent},
		},
		Edges: []*EdgeDecl{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
		},
	}

	yPos := ComputeYPositions(doc)

	// A should be at top (rank 0)
	if yPos["A"] >= yPos["B"] {
		t.Errorf("A (Y=%d) should be above B (Y=%d)", yPos["A"], yPos["B"])
	}
	if yPos["B"] >= yPos["C"] {
		t.Errorf("B (Y=%d) should be above C (Y=%d)", yPos["B"], yPos["C"])
	}
}

func TestComputeYPositions_Diamond(t *testing.T) {
	doc := &Document{
		Nodes: []*NodeDecl{
			{Name: "A", Kind: KindAnchor},
			{Name: "B", Kind: KindComponent},
			{Name: "C", Kind: KindComponent},
			{Name: "D", Kind: KindComponent},
		},
		Edges: []*EdgeDecl{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "B", To: "D"},
			{From: "C", To: "D"},
		},
	}

	yPos := ComputeYPositions(doc)

	// B and C should be at the same rank
	if yPos["B"] != yPos["C"] {
		t.Errorf("B (Y=%d) and C (Y=%d) should be at the same level", yPos["B"], yPos["C"])
	}
	// D should be below B and C
	if yPos["D"] <= yPos["B"] {
		t.Errorf("D (Y=%d) should be below B (Y=%d)", yPos["D"], yPos["B"])
	}
}

func TestComputeYPositions_PipelineMembers(t *testing.T) {
	doc := &Document{
		Nodes: []*NodeDecl{
			{Name: "A", Kind: KindAnchor},
			{Name: "Engine", Kind: KindComponent},
		},
		Edges: []*EdgeDecl{
			{From: "A", To: "Engine"},
		},
		Pipelines: []*PipelineDecl{
			{
				Name: "Engine",
				Members: []*PipelineMemberDecl{
					{Name: "AlgoA", Position: "III.5"},
					{Name: "AlgoB", Position: "II.3"},
				},
			},
		},
	}

	yPos := ComputeYPositions(doc)

	// Pipeline members should share parent Y
	if yPos["AlgoA"] != yPos["Engine"] {
		t.Errorf("AlgoA (Y=%d) should have same Y as Engine (Y=%d)", yPos["AlgoA"], yPos["Engine"])
	}
	if yPos["AlgoB"] != yPos["Engine"] {
		t.Errorf("AlgoB (Y=%d) should have same Y as Engine (Y=%d)", yPos["AlgoB"], yPos["Engine"])
	}
}

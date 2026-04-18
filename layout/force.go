package layout

import (
	"math"
	"sort"
)

// forceSpread refines vertical positions using a force-directed simulation.
//
// Starting from the initial linear placement (rank → Y), the algorithm
// iterates [Options.ForceIterations] times, applying two classes of forces
// each pass:
//
//   - Repulsion: every pair of non-pipeline-member nodes that are closer
//     than [Options.MinSpacing] pushes each other apart. The repulsive
//     force magnitude is RepulsionStrength / dy², producing strong
//     short-range separation. When two nodes are nearly coincident
//     (|dy| < 0.1), dy is clamped to ±0.1 to avoid division by zero.
//
//   - Edge attraction: for each edge, the connected nodes are pulled
//     toward their ideal separation (rank difference × layer height).
//     The attraction coefficient is 0.1, weaker than repulsion so that
//     overlap avoidance takes priority.
//
// After each iteration the computed displacements are scaled by damping
// (which decays by 2 % per iteration for convergence), positions are
// clamped to [MinY, MaxY], and [enforceRankOrder] restores the
// topological ordering invariant.
//
// Pipeline members are excluded entirely — they inherit their parent's
// Y coordinate in the post-processing phase of [defaultLayouter.Layout].
// Anchor nodes (rank 0) are pinned: they contribute to repulsion but
// their own positions are never modified.
func forceSpread(
	positions map[string]float64,
	edges []Edge,
	pipelineMembers map[string]bool,
	pipelineParents map[string]bool,
	ranks map[string]int,
	maxRank int,
	opts Options,
) map[string]float64 {
	if len(positions) <= 1 {
		return positions
	}

	// Collect names of movable nodes (not pipeline members, not anchors)
	pinned := make(map[string]bool)
	names := make([]string, 0, len(positions))
	allNames := make([]string, 0, len(positions))
	for name := range positions {
		if pipelineMembers[name] {
			continue
		}
		allNames = append(allNames, name)
		if ranks[name] == 0 {
			pinned[name] = true
		} else {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	sort.Strings(allNames)

	if len(names) == 0 {
		return positions
	}

	minY := float64(opts.MinY)
	maxY := float64(opts.MaxY)
	minSpacing := float64(opts.MinSpacing)
	damping := opts.Damping

	// Compute ideal layer height for edge attraction
	layerHeight := (maxY - minY) / math.Max(float64(maxRank), 1)

	// Build index for O(1) lookup instead of map[string]float64 per iteration.
	nameIndex := make(map[string]int, len(names))
	for i, n := range names {
		nameIndex[n] = i
	}
	allNameIndex := make(map[string]int, len(allNames))
	for i, n := range allNames {
		allNameIndex[n] = i
	}
	displacements := make([]float64, len(names))

	// Pre-sort names by rank for enforceRankOrder (rank is stable across iterations).
	sortedByRank := make([]string, len(allNames))
	copy(sortedByRank, allNames)
	sort.Slice(sortedByRank, func(i, j int) bool {
		ri, rj := ranks[sortedByRank[i]], ranks[sortedByRank[j]]
		if ri != rj {
			return ri < rj
		}
		return sortedByRank[i] < sortedByRank[j]
	})

	for iter := 0; iter < opts.ForceIterations; iter++ {
		// Reset displacements without re-allocating.
		for i := range displacements {
			displacements[i] = 0
		}

		// Pairwise repulsion between all active nodes (including pinned for repulsion source).
		// Same-rank nodes get a weaker repulsion (30%) to separate them
		// without compressing the overall layout.
		// Pipeline parents use 3× MinSpacing because their visual footprint
		// (the pipeline rectangle) is much taller than a regular component circle.
		for i := 0; i < len(allNames); i++ {
			for j := i + 1; j < len(allNames); j++ {
				a, b := allNames[i], allNames[j]
				dy := positions[b] - positions[a]
				absDy := math.Abs(dy)

				effectiveSpacing := minSpacing
				if pipelineParents[a] || pipelineParents[b] {
					effectiveSpacing = minSpacing * 3
				}

				if absDy < effectiveSpacing {
					if absDy < 0.1 {
						absDy = 0.1
						if dy >= 0 {
							dy = 0.1
						} else {
							dy = -0.1
						}
					}
					strength := opts.RepulsionStrength
					if ranks[a] == ranks[b] && !pipelineParents[a] && !pipelineParents[b] {
						strength *= 0.3
					}
					force := strength / (absDy * absDy)
					sign := 1.0
					if dy < 0 {
						sign = -1.0
					}
					// Only apply displacement to non-pinned nodes
					if !pinned[a] {
						if idx, ok := nameIndex[a]; ok {
							displacements[idx] -= force * sign
						}
					}
					if !pinned[b] {
						if idx, ok := nameIndex[b]; ok {
							displacements[idx] += force * sign
						}
					}
				}
			}
		}

		// Edge attraction: connected nodes pulled toward ideal separation
		for _, e := range edges {
			from := resolveMember(e.From)
			to := resolveMember(e.To)
			if pipelineMembers[from] || pipelineMembers[to] {
				continue
			}
			if _, ok := positions[from]; !ok {
				continue
			}
			if _, ok := positions[to]; !ok {
				continue
			}

			idealDy := float64(ranks[to]-ranks[from]) * layerHeight
			actualDy := positions[to] - positions[from]
			delta := (actualDy - idealDy) * 0.1
			if !pinned[from] {
				if idx, ok := nameIndex[from]; ok {
					displacements[idx] += delta
				}
			}
			if !pinned[to] {
				if idx, ok := nameIndex[to]; ok {
					displacements[idx] -= delta
				}
			}
		}

		// Apply displacements with damping, capping per-iteration movement
		// to minSpacing. Without the cap, coincident same-rank nodes
		// generate enormous repulsion (strength/0.01 ≈ 240) that throws
		// them to the [minY,maxY] boundaries in a single step.
		// enforceRankOrder then uses the outlier's position as the floor
		// for subsequent ranks, cascading all deeper nodes to maxY.
		maxDisp := minSpacing
		for i, name := range names {
			d := displacements[i] * damping
			if d > maxDisp {
				d = maxDisp
			}
			if d < -maxDisp {
				d = -maxDisp
			}
			positions[name] += d
			if positions[name] < minY {
				positions[name] = minY
			}
			if positions[name] > maxY {
				positions[name] = maxY
			}
		}

		// Enforce rank ordering: nodes with higher ranks must have higher Y.
		// Re-sort the pre-sorted slice by position within same rank groups.
		sort.SliceStable(sortedByRank, func(i, j int) bool {
			ri, rj := ranks[sortedByRank[i]], ranks[sortedByRank[j]]
			if ri != rj {
				return ri < rj
			}
			return positions[sortedByRank[i]] < positions[sortedByRank[j]]
		})
		enforceRankOrder(positions, sortedByRank, ranks, minY, maxY, minSpacing, pipelineParents)

		damping *= 0.98
	}

	return positions
}

// enforceRankOrder restores the topological ordering invariant after
// force displacements.
//
// Nodes are sorted by rank, then swept forward per rank group: all
// nodes in a rank group are raised to at least (previous group's max Y
// + minSpacing). Within a rank group, the relative ordering produced by
// repulsion is preserved but no minimum spacing is enforced — this
// prevents same-rank separation from compressing deeper nodes toward
// the bottom.
//
// Anchor nodes (rank 0) are left untouched because they are already
// pinned at the top of the map.

// spreadPipelineParents is a post-processing pass that ensures pipeline
// parents have sufficient vertical separation. It sorts all non-pipeline-member
// nodes by Y position, then sweeps forward: when two consecutive nodes include
// at least one pipeline parent, they must be at least pipelineSpacing apart.
// Pushed nodes are clamped to maxY.
func spreadPipelineParents(positions map[string]float64, pipelineMembers, pipelineParents map[string]bool, minSpacing, maxY float64) {
	pipelineSpacing := minSpacing * 3

	// Collect non-pipeline-member nodes sorted by Y position.
	sorted := make([]string, 0, len(positions))
	for name := range positions {
		if !pipelineMembers[name] {
			sorted = append(sorted, name)
		}
	}
	sort.Strings(sorted)
	sort.SliceStable(sorted, func(i, j int) bool {
		return positions[sorted[i]] < positions[sorted[j]]
	})

	// Sweep forward: push nodes down when a pipeline parent pair is too close.
	for i := 1; i < len(sorted); i++ {
		prev := sorted[i-1]
		cur := sorted[i]
		if !pipelineParents[prev] && !pipelineParents[cur] {
			continue
		}
		gap := positions[cur] - positions[prev]
		if gap < pipelineSpacing {
			positions[cur] = positions[prev] + pipelineSpacing
			if positions[cur] > maxY {
				positions[cur] = maxY
			}
		}
	}
}

func enforceRankOrder(positions map[string]float64, sorted []string, ranks map[string]int, minY, _ float64, minSpacing float64, pipelineParents map[string]bool) {
	// sorted is expected to be pre-sorted by rank then by position.
	// Sweep forward per rank group: enforce minSpacing only between
	// different rank groups, not between siblings at the same rank.
	// Pipeline parents use 3× minSpacing because their visual rectangle
	// is taller than a regular component.
	prevRank := -1
	prevRankMaxY := minY - minSpacing
	prevRankMaxName := ""
	currentRankMinAllowed := minY

	for _, name := range sorted {
		r := ranks[name]
		if r != prevRank {
			// Entering a new rank group — compute minimum allowed Y
			// based on the previous group's maximum.
			spacing := minSpacing
			if pipelineParents[prevRankMaxName] || pipelineParents[name] {
				spacing = minSpacing * 3
			}
			currentRankMinAllowed = prevRankMaxY + spacing
			prevRank = r
		}
		if r > 0 && positions[name] < currentRankMinAllowed {
			positions[name] = currentRankMinAllowed
		}
		// Track the max Y across the current rank group. When we
		// transition to the next rank, this becomes prevRankMaxY.
		if positions[name] > prevRankMaxY {
			prevRankMaxY = positions[name]
			prevRankMaxName = name
		}
	}
}

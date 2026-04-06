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

	if len(names) == 0 {
		return positions
	}

	minY := float64(opts.MinY)
	maxY := float64(opts.MaxY)
	minSpacing := float64(opts.MinSpacing)
	damping := opts.Damping

	// Compute ideal layer height for edge attraction
	layerHeight := (maxY - minY) / math.Max(float64(maxRank), 1)

	for iter := 0; iter < opts.ForceIterations; iter++ {
		displacements := make(map[string]float64, len(names))

		// Pairwise repulsion between all active nodes (including pinned for repulsion source)
		for i := 0; i < len(allNames); i++ {
			for j := i + 1; j < len(allNames); j++ {
				a, b := allNames[i], allNames[j]
				dy := positions[b] - positions[a]
				absDy := math.Abs(dy)

				if absDy < minSpacing {
					if absDy < 0.1 {
						absDy = 0.1
						if dy >= 0 {
							dy = 0.1
						} else {
							dy = -0.1
						}
					}
					force := opts.RepulsionStrength / (absDy * absDy)
					sign := 1.0
					if dy < 0 {
						sign = -1.0
					}
					// Only apply displacement to non-pinned nodes
					if !pinned[a] {
						displacements[a] -= force * sign
					}
					if !pinned[b] {
						displacements[b] += force * sign
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
				displacements[from] += delta
			}
			if !pinned[to] {
				displacements[to] -= delta
			}
		}

		// Apply displacements with damping
		for _, name := range names {
			positions[name] += displacements[name] * damping
			if positions[name] < minY {
				positions[name] = minY
			}
			if positions[name] > maxY {
				positions[name] = maxY
			}
		}

		// Enforce rank ordering: nodes with higher ranks must have higher Y
		enforceRankOrder(positions, allNames, ranks, minY, maxY, minSpacing)

		damping *= 0.98
	}

	return positions
}

// enforceRankOrder restores the topological ordering invariant after
// force displacements.
//
// Nodes are sorted by rank, then swept forward: each node's Y is raised
// to at least (previous node's Y + minSpacing) if its rank is higher than
// the previous node's rank. This prevents crossings where a downstream
// node ends up above an upstream one due to accumulated repulsive forces.
//
// Anchor nodes (rank 0) are left untouched because they are already
// pinned at the top of the map.
func enforceRankOrder(positions map[string]float64, names []string, ranks map[string]int, minY, maxY, minSpacing float64) {
	// Sort names by rank
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Slice(sorted, func(i, j int) bool {
		return ranks[sorted[i]] < ranks[sorted[j]]
	})

	// Sweep forward: ensure each node is at least minSpacing below the previous rank group
	prevRank := -1
	prevMaxY := minY - minSpacing
	for _, name := range sorted {
		r := ranks[name]
		if r != prevRank {
			prevRank = r
		}
		minAllowed := prevMaxY + minSpacing
		if r > 0 && positions[name] < minAllowed {
			positions[name] = minAllowed
		}
		if positions[name] > prevMaxY {
			prevMaxY = positions[name]
		}
	}
}

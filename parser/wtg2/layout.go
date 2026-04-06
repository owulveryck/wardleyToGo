package wtg2

// ComputeYPositions assigns Y coordinates based on dependency depth.
// Anchors get rank 0 (top of map, low Y). Each dependency level pushes Y downward.
// The result maps node names to Y coordinates in the 0-100 range.
func ComputeYPositions(doc *Document) map[string]int {
	// Build adjacency and find all node names
	nodeSet := make(map[string]bool)
	children := make(map[string][]string) // parent -> children
	hasParent := make(map[string]bool)

	for _, n := range doc.Nodes {
		nodeSet[n.Name] = true
	}
	for _, pl := range doc.Pipelines {
		for _, m := range pl.Members {
			nodeSet[m.Name] = true
		}
	}

	for _, e := range doc.Edges {
		from := e.From
		to := e.To
		// Pipeline:Member refs - use just the member name for lookup
		if idx := colonMemberIndex(from); idx >= 0 {
			from = from[idx+1:]
		}
		if idx := colonMemberIndex(to); idx >= 0 {
			to = to[idx+1:]
		}
		if nodeSet[from] && nodeSet[to] {
			children[from] = append(children[from], to)
			hasParent[to] = true
		}
	}

	// Find roots: anchors first, then any node with no incoming edges
	var roots []string
	for _, n := range doc.Nodes {
		if n.Kind == KindAnchor {
			roots = append(roots, n.Name)
		}
	}
	if len(roots) == 0 {
		for name := range nodeSet {
			if !hasParent[name] {
				roots = append(roots, name)
			}
		}
	}

	// BFS to assign ranks
	rank := make(map[string]int)
	for name := range nodeSet {
		rank[name] = -1
	}
	queue := make([]string, 0, len(roots))
	for _, r := range roots {
		rank[r] = 0
		queue = append(queue, r)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, child := range children[current] {
			newRank := rank[current] + 1
			if newRank > rank[child] {
				rank[child] = newRank
				queue = append(queue, child)
			}
		}
	}

	// Assign unvisited nodes a default rank
	maxRank := 0
	for _, r := range rank {
		if r > maxRank {
			maxRank = r
		}
	}
	for name, r := range rank {
		if r < 0 {
			maxRank++
			rank[name] = maxRank
		}
	}
	// Recalculate maxRank after defaults
	for _, r := range rank {
		if r > maxRank {
			maxRank = r
		}
	}

	// Convert rank to Y coordinate: rank 0 = top (Y=3), max rank = bottom (Y=97)
	yPos := make(map[string]int)
	for name, r := range rank {
		if maxRank == 0 {
			yPos[name] = 3
		} else {
			yPos[name] = r*94/maxRank + 3
		}
	}

	// Pipeline members inherit their pipeline parent's Y
	for _, pl := range doc.Pipelines {
		if parentY, ok := yPos[pl.Name]; ok {
			for _, m := range pl.Members {
				yPos[m.Name] = parentY
			}
		}
	}

	return yPos
}

// colonMemberIndex returns the index of ':' in a "Pipeline:Member" reference,
// but only if there are no spaces around the colon (distinguishing from " : ").
func colonMemberIndex(s string) int {
	idx := -1
	for i, c := range s {
		if c == ':' {
			// Check no space before or after
			if i > 0 && s[i-1] != ' ' && i < len(s)-1 && s[i+1] != ' ' {
				idx = i
				break
			}
		}
	}
	return idx
}

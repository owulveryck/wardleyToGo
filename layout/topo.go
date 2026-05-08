package layout

import (
	"fmt"
	"sort"
)

// topoRanks computes a longest-path rank assignment for every node in g
// using a breadth-first search (BFS).
//
// The algorithm works as follows:
//
//  1. Build a forward-adjacency list (children) from g.Edges.
//     Edge names that use the "Pipeline:Member" syntax are resolved
//     to the member name via [resolveMember].
//
//  2. Identify root nodes. If any node has [KindAnchor], all anchors
//     become roots. Otherwise, every node with no incoming edges is a root.
//
//  3. BFS from every root with initial rank 0. For each edge parent→child,
//     set rank[child] = max(rank[child], rank[parent]+1). Because a node
//     may be reachable through paths of different lengths, it can be
//     enqueued multiple times — the longest path wins.
//
//  4. Any node still unvisited after BFS (disconnected from all roots)
//     is assigned a unique rank beyond the current maximum, ensuring it
//     appears at the bottom of the map.
//
// Returns (ranks, maxRank) where ranks maps node name → rank and maxRank
// is the highest rank value in the map.
//
// Complexity: O(V + E) for the BFS pass, where V = |Nodes| and E = |Edges|.
// The disconnected-node fallback adds at most V extra assignments.
func topoRanks(g *Graph) (map[string]int, int, error) {
	// Build adjacency and collect all node names
	nodeSet := make(map[string]bool)
	children := make(map[string][]string)
	hasParent := make(map[string]bool)

	for _, n := range g.Nodes {
		nodeSet[n.Name] = true
	}
	for _, pl := range g.Pipelines {
		for _, m := range pl.Members {
			nodeSet[m] = true
			children[pl.Parent] = append(children[pl.Parent], m)
			hasParent[m] = true
		}
	}

	for _, e := range g.Edges {
		from := resolveMember(e.From)
		to := resolveMember(e.To)
		if nodeSet[from] && nodeSet[to] {
			children[from] = append(children[from], to)
			hasParent[to] = true
		}
	}

	// Detect cycles before BFS — a cycle would cause infinite looping.
	if err := detectCycle(nodeSet, children); err != nil {
		return nil, 0, err
	}

	// Find roots: anchors first, then nodes with no incoming edges
	anchorSet := make(map[string]bool)
	var roots []string
	for _, n := range g.Nodes {
		if n.Kind == KindAnchor {
			roots = append(roots, n.Name)
			anchorSet[n.Name] = true
		}
	}
	if len(roots) == 0 {
		for name := range nodeSet {
			if !hasParent[name] {
				roots = append(roots, name)
			}
		}
		sort.Strings(roots)
	}

	// BFS longest-path rank assignment
	rank := make(map[string]int, len(nodeSet))
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
			if anchorSet[child] {
				continue
			}
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
	var unvisited []string
	for name, r := range rank {
		if r < 0 {
			unvisited = append(unvisited, name)
		}
	}
	sort.Strings(unvisited)
	for _, name := range unvisited {
		minChildRank := -1
		for _, child := range children[name] {
			if rank[child] >= 0 && (minChildRank < 0 || rank[child] < minChildRank) {
				minChildRank = rank[child]
			}
		}
		if minChildRank > 0 {
			rank[name] = minChildRank - 1
		} else if minChildRank == 0 {
			rank[name] = 0
		} else {
			maxRank++
			rank[name] = maxRank
		}
	}
	// Recalculate maxRank
	for _, r := range rank {
		if r > maxRank {
			maxRank = r
		}
	}

	return rank, maxRank, nil
}

// detectCycle uses DFS with white/gray/black coloring to find cycles.
// A back edge (current → gray node) proves a cycle exists.
func detectCycle(nodes map[string]bool, children map[string][]string) error {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int, len(nodes))

	// Sort node names for deterministic error messages.
	sorted := make([]string, 0, len(nodes))
	for n := range nodes {
		sorted = append(sorted, n)
	}
	sort.Strings(sorted)

	var dfs func(string) error
	dfs = func(node string) error {
		color[node] = gray
		for _, child := range children[node] {
			if !nodes[child] {
				continue
			}
			switch color[child] {
			case gray:
				return fmt.Errorf("cycle detected: %s -> %s", node, child)
			case white:
				if err := dfs(child); err != nil {
					return err
				}
			}
		}
		color[node] = black
		return nil
	}

	for _, n := range sorted {
		if color[n] == white {
			if err := dfs(n); err != nil {
				return err
			}
		}
	}
	return nil
}

// resolveMember extracts the member name from a "Pipeline:Member" reference.
//
// If s contains a colon with no spaces on either side (e.g., "Engine:AlgoA"),
// the substring after the colon is returned ("AlgoA"). If the colon has
// adjacent spaces (e.g., "A : B") or no colon is present, s is returned
// unchanged.
//
// This convention allows edges to target pipeline members without
// requiring callers to flatten the pipeline structure first.
func resolveMember(s string) string {
	for i, c := range s {
		if c == ':' && i > 0 && s[i-1] != ' ' && i < len(s)-1 && s[i+1] != ' ' {
			return s[i+1:]
		}
	}
	return s
}

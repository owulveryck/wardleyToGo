// Package layout provides vertical placement algorithms for Wardley Map
// components.
//
// A Wardley Map positions components along two axes: evolution (horizontal,
// X) and value chain / visibility (vertical, Y). This package handles the
// vertical axis — it assigns Y coordinates to nodes based on their
// dependency relationships, ensuring that upstream components (anchors,
// user needs) sit at the top and downstream components (infrastructure,
// commodities) sit at the bottom.
//
// The default algorithm works in three phases:
//
//  1. Topological rank assignment — a longest-path BFS from root nodes
//     (anchors or nodes with no incoming edges) assigns each node a rank
//     equal to its maximum depth in the dependency DAG.
//
//  2. Force-directed spacing — an iterative simulation applies repulsive
//     forces between nodes that are too close and attractive forces along
//     edges to preserve ideal layer separation. Anchors are pinned in
//     place; rank ordering is enforced after each iteration to prevent
//     crossings.
//
//  3. Normalization — floating-point positions are rounded to integers,
//     clamped to the [MinY, MaxY] range, and pipeline members inherit
//     their parent's Y coordinate.
//
// # Usage
//
//	g := &layout.Graph{
//	    Nodes: []layout.Node{
//	        {Name: "User Need", Kind: layout.KindAnchor},
//	        {Name: "Web App"},
//	        {Name: "Database"},
//	    },
//	    Edges: []layout.Edge{
//	        {From: "User Need", To: "Web App"},
//	        {From: "Web App", To: "Database"},
//	    },
//	}
//	l := layout.New(layout.DefaultOptions())
//	positions := l.Layout(g) // map["User Need"]→3, "Web App"→50, "Database"→97
//
// # Custom options
//
// Call [DefaultOptions] and override individual fields:
//
//	opts := layout.DefaultOptions()
//	opts.MinSpacing = 10          // wider gaps between layers
//	opts.ForceIterations = 100    // more simulation passes
//	l := layout.New(opts)
package layout

import "math"

// NodeKind classifies nodes for layout purposes.
//
// The layout algorithm treats anchors specially: they are forced to rank 0
// (topmost position) and pinned during force simulation.
type NodeKind int

const (
	// KindRegular represents an ordinary component that is positioned
	// by the layout algorithm based on its dependencies.
	KindRegular NodeKind = iota

	// KindAnchor represents a user need or market anchor. Anchors are
	// always placed at rank 0 (top of the map) and their positions are
	// pinned during force-directed spacing.
	KindAnchor
)

// Node is a lightweight, layout-only representation of a graph node.
// It carries only the information needed for vertical placement:
// a unique name and a kind.
type Node struct {
	// Name uniquely identifies this node within the graph.
	// It is used as the key in the returned position map.
	Name string

	// Kind determines how the layout algorithm treats this node.
	// Anchors are pinned at the top; regular nodes are placed by
	// dependency depth.
	Kind NodeKind
}

// Edge represents a directed dependency between two nodes, identified
// by name. The direction is from a higher-visibility (upstream) node
// to a lower-visibility (downstream) node: From depends on To, so To
// appears below From on the map.
//
// Edge names may use the "Pipeline:Member" syntax (e.g., "Engine:AlgoA")
// to reference a member inside a pipeline. The layout engine resolves
// these references automatically.
type Edge struct {
	From string
	To   string
}

// Pipeline groups a parent component with its internal members.
// After layout, all members receive the same Y coordinate as their
// parent, since pipeline members share a single horizontal band.
type Pipeline struct {
	// Parent is the name of the pipeline component.
	Parent string
	// Members lists the names of components inside the pipeline.
	Members []string
}

// Graph is the input to a [Layouter]. It describes the dependency
// structure of a Wardley Map without coupling to any specific parser
// or AST.
//
// Callers should populate Nodes with every component (including
// pipeline members if they also appear as standalone nodes), Edges
// with all dependency relationships, and Pipelines with grouping
// information.
type Graph struct {
	Nodes     []Node
	Edges     []Edge
	Pipelines []Pipeline
}

// Options controls the behaviour of the layout algorithm.
//
// All fields have sensible defaults returned by [DefaultOptions].
// Adjust them to trade off between spacing quality and computation time.
type Options struct {
	// MinY is the smallest Y coordinate that will be assigned.
	// It provides a top margin on the map. Default: 3.
	MinY int

	// MaxY is the largest Y coordinate that will be assigned.
	// It provides a bottom margin on the map. Default: 97.
	MaxY int

	// MinSpacing is the minimum vertical gap between adjacent layers,
	// expressed in the [MinY, MaxY] coordinate space. The force-directed
	// pass pushes nodes apart until this minimum is satisfied. Default: 5.
	MinSpacing int

	// ForceIterations is the number of repulsion/attraction passes.
	// More iterations produce better spacing at the cost of computation
	// time. For typical maps (< 50 nodes), 50 iterations complete in
	// under 10 ms. Default: 50.
	ForceIterations int

	// RepulsionStrength controls how strongly nodes push each other
	// apart when they are closer than MinSpacing. Higher values produce
	// faster convergence but may cause oscillation. Default: 8.0.
	RepulsionStrength float64

	// Damping is the fraction of computed displacement that is actually
	// applied each iteration. It decays by 2 % per iteration to help
	// convergence. Values close to 1.0 converge faster but risk
	// oscillation; values close to 0.0 are very stable but slow.
	// Default: 0.85.
	Damping float64
}

// DefaultOptions returns sensible defaults for Wardley Map layout.
//
// The returned options produce good results for maps with up to 50 nodes
// on a 100-unit vertical axis with 3-unit margins at top and bottom.
func DefaultOptions() Options {
	return Options{
		MinY:              3,
		MaxY:              97,
		MinSpacing:        5,
		ForceIterations:   50,
		RepulsionStrength: 8.0,
		Damping:           0.85,
	}
}

// Layouter computes vertical (Y) positions for all nodes in a [Graph].
//
// Implementations must return a map from node name to Y coordinate,
// where every node in the graph (including pipeline members) has an
// entry. All Y values must be within the [Options.MinY, Options.MaxY]
// range.
type Layouter interface {
	// Layout assigns Y coordinates to every node in g.
	// The returned map uses node names as keys and Y positions as values.
	// It returns an error if the graph contains cycles.
	Layout(g *Graph) (map[string]int, error)
}

// defaultLayouter is the standard implementation combining topological
// sort with force-directed spacing.
type defaultLayouter struct {
	opts Options
}

// New returns the default [Layouter] that combines topological sort
// with force-directed spacing.
//
// The algorithm is deterministic for a given graph and options.
func New(opts Options) Layouter {
	return &defaultLayouter{opts: opts}
}

// Layout implements [Layouter] using topological rank assignment followed
// by force-directed spacing.
func (l *defaultLayouter) Layout(g *Graph) (map[string]int, error) {
	// Phase 1: topological rank assignment
	ranks, _, err := topoRanks(g)
	if err != nil {
		return nil, err
	}

	// Collect pipeline members early so we can compute effective max rank
	pipelineMembers := make(map[string]bool)
	for _, pl := range g.Pipelines {
		for _, m := range pl.Members {
			pipelineMembers[m] = true
		}
	}

	// Collect pipeline parents for force-directed spacing: pipeline
	// parents need extra vertical clearance because they render as
	// rectangles taller than regular component circles.
	pipelineParents := make(map[string]bool)
	for _, pl := range g.Pipelines {
		pipelineParents[pl.Parent] = true
	}

	// Compute effectiveMaxRank: highest rank among non-pipeline-member
	// nodes. Disconnected pipeline members get inflated ranks from
	// topoRanks, but they inherit their parent's Y later, so their
	// ranks should not compress the layout for other nodes.
	effectiveMaxRank := 0
	for name, r := range ranks {
		if !pipelineMembers[name] && r > effectiveMaxRank {
			effectiveMaxRank = r
		}
	}

	// Phase 2: initial linear placement then force-directed spacing
	positions := make(map[string]float64, len(ranks))
	minY := float64(l.opts.MinY)
	maxY := float64(l.opts.MaxY)

	// Adaptive pipeline spacing: scale down when many pipeline parents
	// would otherwise exceed available vertical space.
	pipelineSpacing := float64(l.opts.MinSpacing) * 3.0
	if len(pipelineParents) > 0 {
		available := maxY - minY
		maxConsumption := available * 0.5
		idealTotal := float64(len(pipelineParents)) * pipelineSpacing
		if idealTotal > maxConsumption {
			pipelineSpacing = maxConsumption / float64(len(pipelineParents))
			if pipelineSpacing < float64(l.opts.MinSpacing) {
				pipelineSpacing = float64(l.opts.MinSpacing)
			}
		}
	}

	for name, r := range ranks {
		if pipelineMembers[name] {
			continue // pipeline members inherit parent Y in post-processing
		}
		if effectiveMaxRank == 0 {
			positions[name] = minY
		} else {
			positions[name] = minY + float64(r)*(maxY-minY)/float64(effectiveMaxRank)
		}
	}

	positions = forceSpread(positions, g.Edges, pipelineMembers, pipelineParents, ranks, effectiveMaxRank, l.opts, pipelineSpacing)

	// Phase 2.5: separate pipeline parents that are too close.
	// Pipeline parents render as rectangles taller than regular component
	// circles, so they need extra clearance. Sort all non-pipeline-member
	// nodes by Y, then push apart consecutive nodes where either is a
	// pipeline parent.
	if len(pipelineParents) > 1 {
		spreadPipelineParents(positions, pipelineMembers, pipelineParents, maxY, pipelineSpacing)
	}

	// Phase 3: round to int, apply pipeline inheritance, clamp
	result := make(map[string]int, len(positions))
	for name, pos := range positions {
		y := int(math.Round(pos))
		if y < l.opts.MinY {
			y = l.opts.MinY
		}
		if y > l.opts.MaxY {
			y = l.opts.MaxY
		}
		result[name] = y
	}

	// Pipeline members inherit their parent's Y
	for _, pl := range g.Pipelines {
		if parentY, ok := result[pl.Parent]; ok {
			for _, m := range pl.Members {
				result[m] = parentY
			}
		}
	}

	return result, nil
}

// TopoRanks returns the topological depth of each node from root anchors.
// Anchors have rank 0, their direct children rank 1, etc.
func TopoRanks(g *Graph) (map[string]int, error) {
	ranks, _, err := topoRanks(g)
	return ranks, err
}

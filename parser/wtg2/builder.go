package wtg2

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/layout"
	"github.com/owulveryck/wardleyToGo/layout/labels"
)

// BuildResult holds the map and evolution stages produced by the builder.
type BuildResult struct {
	Map         *wardleyToGo.Map
	Stages      []svgmap.Evolution
	Legend      bool
	LegendItems []svgmap.LegendItem
	Focus       *FocusSet // nil if no focus directive is present
}

// FocusSet identifies which elements should remain at full opacity when focus is active.
type FocusSet struct {
	ComponentIDs map[int64]bool    // Focused component IDs
	EdgeKeys     map[[2]int64]bool // Focused edges as (fromID, toID) pairs
	GroupIDs     map[int64]bool    // Groups containing focused components
}

type nodeEntry struct {
	node     wardleyToGo.Component
	decl     *NodeDecl
	isAnchor bool
}

// BuildMap converts a parsed Document AST into a wardleyToGo.Map.
func BuildMap(doc *Document) (*BuildResult, error) {
	m := wardleyToGo.NewMap(0)
	m.Title = doc.Title

	var nextID int64
	newID := func() int64 {
		id := nextID
		nextID++
		return id
	}

	nodeDict := make(map[string]*nodeEntry)
	evolvedMap := make(map[int64]int64) // original component ID → evolved component ID

	// Compute Y positions from dependency graph
	lg := documentToLayoutGraph(doc)
	layouter := layout.New(layout.DefaultOptions())
	yPositions, err := layouter.Layout(lg)
	if err != nil {
		return nil, fmt.Errorf("layout: %w", err)
	}

	// Phase 1: Create all declared nodes
	for _, nd := range doc.Nodes {
		y := yPositions[nd.Name]

		switch nd.Kind {
		case KindAnchor:
			anchor := wardley.NewAnchor(newID())
			anchor.Label = nd.Name
			x := 50 // default center
			if nd.Evolution != "" {
				var err error
				x, err = ParsePosition(nd.Evolution)
				if err != nil {
					return nil, fmt.Errorf("anchor %q: %w", nd.Name, err)
				}
			}
			anchor.Placement = image.Pt(x, y)
			if err := m.AddComponent(anchor); err != nil {
				return nil, fmt.Errorf("anchor %q: %w", nd.Name, err)
			}
			nodeDict[nd.Name] = &nodeEntry{node: anchor, decl: nd, isAnchor: true}

		case KindComponent, KindSubmap:
			comp := wardley.NewComponent(newID())
			comp.Label = nd.Name
			comp.Configured = true

			x := 50
			if nd.Evolution != "" {
				var err error
				x, err = ParsePosition(nd.Evolution)
				if err != nil {
					return nil, fmt.Errorf("component %q: %w", nd.Name, err)
				}
			}
			comp.Placement = image.Pt(x, y)

			// Set type
			switch strings.ToLower(nd.Type) {
			case "build":
				comp.Type = wardley.BuildComponent
			case "buy":
				comp.Type = wardley.BuyComponent
			case "outsource":
				comp.Type = wardley.OutsourceComponent
			}

			// Set color
			if nd.Color != "" {
				if c, err := parseHexColor(nd.Color); err == nil {
					comp.Color = c
				}
			}

			// Build description from note, cost, and asset for SVG tooltip
			var descParts []string
			if nd.Note != "" {
				descParts = append(descParts, nd.Note)
			}
			if nd.Asset != "" {
				comp.Asset = nd.Asset
				descParts = append(descParts, "Asset: "+nd.Asset)
			}
			if nd.Cost != "" {
				comp.Cost = nd.Cost
				descParts = append(descParts, "Cost: "+nd.Cost)
			}
			if len(descParts) > 0 {
				comp.Description = strings.Join(descParts, " | ")
			}

			comp.Inertia = nd.Inertia
			comp.InertiaKinds = nd.InertiaKinds

			if err := m.AddComponent(comp); err != nil {
				return nil, fmt.Errorf("component %q: %w", nd.Name, err)
			}
			nodeDict[nd.Name] = &nodeEntry{node: comp, decl: nd}

			// Handle evolution movement (>>)
			if nd.EvolvedTo != "" {
				evolvedX, err := ParsePosition(nd.EvolvedTo)
				if err != nil {
					return nil, fmt.Errorf("component %q evolved position: %w", nd.Name, err)
				}
				evolved := wardley.NewEvolvedComponent(newID())
				evolved.Label = nd.Name
				evolved.Placement = image.Pt(evolvedX, y)
				evolved.Configured = true
				if err := m.AddComponent(evolved); err != nil {
					return nil, fmt.Errorf("evolved component %q: %w", nd.Name, err)
				}
				evolvedMap[comp.ID()] = evolved.ID()
				comp.Label = ""

				// Create evolution edge
				inertiaX := 0
				if nd.Inertia > 0 {
					inertiaX = x + (evolvedX-x)/3
				}
				collab := &wardley.Collaboration{
					F:            comp,
					T:            evolved,
					Type:         wardley.EvolvedComponentEdge,
					Inertia:      image.Pt(inertiaX, 0),
					InertiaKinds: nd.InertiaKinds,
				}
				if err := m.SetCollaboration(collab); err != nil {
					return nil, fmt.Errorf("evolution edge for %q: %w", nd.Name, err)
				}
			}
		}
	}

	// Phase 2: Create pipeline components
	for _, pl := range doc.Pipelines {
		entry, ok := nodeDict[pl.Name]
		if !ok {
			continue
		}
		comp, ok := entry.node.(*wardley.Component)
		if !ok {
			continue
		}
		comp.Type = wardley.PipelineComponent

		parentY := comp.Placement.Y
		for _, member := range pl.Members {
			memberComp := wardley.NewComponent(newID())
			memberComp.Label = member.Name
			memberComp.Configured = true
			x, err := ParsePosition(member.Position)
			if err != nil {
				return nil, fmt.Errorf("pipeline member %q: %w", member.Name, err)
			}
			memberComp.Placement = image.Pt(x, parentY)
			memberComp.PipelineReference = comp
			comp.PipelinedComponents = append(comp.PipelinedComponents, memberComp)

			if err := m.AddComponent(memberComp); err != nil {
				return nil, fmt.Errorf("pipeline member %q: %w", member.Name, err)
			}
			nodeDict[member.Name] = &nodeEntry{node: memberComp}
		}
	}

	// Phase 2.5: Automatic label placement.
	// Collect all components (including anchors and evolved ones).
	type labelTarget struct {
		comp    *wardley.Component
		evolved *wardley.EvolvedComponent
		anchor  *wardley.Anchor
	}
	allComponents := m.Components()
	labelTargets := make(map[string]labelTarget)
	labelComps := make([]labels.Component, 0, len(allComponents))

	for _, n := range allComponents {
		switch c := n.(type) {
		case *wardley.Anchor:
			if c.LabelPlacement.X != components.UndefinedCoord {
				continue
			}
			key := "anchor:" + c.Label + "@" + strconv.Itoa(c.Placement.X) + "," + strconv.Itoa(c.Placement.Y)
			labelTargets[key] = labelTarget{anchor: c}
			labelComps = append(labelComps, labels.Component{
				Name:     key,
				Position: c.Placement,
				Label:    c.Label,
			})
		case *wardley.Component:
			if c.LabelPlacement.X != components.UndefinedCoord {
				continue
			}
			key := c.Label + "@" + strconv.Itoa(c.Placement.X) + "," + strconv.Itoa(c.Placement.Y)
			labelTargets[key] = labelTarget{comp: c}
			labelComps = append(labelComps, labels.Component{
				Name:       key,
				Position:   c.Placement,
				Label:      c.Label,
				IsPipeline: c.Type == wardley.PipelineComponent,
			})
		case *wardley.EvolvedComponent:
			if c.LabelPlacement.X != components.UndefinedCoord {
				continue
			}
			key := c.Label + "@" + strconv.Itoa(c.Placement.X) + "," + strconv.Itoa(c.Placement.Y)
			labelTargets[key] = labelTarget{evolved: c}
			labelComps = append(labelComps, labels.Component{
				Name:     key,
				Position: c.Placement,
				Label:    c.Label,
			})
		}
	}

	if len(labelComps) > 1 {
		placements := labels.PlaceLabels(labelComps, labels.DefaultOptions())
		for key, target := range labelTargets {
			r, ok := placements[key]
			if !ok {
				continue
			}
			// Map labels.Anchor* to wardley.Adjust* constants.
			adjustAnchor := wardley.AdjustUndefined
			switch r.Anchor {
			case labels.AnchorStart:
				adjustAnchor = wardley.AdjustStart
			case labels.AnchorMiddle:
				adjustAnchor = wardley.AdjustMiddle
			case labels.AnchorEnd:
				adjustAnchor = wardley.AdjustEnd
			}
			if target.comp != nil {
				target.comp.LabelPlacement = r.Offset
				if adjustAnchor != wardley.AdjustUndefined {
					target.comp.Anchor = adjustAnchor
				}
			}
			if target.evolved != nil {
				target.evolved.LabelPlacement = r.Offset
				if adjustAnchor != wardley.AdjustUndefined {
					target.evolved.Anchor = adjustAnchor
				}
			}
			if target.anchor != nil {
				target.anchor.LabelPlacement = r.Offset
				if adjustAnchor != wardley.AdjustUndefined {
					target.anchor.Anchor = adjustAnchor
				}
			}
		}
	}

	// Phase 3: Create edges
	for _, ed := range doc.Edges {
		fromNode := resolveNodeRef(ed.From, nodeDict)
		toNode := resolveNodeRef(ed.To, nodeDict)
		if fromNode == nil || toNode == nil {
			continue // skip unresolved references
		}

		collab := &wardley.Collaboration{
			F:    fromNode,
			T:    toNode,
			Type: wardley.RegularEdge,
		}
		if ed.Label != "" {
			collab.Label = ed.Label
		}
		if err := m.SetCollaboration(collab); err != nil {
			// Skip self-loops silently
			continue
		}
	}

	// Phase 3.5: Detect near-parallel edges and assign curve offsets.
	spreadOverlappingEdges(m)

	// Phase 3.7: Create group components
	defaultGroupColors := []color.RGBA{
		{0x34, 0x98, 0xDB, 0x30}, // blue
		{0x2E, 0xCC, 0x71, 0x30}, // green
		{0xE7, 0x4C, 0x3C, 0x30}, // red
		{0x9B, 0x59, 0xB6, 0x30}, // purple
		{0xE6, 0x7E, 0x22, 0x30}, // orange
		{0x1A, 0xBC, 0x9C, 0x30}, // teal
		{0xF3, 0x9C, 0x12, 0x30}, // yellow
		{0x34, 0x49, 0x5E, 0x30}, // dark blue
	}
	for i, gd := range doc.Groups {
		var memberPoints []image.Point
		for _, memberName := range gd.Members {
			entry, ok := nodeDict[memberName]
			if !ok {
				continue
			}
			memberPoints = append(memberPoints, entry.node.GetPosition())
		}
		if len(memberPoints) == 0 {
			continue
		}

		var fillColor color.Color
		if gd.Color != "" {
			c, err := parseHexColor(gd.Color)
			if err == nil {
				// Set alpha to 0x30 for translucency
				r, g, b, _ := c.RGBA()
				fillColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0x30}
			}
		}
		if fillColor == nil {
			fillColor = defaultGroupColors[i%len(defaultGroupColors)]
		}

		group := wardley.NewGroup(newID(), gd.Name, memberPoints, fillColor)
		group.TeamType = gd.Team
		if err := m.AddComponent(group); err != nil {
			return nil, fmt.Errorf("group %q: %w", gd.Name, err)
		}
	}

	// Phase 3.8: Attach signals, annotations, and gameplays to components
	for _, sig := range doc.Signals {
		entry, ok := nodeDict[sig.Target]
		if !ok {
			continue
		}
		if comp, ok := entry.node.(*wardley.Component); ok {
			comp.Signals = append(comp.Signals, wardley.ComponentSignal{Type: sig.Type})
		}
	}
	for _, ann := range doc.Annotations {
		entry, ok := nodeDict[ann.Target]
		if !ok {
			continue
		}
		if comp, ok := entry.node.(*wardley.Component); ok {
			comp.Annotations = append(comp.Annotations, wardley.ComponentAnnotation{
				Kind: ann.Kind,
				Text: ann.Text,
			})
		}
	}
	for _, gp := range doc.Gameplays {
		entry, ok := nodeDict[gp.Target]
		if !ok {
			continue
		}
		if comp, ok := entry.node.(*wardley.Component); ok {
			comp.Gameplays = append(comp.Gameplays, wardley.ComponentGameplay{
				Type: gp.Type,
				Text: gp.Text,
			})
		}
	}

	// Phase 4: Build evolution stages
	stages := buildEvolutionStages(doc.Stages)

	result := &BuildResult{
		Map:         m,
		Stages:      stages,
		Legend:      doc.Legend,
		LegendItems: buildLegendItems(doc),
	}

	// Phase 5: Compute focus set
	if len(doc.Focuses) > 0 {
		fs := &FocusSet{
			ComponentIDs: make(map[int64]bool),
			EdgeKeys:     make(map[[2]int64]bool),
			GroupIDs:     make(map[int64]bool),
		}

		// Collect focus root IDs
		var roots []int64
		for _, fd := range doc.Focuses {
			entry, ok := nodeDict[fd.Target]
			if !ok {
				continue
			}
			roots = append(roots, entry.node.ID())

			// If root is a pipeline parent, include members
			if comp, ok := entry.node.(*wardley.Component); ok {
				for _, member := range comp.PipelinedComponents {
					roots = append(roots, member.ID())
				}
				// If root is a pipeline member, include parent
				if comp.PipelineReference != nil {
					roots = append(roots, comp.PipelineReference.ID())
					for _, sibling := range comp.PipelineReference.PipelinedComponents {
						roots = append(roots, sibling.ID())
					}
				}
			}
		}

		// DFS from each root to collect descendants
		var dfs func(id int64)
		dfs = func(id int64) {
			if fs.ComponentIDs[id] {
				return
			}
			fs.ComponentIDs[id] = true

			// Include evolved counterpart if present
			if evolvedID, ok := evolvedMap[id]; ok {
				fs.ComponentIDs[evolvedID] = true
				fs.EdgeKeys[[2]int64{id, evolvedID}] = true
			}

			for _, succ := range m.From(id) {
				fs.EdgeKeys[[2]int64{id, succ.ID()}] = true
				dfs(succ.ID())
			}
		}
		for _, rootID := range roots {
			dfs(rootID)
		}

		// Identify groups containing focused components
		for _, gd := range doc.Groups {
			groupFocused := false
			for _, memberName := range gd.Members {
				entry, ok := nodeDict[memberName]
				if !ok {
					continue
				}
				if fs.ComponentIDs[entry.node.ID()] {
					groupFocused = true
					break
				}
			}
			if groupFocused {
				// Find the Group component by label
				for _, c := range m.Components() {
					if g, ok := c.(*wardley.Group); ok && g.Label == gd.Name {
						fs.GroupIDs[g.ID()] = true
						break
					}
				}
			}
		}

		if len(fs.ComponentIDs) > 0 {
			result.Focus = fs
		}
	}

	return result, nil
}

// resolveNodeRef looks up a node by name, handling "Pipeline:Member" syntax.
func resolveNodeRef(ref string, nodeDict map[string]*nodeEntry) wardleyToGo.Component {
	// Check for Pipeline:Member syntax (colon without spaces)
	if idx := colonMemberIndex(ref); idx >= 0 {
		memberName := strings.TrimSpace(ref[idx+1:])
		if entry, ok := nodeDict[memberName]; ok {
			return entry.node
		}
		return nil
	}
	ref = strings.TrimSpace(ref)
	if entry, ok := nodeDict[ref]; ok {
		return entry.node
	}
	return nil
}

func buildEvolutionStages(stages [4]string) []svgmap.Evolution {
	roman := [4]string{"I", "II", "III", "IV"}
	return []svgmap.Evolution{
		{Position: 0.0, Label: roman[0], ZoneLabel: stages[0]},
		{Position: 0.25, Label: roman[1], ZoneLabel: stages[1]},
		{Position: 0.50, Label: roman[2], ZoneLabel: stages[2]},
		{Position: 0.75, Label: roman[3], ZoneLabel: stages[3]},
	}
}

// documentToLayoutGraph converts a Document AST to the abstract layout.Graph.
func documentToLayoutGraph(doc *Document) *layout.Graph {
	g := &layout.Graph{}
	for _, n := range doc.Nodes {
		kind := layout.KindRegular
		if n.Kind == KindAnchor {
			kind = layout.KindAnchor
		}
		g.Nodes = append(g.Nodes, layout.Node{Name: n.Name, Kind: kind})
	}
	for _, e := range doc.Edges {
		g.Edges = append(g.Edges, layout.Edge{From: e.From, To: e.To})
	}
	for _, pl := range doc.Pipelines {
		p := layout.Pipeline{Parent: pl.Name}
		for _, m := range pl.Members {
			p.Members = append(p.Members, m.Name)
			g.Nodes = append(g.Nodes, layout.Node{Name: m.Name, Kind: layout.KindRegular})
		}
		g.Pipelines = append(g.Pipelines, p)
	}
	return g
}

// colonMemberIndex returns the index of ':' in a "Pipeline:Member" reference,
// but only if there are no spaces around the colon (distinguishing from " : ").
func colonMemberIndex(s string) int {
	for i, c := range s {
		if c == ':' && i > 0 && s[i-1] != ' ' && i < len(s)-1 && s[i+1] != ' ' {
			return i
		}
	}
	return -1
}

// spreadOverlappingEdges detects edges that share an endpoint and have
// nearly parallel paths, then assigns CurveOffset to spread them apart
// visually using Bézier curves.
func spreadOverlappingEdges(m *wardleyToGo.Map) {
	collabs := m.Collaborations()

	// Collect concrete collaborations.
	edges := make([]*wardley.Collaboration, 0, len(collabs))
	for _, c := range collabs {
		if wc, ok := c.(*wardley.Collaboration); ok {
			edges = append(edges, wc)
		}
	}

	// For each pair of edges, check if they share an endpoint and
	// their non-shared endpoints are close enough that the lines
	// would visually overlap.
	const angleThreshold = 0.26   // ~15 degrees in radians
	const proximityThreshold = 15 // max distance between non-shared endpoints (100-unit space)
	const curvePixels = 20        // perpendicular offset in SVG pixels

	for i := 0; i < len(edges); i++ {
		for j := i + 1; j < len(edges); j++ {
			a := edges[i]
			b := edges[j]

			af := a.F.GetPosition()
			at := a.T.GetPosition()
			bf := b.F.GetPosition()
			bt := b.T.GetPosition()

			// Check if they share an endpoint.
			var shared, otherA, otherB image.Point
			switch {
			case af == bf:
				shared, otherA, otherB = af, at, bt
			case af == bt:
				shared, otherA, otherB = af, at, bf
			case at == bf:
				shared, otherA, otherB = at, af, bt
			case at == bt:
				shared, otherA, otherB = at, af, bf
			default:
				continue
			}

			// Check if the non-shared endpoints are close to each other.
			dx := float64(otherA.X - otherB.X)
			dy := float64(otherA.Y - otherB.Y)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > proximityThreshold {
				continue
			}

			// Compute angle between the two edges from the shared point.
			dxA := float64(otherA.X - shared.X)
			dyA := float64(otherA.Y - shared.Y)
			dxB := float64(otherB.X - shared.X)
			dyB := float64(otherB.Y - shared.Y)
			lenA := math.Sqrt(dxA*dxA + dyA*dyA)
			lenB := math.Sqrt(dxB*dxB + dyB*dyB)
			if lenA < 1 || lenB < 1 {
				continue
			}
			cosAngle := (dxA*dxB + dyA*dyB) / (lenA * lenB)
			if cosAngle > 1 {
				cosAngle = 1
			}
			if cosAngle < -1 {
				cosAngle = -1
			}
			angle := math.Acos(cosAngle)

			if angle < angleThreshold {
				if a.CurveOffset == 0 {
					a.CurveOffset = curvePixels
				}
				if b.CurveOffset == 0 {
					b.CurveOffset = -curvePixels
				}
			}
		}
	}
}

// buildLegendItems scans the Document AST and returns legend items for
// element types that are actually present in the map.
func buildLegendItems(doc *Document) []svgmap.LegendItem {
	var items []svgmap.LegendItem

	// Scan component types
	hasRegular := false
	hasBuild := false
	hasBuy := false
	hasOutsource := false
	hasEvolved := false
	hasUnqualifiedInertia := false
	seenInertiaKinds := make(map[string]bool)

	for _, nd := range doc.Nodes {
		if nd.Kind == KindAnchor {
			continue
		}
		switch strings.ToLower(nd.Type) {
		case "build":
			hasBuild = true
		case "buy":
			hasBuy = true
		case "outsource":
			hasOutsource = true
		default:
			hasRegular = true
		}
		if nd.EvolvedTo != "" {
			hasEvolved = true
		}
		if nd.Inertia > 0 {
			if len(nd.InertiaKinds) > 0 {
				for _, kind := range nd.InertiaKinds {
					seenInertiaKinds[kind] = true
				}
			} else {
				hasUnqualifiedInertia = true
			}
		}
	}

	if hasRegular {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Component", Type: "component"})
	}
	if hasBuild {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Build", Type: "build"})
	}
	if hasBuy {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Buy", Type: "buy"})
	}
	if hasOutsource {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Outsource", Type: "outsource"})
	}
	if hasEvolved {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Evolved", Type: "evolved"})
	}

	// Pipelines
	if len(doc.Pipelines) > 0 {
		items = append(items, svgmap.LegendItem{Category: "Components", Label: "Pipeline", Type: "pipeline"})
	}

	// Edges
	if len(doc.Edges) > 0 {
		items = append(items, svgmap.LegendItem{Category: "Edges", Label: "Dependency", Type: "edge"})
	}
	if hasEvolved {
		items = append(items, svgmap.LegendItem{Category: "Edges", Label: "Evolution", Type: "evolved_edge"})
	}
	if hasUnqualifiedInertia {
		items = append(items, svgmap.LegendItem{Category: "Edges", Label: "Inertia", Type: "inertia"})
	}
	for _, kind := range []string{"tech", "financial", "human", "relational", "social"} {
		if seenInertiaKinds[kind] {
			items = append(items, svgmap.LegendItem{
				Category: "Edges",
				Label:    "Inertia (" + inertiaKindLabel(kind) + ")",
				Type:     "inertia_" + kind,
				Color:    inertiaKindColorForLegend(kind),
			})
		}
	}

	// Groups — one entry per group with its name and color
	defaultGroupColors := []color.RGBA{
		{0x34, 0x98, 0xDB, 0xFF}, // blue
		{0x2E, 0xCC, 0x71, 0xFF}, // green
		{0xE7, 0x4C, 0x3C, 0xFF}, // red
		{0x9B, 0x59, 0xB6, 0xFF}, // purple
		{0xE6, 0x7E, 0x22, 0xFF}, // orange
		{0x1A, 0xBC, 0x9C, 0xFF}, // teal
		{0xF3, 0x9C, 0x12, 0xFF}, // yellow
		{0x34, 0x49, 0x5E, 0xFF}, // dark blue
	}
	for i, gd := range doc.Groups {
		var groupColor color.Color
		if gd.Color != "" {
			c, err := parseHexColor(gd.Color)
			if err == nil {
				groupColor = c
			}
		}
		if groupColor == nil {
			groupColor = defaultGroupColors[i%len(defaultGroupColors)]
		}
		items = append(items, svgmap.LegendItem{Category: "Groups", Label: gd.Name, Type: "group", Color: groupColor})
	}

	// Signals — one entry per distinct signal type
	seenSignals := make(map[string]bool)
	for _, s := range doc.Signals {
		if !seenSignals[s.Type] {
			seenSignals[s.Type] = true
			items = append(items, svgmap.LegendItem{Category: "Signals", Label: signalLabel(s.Type), Type: "signal_" + s.Type})
		}
	}

	// Gameplays
	if len(doc.Gameplays) > 0 {
		items = append(items, svgmap.LegendItem{Category: "Gameplays", Label: "Gameplay", Type: "gameplay"})
	}

	// Annotations
	if len(doc.Annotations) > 0 {
		items = append(items, svgmap.LegendItem{Category: "Other", Label: "Annotation", Type: "annotation"})
	}

	return items
}

// signalLabel returns a human-readable label for a signal type.
func signalLabel(signalType string) string {
	switch signalType {
	case "accelerating":
		return "Accelerating"
	case "declining":
		return "Declining"
	case "stagnating":
		return "Stagnating"
	case "co-evolution":
		return "Co-evolution"
	case "red-queen":
		return "Red Queen"
	case "commoditization":
		return "Commoditization"
	case "network-effects":
		return "Network effects"
	case "economies-of-scale":
		return "Economies of scale"
	default:
		if len(signalType) == 0 {
			return signalType
		}
		return strings.ToUpper(signalType[:1]) + signalType[1:]
	}
}

func inertiaKindLabel(kind string) string {
	switch kind {
	case "tech":
		return "Tech"
	case "financial":
		return "Financial"
	case "human":
		return "Human"
	case "relational":
		return "Relational"
	case "social":
		return "Social"
	default:
		if len(kind) == 0 {
			return kind
		}
		return strings.ToUpper(kind[:1]) + kind[1:]
	}
}

func inertiaKindColorForLegend(kind string) color.RGBA {
	switch kind {
	case "tech":
		return color.RGBA{0x29, 0x80, 0xB9, 0xFF}
	case "financial":
		return color.RGBA{0x27, 0xAE, 0x60, 0xFF}
	case "human":
		return color.RGBA{0xE6, 0x7E, 0x22, 0xFF}
	case "relational":
		return color.RGBA{0x8E, 0x44, 0xAD, 0xFF}
	case "social":
		return color.RGBA{0x16, 0xA0, 0x85, 0xFF}
	default:
		return color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}
}

// parseHexColor parses a hex color string like "#3498DB" or "#abc".
func parseHexColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 6:
		b, err := hex.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("invalid hex color: %w", err)
		}
		return color.RGBA{R: b[0], G: b[1], B: b[2], A: 255}, nil
	case 3:
		b, err := hex.DecodeString(s[0:1] + s[0:1] + s[1:2] + s[1:2] + s[2:3] + s[2:3])
		if err != nil {
			return nil, fmt.Errorf("invalid hex color: %w", err)
		}
		return color.RGBA{R: b[0], G: b[1], B: b[2], A: 255}, nil
	default:
		return nil, fmt.Errorf("invalid hex color: %q", s)
	}
}

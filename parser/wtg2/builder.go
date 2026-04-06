package wtg2

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/owulveryck/wardleyToGo"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

// BuildResult holds the map and evolution stages produced by the builder.
type BuildResult struct {
	Map    *wardleyToGo.Map
	Stages []svgmap.Evolution
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

	// Compute Y positions from dependency graph
	yPositions := ComputeYPositions(doc)

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

			// Set description from note
			if nd.Note != "" {
				comp.Description = nd.Note
			}

			comp.Inertia = nd.Inertia

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

				// Create evolution edge
				inertiaX := 0
				if nd.Inertia > 0 {
					inertiaX = x + (evolvedX-x)/3
				}
				collab := &wardley.Collaboration{
					F:       comp,
					T:       evolved,
					Type:    wardley.EvolvedComponentEdge,
					Inertia: image.Pt(inertiaX, 0),
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

	// Phase 4: Build evolution stages
	stages := buildEvolutionStages(doc.Stages)

	return &BuildResult{
		Map:    m,
		Stages: stages,
	}, nil
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
	return []svgmap.Evolution{
		{Position: 0.0, Label: stages[0]},
		{Position: 0.25, Label: stages[1]},
		{Position: 0.50, Label: stages[2]},
		{Position: 0.75, Label: stages[3]},
	}
}

// parseHexColor parses a hex color string like "#3498DB" or "#abc".
func parseHexColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")
	var r, g, b uint8
	switch len(s) {
	case 6:
		_, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return nil, err
		}
	case 3:
		_, err := fmt.Sscanf(s, "%1x%1x%1x", &r, &g, &b)
		if err != nil {
			return nil, err
		}
		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
	default:
		return nil, fmt.Errorf("invalid hex color: %q", s)
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

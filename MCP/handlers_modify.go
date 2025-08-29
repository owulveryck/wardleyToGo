package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph"
)

// MODIFY HANDLERS - Tools for modifying existing maps (positioning, layout, configuration, removal)

func moveElementsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	movesJSON, err := request.RequireString("moves")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get moves parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Parse moves array
	type MoveOperation struct {
		Name string `json:"name"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	}

	var moves []MoveOperation
	if err := json.Unmarshal([]byte(movesJSON), &moves); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse moves JSON array", err), nil
	}

	// Process each move
	for _, move := range moves {
		// Validate coordinates
		if move.X < 0 || move.X > 100 || move.Y < 0 || move.Y > 100 {
			return mcp.NewToolResultErrorf("Invalid coordinates (%d, %d) for element '%s': both x and y must be between 0 and 100", move.X, move.Y, move.Name), nil
		}

		// Find the element
		element := findComponentByName(m, move.Name)
		if element == nil {
			return mcp.NewToolResultErrorf("Element '%s' not found in map", move.Name), nil
		}

		// Update position based on type
		switch e := element.(type) {
		case *wardley.Component:
			e.Placement = image.Pt(move.X, move.Y)
		case *wardley.Anchor:
			e.Placement = image.Pt(move.X, move.Y)
		default:
			return mcp.NewToolResultErrorf("Unknown element type for '%s'", move.Name), nil
		}
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func autoLayoutHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Check if map has any components
	nodeCount := 0
	nodes := m.Nodes()
	for nodes.Next() {
		nodeCount++
	}

	if nodeCount == 0 {
		return mcp.NewToolResultError("Cannot auto-position value chain: map has no components"), nil
	}

	// Apply value chain positioning algorithm
	positionComponentsInValueChain(m)

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output after value chain positioning", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func configureEvolutionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	stage1, err := request.RequireString("stage1")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get stage1 parameter", err), nil
	}

	stage2, err := request.RequireString("stage2")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get stage2 parameter", err), nil
	}

	stage3, err := request.RequireString("stage3")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get stage3 parameter", err), nil
	}

	stage4, err := request.RequireString("stage4")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get stage4 parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Create custom stages
	customStages := []JSONEvolution{
		{Position: 0, Label: stage1},
		{Position: 0.174, Label: stage2},
		{Position: 0.4, Label: stage3},
		{Position: 0.7, Label: stage4},
	}

	// Store the custom stages for this map
	mapStages[m.ID()] = customStages

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output with custom stages", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func removeElementsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	removeType, err := request.RequireString("remove_type")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get remove_type parameter", err), nil
	}

	itemsJSON, err := request.RequireString("items")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get items parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	if removeType == "elements" {
		// Remove elements (components or anchors)
		type RemoveElement struct {
			Name string `json:"name"`
		}

		var items []RemoveElement
		if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to parse items JSON array for elements", err), nil
		}

		for _, item := range items {
			element := findComponentByName(m, item.Name)
			if element == nil {
				return mcp.NewToolResultErrorf("Element '%s' not found in map", item.Name), nil
			}

			// Remove all edges connected to this element first
			edgesToRemove := make([]graph.Edge, 0)
			edges := m.Edges()
			for edges.Next() {
				edge := edges.Edge()
				if edge.From().ID() == element.ID() || edge.To().ID() == element.ID() {
					edgesToRemove = append(edgesToRemove, edge)
				}
			}

			for _, edge := range edgesToRemove {
				m.RemoveEdge(edge.From().ID(), edge.To().ID())
			}

			// Remove the element itself
			m.RemoveNode(element.ID())
		}

	} else if removeType == "links" {
		// Remove links
		type RemoveLink struct {
			From string `json:"from"`
			To   string `json:"to"`
		}

		var items []RemoveLink
		if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to parse items JSON array for links", err), nil
		}

		for _, item := range items {
			fromElement := findComponentByName(m, item.From)
			if fromElement == nil {
				return mcp.NewToolResultErrorf("Source element '%s' not found in map", item.From), nil
			}

			toElement := findComponentByName(m, item.To)
			if toElement == nil {
				return mcp.NewToolResultErrorf("Target element '%s' not found in map", item.To), nil
			}

			// Check if link exists and remove it
			if !linkExists(m, item.From, item.To) {
				return mcp.NewToolResultErrorf("No link found between '%s' and '%s'", item.From, item.To), nil
			}

			m.RemoveEdge(fromElement.ID(), toElement.ID())
		}

	} else {
		return mcp.NewToolResultErrorf("Invalid remove_type '%s': must be 'elements' or 'links'", removeType), nil
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

// Legacy handlers for backward compatibility

func moveComponentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	componentName, err := request.RequireString("component_name")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get component_name parameter", err), nil
	}

	x, err := request.RequireInt("x")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get x coordinate parameter", err), nil
	}

	y, err := request.RequireInt("y")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get y coordinate parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Validate coordinates
	if x < 0 || x > 100 || y < 0 || y > 100 {
		return mcp.NewToolResultErrorf("Invalid coordinates (%d, %d) for component '%s': both x and y must be between 0 and 100", x, y, componentName), nil
	}

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Find the component or anchor
	element := findComponentByName(m, componentName)
	if element == nil {
		return mcp.NewToolResultErrorf("Component or anchor '%s' not found in map", componentName), nil
	}

	// Update position based on type
	switch e := element.(type) {
	case *wardley.Component:
		e.Placement = image.Pt(x, y)
	case *wardley.Anchor:
		e.Placement = image.Pt(x, y)
	default:
		return mcp.NewToolResultErrorf("Unknown element type for '%s'", componentName), nil
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func linkComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	fromComponent, err := request.RequireString("from_component")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get from_component parameter", err), nil
	}

	toComponent, err := request.RequireString("to_component")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get to_component parameter", err), nil
	}

	linkType := request.GetString("link_type", "regular")
	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Find the components
	fromComp := findComponentByName(m, fromComponent)
	if fromComp == nil {
		return mcp.NewToolResultErrorf("Source component '%s' not found in map", fromComponent), nil
	}

	toComp := findComponentByName(m, toComponent)
	if toComp == nil {
		return mcp.NewToolResultErrorf("Target component '%s' not found in map", toComponent), nil
	}

	// Create collaboration
	collab := &wardley.Collaboration{
		F: fromComp,
		T: toComp,
	}

	// Set collaboration type
	switch linkType {
	case "evolved_component":
		collab.Type = wardley.EvolvedComponentEdge
	case "evolved":
		collab.Type = wardley.EvolvedEdge
	default:
		collab.Type = wardley.RegularEdge
	}

	// Add collaboration to map
	if err := m.SetCollaboration(collab); err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add link from '%s' to '%s'", fromComponent, toComponent), err), nil
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func unlinkComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	fromComponent, err := request.RequireString("from_component")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get from_component parameter", err), nil
	}

	toComponent, err := request.RequireString("to_component")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get to_component parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Find the components
	fromComp := findComponentByName(m, fromComponent)
	if fromComp == nil {
		return mcp.NewToolResultErrorf("Source component '%s' not found in map", fromComponent), nil
	}

	toComp := findComponentByName(m, toComponent)
	if toComp == nil {
		return mcp.NewToolResultErrorf("Target component '%s' not found in map", toComponent), nil
	}

	// Find and remove the collaboration
	edges := m.Edges()
	var edgeToRemove graph.Edge
	for edges.Next() {
		edge := edges.Edge()
		if edge.From().ID() == fromComp.ID() && edge.To().ID() == toComp.ID() {
			edgeToRemove = edge
			break
		}
	}

	if edgeToRemove == nil {
		return mcp.NewToolResultErrorf("No link found between '%s' and '%s'", fromComponent, toComponent), nil
	}

	// Remove the edge
	m.RemoveEdge(edgeToRemove.From().ID(), edgeToRemove.To().ID())

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

// BUILD HANDLERS - Tools for creating new maps and adding initial elements

func createMapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract optional parameters with defaults
	title := request.GetString("title", "New Wardley Map")
	mapID := request.GetInt("map_id", 1)
	output := request.GetString("output", "svg")

	// Create empty map
	m := wardleyToGo.NewMap(int64(mapID))
	m.Title = title

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output for empty map", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func addElementsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	elementsJSON, err := request.RequireString("elements")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get elements parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the existing map or create a new one
	var m *wardleyToGo.Map
	if mapJSON == "" || mapJSON == "{}" {
		// Create new map
		m = wardleyToGo.NewMap(1)
		m.Title = "New Wardley Map"
	} else {
		// Parse existing map
		m, err = UnmarshalMap([]byte(mapJSON))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
		}
	}

	// Parse elements array
	var inputElements []InputElement
	if err := json.Unmarshal([]byte(elementsJSON), &inputElements); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse elements JSON array", err), nil
	}

	// Find the next available ID
	nextID := int64(1)
	nodes := m.Nodes()
	for nodes.Next() {
		if nodes.Node().ID() >= nextID {
			nextID = nodes.Node().ID() + 1
		}
	}

	// Process each element
	for _, inputElement := range inputElements {
		// Validate coordinates
		if inputElement.X < 0 || inputElement.X > 100 || inputElement.Y < 0 || inputElement.Y > 100 {
			return mcp.NewToolResultErrorf("Invalid coordinates (%d, %d) for element '%s': both x and y must be between 0 and 100", inputElement.X, inputElement.Y, inputElement.Name), nil
		}

		// Default element type is component
		elementType := inputElement.ElementType
		if elementType == "" {
			elementType = "component"
		}

		// Check if element already exists
		existingElement := findComponentByName(m, inputElement.Name)

		if existingElement != nil {
			// Update existing element position
			switch e := existingElement.(type) {
			case *wardley.Component:
				e.Placement = image.Pt(inputElement.X, inputElement.Y)
				// Update type if provided and it's a component
				if elementType == "component" && inputElement.Type != "" {
					switch inputElement.Type {
					case "build":
						e.Type = wardley.BuildComponent
					case "buy":
						e.Type = wardley.BuyComponent
					case "outsource":
						e.Type = wardley.OutsourceComponent
					case "dataproduct":
						e.Type = wardley.DataProductComponent
					default:
						e.Type = wardley.RegularComponent
					}
				}
			case *wardley.Anchor:
				e.Placement = image.Pt(inputElement.X, inputElement.Y)
			}
		} else {
			// Create new element
			if elementType == "anchor" {
				// Create new anchor
				anchor := wardley.NewAnchor(nextID)
				anchor.Label = inputElement.Name
				anchor.Placement = image.Pt(inputElement.X, inputElement.Y)

				// Add anchor to map
				if err := m.AddComponent(anchor); err != nil {
					return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add anchor '%s' to map", anchor.Label), err), nil
				}
			} else {
				// Create new component
				comp := wardley.NewComponent(nextID)
				comp.Label = inputElement.Name
				comp.Placement = image.Pt(inputElement.X, inputElement.Y)

				// Set component type
				switch inputElement.Type {
				case "build":
					comp.Type = wardley.BuildComponent
				case "buy":
					comp.Type = wardley.BuyComponent
				case "outsource":
					comp.Type = wardley.OutsourceComponent
				case "dataproduct":
					comp.Type = wardley.DataProductComponent
				default:
					comp.Type = wardley.RegularComponent
				}

				// Add component to map
				if err := m.AddComponent(comp); err != nil {
					return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add component '%s' to map", comp.Label), err), nil
				}
			}

			nextID++
		}
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func addLinksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	linksJSON, err := request.RequireString("links")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get links parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
	}

	// Parse links array
	var inputLinks []InputLink
	if err := json.Unmarshal([]byte(linksJSON), &inputLinks); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse links JSON array", err), nil
	}

	// Process each link
	for _, inputLink := range inputLinks {
		// Find the components
		fromComp := findComponentByName(m, inputLink.From)
		if fromComp == nil {
			return mcp.NewToolResultErrorf("Source component '%s' not found in map for link", inputLink.From), nil
		}

		toComp := findComponentByName(m, inputLink.To)
		if toComp == nil {
			return mcp.NewToolResultErrorf("Target component '%s' not found in map for link", inputLink.To), nil
		}

		// Check if link already exists
		if linkExists(m, inputLink.From, inputLink.To) {
			continue // Skip existing links
		}

		// Create collaboration
		collab := &wardley.Collaboration{
			F: fromComp,
			T: toComp,
		}

		// Set collaboration type
		switch inputLink.Type {
		case "evolved_component":
			collab.Type = wardley.EvolvedComponentEdge
		case "evolved":
			collab.Type = wardley.EvolvedEdge
		default:
			collab.Type = wardley.RegularEdge
		}

		// Add collaboration to map
		if err := m.SetCollaboration(collab); err != nil {
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add link from '%s' to '%s'", inputLink.From, inputLink.To), err), nil
		}
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

// Legacy handler for backward compatibility
func addComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	componentsJSON, err := request.RequireString("components")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get components parameter", err), nil
	}

	output := request.GetString("output", "svg")

	// Parse the existing map or create a new one
	var m *wardleyToGo.Map
	if mapJSON == "" || mapJSON == "{}" {
		// Create new map
		m = wardleyToGo.NewMap(1)
		m.Title = "New Wardley Map"
	} else {
		// Parse existing map
		m, err = UnmarshalMap([]byte(mapJSON))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
		}
	}

	// Parse components array
	var inputComponents []InputComponent
	if err := json.Unmarshal([]byte(componentsJSON), &inputComponents); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to parse components JSON array", err), nil
	}

	// Find the next available ID
	nextID := int64(1)
	nodes := m.Nodes()
	for nodes.Next() {
		if nodes.Node().ID() >= nextID {
			nextID = nodes.Node().ID() + 1
		}
	}

	// Process each component
	for _, inputComp := range inputComponents {
		// Validate coordinates
		if inputComp.X < 0 || inputComp.X > 100 || inputComp.Y < 0 || inputComp.Y > 100 {
			return mcp.NewToolResultErrorf("Invalid coordinates (%d, %d) for component '%s': both x and y must be between 0 and 100", inputComp.X, inputComp.Y, inputComp.Name), nil
		}

		// Check if component already exists
		existingElement := findComponentByName(m, inputComp.Name)

		if existingElement != nil {
			// Only allow updating existing components, not anchors
			if existingComp, ok := existingElement.(*wardley.Component); ok {
				// Update existing component position
				existingComp.Placement = image.Pt(inputComp.X, inputComp.Y)

				// Update type if provided
				if inputComp.Type != "" {
					switch inputComp.Type {
					case "build":
						existingComp.Type = wardley.BuildComponent
					case "buy":
						existingComp.Type = wardley.BuyComponent
					case "outsource":
						existingComp.Type = wardley.OutsourceComponent
					case "dataproduct":
						existingComp.Type = wardley.DataProductComponent
					default:
						existingComp.Type = wardley.RegularComponent
					}
				}
			} else {
				// If it's an anchor, skip updating it in add_components
				continue
			}
		} else {
			// Create new component
			comp := wardley.NewComponent(nextID)
			comp.Label = inputComp.Name
			comp.Placement = image.Pt(inputComp.X, inputComp.Y)

			// Set component type
			switch inputComp.Type {
			case "build":
				comp.Type = wardley.BuildComponent
			case "buy":
				comp.Type = wardley.BuyComponent
			case "outsource":
				comp.Type = wardley.OutsourceComponent
			case "dataproduct":
				comp.Type = wardley.DataProductComponent
			default:
				comp.Type = wardley.RegularComponent
			}

			// Add component to map
			if err := m.AddComponent(comp); err != nil {
				return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add component '%s' to map", comp.Label), err), nil
			}

			nextID++
		}
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

func addAnchorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get map_json parameter", err), nil
	}

	anchorName, err := request.RequireString("anchor_name")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get anchor_name parameter", err), nil
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
		return mcp.NewToolResultErrorf("Invalid coordinates (%d, %d) for anchor '%s': both x and y must be between 0 and 100", x, y, anchorName), nil
	}

	// Parse the existing map or create a new one
	var m *wardleyToGo.Map
	if mapJSON == "" || mapJSON == "{}" {
		// Create new map
		m = wardleyToGo.NewMap(1)
		m.Title = "New Wardley Map"
	} else {
		// Parse existing map
		m, err = UnmarshalMap([]byte(mapJSON))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to parse map JSON", err), nil
		}
	}

	// Find the next available ID
	nextID := int64(1)
	nodes := m.Nodes()
	for nodes.Next() {
		if nodes.Node().ID() >= nextID {
			nextID = nodes.Node().ID() + 1
		}
	}

	// Create new anchor
	anchor := wardley.NewAnchor(nextID)
	anchor.Label = anchorName
	anchor.Placement = image.Pt(x, y)

	// Add anchor to map
	if err := m.AddComponent(anchor); err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to add anchor '%s' to map", anchorName), err), nil
	}

	// Generate output in requested format
	content, err := generateOutput(m, output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to generate output", err), nil
	}

	return mcp.NewToolResultText(content), nil
}

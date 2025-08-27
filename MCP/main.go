package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgencoding "github.com/owulveryck/wardleyToGo/encoding/svg"
	"gonum.org/v1/gonum/graph"
)

// JSONComponent represents a component in JSON format
type JSONComponent struct {
	ID    int64      `json:"id"`
	Name  string     `json:"name"`
	X     int        `json:"x"`
	Y     int        `json:"y"`
	Type  string     `json:"type,omitempty"`
	Color *JSONColor `json:"color,omitempty"`
}

// JSONColor represents a color in JSON format
type JSONColor struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

// JSONCollaboration represents a collaboration/edge in JSON format
type JSONCollaboration struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type,omitempty"`
}

// JSONEvolution represents an evolution stage in JSON format
type JSONEvolution struct {
	Position float64 `json:"position"`
	Label    string  `json:"label"`
}

// JSONAnchor represents an anchor in JSON format
type JSONAnchor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// JSONMap represents a Wardley map in JSON format
type JSONMap struct {
	ID             int64               `json:"id"`
	Title          string              `json:"title"`
	Components     []JSONComponent     `json:"components"`
	Collaborations []JSONCollaboration `json:"collaborations"`
	Anchors        []JSONAnchor        `json:"anchors,omitempty"`
	Stages         []JSONEvolution     `json:"stages,omitempty"`
}

// InputComponent represents a component in the input for add_components
type InputComponent struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type,omitempty"`
}

// InputLink represents a link in the input for add_links
type InputLink struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type,omitempty"`
}

// Global storage for map stages (in production, this could be a database)
var mapStages = make(map[int64][]JSONEvolution)

// getDefaultStages returns the default evolution stages
func getDefaultStages() []JSONEvolution {
	return []JSONEvolution{
		{Position: 0, Label: "Genesis / Concept"},
		{Position: 0.174, Label: "Custom / Emerging"},
		{Position: 0.4, Label: "Product / Converging"},
		{Position: 0.7, Label: "Commodity / Accepted"},
	}
}

// MarshalJSON converts a wardleyToGo.Map to JSON
func MarshalMap(m *wardleyToGo.Map) ([]byte, error) {
	// Get stages for this map, or use default if none set
	stages, exists := mapStages[m.ID()]
	if !exists {
		stages = getDefaultStages()
		mapStages[m.ID()] = stages
	}

	jsonMap := JSONMap{
		ID:             m.ID(),
		Title:          m.Title,
		Components:     make([]JSONComponent, 0),
		Collaborations: make([]JSONCollaboration, 0),
		Anchors:        make([]JSONAnchor, 0),
		Stages:         stages,
	}

	nodes := m.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		if comp, ok := node.(*wardley.Component); ok {
			jsonComp := JSONComponent{
				ID:   comp.ID(),
				Name: comp.Label,
				X:    comp.GetPosition().X,
				Y:    comp.GetPosition().Y,
			}

			// Add type if it's a specific type
			switch comp.Type {
			case wardley.BuildComponent:
				jsonComp.Type = "build"
			case wardley.BuyComponent:
				jsonComp.Type = "buy"
			case wardley.OutsourceComponent:
				jsonComp.Type = "outsource"
			case wardley.DataProductComponent:
				jsonComp.Type = "dataproduct"
			}

			// Add color if it's not black
			if comp.Color != nil {
				r, g, b, a := comp.Color.RGBA()
				if r != 0 || g != 0 || b != 0 || a != 65535 {
					jsonComp.Color = &JSONColor{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a >> 8),
					}
				}
			}

			jsonMap.Components = append(jsonMap.Components, jsonComp)
		}
		if anchor, ok := node.(*wardley.Anchor); ok {
			jsonAnchor := JSONAnchor{
				ID:   anchor.ID(),
				Name: anchor.Label,
				X:    anchor.GetPosition().X,
				Y:    anchor.GetPosition().Y,
			}
			jsonMap.Anchors = append(jsonMap.Anchors, jsonAnchor)
		}
	}

	// Add collaborations/edges
	edges := m.Edges()
	for edges.Next() {
		edge := edges.Edge()
		if collab, ok := edge.(*wardley.Collaboration); ok {
			// Find component names by ID
			fromName := findComponentNameByID(m, collab.From().ID())
			toName := findComponentNameByID(m, collab.To().ID())

			if fromName != "" && toName != "" {
				jsonCollab := JSONCollaboration{
					From: fromName,
					To:   toName,
				}

				// Add type if it's not regular
				switch collab.Type {
				case wardley.EvolvedComponentEdge:
					jsonCollab.Type = "evolved_component"
				case wardley.EvolvedEdge:
					jsonCollab.Type = "evolved"
				default:
					jsonCollab.Type = "regular"
				}

				jsonMap.Collaborations = append(jsonMap.Collaborations, jsonCollab)
			}
		}
	}

	return json.Marshal(jsonMap)
}

// Helper function to find component name by ID
func findComponentNameByID(m *wardleyToGo.Map, id int64) string {
	nodes := m.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		if comp, ok := node.(*wardley.Component); ok && comp.ID() == id {
			return comp.Label
		}
	}
	return ""
}

// Helper function to find component by name
func findComponentByName(m *wardleyToGo.Map, name string) *wardley.Component {
	nodes := m.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		if comp, ok := node.(*wardley.Component); ok && comp.Label == name {
			return comp
		}
	}
	return nil
}

// Helper function to check if a link already exists between two components
func linkExists(m *wardleyToGo.Map, fromName, toName string) bool {
	fromComp := findComponentByName(m, fromName)
	toComp := findComponentByName(m, toName)

	if fromComp == nil || toComp == nil {
		return false
	}

	edges := m.Edges()
	for edges.Next() {
		edge := edges.Edge()
		if edge.From().ID() == fromComp.ID() && edge.To().ID() == toComp.ID() {
			return true
		}
	}
	return false
}

// UnmarshalMap converts JSON to a wardleyToGo.Map
func UnmarshalMap(data []byte) (*wardleyToGo.Map, error) {
	var jsonMap JSONMap
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	m := wardleyToGo.NewMap(jsonMap.ID)
	m.Title = jsonMap.Title

	// Store stages for this map
	if len(jsonMap.Stages) > 0 {
		mapStages[jsonMap.ID] = jsonMap.Stages
	} else {
		// Use default stages if none provided
		mapStages[jsonMap.ID] = getDefaultStages()
	}

	for _, jsonComp := range jsonMap.Components {
		comp := wardley.NewComponent(jsonComp.ID)
		comp.Label = jsonComp.Name
		comp.Placement = image.Pt(jsonComp.X, jsonComp.Y)

		// Set type
		switch jsonComp.Type {
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

		// Set color if provided
		if jsonComp.Color != nil {
			comp.Color = color.RGBA{
				R: jsonComp.Color.R,
				G: jsonComp.Color.G,
				B: jsonComp.Color.B,
				A: jsonComp.Color.A,
			}
		}

		if err := m.AddComponent(comp); err != nil {
			return nil, fmt.Errorf("failed to add component %s: %w", comp.Label, err)
		}
	}

	// Add anchors
	for _, jsonAnchor := range jsonMap.Anchors {
		anchor := wardley.NewAnchor(jsonAnchor.ID)
		anchor.Label = jsonAnchor.Name
		anchor.Placement = image.Pt(jsonAnchor.X, jsonAnchor.Y)

		if err := m.AddComponent(anchor); err != nil {
			return nil, fmt.Errorf("failed to add anchor %s: %w", anchor.Label, err)
		}
	}

	// Add collaborations
	for _, jsonCollab := range jsonMap.Collaborations {
		fromComp := findComponentByName(m, jsonCollab.From)
		toComp := findComponentByName(m, jsonCollab.To)

		if fromComp == nil {
			return nil, fmt.Errorf("component '%s' not found for collaboration", jsonCollab.From)
		}
		if toComp == nil {
			return nil, fmt.Errorf("component '%s' not found for collaboration", jsonCollab.To)
		}

		collab := &wardley.Collaboration{
			F: fromComp,
			T: toComp,
		}

		// Set collaboration type
		switch jsonCollab.Type {
		case "evolved_component":
			collab.Type = wardley.EvolvedComponentEdge
		case "evolved":
			collab.Type = wardley.EvolvedEdge
		default:
			collab.Type = wardley.RegularEdge
		}

		if err := m.SetCollaboration(collab); err != nil {
			return nil, fmt.Errorf("failed to add collaboration from %s to %s: %w", jsonCollab.From, jsonCollab.To, err)
		}
	}

	return m, nil
}

// GenerateSVG creates an SVG representation of the map with embedded JSON data
// Uses the same configuration as wtg2svg example
func GenerateSVG(m *wardleyToGo.Map) (string, error) {
	var buf bytes.Buffer

	// Create encoder with laptop screen ratio (16:10)
	width, height := 1280, 800
	imageSize := image.Rect(0, 0, width, height)
	mapSize := image.Rect(30, 50, width-30, height-50)

	encoder, err := svgencoding.NewEncoder(&buf, imageSize, mapSize)
	if err != nil {
		return "", fmt.Errorf("failed to create SVG encoder: %w", err)
	}
	defer encoder.Close()

	// Get stages for this map and convert to Evolution structs
	stages, exists := mapStages[m.ID()]
	if !exists {
		stages = getDefaultStages()
		mapStages[m.ID()] = stages
	}

	evolutionSteps := make([]svgencoding.Evolution, len(stages))
	for i, stage := range stages {
		evolutionSteps[i] = svgencoding.Evolution{
			Position: stage.Position,
			Label:    stage.Label,
		}
	}

	// Use OctoStyle with evolution stages
	style := svgencoding.NewOctoStyle(evolutionSteps)
	style.WithSpace = true
	style.WithControls = false
	style.WithValueChain = true
	// WithIndicators defaults to false (no indicators)

	encoder.Init(style)

	if err := encoder.Encode(m); err != nil {
		return "", fmt.Errorf("failed to encode SVG: %w", err)
	}

	svgContent := buf.String()

	// Generate JSON representation
	jsonData, err := MarshalMap(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	// Embed JSON as comment in SVG
	// Find the position right after the opening <svg> tag to insert the comment
	svgTag := "<svg"
	tagIndex := strings.Index(svgContent, svgTag)
	if tagIndex == -1 {
		return "", fmt.Errorf("invalid SVG: could not find opening svg tag")
	}

	// Find the end of the opening tag
	tagEndIndex := strings.Index(svgContent[tagIndex:], ">")
	if tagEndIndex == -1 {
		return "", fmt.Errorf("invalid SVG: could not find end of opening svg tag")
	}
	tagEndIndex += tagIndex + 1

	// Insert the JSON comment right after the opening SVG tag
	comment := fmt.Sprintf("\n<!-- WARDLEY_MAP_DATA: %s -->\n", string(jsonData))
	result := svgContent[:tagEndIndex] + comment + svgContent[tagEndIndex:]

	return result, nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Wardley Map SVG Generator 🗺️",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add add_component tool
	addComponentTool := mcp.NewTool("add_component",
		mcp.WithDescription("Add a component to a Wardley map and return the updated map with SVG representation"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("component_name",
			mcp.Required(),
			mcp.Description("Name of the component to add"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("X coordinate (0-100)"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Y coordinate (0-100)"),
		),
		mcp.WithString("component_type",
			mcp.Description("Type of component: 'regular', 'build', 'buy', 'outsource', or 'dataproduct' (default: 'regular')"),
		),
	)

	// Add link_components tool
	linkComponentsTool := mcp.NewTool("link_components",
		mcp.WithDescription("Link two components in a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("from_component",
			mcp.Required(),
			mcp.Description("Name of the source component"),
		),
		mcp.WithString("to_component",
			mcp.Required(),
			mcp.Description("Name of the target component"),
		),
		mcp.WithString("link_type",
			mcp.Description("Type of link: 'regular', 'evolved_component', or 'evolved' (default: 'regular')"),
		),
	)

	// Add unlink_components tool
	unlinkComponentsTool := mcp.NewTool("unlink_components",
		mcp.WithDescription("Remove a link between two components in a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("from_component",
			mcp.Required(),
			mcp.Description("Name of the source component"),
		),
		mcp.WithString("to_component",
			mcp.Required(),
			mcp.Description("Name of the target component"),
		),
	)

	// Add move_component tool
	moveComponentTool := mcp.NewTool("move_component",
		mcp.WithDescription("Move a component to new coordinates in a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("component_name",
			mcp.Required(),
			mcp.Description("Name of the component to move"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("New X coordinate (0-100)"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("New Y coordinate (0-100)"),
		),
	)

	// Add get_empty_map tool
	getEmptyMapTool := mcp.NewTool("get_empty_map",
		mcp.WithDescription("Get an empty Wardley map with no components"),
		mcp.WithString("title",
			mcp.Description("Title for the map (default: 'New Wardley Map')"),
		),
		mcp.WithNumber("map_id",
			mcp.Description("ID for the map (default: 1)"),
		),
	)

	// Add add_components tool
	addComponentsTool := mcp.NewTool("add_components",
		mcp.WithDescription("Add or update multiple components in a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("components",
			mcp.Required(),
			mcp.Description("JSON array of components with name, x, y, and optional type"),
		),
	)

	// Add add_links tool
	addLinksTool := mcp.NewTool("add_links",
		mcp.WithDescription("Add multiple links between components in a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("links",
			mcp.Required(),
			mcp.Description("JSON array of links with from, to, and optional type"),
		),
	)

	// Add set_stages tool
	setStagesTool := mcp.NewTool("set_stages",
		mcp.WithDescription("Set evolution stages for a Wardley map"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("stage1",
			mcp.Required(),
			mcp.Description("Label for stage 1 (Genesis/Concept)"),
		),
		mcp.WithString("stage2",
			mcp.Required(),
			mcp.Description("Label for stage 2 (Custom/Emerging)"),
		),
		mcp.WithString("stage3",
			mcp.Required(),
			mcp.Description("Label for stage 3 (Product/Converging)"),
		),
		mcp.WithString("stage4",
			mcp.Required(),
			mcp.Description("Label for stage 4 (Commodity/Accepted)"),
		),
	)

	// Add add_anchor tool
	addAnchorTool := mcp.NewTool("add_anchor",
		mcp.WithDescription("Add an anchor to a Wardley map and return the updated map with SVG representation"),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("anchor_name",
			mcp.Required(),
			mcp.Description("Name/label of the anchor to add"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("X coordinate (0-100)"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Y coordinate (0-100)"),
		),
	)

	// Add tool handlers
	s.AddTool(addComponentTool, addComponentHandler)
	s.AddTool(addComponentsTool, addComponentsHandler)
	s.AddTool(addLinksTool, addLinksHandler)
	s.AddTool(linkComponentsTool, linkComponentsHandler)
	s.AddTool(unlinkComponentsTool, unlinkComponentsHandler)
	s.AddTool(moveComponentTool, moveComponentHandler)
	s.AddTool(getEmptyMapTool, getEmptyMapHandler)
	s.AddTool(setStagesTool, setStagesHandler)
	s.AddTool(addAnchorTool, addAnchorHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func addComponentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	componentName, err := request.RequireString("component_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("component_name parameter required: %v", err)), nil
	}

	x, err := request.RequireInt("x")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("x parameter required: %v", err)), nil
	}

	y, err := request.RequireInt("y")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("y parameter required: %v", err)), nil
	}

	componentType := request.GetString("component_type", "regular")

	// Validate coordinates
	if x < 0 || x > 100 || y < 0 || y > 100 {
		return mcp.NewToolResultError("coordinates must be between 0 and 100"), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
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

	// Create new component
	comp := wardley.NewComponent(nextID)
	comp.Label = componentName
	comp.Placement = image.Pt(x, y)

	// Set component type
	switch componentType {
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
		return mcp.NewToolResultError(fmt.Sprintf("failed to add component: %v", err)), nil
	}

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func linkComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	fromComponent, err := request.RequireString("from_component")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("from_component parameter required: %v", err)), nil
	}

	toComponent, err := request.RequireString("to_component")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("to_component parameter required: %v", err)), nil
	}

	linkType := request.GetString("link_type", "regular")

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
	}

	// Find the components
	fromComp := findComponentByName(m, fromComponent)
	if fromComp == nil {
		return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found", fromComponent)), nil
	}

	toComp := findComponentByName(m, toComponent)
	if toComp == nil {
		return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found", toComponent)), nil
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
		return mcp.NewToolResultError(fmt.Sprintf("failed to add collaboration: %v", err)), nil
	}

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func unlinkComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	fromComponent, err := request.RequireString("from_component")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("from_component parameter required: %v", err)), nil
	}

	toComponent, err := request.RequireString("to_component")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("to_component parameter required: %v", err)), nil
	}

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
	}

	// Find the components
	fromComp := findComponentByName(m, fromComponent)
	if fromComp == nil {
		return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found", fromComponent)), nil
	}

	toComp := findComponentByName(m, toComponent)
	if toComp == nil {
		return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found", toComponent)), nil
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
		return mcp.NewToolResultError(fmt.Sprintf("no link found between '%s' and '%s'", fromComponent, toComponent)), nil
	}

	// Remove the edge
	m.RemoveEdge(edgeToRemove.From().ID(), edgeToRemove.To().ID())

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func moveComponentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	componentName, err := request.RequireString("component_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("component_name parameter required: %v", err)), nil
	}

	x, err := request.RequireInt("x")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("x parameter required: %v", err)), nil
	}

	y, err := request.RequireInt("y")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("y parameter required: %v", err)), nil
	}

	// Validate coordinates
	if x < 0 || x > 100 || y < 0 || y > 100 {
		return mcp.NewToolResultError("coordinates must be between 0 and 100"), nil
	}

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
	}

	// Find the component
	comp := findComponentByName(m, componentName)
	if comp == nil {
		return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found", componentName)), nil
	}

	// Update component position
	comp.Placement = image.Pt(x, y)

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func getEmptyMapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract optional parameters with defaults
	title := request.GetString("title", "New Wardley Map")
	mapID := request.GetInt("map_id", 1)

	// Create empty map
	m := wardleyToGo.NewMap(int64(mapID))
	m.Title = title

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func addComponentsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	componentsJSON, err := request.RequireString("components")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("components parameter required: %v", err)), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
		}
	}

	// Parse components array
	var inputComponents []InputComponent
	if err := json.Unmarshal([]byte(componentsJSON), &inputComponents); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse components JSON: %v", err)), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("coordinates for component '%s' must be between 0 and 100", inputComp.Name)), nil
		}

		// Check if component already exists
		existingComp := findComponentByName(m, inputComp.Name)

		if existingComp != nil {
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
				return mcp.NewToolResultError(fmt.Sprintf("failed to add component '%s': %v", comp.Label, err)), nil
			}

			nextID++
		}
	}

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func addLinksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	linksJSON, err := request.RequireString("links")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("links parameter required: %v", err)), nil
	}

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
	}

	// Parse links array
	var inputLinks []InputLink
	if err := json.Unmarshal([]byte(linksJSON), &inputLinks); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse links JSON: %v", err)), nil
	}

	// Process each link
	for _, inputLink := range inputLinks {
		// Find the components
		fromComp := findComponentByName(m, inputLink.From)
		if fromComp == nil {
			return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found for link", inputLink.From)), nil
		}

		toComp := findComponentByName(m, inputLink.To)
		if toComp == nil {
			return mcp.NewToolResultError(fmt.Sprintf("component '%s' not found for link", inputLink.To)), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("failed to add collaboration from %s to %s: %v", inputLink.From, inputLink.To, err)), nil
		}
	}

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func setStagesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	stage1, err := request.RequireString("stage1")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("stage1 parameter required: %v", err)), nil
	}

	stage2, err := request.RequireString("stage2")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("stage2 parameter required: %v", err)), nil
	}

	stage3, err := request.RequireString("stage3")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("stage3 parameter required: %v", err)), nil
	}

	stage4, err := request.RequireString("stage4")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("stage4 parameter required: %v", err)), nil
	}

	// Parse the map
	m, err := UnmarshalMap([]byte(mapJSON))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
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

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

func addAnchorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	mapJSON, err := request.RequireString("map_json")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("map_json parameter required: %v", err)), nil
	}

	anchorName, err := request.RequireString("anchor_name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("anchor_name parameter required: %v", err)), nil
	}

	x, err := request.RequireInt("x")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("x parameter required: %v", err)), nil
	}

	y, err := request.RequireInt("y")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("y parameter required: %v", err)), nil
	}

	// Validate coordinates
	if x < 0 || x > 100 || y < 0 || y > 100 {
		return mcp.NewToolResultError("coordinates must be between 0 and 100"), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse map JSON: %v", err)), nil
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
		return mcp.NewToolResultError(fmt.Sprintf("failed to add anchor: %v", err)), nil
	}

	// Generate SVG with embedded JSON
	svgContent, err := GenerateSVG(m)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate SVG: %v", err)), nil
	}

	return mcp.NewToolResultText(svgContent), nil
}

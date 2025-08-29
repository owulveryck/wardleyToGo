package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgencoding "github.com/owulveryck/wardleyToGo/encoding/svg"
)

// Global storage for map stages (in production, this could be a database)
var mapStages = make(map[int64][]JSONEvolution)

// URI generation environment variables
var (
	uriScheme = getEnvWithDefault("WARDLEY_URI_SCHEME", "http")
	uriHost   = getEnvWithDefault("WARDLEY_URI_HOST", "localhost")
	uriPort   = getEnvWithDefault("WARDLEY_URI_PORT", "8585")
)

// getEnvWithDefault returns the environment variable value or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDefaultStages returns the default evolution stages
func getDefaultStages() []JSONEvolution {
	return []JSONEvolution{
		{Position: 0, Label: "Genesis / Concept"},
		{Position: 0.174, Label: "Custom / Emerging"},
		{Position: 0.4, Label: "Product / Converging"},
		{Position: 0.7, Label: "Commodity / Accepted"},
	}
}

// Helper function to find component name by ID
func findComponentNameByID(m *wardleyToGo.Map, id int64) string {
	nodes := m.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		// Check for regular components
		if comp, ok := node.(*wardley.Component); ok && comp.ID() == id {
			return comp.Label
		}
		// Check for anchors
		if anchor, ok := node.(*wardley.Anchor); ok && anchor.ID() == id {
			return anchor.Label
		}
	}
	return ""
}

// Helper function to find component by name (includes both components and anchors)
func findComponentByName(m *wardleyToGo.Map, name string) wardleyToGo.Component {
	nodes := m.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		// Check for regular components
		if comp, ok := node.(*wardley.Component); ok && comp.Label == name {
			return comp
		}
		// Check for anchors
		if anchor, ok := node.(*wardley.Anchor); ok && anchor.Label == name {
			return anchor
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

// encodeMapToGzippedBase64 compresses the JSON representation of a map using gzip and encodes it in base64
func encodeMapToGzippedBase64(m *wardleyToGo.Map) (string, error) {
	// Generate JSON representation
	jsonData, err := MarshalMap(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	// Compress JSON data using gzip
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(jsonData); err != nil {
		return "", fmt.Errorf("failed to compress JSON data: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Encode compressed data to base64
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

// generateURI creates a URI with the gzipped and base64-encoded map data
func generateURI(m *wardleyToGo.Map) (string, error) {
	// Encode map data
	encodedData, err := encodeMapToGzippedBase64(m)
	if err != nil {
		return "", fmt.Errorf("failed to encode map data: %w", err)
	}

	// Construct URI
	uri := fmt.Sprintf("%s://%s:%s/map?wardley_map_json_base64=%s", uriScheme, uriHost, uriPort, encodedData)
	return uri, nil
}

// decodeMapFromGzippedBase64 decompresses and decodes a base64-encoded gzipped JSON map
func decodeMapFromGzippedBase64(encodedData string) (*wardleyToGo.Map, error) {
	// Decode base64
	compressedData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Decompress gzip
	gzReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gzReader); err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	// Parse JSON and convert to map
	m, err := UnmarshalMap(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map JSON: %w", err)
	}

	return m, nil
}

// extractBase64FromURI extracts the base64 encoded map data from a URI
func extractBase64FromURI(uri string) (string, error) {
	// Parse the URI to extract the base64 data
	queryParamPrefix := "wardley_map_json_base64="
	idx := strings.Index(uri, queryParamPrefix)
	if idx == -1 {
		return "", fmt.Errorf("URI does not contain wardley_map_json_base64 parameter")
	}

	// Extract everything after the parameter name
	encodedData := uri[idx+len(queryParamPrefix):]

	// Remove any additional query parameters that might follow
	if ampIdx := strings.Index(encodedData, "&"); ampIdx != -1 {
		encodedData = encodedData[:ampIdx]
	}

	if encodedData == "" {
		return "", fmt.Errorf("empty base64 data in URI")
	}

	return encodedData, nil
}

// generateOutput creates either SVG, JSON, or URI representation based on output format
func generateOutput(m *wardleyToGo.Map, format string) (string, error) {
	switch format {
	case "json":
		jsonData, err := MarshalMap(m)
		if err != nil {
			return "", fmt.Errorf("failed to marshal map to JSON: %w", err)
		}
		return string(jsonData), nil
	case "uri":
		return generateURI(m)
	case "svg":
		fallthrough
	default:
		return GenerateSVG(m)
	}
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
		encoder.Close()
		return "", fmt.Errorf("failed to encode SVG: %w", err)
	}

	encoder.Close()
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

// calculateValueChainDepths calculates the depth of each node in the value chain
func calculateValueChainDepths(m *wardleyToGo.Map) (map[int64]*ValueChainNode, int) {
	nodes := make(map[int64]*ValueChainNode)
	inDegree := make(map[int64]int)
	maxDepth := 0

	// Initialize all nodes and calculate in-degrees
	mapNodes := m.Nodes()
	for mapNodes.Next() {
		node := mapNodes.Node()
		nodeID := node.ID()
		vcNode := &ValueChainNode{
			Node:     node,
			Depth:    -1, // -1 means unvisited
			IsAnchor: false,
		}
		// Check if this is an anchor
		if _, ok := node.(*wardley.Anchor); ok {
			vcNode.IsAnchor = true
		}
		nodes[nodeID] = vcNode
		inDegree[nodeID] = 0
	}

	// Calculate in-degrees (count incoming edges)
	edges := m.Edges()
	for edges.Next() {
		edge := edges.Edge()
		toID := edge.To().ID()
		inDegree[toID]++
	}

	// Find all root nodes (nodes with no incoming edges)
	queue := make([]int64, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			nodes[nodeID].Depth = 0
			queue = append(queue, nodeID)
		}
	}

	// If no root nodes found (cycle or isolated components), treat anchors as roots
	if len(queue) == 0 {
		for nodeID, vcNode := range nodes {
			if vcNode.IsAnchor {
				vcNode.Depth = 0
				queue = append(queue, nodeID)
			}
		}
	}

	// If still no roots, treat all nodes as depth 0
	if len(queue) == 0 {
		for _, vcNode := range nodes {
			vcNode.Depth = 0
		}
		return nodes, 0
	}

	// BFS to calculate depths (topological sort)
	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]
		currentDepth := nodes[currentID].Depth

		// Visit all nodes that depend on this one (outgoing edges)
		outgoing := m.From(currentID)
		for outgoing.Next() {
			neighborID := outgoing.Node().ID()
			neighbor := nodes[neighborID]
			newDepth := currentDepth + 1

			// Update depth if this is a longer path (we want the maximum depth for each node)
			if neighbor.Depth < newDepth {
				neighbor.Depth = newDepth
				queue = append(queue, neighborID)
				if newDepth > maxDepth {
					maxDepth = newDepth
				}
			}
		}
	}

	// Handle any remaining unvisited nodes
	for _, vcNode := range nodes {
		if vcNode.Depth == -1 {
			vcNode.Depth = 0
		}
	}

	return nodes, maxDepth
}

// positionComponentsInValueChain repositions components based on value chain analysis
func positionComponentsInValueChain(m *wardleyToGo.Map) {
	vcNodes, maxDepth := calculateValueChainDepths(m)

	// If maxDepth is 0, all components are at the same level
	if maxDepth == 0 {
		distributeComponentsEvenly(m, vcNodes)
		return
	}

	// Define N+1 zones where N is the maximum depth
	numZones := maxDepth + 1
	zones := make([][]int64, numZones)
	zoneHeight := 100 / numZones

	// Assign components to zones based on their depth
	for nodeID, vcNode := range vcNodes {
		// Direct mapping: depth 0 -> zone 0, depth 1 -> zone 1, etc.
		zoneIndex := vcNode.Depth
		if zoneIndex >= numZones {
			zoneIndex = numZones - 1
		}
		zones[zoneIndex] = append(zones[zoneIndex], nodeID)
	}

	// Position components within each zone
	for zoneIndex, nodeIDs := range zones {
		if len(nodeIDs) == 0 {
			continue
		}

		// Calculate Y position for this zone
		zoneStartY := zoneIndex * zoneHeight
		zoneEndY := (zoneIndex + 1) * zoneHeight

		// Distribute components within the zone
		distributeComponentsInZone(m, nodeIDs, zoneStartY, zoneEndY)
	}
}

// distributeComponentsEvenly distributes components evenly when no depth hierarchy exists
func distributeComponentsEvenly(m *wardleyToGo.Map, vcNodes map[int64]*ValueChainNode) {
	nodeIDs := make([]int64, 0, len(vcNodes))
	for nodeID := range vcNodes {
		nodeIDs = append(nodeIDs, nodeID)
	}
	distributeComponentsInZone(m, nodeIDs, 5, 95) // Use almost full height
}

// distributeComponentsInZone distributes components within a specific Y range
func distributeComponentsInZone(m *wardleyToGo.Map, nodeIDs []int64, startY, endY int) {
	if len(nodeIDs) == 0 {
		return
	}

	// Sort components by their current X coordinate to maintain horizontal order
	type componentPos struct {
		ID int64
		X  int
	}

	components := make([]componentPos, 0, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		node := m.Node(nodeID)
		if comp, ok := node.(*wardley.Component); ok {
			components = append(components, componentPos{
				ID: nodeID,
				X:  comp.GetPosition().X,
			})
		} else if anchor, ok := node.(*wardley.Anchor); ok {
			components = append(components, componentPos{
				ID: nodeID,
				X:  anchor.GetPosition().X,
			})
		}
	}

	// Sort by X coordinate
	for i := 0; i < len(components)-1; i++ {
		for j := i + 1; j < len(components); j++ {
			if components[i].X > components[j].X {
				components[i], components[j] = components[j], components[i]
			}
		}
	}

	// Calculate Y positions to avoid overlapping with padding
	const zonePadding = 5 // 5% padding at top and bottom of each zone
	zoneHeight := endY - startY
	paddedStartY := startY + zonePadding
	paddedEndY := endY - zonePadding
	paddedHeight := paddedEndY - paddedStartY

	if len(components) == 1 {
		// Single component goes in the middle of the padded zone
		targetY := paddedStartY + paddedHeight/2
		updateComponentY(m, components[0].ID, targetY)
	} else {
		// Multiple components distributed evenly within padded zone
		if paddedHeight <= 0 {
			// If zone is too small for padding, just center all components
			centerY := startY + zoneHeight/2
			for _, comp := range components {
				updateComponentY(m, comp.ID, centerY)
			}
		} else {
			yStep := paddedHeight / (len(components) - 1)
			for i, comp := range components {
				targetY := paddedStartY + i*yStep
				updateComponentY(m, comp.ID, targetY)
			}
		}
	}
}

// updateComponentY updates the Y coordinate of a component or anchor
func updateComponentY(m *wardleyToGo.Map, nodeID int64, newY int) {
	node := m.Node(nodeID)
	if comp, ok := node.(*wardley.Component); ok {
		comp.Placement = image.Pt(comp.Placement.X, newY)
	} else if anchor, ok := node.(*wardley.Anchor); ok {
		anchor.Placement = image.Pt(anchor.Placement.X, newY)
	}
}

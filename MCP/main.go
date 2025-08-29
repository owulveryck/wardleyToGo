package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"
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

// InputElement represents an element (component or anchor) in the input for add_elements
type InputElement struct {
	Name        string `json:"name"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	ElementType string `json:"element_type,omitempty"` // "component" or "anchor", default: "component"
	Type        string `json:"type,omitempty"`         // for components: "regular", "build", "buy", "outsource", "dataproduct"
}

// InputLink represents a link in the input for add_links
type InputLink struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type,omitempty"`
}

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

// handleMapRequest handles HTTP requests for the /map endpoint
func handleMapRequest(w http.ResponseWriter, r *http.Request) {
	// Get the base64 encoded map data from query parameter
	encodedData := r.URL.Query().Get("wardley_map_json_base64")
	if encodedData == "" {
		http.Error(w, "Missing wardley_map_json_base64 query parameter", http.StatusBadRequest)
		return
	}

	// Get the output format (default to SVG)
	outputFormat := r.URL.Query().Get("output")
	if outputFormat == "" {
		outputFormat = "svg"
	}

	// Validate output format
	if outputFormat != "svg" && outputFormat != "json" {
		http.Error(w, "Invalid output format. Must be 'svg' or 'json'", http.StatusBadRequest)
		return
	}

	// Decode the map data
	m, err := decodeMapFromGzippedBase64(encodedData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode map data: %v", err), http.StatusBadRequest)
		return
	}

	// Generate output in requested format
	content, err := generateOutput(m, outputFormat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate output: %v", err), http.StatusInternalServerError)
		return
	}

	// Set appropriate content type
	switch outputFormat {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	}

	// Allow CORS for web applications
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request for CORS preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Write the content
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

// startWebServer starts the HTTP server for serving Wardley maps
func startWebServer() {
	port, err := strconv.Atoi(uriPort)
	if err != nil {
		log.Printf("Invalid port number for web server: %s, web server disabled", uriPort)
		return
	}

	addr := fmt.Sprintf("%s:%d", uriHost, port)

	// Set up routes
	http.HandleFunc("/map", handleMapRequest)

	// Add a simple health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"wardley-map-server"}`))
	})

	// Add a root endpoint that shows usage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		usage := `<!DOCTYPE html>
<html>
<head>
    <title>Wardley Map Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        code { background-color: #f4f4f4; padding: 2px 4px; border-radius: 3px; }
        pre { background-color: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Wardley Map Server</h1>
    <p>This server renders Wardley maps from base64-encoded, gzipped JSON data.</p>
    
    <h2>Usage</h2>
    <p>Send a GET request to:</p>
    <pre><code>/map?wardley_map_json_base64=&lt;encoded_data&gt;&amp;output=&lt;format&gt;</code></pre>
    
    <h3>Parameters:</h3>
    <ul>
        <li><strong>wardley_map_json_base64</strong> (required): Base64-encoded, gzipped JSON representation of the Wardley map</li>
        <li><strong>output</strong> (optional): Output format - "svg" (default) or "json"</li>
    </ul>
    
    <h3>Examples:</h3>
    <ul>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;</code> - Returns SVG</li>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;&amp;output=svg</code> - Returns SVG</li>
        <li><code>/map?wardley_map_json_base64=&lt;data&gt;&amp;output=json</code> - Returns JSON</li>
    </ul>
    
    <h3>Health Check:</h3>
    <p><a href="/health">/health</a> - Returns server status</p>
</body>
</html>`
		w.Write([]byte(usage))
	})

	log.Printf("Starting Wardley Map web server on %s", addr)
	log.Printf("Server endpoints:")
	log.Printf("  GET /map?wardley_map_json_base64=<data>&output=<format>")
	log.Printf("  GET /health")
	log.Printf("  GET / (usage information)")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("Web server failed to start: %v", err)
	}
}

// ValueChainNode represents a node in the value chain analysis
type ValueChainNode struct {
	Node     graph.Node
	Depth    int
	IsAnchor bool
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

func main() {
	// Parse command line flags
	disableWeb := flag.Bool("no-web", false, "Disable web server (MCP server only)")
	enablePrompts := flag.Bool("prompt", false, "Enable prompt capabilities in MCP server")
	flag.Parse()

	// Start the web server in a goroutine unless disabled
	if !*disableWeb {
		go startWebServer()
	}

	// Create a new MCP server
	/*
		🤖 AI WORKFLOW GUIDE 🤖

		This MCP server provides tools for creating and editing Wardley Maps with a complete workflow designed for AI agents.
		The server also runs a web server (localhost:8585) that serves maps from shareable URIs.

		📋 CORE WORKFLOW PRINCIPLES:
		1. ALWAYS use 'json' output when you plan to make more changes (intermediate operations)
		2. Use 'uri' output to create shareable links that display SVG maps on the web server
		3. Use 'svg' output only for final display when no more changes are needed
		4. Each step should return JSON that is understandable and can be passed to the next tool

		🔄 COMPLETE WORKFLOWS:

		A) NEW MAP CREATION:
		   1. create_map(output="json") → Get starting JSON
		   2. add_elements(map_json=..., output="json") → Add components and anchors
		   3. add_links(map_json=..., output="json") → Add dependencies
		   4. auto_layout(map_json=..., output="json") → Position elements (optional)
		   5. Final step: Use output="uri" to create shareable link

		B) EDITING EXISTING MAP:
		   1. decode_uri(uri="...") → Extract JSON from shareable URI
		   2. Modify using any tools with output="json"
		   3. Final step: Use output="uri" to create new shareable link

		C) QUICK DISPLAY:
		   - Any tool with output="svg" for immediate visualization
		   - Visit the URI in a browser to see the interactive map

		🎯 KEY TOOLS BY PURPOSE:
		- 🚀 CREATE: create_map
		- 🔧 ADD: add_elements, add_links
		- 📐 POSITION: move_elements, auto_layout
		- 🎨 CONFIGURE: configure_evolution
		- 🔄 CONVERT: decode_uri
		- 🗑️ REMOVE: remove_elements

		⚠️ IMPORTANT: The JSON format is the universal interchange format. All tools understand it.
		   URIs contain compressed JSON and can be decoded back to JSON for further editing.
	*/
	s := server.NewMCPServer(
		"Wardley Map Generator & Web Server 🗺️",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// 📐 POSITION: move_elements tool
	moveElementsTool := mcp.NewTool("move_elements",
		mcp.WithDescription("📐 POSITION: Reposition specific elements by name. Handles single or multiple moves. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("moves",
			mcp.Required(),
			mcp.Description("JSON array of move operations. Each move must have: name (string, element name), x (0-100, evolution: 0=genesis/left, 100=commodity/right), y (0-100, visibility: 0=visible/bottom, 100=invisible/top). Example: [{'name':'Customer','x':20,'y':15}, {'name':'Service','x':60,'y':50}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🚀 CREATE: create_map tool
	createMapTool := mcp.NewTool("create_map",
		mcp.WithDescription("🚀 CREATE: Create a new empty Wardley map. Starting point for all workflows. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("title",
			mcp.Description("Title for the map (default: 'New Wardley Map')"),
		),
		mcp.WithNumber("map_id",
			mcp.Description("Unique identifier for the map (default: 1). Use different IDs to create multiple independent maps."),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔧 ADD: add_elements tool (unified components and anchors)
	addElementsTool := mcp.NewTool("add_elements",
		mcp.WithDescription("🔧 ADD: Add or update components and anchors. Handles single or multiple elements. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("elements",
			mcp.Required(),
			mcp.Description("JSON array of elements. Each element must have: name (string), x (0-100, evolution: 0=genesis/left, 100=commodity/right), y (0-100, visibility: 0=visible/bottom, 100=invisible/top), element_type ('component' or 'anchor', default: 'component'), type (optional for components: 'regular', 'build', 'buy', 'outsource', 'dataproduct'). Example: [{'name':'Customer','x':15,'y':15,'element_type':'anchor'}, {'name':'Service','x':50,'y':50,'element_type':'component','type':'regular'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔧 ADD: add_links tool
	addLinksTool := mcp.NewTool("add_links",
		mcp.WithDescription("🔧 ADD: Add dependency relationships between elements. Handles single or multiple links. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("links",
			mcp.Required(),
			mcp.Description("JSON array of dependency links. Each link must have: from (string, source element name), to (string, target element name), type (optional: 'regular', 'evolved_component', 'evolved'). The 'from' element depends on the 'to' element. Example: [{'from':'Customer','to':'Service','type':'regular'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🎨 CONFIGURE: configure_evolution tool
	configureEvolutionTool := mcp.NewTool("configure_evolution",
		mcp.WithDescription("🎨 CONFIGURE: Customize the four evolution stage labels displayed on the X-axis. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("stage1",
			mcp.Required(),
			mcp.Description("Label for stage 1 (Genesis/Concept): Novel, experimental, uncertain"),
		),
		mcp.WithString("stage2",
			mcp.Required(),
			mcp.Description("Label for stage 2 (Custom/Emerging): Bespoke, tailored solutions"),
		),
		mcp.WithString("stage3",
			mcp.Required(),
			mcp.Description("Label for stage 3 (Product/Converging): Packaged, feature-complete products"),
		),
		mcp.WithString("stage4",
			mcp.Required(),
			mcp.Description("Label for stage 4 (Commodity/Accepted): Most evolved, standardized, utility-like"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 📐 POSITION: auto_layout tool
	autoLayoutTool := mcp.NewTool("auto_layout",
		mcp.WithDescription("📐 POSITION: Automatically arrange all elements based on value chain depth analysis. Calculates dependency paths and positions elements in horizontal layers. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// 🔄 CONVERT: decode_uri tool
	decodeUriTool := mcp.NewTool("decode_uri",
		mcp.WithDescription("🔄 CONVERT: Extract map JSON from a shareable URI for editing. Essential for modifying existing maps. Returns JSON that can be passed to other tools."),
		mcp.WithString("uri",
			mcp.Required(),
			mcp.Description("The complete URI containing the base64-encoded map data (e.g., 'http://localhost:8585/map?wardley_map_json_base64=...')"),
		),
	)

	// 🗑️ REMOVE: remove_elements tool
	removeElementsTool := mcp.NewTool("remove_elements",
		mcp.WithDescription("🗑️ REMOVE: Remove components, anchors, or links from the map. Handles single or multiple removals. Use 'json' output for building workflows, 'uri' for shareable links, 'svg' for final display."),
		mcp.WithString("map_json",
			mcp.Required(),
			mcp.Description("JSON representation of the current map"),
		),
		mcp.WithString("remove_type",
			mcp.Required(),
			mcp.Description("Type of removal: 'elements' to remove components/anchors, 'links' to remove connections"),
		),
		mcp.WithString("items",
			mcp.Required(),
			mcp.Description("JSON array of items to remove. For elements: [{'name':'ElementName'}]. For links: [{'from':'SourceName','to':'TargetName'}]. Example: [{'name':'Customer'}] or [{'from':'Customer','to':'Service'}]"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: 'json' for intermediate operations, 'svg' for final display (default: 'svg'), or 'uri' for shareable link. Always use 'json' when planning to call more tools afterward."),
		),
	)

	// Add tool handlers
	s.AddTool(createMapTool, createMapHandler)
	s.AddTool(addElementsTool, addElementsHandler)
	s.AddTool(addLinksTool, addLinksHandler)
	s.AddTool(moveElementsTool, moveElementsHandler)
	s.AddTool(autoLayoutTool, autoLayoutHandler)
	s.AddTool(configureEvolutionTool, configureEvolutionHandler)
	s.AddTool(decodeUriTool, decodeUriHandler)
	s.AddTool(removeElementsTool, removeElementsHandler)

	// Add workflow prompts only if -prompt flag is set
	if *enablePrompts {
		addWorkflowPrompts(s)
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
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

func decodeUriHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	uri, err := request.RequireString("uri")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get uri parameter", err), nil
	}

	// Extract base64 data from URI
	encodedData, err := extractBase64FromURI(uri)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to extract base64 data from URI", err), nil
	}

	// Decode the map data
	m, err := decodeMapFromGzippedBase64(encodedData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to decode map data from URI", err), nil
	}

	// Return JSON representation
	jsonData, err := MarshalMap(m)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal decoded map to JSON", err), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
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

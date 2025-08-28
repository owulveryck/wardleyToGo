package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestURIGeneration(t *testing.T) {
	// Create a simple test map
	m := wardleyToGo.NewMap(1)
	m.Title = "Test Map"

	// Add a simple component
	comp := wardley.NewComponent(1)
	comp.Label = "Test Component"
	comp.Placement.X = 50
	comp.Placement.Y = 50
	m.AddComponent(comp)

	// Test URI generation
	uri, err := generateURI(m)
	if err != nil {
		t.Fatalf("Failed to generate URI: %v", err)
	}

	// Verify URI format
	expectedPrefix := "http://localhost:8585/map?wardley_map_json_base64="
	if !strings.HasPrefix(uri, expectedPrefix) {
		t.Errorf("Expected URI to start with %s, got: %s", expectedPrefix, uri)
	}

	// Extract and verify the base64 encoded data
	queryParamPrefix := "wardley_map_json_base64="
	idx := strings.Index(uri, queryParamPrefix)
	if idx == -1 {
		t.Fatalf("Expected URI to contain query parameter %s", queryParamPrefix)
	}

	encodedData := uri[idx+len(queryParamPrefix):]

	// Decode base64
	compressedData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		t.Fatalf("Failed to decode base64 data: %v", err)
	}

	// Decompress gzip
	gzReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gzReader); err != nil {
		t.Fatalf("Failed to decompress data: %v", err)
	}

	// Parse JSON
	var jsonMap JSONMap
	if err := json.Unmarshal(buf.Bytes(), &jsonMap); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify the data
	if jsonMap.Title != "Test Map" {
		t.Errorf("Expected title 'Test Map', got: %s", jsonMap.Title)
	}

	if len(jsonMap.Components) != 1 {
		t.Errorf("Expected 1 component, got: %d", len(jsonMap.Components))
	}

	if jsonMap.Components[0].Name != "Test Component" {
		t.Errorf("Expected component name 'Test Component', got: %s", jsonMap.Components[0].Name)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test default values
	if uriScheme != "http" {
		t.Errorf("Expected default scheme 'http', got: %s", uriScheme)
	}
	if uriHost != "localhost" {
		t.Errorf("Expected default host 'localhost', got: %s", uriHost)
	}
	if uriPort != "8585" {
		t.Errorf("Expected default port '8585', got: %s", uriPort)
	}

	// Test environment variable override
	os.Setenv("WARDLEY_URI_SCHEME", "https")
	os.Setenv("WARDLEY_URI_HOST", "example.com")
	os.Setenv("WARDLEY_URI_PORT", "9090")

	// Test that our getEnvWithDefault function works
	if getEnvWithDefault("WARDLEY_URI_SCHEME", "http") != "https" {
		t.Error("Environment variable override not working for scheme")
	}
	if getEnvWithDefault("WARDLEY_URI_HOST", "localhost") != "example.com" {
		t.Error("Environment variable override not working for host")
	}
	if getEnvWithDefault("WARDLEY_URI_PORT", "8585") != "9090" {
		t.Error("Environment variable override not working for port")
	}

	// Clean up
	os.Unsetenv("WARDLEY_URI_SCHEME")
	os.Unsetenv("WARDLEY_URI_HOST")
	os.Unsetenv("WARDLEY_URI_PORT")
}

func TestOutputGeneration(t *testing.T) {
	// Create a simple test map
	m := wardleyToGo.NewMap(1)
	m.Title = "Test Map"

	// Test URI output
	output, err := generateOutput(m, "uri")
	if err != nil {
		t.Fatalf("Failed to generate URI output: %v", err)
	}

	if !strings.HasPrefix(output, "http://localhost:8585/map?wardley_map_json_base64=") {
		t.Errorf("URI output format incorrect: %s", output)
	}

	// Test JSON output still works
	jsonOutput, err := generateOutput(m, "json")
	if err != nil {
		t.Fatalf("Failed to generate JSON output: %v", err)
	}

	var jsonMap JSONMap
	if err := json.Unmarshal([]byte(jsonOutput), &jsonMap); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v", err)
	}

	// Test SVG output still works
	svgOutput, err := generateOutput(m, "svg")
	if err != nil {
		t.Fatalf("Failed to generate SVG output: %v", err)
	}

	if !strings.Contains(svgOutput, "<svg") {
		t.Error("SVG output doesn't contain SVG tag")
	}
}

func TestWebServerDecoding(t *testing.T) {
	// Create a simple test map
	m := wardleyToGo.NewMap(1)
	m.Title = "Test Web Server Map"

	// Add a component
	comp := wardley.NewComponent(1)
	comp.Label = "Web Component"
	comp.Placement.X = 30
	comp.Placement.Y = 70
	m.AddComponent(comp)

	// Generate base64 encoded data
	encodedData, err := encodeMapToGzippedBase64(m)
	if err != nil {
		t.Fatalf("Failed to encode map data: %v", err)
	}

	// Test decoding the data (simulating what the web server does)
	decodedMap, err := decodeMapFromGzippedBase64(encodedData)
	if err != nil {
		t.Fatalf("Failed to decode map data: %v", err)
	}

	// Verify the decoded map
	if decodedMap.Title != "Test Web Server Map" {
		t.Errorf("Expected title 'Test Web Server Map', got: %s", decodedMap.Title)
	}

	// Count components
	nodeCount := 0
	nodes := decodedMap.Nodes()
	for nodes.Next() {
		nodeCount++
	}

	if nodeCount != 1 {
		t.Errorf("Expected 1 component, got: %d", nodeCount)
	}

	// Test that we can generate outputs from the decoded map
	jsonOutput, err := generateOutput(decodedMap, "json")
	if err != nil {
		t.Fatalf("Failed to generate JSON from decoded map: %v", err)
	}

	var jsonMap JSONMap
	if err := json.Unmarshal([]byte(jsonOutput), &jsonMap); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v", err)
	}

	if jsonMap.Components[0].Name != "Web Component" {
		t.Errorf("Expected component name 'Web Component', got: %s", jsonMap.Components[0].Name)
	}

	// Test SVG generation from decoded map
	svgOutput, err := generateOutput(decodedMap, "svg")
	if err != nil {
		t.Fatalf("Failed to generate SVG from decoded map: %v", err)
	}

	if !strings.Contains(svgOutput, "<svg") {
		t.Error("SVG output doesn't contain SVG tag")
	}
}

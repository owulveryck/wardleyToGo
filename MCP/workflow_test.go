package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestCompleteWorkflow(t *testing.T) {
	t.Run("New Map Creation Workflow", func(t *testing.T) {
		// Step 1: Create empty map
		emptyMap := wardleyToGo.NewMap(1)
		emptyMap.Title = "Workflow Test Map"

		jsonOutput1, err := generateOutput(emptyMap, "json")
		if err != nil {
			t.Fatalf("Failed to generate JSON for empty map: %v", err)
		}

		// Verify we have valid JSON
		var jsonMap1 JSONMap
		if err := json.Unmarshal([]byte(jsonOutput1), &jsonMap1); err != nil {
			t.Fatalf("Empty map JSON is invalid: %v", err)
		}

		if jsonMap1.Title != "Workflow Test Map" {
			t.Errorf("Expected title 'Workflow Test Map', got: %s", jsonMap1.Title)
		}

		// Step 2: Add components
		m, err := UnmarshalMap([]byte(jsonOutput1))
		if err != nil {
			t.Fatalf("Failed to unmarshal map: %v", err)
		}

		// Add components
		comp1 := wardley.NewComponent(1)
		comp1.Label = "User"
		comp1.Placement.X = 10
		comp1.Placement.Y = 20
		comp1.Type = wardley.RegularComponent
		m.AddComponent(comp1)

		comp2 := wardley.NewComponent(2)
		comp2.Label = "Service"
		comp2.Placement.X = 50
		comp2.Placement.Y = 60
		comp2.Type = wardley.BuildComponent
		m.AddComponent(comp2)

		comp3 := wardley.NewComponent(3)
		comp3.Label = "Database"
		comp3.Placement.X = 80
		comp3.Placement.Y = 80
		comp3.Type = wardley.BuyComponent
		m.AddComponent(comp3)

		jsonOutput2, err := generateOutput(m, "json")
		if err != nil {
			t.Fatalf("Failed to generate JSON after adding components: %v", err)
		}

		// Step 3: Add links
		m, err = UnmarshalMap([]byte(jsonOutput2))
		if err != nil {
			t.Fatalf("Failed to unmarshal map after components: %v", err)
		}

		// Add collaborations
		userComp := findComponentByName(m, "User")
		serviceComp := findComponentByName(m, "Service")
		dbComp := findComponentByName(m, "Database")

		if userComp == nil || serviceComp == nil || dbComp == nil {
			t.Fatal("Could not find all components")
		}

		// User depends on Service
		collab1 := &wardley.Collaboration{
			F:    userComp,
			T:    serviceComp,
			Type: wardley.RegularEdge,
		}
		m.SetCollaboration(collab1)

		// Service depends on Database
		collab2 := &wardley.Collaboration{
			F:    serviceComp,
			T:    dbComp,
			Type: wardley.RegularEdge,
		}
		m.SetCollaboration(collab2)

		jsonOutput3, err := generateOutput(m, "json")
		if err != nil {
			t.Fatalf("Failed to generate JSON after adding links: %v", err)
		}

		// Step 4: Generate URI for sharing
		uriOutput, err := generateOutput(m, "uri")
		if err != nil {
			t.Fatalf("Failed to generate URI: %v", err)
		}

		if !strings.HasPrefix(uriOutput, "http://localhost:8585/map?wardley_map_json_base64=") {
			t.Errorf("URI format incorrect: %s", uriOutput)
		}

		// Step 5: Test editing workflow - decode URI back to JSON
		encodedData, err := extractBase64FromURI(uriOutput)
		if err != nil {
			t.Fatalf("Failed to extract base64 from URI: %v", err)
		}

		decodedMap, err := decodeMapFromGzippedBase64(encodedData)
		if err != nil {
			t.Fatalf("Failed to decode map from base64: %v", err)
		}

		decodedJSON, err := generateOutput(decodedMap, "json")
		if err != nil {
			t.Fatalf("Failed to generate JSON from decoded map: %v", err)
		}

		// Verify the decoded map matches the original
		var originalMap JSONMap
		var decodedMapData JSONMap

		if err := json.Unmarshal([]byte(jsonOutput3), &originalMap); err != nil {
			t.Fatalf("Failed to parse original JSON: %v", err)
		}

		if err := json.Unmarshal([]byte(decodedJSON), &decodedMapData); err != nil {
			t.Fatalf("Failed to parse decoded JSON: %v", err)
		}

		// Check that key data matches
		if originalMap.Title != decodedMapData.Title {
			t.Errorf("Titles don't match: original=%s, decoded=%s", originalMap.Title, decodedMapData.Title)
		}

		if len(originalMap.Components) != len(decodedMapData.Components) {
			t.Errorf("Component count doesn't match: original=%d, decoded=%d",
				len(originalMap.Components), len(decodedMapData.Components))
		}

		if len(originalMap.Collaborations) != len(decodedMapData.Collaborations) {
			t.Errorf("Collaboration count doesn't match: original=%d, decoded=%d",
				len(originalMap.Collaborations), len(decodedMapData.Collaborations))
		}

		// Verify component names exist
		componentNames := make(map[string]bool)
		for _, comp := range decodedMapData.Components {
			componentNames[comp.Name] = true
		}

		expectedComponents := []string{"User", "Service", "Database"}
		for _, expected := range expectedComponents {
			if !componentNames[expected] {
				t.Errorf("Expected component '%s' not found in decoded map", expected)
			}
		}

		// Step 6: Test SVG generation from decoded map
		svgOutput, err := generateOutput(decodedMap, "svg")
		if err != nil {
			t.Fatalf("Failed to generate SVG from decoded map: %v", err)
		}

		if !strings.Contains(svgOutput, "<svg") {
			t.Error("SVG output doesn't contain SVG tag")
		}

		// Check that component names appear in SVG
		for _, name := range expectedComponents {
			if !strings.Contains(svgOutput, name) {
				t.Errorf("Component '%s' not found in SVG output", name)
			}
		}
	})
}

func TestURIDecodeFunction(t *testing.T) {
	// Create a test map
	m := wardleyToGo.NewMap(42)
	m.Title = "URI Decode Test"

	comp := wardley.NewComponent(1)
	comp.Label = "Test Component"
	comp.Placement.X = 25
	comp.Placement.Y = 75
	m.AddComponent(comp)

	// Generate URI
	uri, err := generateURI(m)
	if err != nil {
		t.Fatalf("Failed to generate URI: %v", err)
	}

	// Test the extractBase64FromURI function
	extractedData, err := extractBase64FromURI(uri)
	if err != nil {
		t.Fatalf("Failed to extract base64 from URI: %v", err)
	}

	// Test decoding
	decodedMap, err := decodeMapFromGzippedBase64(extractedData)
	if err != nil {
		t.Fatalf("Failed to decode map: %v", err)
	}

	if decodedMap.Title != "URI Decode Test" {
		t.Errorf("Expected title 'URI Decode Test', got: %s", decodedMap.Title)
	}

	if decodedMap.ID() != 42 {
		t.Errorf("Expected map ID 42, got: %d", decodedMap.ID())
	}

	// Test that we can extract base64 from URIs with additional query parameters
	uriWithParams := uri + "&output=svg&other=param"
	extractedData2, err := extractBase64FromURI(uriWithParams)
	if err != nil {
		t.Fatalf("Failed to extract base64 from URI with params: %v", err)
	}

	if extractedData != extractedData2 {
		t.Error("Extracted data differs when URI has additional parameters")
	}

	// Test error cases
	_, err = extractBase64FromURI("http://example.com/map?other=param")
	if err == nil {
		t.Error("Expected error for URI without base64 parameter")
	}

	_, err = extractBase64FromURI("http://example.com/map?wardley_map_json_base64=")
	if err == nil {
		t.Error("Expected error for URI with empty base64 parameter")
	}
}

func TestWorkflowJSONValidation(t *testing.T) {
	// Test that all workflow steps produce valid JSON that can be parsed
	emptyMap := wardleyToGo.NewMap(1)
	emptyMap.Title = "JSON Validation Test"

	// Test empty map JSON
	jsonOutput, err := generateOutput(emptyMap, "json")
	if err != nil {
		t.Fatalf("Failed to generate JSON: %v", err)
	}

	var jsonMap JSONMap
	if err := json.Unmarshal([]byte(jsonOutput), &jsonMap); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Test that JSON can be round-tripped through UnmarshalMap
	recreatedMap, err := UnmarshalMap([]byte(jsonOutput))
	if err != nil {
		t.Fatalf("Failed to unmarshal generated JSON: %v", err)
	}

	if recreatedMap.Title != emptyMap.Title {
		t.Errorf("Round-trip failed: titles don't match")
	}

	if recreatedMap.ID() != emptyMap.ID() {
		t.Errorf("Round-trip failed: IDs don't match")
	}

	// Test JSON structure includes required fields
	if jsonMap.ID == 0 {
		t.Error("JSON missing ID field")
	}

	if jsonMap.Title == "" {
		t.Error("JSON missing Title field")
	}

	if jsonMap.Components == nil {
		t.Error("JSON missing Components field")
	}

	if jsonMap.Collaborations == nil {
		t.Error("JSON missing Collaborations field")
	}

	if jsonMap.Stages == nil {
		t.Error("JSON missing Stages field")
	}

	// Verify default stages are present
	if len(jsonMap.Stages) != 4 {
		t.Errorf("Expected 4 default stages, got %d", len(jsonMap.Stages))
	}
}

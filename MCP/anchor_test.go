package main

import (
	"encoding/json"
	"image"
	"testing"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func TestAnchorLinking(t *testing.T) {
	// Create a map with both a component and an anchor
	m := wardleyToGo.NewMap(1)
	m.Title = "Anchor Linking Test"

	// Add an anchor
	anchor := wardley.NewAnchor(1)
	anchor.Label = "Customer"
	anchor.Placement = image.Pt(15, 15)
	m.AddComponent(anchor)

	// Add a component
	comp := wardley.NewComponent(2)
	comp.Label = "Cup of Tea"
	comp.Placement = image.Pt(50, 50)
	m.AddComponent(comp)

	// Test that we can find both the anchor and component
	foundAnchor := findComponentByName(m, "Customer")
	if foundAnchor == nil {
		t.Fatal("Could not find anchor 'Customer'")
	}

	foundComponent := findComponentByName(m, "Cup of Tea")
	if foundComponent == nil {
		t.Fatal("Could not find component 'Cup of Tea'")
	}

	// Test that both implement wardleyToGo.Component interface
	if foundAnchor.GetPosition().X != 15 || foundAnchor.GetPosition().Y != 15 {
		t.Errorf("Anchor position incorrect: expected (15,15), got (%d,%d)",
			foundAnchor.GetPosition().X, foundAnchor.GetPosition().Y)
	}

	if foundComponent.GetPosition().X != 50 || foundComponent.GetPosition().Y != 50 {
		t.Errorf("Component position incorrect: expected (50,50), got (%d,%d)",
			foundComponent.GetPosition().X, foundComponent.GetPosition().Y)
	}

	// Test creating a collaboration from anchor to component
	collab := &wardley.Collaboration{
		F:    foundAnchor, // This should work now!
		T:    foundComponent,
		Type: wardley.RegularEdge,
	}

	// Add the collaboration to the map
	if err := m.SetCollaboration(collab); err != nil {
		t.Fatalf("Failed to create collaboration from anchor to component: %v", err)
	}

	// Verify the collaboration exists in the map
	edges := m.Edges()
	collaborationFound := false
	for edges.Next() {
		edge := edges.Edge()
		if edge.From().ID() == anchor.ID() && edge.To().ID() == comp.ID() {
			collaborationFound = true
			break
		}
	}

	if !collaborationFound {
		t.Error("Collaboration between anchor and component not found in map")
	}

	// Test JSON marshaling includes both anchor and component
	jsonData, err := MarshalMap(m)
	if err != nil {
		t.Fatalf("Failed to marshal map to JSON: %v", err)
	}

	var jsonMap JSONMap
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check that we have both anchors and components
	if len(jsonMap.Anchors) != 1 {
		t.Errorf("Expected 1 anchor, got %d", len(jsonMap.Anchors))
	}

	if len(jsonMap.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(jsonMap.Components))
	}

	if len(jsonMap.Collaborations) != 1 {
		t.Errorf("Expected 1 collaboration, got %d", len(jsonMap.Collaborations))
	}

	// Check collaboration details
	if jsonMap.Collaborations[0].From != "Customer" {
		t.Errorf("Expected collaboration from 'Customer', got '%s'", jsonMap.Collaborations[0].From)
	}

	if jsonMap.Collaborations[0].To != "Cup of Tea" {
		t.Errorf("Expected collaboration to 'Cup of Tea', got '%s'", jsonMap.Collaborations[0].To)
	}
}

func TestMoveAnchor(t *testing.T) {
	// Create a map with an anchor
	m := wardleyToGo.NewMap(1)
	m.Title = "Move Anchor Test"

	// Add an anchor
	anchor := wardley.NewAnchor(1)
	anchor.Label = "User"
	anchor.Placement = image.Pt(20, 30)
	m.AddComponent(anchor)

	// Test that findComponentByName can find and move the anchor
	element := findComponentByName(m, "User")
	if element == nil {
		t.Fatal("Could not find anchor 'User'")
	}

	// Test moving the anchor (simulating move_component tool)
	switch e := element.(type) {
	case *wardley.Anchor:
		e.Placement = image.Pt(40, 60)
	default:
		t.Fatal("Element is not an anchor")
	}

	// Verify the anchor was moved
	movedElement := findComponentByName(m, "User")
	if movedElement.GetPosition().X != 40 || movedElement.GetPosition().Y != 60 {
		t.Errorf("Anchor not moved correctly: expected (40,60), got (%d,%d)",
			movedElement.GetPosition().X, movedElement.GetPosition().Y)
	}
}

func TestAnchorAndComponentInSameMap(t *testing.T) {
	// Test the tea shop scenario: Customer anchor -> Cup of Tea component
	m := wardleyToGo.NewMap(1)
	m.Title = "Tea Shop"

	// Add Customer anchor
	customer := wardley.NewAnchor(1)
	customer.Label = "Customer"
	customer.Placement = image.Pt(15, 15)
	m.AddComponent(customer)

	// Add Cup of Tea component
	cupOfTea := wardley.NewComponent(2)
	cupOfTea.Label = "Cup of Tea"
	cupOfTea.Placement = image.Pt(80, 60)
	m.AddComponent(cupOfTea)

	// Create link from Customer to Cup of Tea
	customerElement := findComponentByName(m, "Customer")
	cupOfTeaElement := findComponentByName(m, "Cup of Tea")

	if customerElement == nil {
		t.Fatal("Customer anchor not found")
	}
	if cupOfTeaElement == nil {
		t.Fatal("Cup of Tea component not found")
	}

	// This should now work!
	collab := &wardley.Collaboration{
		F:    customerElement,
		T:    cupOfTeaElement,
		Type: wardley.RegularEdge,
	}

	if err := m.SetCollaboration(collab); err != nil {
		t.Fatalf("Failed to link Customer anchor to Cup of Tea component: %v", err)
	}

	// Generate SVG to ensure it renders correctly
	svgOutput, err := GenerateSVG(m)
	if err != nil {
		t.Fatalf("Failed to generate SVG: %v", err)
	}

	// Check that both anchor and component appear in SVG
	if !containsSubstring(svgOutput, "Customer") {
		t.Error("SVG doesn't contain 'Customer' anchor")
	}

	if !containsSubstring(svgOutput, "Cup of Tea") {
		t.Error("SVG doesn't contain 'Cup of Tea' component")
	}

	t.Logf("Successfully created link from Customer anchor to Cup of Tea component!")
}

// Helper function to check if string contains substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) != -1
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

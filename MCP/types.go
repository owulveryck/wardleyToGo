package main

import (
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

// ValueChainNode represents a node in the value chain analysis
type ValueChainNode struct {
	Node     graph.Node
	Depth    int
	IsAnchor bool
}

# Wardley Map MCP Server

This is a Model Context Protocol (MCP) server that provides tools for generating and manipulating Wardley Maps as SVG images. All tools return SVG content with the map's JSON data embedded as comments for easy parsing by LLMs.

## Features

- **get_empty_map**: Get an empty Wardley map with no components as SVG
- **add_component**: Add a single component to a Wardley map and return the updated map as SVG
- **add_components**: Add or update multiple components in a Wardley map and return the updated map as SVG
- **link_components**: Link two components with a collaboration/dependency and return the updated map as SVG
- **add_links**: Add multiple links between components in a Wardley map and return the updated map as SVG
- **unlink_components**: Remove a link between two components and return the updated map as SVG
- **move_component**: Move a component to new coordinates and return the updated map as SVG
- **set_stages**: Set custom evolution stage labels for a Wardley map and return the updated map as SVG
- **add_anchor**: Add an anchor (text label) to a Wardley map and return the updated map as SVG

## Installation

```bash
cd MCP
go build
```

## Usage

The server runs as an MCP server that communicates via stdio. It accepts JSON-RPC requests and returns responses.

### get_empty_map Tool

Creates an empty Wardley map with no components.

**Parameters:**
- `title` (optional): Title for the map (default: "New Wardley Map")
- `map_id` (optional): ID for the map (default: 1)

**Returns:**
SVG content with embedded JSON data as comment:
```xml
<svg ...>
<!-- WARDLEY_MAP_DATA: {"id":1,"title":"New Wardley Map","components":[],"collaborations":[]} -->
<!-- SVG content follows -->
</svg>
```

### add_component Tool

Adds a component to a Wardley map.

**Parameters:**
- `map_json` (required): JSON representation of the current map (use `"{}"` for empty map)
- `component_name` (required): Name of the component to add
- `x` (required): X coordinate (0-100)
- `y` (required): Y coordinate (0-100)  
- `component_type` (optional): Type of component - "regular", "build", "buy", "outsource", or "dataproduct" (default: "regular")

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with the new component.

### add_components Tool

Adds or updates multiple components in a Wardley map. If a component with the same name already exists, only its coordinates (and optionally type) are updated without returning an error.

**Parameters:**
- `map_json` (required): JSON representation of the current map (use `"{}"` for empty map)
- `components` (required): JSON array of components to add/update

**Component format:**
```json
[
  {
    "name": "Component Name",
    "x": 50,
    "y": 60,
    "type": "regular"
  }
]
```

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with all components added/updated.

### link_components Tool

Links two components with a collaboration/dependency.

**Parameters:**
- `map_json` (required): JSON representation of the current map
- `from_component` (required): Name of the source component
- `to_component` (required): Name of the target component
- `link_type` (optional): Type of link - "regular", "evolved_component", or "evolved" (default: "regular")

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with the new link.

### add_links Tool

Adds multiple links between components in a Wardley map. If a link between two components already exists, it is skipped without returning an error.

**Parameters:**
- `map_json` (required): JSON representation of the current map
- `links` (required): JSON array of links to add

**Link format:**
```json
[
  {
    "from": "Component A",
    "to": "Component B", 
    "type": "regular"
  }
]
```

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with all links added.

### unlink_components Tool

Removes a link between two components.

**Parameters:**
- `map_json` (required): JSON representation of the current map
- `from_component` (required): Name of the source component
- `to_component` (required): Name of the target component

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with the link removed.

### move_component Tool

Moves a component to new coordinates.

**Parameters:**
- `map_json` (required): JSON representation of the current map
- `component_name` (required): Name of the component to move
- `x` (required): New X coordinate (0-100)
- `y` (required): New Y coordinate (0-100)

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with the component moved to its new position.

### set_stages Tool

Sets custom evolution stage labels for a Wardley map. This tool allows you to customize the four evolution stages that appear on the map.

**Parameters:**
- `map_json` (required): JSON representation of the current map
- `stage1` (required): Label for stage 1 (Genesis/Concept position, default: "Genesis / Concept")
- `stage2` (required): Label for stage 2 (Custom/Emerging position, default: "Custom / Emerging")
- `stage3` (required): Label for stage 3 (Product/Converging position, default: "Product / Converging")
- `stage4` (required): Label for stage 4 (Commodity/Accepted position, default: "Commodity / Accepted")

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with custom evolution stage labels.

**Default Stages:**
- Stage 1 (Position 0): "Genesis / Concept"
- Stage 2 (Position 0.174): "Custom / Emerging"
- Stage 3 (Position 0.4): "Product / Converging"
- Stage 4 (Position 0.7): "Commodity / Accepted"

### add_anchor Tool

Adds an anchor (text label) to a Wardley map. Anchors are reference points that appear as plain text on the map without visual decoration.

**Parameters:**
- `map_json` (required): JSON representation of the current map (use `"{}"` for empty map)
- `anchor_name` (required): Name/label of the anchor to add
- `x` (required): X coordinate (0-100)
- `y` (required): Y coordinate (0-100)

**Returns:**
SVG content with embedded JSON data as comment containing the updated map with the new anchor.

**Usage Example:**
Anchors are commonly used to mark reference points like "Business" or "Public" on Wardley maps to indicate different perspectives or user types.

## Embedded JSON Format

All SVG outputs contain the map data as an embedded comment in this format:

```xml
<!-- WARDLEY_MAP_DATA: {"id":1,"title":"My Wardley Map","components":[...],"collaborations":[...]} -->
```

The JSON structure follows this format:

```json
{
  "id": 1,
  "title": "My Wardley Map", 
  "components": [
    {
      "id": 1,
      "name": "Customer",
      "x": 10,
      "y": 80,
      "type": "regular"
    },
    {
      "id": 2,
      "name": "Product",
      "x": 50,
      "y": 60,
      "type": "regular"
    }
  ],
  "collaborations": [
    {
      "from": "Customer",
      "to": "Product", 
      "type": "regular"
    }
  ],
  "anchors": [
    {
      "id": 3,
      "name": "Business",
      "x": 94,
      "y": 55
    }
  ],
  "stages": [
    {
      "position": 0,
      "label": "Genesis / Concept"
    },
    {
      "position": 0.174,
      "label": "Custom / Emerging"
    },
    {
      "position": 0.4,
      "label": "Product / Converging"
    },
    {
      "position": 0.7,
      "label": "Commodity / Accepted"
    }
  ]
}
```

**Extracting JSON from SVG:**
LLMs can extract the JSON data by finding the comment that starts with `<!-- WARDLEY_MAP_DATA:` and parsing the JSON content.

## Component Types

- `regular`: Standard component (default)
- `build`: Build component (custom built)
- `buy`: Buy component (off-the-shelf)
- `outsource`: Outsourced component
- `dataproduct`: Data product component

## Link Types

- `regular`: Standard collaboration/dependency (default)
- `evolved_component`: Evolution edge with arrow (red dashed line)
- `evolved`: Evolution edge (red solid line)

## Coordinate System

The map uses a 100x100 coordinate system where:
- (0,0) is top-left
- (100,100) is bottom-right
- X represents evolution (0 = genesis, 100 = commodity)
- Y represents visibility (0 = invisible, 100 = visible)
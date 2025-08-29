# Wardley Map MCP Server 🗺️

A comprehensive Model Context Protocol (MCP) server for creating, editing, and visualizing Wardley Maps. Features a complete workflow system with multiple output formats, integrated web server, and AI-guided prompts.

## 🚀 Features

### Core Capabilities
- **🗺️ Complete Wardley Map Creation & Editing**: Create maps from scratch or modify existing ones
- **🌐 Integrated Web Server**: Shareable URIs that render interactive SVG maps
- **📊 Multiple Output Formats**: JSON (for workflows), SVG (for display), URI (for sharing)
- **🤖 AI Workflow Prompts**: Built-in prompts that guide AI assistants through map creation
- **🔄 Value Chain Auto-Layout**: Intelligent positioning based on dependency analysis
- **📱 CORS-Enabled**: Web server supports cross-origin requests for web applications

### Tool Architecture
The server provides 7 unified tools designed for efficient AI workflows:

- **🚀 CREATE**: `create_map` - Start with empty maps
- **🔧 ADD**: `add_elements`, `add_links` - Build map content
- **📐 POSITION**: `move_elements`, `auto_layout` - Arrange components
- **🎨 CONFIGURE**: `configure_evolution` - Customize evolution stages
- **🔄 CONVERT**: `decode_uri` - Extract maps from shareable links
- **🗑️ REMOVE**: `remove_elements` - Clean up maps

## 📦 Installation

```bash
cd MCP
go build
```

## 🎯 Usage

### Command Line Options

```bash
# Standard mode (tools only, web server enabled)
./mcp-server

# Enable prompts for AI assistants
./mcp-server -prompt

# Disable web server (MCP only)
./mcp-server -no-web

# Both flags
./mcp-server -prompt -no-web
```

### Flags
- **`-prompt`**: Enable prompt capabilities for AI workflow guidance
- **`-no-web`**: Disable the integrated web server (port 8585)

## 🔄 Workflow System

### Output Formats
All tools support three output formats via the `output` parameter:

- **`json`** (default for workflows): Structured data for chaining operations
- **`svg`**: Rendered map for final display
- **`uri`**: Shareable link to interactive map on web server

### Recommended Workflow Pattern
```
1. create_map(output="json")
2. add_elements(map_json=..., output="json") 
3. add_links(map_json=..., output="json")
4. auto_layout(map_json=..., output="json")
5. [final tool](map_json=..., output="uri") → Share the link!
```

## 🌐 Web Server

When enabled (default), the server runs on `http://localhost:8585` with:

- **`/map`**: Render maps from base64-encoded data
- **`/health`**: Health check endpoint
- **`/`**: Usage documentation

**Environment Variables:**
- `WARDLEY_URI_SCHEME`: URL scheme (default: "http")
- `WARDLEY_URI_HOST`: Host (default: "localhost")  
- `WARDLEY_URI_PORT`: Port (default: "8585")

## 🤖 AI Prompts

When started with `-prompt`, the server exposes workflow prompts:

### `create_wardley_map`
**Purpose**: Guide AI through complete map creation from text descriptions
**Parameter**: `description` - Text containing all map information
**Workflow**: Empty map → Add elements → Add links → Auto-layout → Generate URI

### `edit_wardley_map`  
**Purpose**: Guide AI through editing existing maps
**Parameters**: 
- `uri` - Existing map URI to edit
- `changes` - Description of modifications needed
**Workflow**: Decode URI → Apply changes → Auto-layout → Generate new URI

## 🛠️ Tools Reference

### 🚀 CREATE: create_map
Creates a new empty Wardley map.

**Parameters:**
- `title` (optional): Map title (default: "New Wardley Map")
- `map_id` (optional): Unique identifier (default: 1)
- `output` (optional): Format - "json", "svg", or "uri" (default: "svg")

### 🔧 ADD: add_elements  
Add or update components and anchors in bulk.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `elements` (required): JSON array of elements to add/update
- `output` (optional): Output format (default: "svg")

**Element Format:**
```json
[
  {
    "name": "Customer",
    "x": 15,
    "y": 85,
    "element_type": "anchor"
  },
  {
    "name": "Service",
    "x": 50,
    "y": 50,
    "element_type": "component",
    "type": "regular"
  }
]
```

**Element Types:**
- `element_type`: "component" (default) or "anchor"
- `type` (for components): "regular", "build", "buy", "outsource", "dataproduct"

### 🔧 ADD: add_links
Add dependency relationships between elements.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `links` (required): JSON array of dependency links
- `output` (optional): Output format (default: "svg")

**Link Format:**
```json
[
  {
    "from": "Customer",
    "to": "Service",
    "type": "regular"
  }
]
```

**Link Types:**
- `regular`: Standard dependency (default)
- `evolved_component`: Evolution with arrow (red dashed)
- `evolved`: Evolution (red solid)

### 📐 POSITION: move_elements
Reposition specific elements by name.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `moves` (required): JSON array of move operations
- `output` (optional): Output format (default: "svg")

**Move Format:**
```json
[
  {
    "name": "Customer",
    "x": 20,
    "y": 80
  }
]
```

### 📐 POSITION: auto_layout
Automatically arrange elements based on value chain analysis.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `output` (optional): Output format (default: "svg")

**Algorithm:**
1. Analyzes dependency relationships
2. Calculates depth from anchor points
3. Distributes elements in horizontal layers
4. Maintains evolution positioning (X-axis)

### 🎨 CONFIGURE: configure_evolution
Customize evolution stage labels on X-axis.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `stage1` (required): Genesis/Concept label
- `stage2` (required): Custom/Emerging label  
- `stage3` (required): Product/Converging label
- `stage4` (required): Commodity/Accepted label
- `output` (optional): Output format (default: "svg")

### 🔄 CONVERT: decode_uri
Extract map JSON from shareable URI for editing.

**Parameters:**
- `uri` (required): Complete URI with base64 map data

**Returns:** JSON representation for use with other tools

### 🗑️ REMOVE: remove_elements
Remove components, anchors, or links from maps.

**Parameters:**
- `map_json` (required): JSON representation of current map
- `remove_type` (required): "elements" or "links"
- `items` (required): JSON array of items to remove
- `output` (optional): Output format (default: "svg")

**Remove Formats:**
```json
// For elements
[{"name": "ComponentName"}]

// For links  
[{"from": "SourceName", "to": "TargetName"}]
```

## 📊 Data Formats

### JSON Map Structure
```json
{
  "id": 1,
  "title": "My Wardley Map",
  "components": [
    {
      "id": 1,
      "name": "Customer", 
      "x": 15,
      "y": 85,
      "type": "regular"
    }
  ],
  "collaborations": [
    {
      "from": "Customer",
      "to": "Service",
      "type": "regular"
    }
  ],
  "anchors": [
    {
      "id": 2,
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

### SVG Output
All SVG outputs include embedded JSON data as comments:
```xml
<svg ...>
<!-- WARDLEY_MAP_DATA: {"id":1,"title":"My Map",...} -->
<!-- SVG content follows -->
</svg>
```

### URI Format
Shareable links encode compressed map data:
```
http://localhost:8585/map?wardley_map_json_base64=<compressed_data>&output=svg
```

## 📏 Coordinate System

- **Range**: 0-100 for both X and Y coordinates
- **Origin**: (0,0) at top-left, (100,100) at bottom-right
- **X-axis (Evolution)**: 0=Genesis/Novel → 100=Commodity/Utility
- **Y-axis (Visibility)**: 0=Invisible/Infrastructure → 100=Visible/User-facing

### Evolution Guidelines
- **Genesis (X≈10)**: Novel, experimental, uncertain
- **Custom (X≈30)**: Bespoke, tailored solutions
- **Product (X≈60)**: Packaged, feature-complete
- **Commodity (X≈90)**: Standardized, utility-like

## 🔗 Integration Examples

### Claude Desktop MCP Configuration
```json
{
  "mcpServers": {
    "wardley-maps": {
      "command": "/path/to/mcp-server",
      "args": ["-prompt"]
    }
  }
}
```

### Web Application Integration
```javascript
// Fetch map from URI
const response = await fetch('http://localhost:8585/map?wardley_map_json_base64=...');
const svgContent = await response.text();

// Display in browser
document.getElementById('map-container').innerHTML = svgContent;
```

## 🚀 Quick Start Example

```bash
# Start server with prompts
./mcp-server -prompt

# Use the create_wardley_map prompt with your AI assistant:
# "Create a map showing: Customer needs a mobile app that uses an API connected to a database"

# The AI will guide you through:
# 1. Creating empty map
# 2. Adding elements (Customer, Mobile App, API, Database)  
# 3. Adding dependencies
# 4. Auto-positioning
# 5. Generating shareable URI

# Result: Interactive map at http://localhost:8585/map?wardley_map_json_base64=...
```

## 📋 Dependencies

- Go 1.19+
- github.com/mark3labs/mcp-go
- github.com/owulveryck/wardleyToGo (core library)
- gonum.org/v1/gonum/graph (graph operations)

## 🤝 Contributing

This MCP server is part of the WardleyToGo project. See the main project documentation for development guidelines and architecture details.
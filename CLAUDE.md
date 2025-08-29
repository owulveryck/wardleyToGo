# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WardleyToGo is a Go library and set of tools for creating and rendering Wardley Maps. It consists of:

- **Core library** (`github.com/owulveryck/wardleyToGo`): A Go SDK providing primitives to code Wardley Maps as graphs
- **WTG DSL**: A high-level domain-specific language for designing Wardley Maps  
- **OWM parser**: Parser for OnlineWardleyMaps (OWM) syntax
- **Multiple output formats**: SVG, PNG, DOT graph format
- **WebAssembly demo**: Interactive playground for WTG DSL

## Key Architecture

### Core Components
- **Map** (`map.go`): Central structure representing a Wardley Map as a directed graph using `gonum.org/v1/gonum/graph`
- **Component** (`defs.go`): Interface for map elements with positions on a 100x100 coordinate system
- **Collaboration** (`collaboration.go`): Interface for edges/relationships between components
- **Area**: Interface for rectangular regions on the map

### Parsers
- **WTG parser** (`parser/wtg/`): Parses the WTG domain-specific language
- **OWM parser** (`parser/owm/`): Parses OnlineWardleyMaps format

### Encoders/Renderers
- **SVG encoder** (`encoding/svg/`): Renders maps to SVG format with embedded CSS/JS
- **DOT encoder** (`encoding/dot/`): Outputs Graphviz DOT format

### Component Types
- **Wardley components** (`components/wardley/`): Traditional Wardley Map components with evolution stages
- **Team Topologies** (`components/teamtopologies/`): Team collaboration patterns

## Common Development Commands

### Building and Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.txt ./...

# Run a single test
go test -run TestSpecificFunction ./path/to/package

# Build a specific example tool
cd examples/wtg2svg && go build

# Format imports (run after modifying Go files)
goimports -w .

# Clean module cache and dependencies
go clean -modcache
go mod tidy
```

### Working with Examples
```bash
# Convert WTG to SVG
cd examples/wtg2svg
cat ../sample.wtg | go run main.go > output.svg

# Convert OWM to SVG  
cd examples/owm2svg
cat teashop.owm | go run main.go > output.svg

# Build WebAssembly demo
cd examples/webassembly
make all
```

### WebAssembly Development
```bash
cd examples/webassembly

# Build WASM module and copy assets
make all

# Clean build artifacts
make clean

# Install to docs directory
make install
```

## Development Notes

### Coordinate System
- All components use a 100x100 coordinate system where (0,0) is top-left
- Components must be within bounds or `AddComponent` returns an error
- Use `image.Point` for positions and `image.Rectangle` for areas

### Graph Structure
- Maps extend `gonum.org/v1/gonum/graph/simple.DirectedGraph`
- Components are nodes implementing the `Component` interface
- Collaborations are edges implementing the `Collaboration` interface
- Maps can contain other maps (submapping) since Map implements Component

### Rendering Pipeline
- Components and collaborations that implement `draw.Drawer` are automatically rendered
- SVG output includes embedded CSS and JavaScript for interactivity
- Canvas can be customized by setting `Map.Canvas`

### Testing Strategy
- Unit tests use `*_test.go` convention
- Example tests demonstrate usage patterns
- Fuzz testing available in `parser/wtg/testdata/fuzz/`

## MCP Server

The project includes an MCP (Model Context Protocol) server in the `MCP/` directory that provides AI-integrated Wardley Map creation and editing capabilities.

### MCP Development
```bash
# Build MCP server
cd MCP && go build

# Run MCP server with prompts enabled
./mcp-server -prompt

# Run tests for MCP server
cd MCP && go test ./...
```

### MCP Architecture
- **Tools**: 7 unified tools for map creation, editing, and visualization
- **Prompts**: AI workflow guidance for complete map creation from text descriptions
- **Web Server**: Integrated server on port 8585 for sharing interactive maps
- **Output Formats**: JSON (workflows), SVG (display), URI (sharing)
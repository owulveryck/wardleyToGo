# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WardleyToGo is a Go library and set of tools for creating and rendering Wardley Maps. It consists of:

- **Core library** (`github.com/owulveryck/wardleyToGo`): A Go SDK providing primitives to code Wardley Maps as graphs
- **WTG2 DSL**: A high-level domain-specific language for designing Wardley Maps (grammar in `wtg2/grammar.bnf`)
- **SVG renderer**: Renders maps to SVG format with embedded CSS/JS
- **WebAssembly demo**: Interactive playground for WTG2 DSL

## Key Architecture

### Core Components
- **Map** (`map.go`): Central structure representing a Wardley Map as a directed graph using `internal/graph`
- **Component** (`defs.go`): Interface for map elements with positions on a 100x100 coordinate system
- **Collaboration** (`collaboration.go`): Interface for edges/relationships between components
- **Area**: Interface for rectangular regions on the map

### Parser
- **WTG2 parser** (`parser/wtg2/`): Parses the WTG2 domain-specific language into an AST, then builds a `wardleyToGo.Map`

### Encoder
- **SVG encoder** (`encoding/svg/`): Renders maps to SVG format with embedded CSS/JS

### Component Types
- **Wardley components** (`components/wardley/`): Traditional Wardley Map components with evolution stages

### Experimental (`exp/`)
- **OWM parser** (`exp/parser_owm/`): Parser for OnlineWardleyMaps format (experimental)
- **Team Topologies** (`exp/teamtopologies/`): Team collaboration patterns (experimental)

## Common Development Commands

### Building and Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.txt ./...

# Run a single test
go test -run TestSpecificFunction ./path/to/package

# Build the wtg2svg CLI
go build ./cmd/wtg2svg/

# Convert WTG2 to SVG
cat wtg2/example.wtg2 | go run ./cmd/wtg2svg/ > output.svg

# Format imports (run after modifying Go files)
goimports -w .

# Clean module cache and dependencies
go clean -modcache
go mod tidy
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
- Maps embed `*internal/graph.DirectedGraph` (stdlib-only, no external deps)
- Components are nodes implementing the `Component` interface
- Collaborations are edges implementing the `Collaboration` interface
- Maps can contain other maps (submapping) since Map implements Component

### WTG2 Parser Pipeline
- **Lexer** (`parser/wtg2/lexer.go`): Line-oriented tokenization with block accumulation
- **Parser** (`parser/wtg2/parser.go`): Builds an AST (`Document`) from tokens
- **Builder** (`parser/wtg2/builder.go`): Converts AST to `wardleyToGo.Map` with `wardley.*` types
- **Evolution** (`parser/wtg2/evolution.go`): Converts roman numeral positions (I.0-IV.9) to 0-100 coords

### Layout Engine (`layout/`)
Standalone package for vertical placement of components on the value chain axis.
- **Interface**: `Layouter` with abstract `Graph`/`Node`/`Edge` types — decoupled from the parser AST
- **Algorithm**: Topological sort (Kahn) for rank assignment + force-directed repulsion for spacing
- **Options**: Configurable min spacing, repulsion strength, damping, iteration count
- The builder calls `layout.New(layout.DefaultOptions()).Layout(graph)` to compute Y positions

### Rendering Pipeline
- Components and collaborations that implement `draw.Drawer` are automatically rendered
- SVG output includes embedded CSS and JavaScript for interactivity
- Canvas can be customized by setting `Map.Canvas`

### Testing Strategy
- Unit tests use `*_test.go` convention
- Example tests demonstrate usage patterns
- Integration tests parse `parser/wtg2/testdata/example.wtg2` end-to-end

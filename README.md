# WardleyToGo

[![Go Reference](https://pkg.go.dev/badge/github.com/owulveryck/wardleyToGo.svg)](https://pkg.go.dev/github.com/owulveryck/wardleyToGo)
[![Try it online](https://img.shields.io/badge/Try%20it-online-blue?logo=github)](https://owulveryck.github.io/wardleyToGo)

A Go library and tools for creating and rendering [Wardley Maps](https://en.wikipedia.org/wiki/Wardley_map).

A map is a directed graph of components, each with a position on a 100x100 coordinate system representing visibility (Y) and evolution (X). The entry point of this API is the `Map` structure.

![wardleyToGo self-map](examples/self-map/this_repo.png)

*The project maps itself using WTG2 — [read the full strategic analysis](examples/self-map/) or see the [source](this_repo.wtg2).*

## WTG2: the Wardley Map DSL

WTG2 is a human-friendly domain-specific language for designing Wardley Maps. The full grammar is in [`wtg2/grammar.bnf`](wtg2/grammar.bnf).

### Quick start

```bash
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
cat input.wtg2 | wtg2svg > output.svg
```

Or without installing:

```bash
cat wtg2/example.wtg2 | go run ./cmd/wtg2svg/ > output.svg
```

### Example

```
title: Plateforme de Navigation -- Strategie 2026
stages: Genesis, Custom, Product, Commodity

anchor User
anchor Partner

Application : III.5
Routing Engine : II.7 !! >> III.5 {
  type: build
  color: #3498DB
  note: "Key differentiator"
}
Map Data : III.1 (buy)
Cloud Infra : IV.3 (buy)

pipeline Routing Engine {
  Classic Algo : III.5
  Predictive AI : II.3
}

User -> Application -> Routing Engine -> Map Data
Application -> Cloud Infra
Routing Engine -> Cloud Infra

note "Outsource candidate Q4" on Cloud Infra
warning "Single vendor dependency" on Map Data
signal accelerating on Predictive AI
```

### Language features

| Feature | Syntax |
|---------|--------|
| Metadata | `title:`, `date:`, `author:`, `scope:`, `question:` |
| Custom stages | `stages: Genesis, Custom, Product, Commodity` |
| Anchors | `anchor Name` |
| Components | `Name : III.5` or `component Name : III.5` |
| Types | `(build)`, `(buy)`, `(outsource)` |
| Block config | `Name : II.7 { type: build, color: #3498DB, note: "..." }` |
| Evolution | `II.7 >> III.5` (with optional inertia `!!`) |
| Pipelines | `pipeline Name { Member1 : III.5 }` |
| Edges | `A -> B -> C` (chains supported) |
| Annotated edges | `A -[label text]-> B` |
| Submaps | `submap Name : III.6` |
| Groups | `group Name { ... }` (visual only) |
| Annotations | `note "text" on Name`, `warning "text" on Name` |
| Signals | `signal accelerating\|stagnating\|declining on Name` |
| Comments | `// line` or `/* block */` |

### Position notation

Positions use roman numerals (I-IV) for evolution phases, with a decimal for sub-positioning:

- `I.0` = far left (genesis), `IV.9` = far right (commodity)
- Each phase spans 25% of the axis
- Without decimal (e.g. `III`), the center of the phase is used

## Using the library

### Creating a map programmatically

First, create a component type:

[embedmd]:# (example_draw_test.go /type dummyComponent.*/ /d.id }/)
```go
type dummyComponent struct {
	id       int64
	position image.Point
}

func (d *dummyComponent) GetPosition() image.Point { return d.position }

func (d *dummyComponent) String() string { return strconv.FormatInt(d.id, 10) }

func (d *dummyComponent) ID() int64 { return d.id }
```

Then a collaboration structure (an edge):

[embedmd]:# (example_draw_test.go /type dummyCollaboration.*/ /^}$/)
```go
type dummyCollaboration struct{ simple.Edge }

func (d *dummyCollaboration) GetType() wardleyToGo.EdgeType { return 0 }

func (d *dummyCollaboration) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	coordsF := utils.CalcCoords(d.F.(wardleyToGo.Component).GetPosition(), r)
	coordsT := utils.CalcCoords(d.T.(wardleyToGo.Component).GetPosition(), r)
	drawing.Line(dst, coordsF.X, coordsF.Y, coordsT.X, coordsT.Y, color.Gray{Y: 128}, [2]int{})
}
```

And finally create the map:

[embedmd]:# (example_draw_test.go /.*m \:= wardleyToGo.NewMap.*/ /.*c1, c3.*/)
```go
	m := wardleyToGo.NewMap(0)
	c0 := &dummyComponent{id: 0, position: image.Pt(25, 25)}
	c1 := &dummyComponent{id: 1, position: image.Pt(50, 50)}
	c2 := &dummyComponent{id: 2, position: image.Pt(50, 75)}
	c3 := &dummyComponent{id: 3, position: image.Pt(75, 75)}
	m.AddComponent(c0)
	m.AddComponent(c1)
	m.AddComponent(c2)
	m.AddComponent(c3)
	// c0 -> c1
	// c1 -> c2
	// c2 -> c3
	// c1 -> c3
	m.SetCollaboration(newCollaboration(c0, c1))
	m.SetCollaboration(newCollaboration(c1, c2))
	m.SetCollaboration(newCollaboration(c2, c3))
	m.SetCollaboration(newCollaboration(c1, c3))
```

## Project structure

```
wardleyToGo/
  map.go, defs.go, types.go     Core interfaces (Map, Component, Collaboration)
  components/wardley/            Wardley Map component implementations
  encoding/svg/                  SVG renderer with embedded CSS/JS
  parser/wtg2/                   WTG2 DSL parser (lexer, AST, layout, builder)
  cmd/wtg2svg/                   CLI tool: stdin WTG2 -> stdout SVG
  examples/webassembly/          Interactive browser playground
  wtg2/                          Language docs (grammar.bnf, algo.md, example)
  exp/                           Experimental: OWM parser, Team Topologies
  internal/                      Internal utilities (SVG primitives, drawing, coords)
```

## WebAssembly playground

Check the online demo at [https://owulveryck.github.io/wardleyToGo/](https://owulveryck.github.io/wardleyToGo/)

```bash
cd examples/webassembly
make all
```

Credits: inspired by [GraphvizOnline](https://github.com/dreampuf/GraphvizOnline). Uses:
- [ACE-editor](https://github.com/ajaxorg/ace/blob/master/LICENSE) BSD-2
- [svg-pan-zoom](https://github.com/bumbu/svg-pan-zoom) BSD-2

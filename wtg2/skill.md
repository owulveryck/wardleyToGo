# WTG2 — Wardley Map Language

You are generating Wardley Maps in the **WTG2** domain-specific language. Your output must be a valid `.wtg2` file that can be parsed and rendered to SVG.

---

## What is a Wardley Map?

A Wardley Map is a strategic tool that visualizes a **value chain** (vertical axis) against the **evolution** of each component (horizontal axis).

- **Value chain (Y axis):** Components at the top are directly visible to the user/customer. Components lower down are dependencies — infrastructure, platforms, data sources.
- **Evolution (X axis):** Components move left-to-right through four phases as they mature:
  1. **Genesis** (I) — Novel, poorly understood, high uncertainty
  2. **Custom-built** (II) — Understood but bespoke, requires expertise
  3. **Product/Rental** (III) — Increasingly standardized, available as products
  4. **Commodity/Utility** (IV) — Highly standardized, pay-per-use, invisible

Key principles:
- **Anchors** represent users or actors — they sit at the top of the value chain.
- **Components** are connected by dependency edges forming the value chain: `User -> Application -> Database -> Compute`.
- Position a component on the evolution axis based on its maturity, not where you *want* it to be.
- Common infrastructure (cloud, networking, power) belongs in phase IV. Novel R&D belongs in phase I.

---

## Document Structure

A WTG2 document follows this canonical order:

```
1. Metadata       (title, date, author, scope, question)
2. Configuration   (stages)
3. Nodes           (anchors, components, submaps, pipelines)
4. Value chain     (edges / dependencies)
5. Groups          (visual organization)
6. Annotations     (notes, warnings, signals)
```

All sections are optional. Comments can appear anywhere.

---

## Syntax Reference

### Comments

```
// Single-line comment
/* Block comment */
```

### Metadata

```
title: My Wardley Map
date: 2026-01-15
author: Strategy Team
scope: B2C mobile platform, European market
question: "Where should we invest to differentiate?"
```

All metadata fields are optional. The `question` value should be quoted.

### Stage Labels

Override the default evolution axis labels (default: `I`, `II`, `III`, `IV`):

```
stages: Genesis, Custom, Product, Commodity
```

Exactly four labels, comma-separated.

### Nodes

There are three node kinds:

| Keyword     | Purpose                                    |
|-------------|--------------------------------------------|
| `anchor`    | User or actor (always at top of map)       |
| `component` | Regular component (keyword is optional)    |
| `submap`    | Encapsulated sub-map shown as a component  |

#### Shorthand declaration (single line)

```
[kind] <name> : <evolution> [(<type>)] [@<visibility>]
```

The `component` keyword is optional — a bare name with a position is treated as a component.

Examples:

```
anchor User
Application : III.5
Database : III.8 (buy)
Infrastructure : IV.3 (buy) @0.2
submap Payment System : III.6
```

#### Block declaration (multi-line)

```
<name> : <evolution> {
  type: build
  color: #3498DB
  note: "Our key differentiator"
}
```

Block config fields:

| Field        | Values                              |
|--------------|-------------------------------------|
| `type`       | `build`, `buy`, `outsource`         |
| `color`      | `#RRGGBB`, `#RGB`, or CSS color name |
| `visibility` | `0.0` (bottom) to `1.0` (top)      |
| `note`       | Quoted text description             |

### Evolution Positioning

The horizontal position uses **roman numerals** with an optional decimal subdivision:

```
<roman>.<digit>
```

Where `<roman>` is `I`, `II`, `III`, or `IV`, and `<digit>` is `0`-`9`.

Each phase spans 25% of the axis. The decimal subdivides within the phase (0 = start, 9 = end). Without a decimal, the center of the phase is used.

**Position mapping to 0-100 coordinate:**

| Position | Coordinate | Meaning                     |
|----------|------------|-----------------------------|
| `I.0`    | 0          | Start of Genesis            |
| `I.5`    | 12         | Middle of Genesis           |
| `I.9`    | 22         | End of Genesis              |
| `II.0`   | 25         | Start of Custom             |
| `II.5`   | 37         | Middle of Custom            |
| `II.9`   | 47         | End of Custom               |
| `III.0`  | 50         | Start of Product            |
| `III.5`  | 62         | Middle of Product           |
| `III.9`  | 72         | End of Product              |
| `IV.0`   | 75         | Start of Commodity          |
| `IV.5`   | 87         | Middle of Commodity         |
| `IV.9`   | 97         | End of Commodity            |

Formula: `floor((base + digit/10 * 0.25) * 100)` where base is `I=0.00, II=0.25, III=0.50, IV=0.75`.

### Evolution Movement

Show that a component is evolving from one position to another:

```
Component : II.7 >> III.5
```

This renders an arrow from position II.7 to III.5 on the map.

### Inertia

Mark resistance to evolution with `!` (1-3 levels):

```
Component : II.7 ! >> III.5     // moderate inertia
Component : II.7 !! >> III.5    // strong inertia
Component : II.7 !!! >> III.5   // blocking inertia
```

Inertia appears between the current position and the `>>` movement operator.

### Visibility Override

By default, vertical positioning is computed automatically from the dependency graph. Override it with `@`:

```
Component : III.5 @0.9    // near top of map
Component : III.5 @0.1    // near bottom of map
```

Values range from `0.0` (bottom) to `1.0` (top).

### Edges (Value Chain)

Edges define dependencies between components.

```
A -> B                          // A depends on B
A <-> B                         // bidirectional relationship
A -[label text]-> B             // annotated dependency
A <-[label text]-> B            // annotated bidirectional
```

Edges can be chained:

```
User -> App -> API -> Database -> Cloud
```

This creates four edges: User->App, App->API, API->Database, Database->Cloud.

#### Pipeline member references

Target a specific member within a pipeline:

```
Component -> Pipeline:Member
```

### Pipelines

A pipeline shows multiple implementations of a component at different evolution stages:

```
pipeline <component-name> {
  Implementation A : III.5
  Implementation B : II.3
  Implementation C : I.2
}
```

Rules:
- The pipeline name must match an already-declared component.
- Members are positioned on the evolution axis only; their vertical position is derived from the parent component.
- The pipeline's horizontal span covers from its leftmost to rightmost member.

### Groups

Visually group components (purely visual, no scoping):

```
group Team Name {
  Component A
  Component B
  Component C
}
```

Members must reference existing component names.

### Annotations

```
note "Description text" on Component Name
warning "Risk description" on Component Name
```

### Signals

Mark market dynamics on a component:

```
signal accelerating on Component Name    // moving rapidly toward commodity
signal stagnating on Component Name      // evolution has plateaued
signal declining on Component Name       // regression in relevance
```

---

## Identifier Rules

Identifiers (component names, group names, etc.):
- Start with a letter or digit
- May contain letters, digits, `.`, `-`, `'`, `_`, and spaces
- Spaces are allowed inside identifiers (e.g., `Application Mobile`)
- Cannot be a reserved keyword used alone

**Reserved keywords:** `anchor`, `component`, `submap`, `pipeline`, `group`, `note`, `warning`, `signal`, `title`, `date`, `author`, `scope`, `question`, `stages`, `evolution`, `type`, `color`, `visibility`, `build`, `buy`, `outsource`, `accelerating`, `stagnating`, `declining`, `on`

---

## Semantic Rules

1. Every node referenced in an edge, annotation, signal, or pipeline must be declared with a position somewhere in the document.
2. Pipeline names must match a declared component.
3. Pipelines cannot be nested.
4. Groups do not create namespaces — components remain global.
5. Anchors do not need an evolution position (they are placed at the top automatically).
6. The dependency graph must be acyclic (no circular dependencies).

---

## Complete Example

```wtg2
// Wardley Map — GPS Navigation Platform

title: Navigation Platform — 2026 Strategy
date: 2026-01-15
author: Product Strategy Cell
scope: B2C mobile app, European market
question: "Where to invest to differentiate against Google Maps?"

stages: Genesis, Custom, Product, Commodity

// Anchors
anchor Driver
anchor Local Authority

// Visible layer
Application : III.5
Displayed Route : III.2
Real-Time Traffic Alerts : II.3

// Core engine with evolution and inertia
Route Calculation Engine : II.7 !! >> III.5 {
  type: build
  color: #3498DB
  note: "Key differentiator — 12 FTEs, 1.2M/year"
}

// Pipeline: the engine exists in multiple forms
pipeline Route Calculation Engine {
  Classic Dijkstra : III.5
  Predictive AI : II.3
  Quantum Algo : I.2
}

Cartographic Data Model : III.1 (buy)
B2G Partner API : II.1

// Infrastructure
OSM Data : III.8 (buy)
Real-Time Sensor Feed : I.8 ! >> II.5 {
  type: build
  color: #E67E22
  note: "Partnership in progress with Waze/TomTom"
}
Cloud Infrastructure : IV.3 (buy)
CDN : IV.5 (buy)
Mobile Network : IV.7 (outsource)

submap Payment System : III.6

// Value chain
Driver -> Application -> Displayed Route -> Route Calculation Engine
Application -> Real-Time Traffic Alerts -> Real-Time Sensor Feed
Route Calculation Engine -> Cartographic Data Model -> OSM Data
Route Calculation Engine -> Cloud Infrastructure
Real-Time Sensor Feed -> Cloud Infrastructure

Local Authority -> B2G Partner API -[Open Data, annual license]-> Cartographic Data Model
Local Authority -> Real-Time Traffic Alerts

Application -> CDN -> Cloud Infrastructure
Cloud Infrastructure -> Mobile Network
Application -> Payment System

// Link to specific pipeline member
Real-Time Traffic Alerts -> Route Calculation Engine:Predictive AI

// Groups
group Core Navigation Team {
  Route Calculation Engine
  Predictive AI
  Cartographic Data Model
}

group Platform Team {
  Cloud Infrastructure
  CDN
  Payment System
}

group Data Team {
  Real-Time Sensor Feed
  OSM Data
  Quantum Algo
}

// Annotations
warning "SPOF — no fallback if unavailable" on Route Calculation Engine
warning "Vendor lock-in AWS, cost rising 30%/year" on Cloud Infrastructure
warning "Critical dependency on single supplier" on OSM Data

note "Candidate for outsourcing Q4 2026" on Payment System
note "Partnership signed with 12 cities" on B2G Partner API
note "R&D budget 400k, horizon 2028" on Quantum Algo

// Market signals
signal accelerating on Predictive AI
signal accelerating on Real-Time Sensor Feed
signal stagnating on Classic Dijkstra
signal declining on OSM Data
```

---

## Generation Guidelines

When generating a WTG2 map:

1. **Start with the user need.** Identify the anchor(s) — who is the user/customer? Place them as `anchor` declarations.

2. **Build the value chain top-down.** Ask: "What does the user need?" Then for each component: "What does *this* component need?" Continue until you reach infrastructure.

3. **Position evolution realistically:**
   - Phase I (Genesis): Research, experiments, novel tech nobody else has
   - Phase II (Custom): Built in-house, understood but bespoke
   - Phase III (Product): Available as products/services, increasingly standardized
   - Phase IV (Commodity): Utilities, pay-per-use, ubiquitous (cloud, electricity, internet)

4. **Use `(buy)` and `(outsource)`** for components you consume rather than build. Leave untyped or use `(build)` for in-house work.

5. **Mark evolution movement** (`>>`) only for components actively transitioning. Add inertia (`!`, `!!`, `!!!`) when organizational or market resistance slows the transition.

6. **Use pipelines** when a component exists in multiple forms at different evolution stages (e.g., legacy vs. modern implementations).

7. **Add annotations sparingly.** Use `warning` for risks and `note` for strategic observations. Use `signal` to mark market dynamics.

8. **Group by team or domain** to show organizational ownership.

9. **Keep identifiers readable.** Use natural language names with spaces, not camelCase or snake_case.

10. **Follow the canonical section order:** metadata, stages, nodes, edges, groups, annotations.

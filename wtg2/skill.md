---
name: wtg2
description: Generate Wardley Maps in the WTG2 domain-specific language. Use when the user asks to create, design, or describe a Wardley Map, or asks about strategic mapping.
argument-hint: [description of the map to generate]
---

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
- **Anchors** represent users or actors — they sit at the top of the value chain, rendered with a person icon. They have a position on the evolution axis to control horizontal placement.
- **Components** are connected by dependency edges forming the value chain: `User -> Application -> Database -> Compute`.
- Position a component on the evolution axis based on its maturity, not where you *want* it to be.
- Common infrastructure (cloud, networking, power) belongs in phase IV. Novel R&D belongs in phase I.

---

## Strategic Concepts for Richer Maps

### Gameplays (Strategic Maneuvers)

A gameplay is a deliberate action to modify your position on the map. Annotate gameplays on the component that is the *target* of the maneuver.

| Gameplay | Description | Typical context |
|----------|-------------|-----------------|
| `ILC` | Innovate-Leverage-Commoditize: provide infrastructure, observe what thrives, absorb | Platform with ecosystem |
| `open-source` | Commoditize a layer to capture value in an adjacent layer | Competitor with proprietary rent |
| `land-grab` | Sacrifice profitability for rapid market share | New market with strong network effects |
| `embrace-extend` | Adopt a standard, add proprietary extensions, close ecosystem | Standard you want to control |
| `tower-moat` | Erect barriers: patents, lock-in, closed protocols | Protecting an existing rent |
| `FUD` | Spread fear/uncertainty/doubt to slow competitor adoption | Competitor gaining traction |
| `strangler-fig` | Progressively replace a legacy system by commoditizing non-differentiating layers | Legacy system blocking evolution |
| `signal-distortion` | Mislead competitors about strategic intent | Competitive misdirection |

```wtg2
gameplay ILC on Platform API
gameplay open-source "Commoditize compute to capture AI layer" on Cloud Infra
gameplay strangler-fig on Legacy CRM
```

### Five Capitals (Asset Classification)

Components represent different types of organizational capital. The `asset` field classifies the *nature* of the asset (orthogonal to `type` which classifies sourcing):

| Asset type | Description | Example |
|------------|-------------|---------|
| `tech` | Technological capital: code, infrastructure, patents | A routing engine, a data pipeline |
| `financial` | Financial capital: revenue models, pricing power | A billing system, a licensing model |
| `human` | Human capital: expertise, skills, tacit knowledge | An ML engineering team, domain experts |
| `relational` | Relational capital: partnerships, brand, contracts | A partner API, a brand, a patent portfolio |
| `social` | Social/environmental capital: community, regulatory | Open-source community, regulatory compliance |

```wtg2
AI Team : II.3 {
    type: build
    asset: human
    note: "12 ML engineers, hard to replace"
}
```

### Qualified Inertia

Inertia is not just a severity level — it has a *nature*. The book identifies 5 forms:

| Kind | Meaning | Symptom |
|------|---------|---------|
| `tech` | Technology lock-in, infrastructure debt | "We've always used Java" |
| `financial` | Sunk costs, established revenue models | "We've invested 5M in this" |
| `human` | Skills gap, identity threat, expertise obsolescence | "Our team doesn't know cloud-native" |
| `relational` | Contractual obligations, partner dependencies | "We have a 3-year vendor contract" |
| `social` | Cultural resistance, regulatory inertia | "That's not how we do things here" |

```wtg2
Component : II.7 !!(tech,human) >> III.5    // tech and human inertia
Component : II.7 !(financial) >> III.5      // financial inertia only
Component : II.7 !!! >> III.5              // unqualified (backward-compatible)
```

### Climatic Patterns (Extended Signals)

Beyond `accelerating`/`stagnating`/`declining`, mark climatic forces that explain *why* components evolve:

| Signal | Meaning |
|--------|---------|
| `co-evolution` | Technology and practice evolving together (e.g., containers + DevOps) |
| `red-queen` | Must evolve constantly just to maintain position |
| `commoditization` | Gravitational pull toward utility/commodity |
| `network-effects` | Value increases with number of users/participants |
| `economies-of-scale` | Cost advantage from volume, favoring consolidation |

```wtg2
signal co-evolution on DevOps Practices
signal commoditization on Cloud Infrastructure
signal network-effects on Social Platform
```

### EVT/PST Team Alignment

The Explorer-Villager-Town-planner model aligns team types to evolution phases. Use `team:` in groups to declare this alignment:

| Team type | Evolution phase | Mindset |
|-----------|----------------|---------|
| `explorer` / `pioneer` | Genesis (I) | Discovery, intuition, high failure tolerance |
| `settler` / `villager` | Custom-Product (II-III) | Productization, standards, analysis |
| `town-planner` | Commodity (IV) | Industrialization, cost optimization, experience-driven |

A mismatch between team type and component evolution phase is a strategic signal worth highlighting.

```wtg2
group R&D Team {
    team: explorer
    Quantum Algo
}
```

### Carrying Cost

The book emphasizes that 70-80% of IT budgets go to maintenance. Use `cost:` to annotate financial context:

```wtg2
Legacy CRM : III.2 {
    type: buy
    cost: "850k/year, 80% maintenance"
}
```

---

## Document Structure

A WTG2 document follows this canonical order:

```
1. Metadata       (title, date, author, scope, question, doctrine)
2. Configuration   (stages, legend)
3. Nodes           (anchors, components, submaps, pipelines)
4. Value chain     (edges / dependencies)
5. Groups          (visual organization)
6. Annotations     (notes, warnings, signals, gameplays, focus)
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

```
doctrine: context
```

The `doctrine` field indicates organizational maturity phase: `hygiene`, `context`, `excellence`, or `evolution`.

### Stage Labels

Override the default evolution axis labels (default: `I`, `II`, `III`, `IV`):

```
stages: Genesis, Custom, Product, Commodity
```

Exactly four labels, comma-separated.

### Legend

Display an auto-generated legend panel to the right of the map. The legend shows only the element types actually present in the map (component types, edges, signals, gameplays, groups, etc.).

```
legend
```

The keyword is standalone — no configuration needed. The SVG viewBox is automatically widened to accommodate the legend without distorting the map.

### Nodes

There are three node kinds:

| Keyword     | Purpose                                    |
|-------------|--------------------------------------------|
| `anchor`    | User or actor — rendered with a person icon, always at the top of the map. Has a position on the evolution axis to control horizontal placement. |
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
anchor User : III.5
anchor User : II.3 >> III.5
Application : III.5
Database : III.8 (buy)
Infrastructure : IV.3 (buy) @0.2
submap Payment System : III.6
```

Anchors have a position on the evolution axis that controls their horizontal placement on the map. They are always rendered at the top of the map (vertical position is fixed) with a person icon. Anchors can also have movement (`>>`) to show a shift in user positioning.

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
| `asset`      | `tech`, `financial`, `human`, `relational`, `social` |
| `evolution`  | Evolution expression (e.g., `II.7 !! >> III.5`) |
| `color`      | `#RRGGBB` or `#RGB`                  |
| `visibility` | `0.0` (bottom) to `1.0` (top)      |
| `cost`       | Free-text cost annotation           |
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

Mark resistance to evolution with `!` (1-3 levels), optionally qualified by kind:

```
Component : II.7 ! >> III.5               // moderate inertia
Component : II.7 !! >> III.5              // strong inertia
Component : II.7 !!! >> III.5             // blocking inertia
Component : II.7 !!(tech,human) >> III.5  // qualified: tech + human inertia
Component : II.7 !(financial) >> III.5    // qualified: financial inertia
```

Inertia appears between the current position and the `>>` movement operator. The optional `(kind,...)` qualifier specifies the nature of the resistance (see "Qualified Inertia" above).

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

All four forms are supported. Annotated bidirectional edges combine `<->` with `[label]`.

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

Visually group components (purely visual, no scoping). Optionally specify team type for EVT/PST alignment:

```
group Team Name {
  Component A
  Component B
  Component C
}

group R&D Team {
  team: explorer
  color: #E74C3C
  Quantum Algo
  Experimental Cache
}
```

Members must reference existing component names.

Group directives:

| Directive | Values |
|-----------|--------|
| `color:`  | `#RRGGBB` or `#RGB`                  |
| `team:`   | `explorer`, `settler`, `town-planner`, `pioneer`, `villager` |

### Annotations

```
note "Description text" on Component Name
warning "Risk description" on Component Name
```

### Signals

Mark market dynamics and climatic patterns on a component:

```
signal accelerating on Component Name    // moving rapidly toward commodity
signal stagnating on Component Name      // evolution has plateaued
signal declining on Component Name       // regression in relevance
signal co-evolution on Component Name    // technology-practice mutual reinforcement
signal red-queen on Component Name       // must evolve to maintain position
signal commoditization on Component Name // gravitational pull toward utility
signal network-effects on Component Name // value grows with adoption
signal economies-of-scale on Component   // cost advantage from volume
```

### Gameplays

Annotate strategic maneuvers on a component:

```
gameplay ILC on Platform API
gameplay open-source "Commoditize to capture adjacent value" on Database Engine
gameplay strangler-fig on Legacy System
gameplay land-grab on New Market
```

Gameplay types: `ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`.

The quoted description is optional and provides strategic context.

### Focus

The `focus` keyword highlights a component and all its descendants in the value chain. Elements outside the focus are rendered with reduced opacity.

```
focus Recommendation Engine
```

Multiple focus declarations can be combined — their subtrees are merged:

```
focus Recommendation Engine
focus Real-Time Personalisation
```

Focus is useful for presenting a specific part of the map during a discussion or presentation.

---

## Identifier Rules

Identifiers (component names, group names, etc.):
- Start with a letter or digit
- May contain letters, digits, `.`, `-`, `'`, `_`, and spaces
- Spaces are allowed inside identifiers (e.g., `Application Mobile`)
- Cannot be a reserved keyword used alone

**Reserved keywords:** `anchor`, `component`, `submap`, `pipeline`, `group`, `note`, `warning`, `signal`, `gameplay`, `legend`, `focus`, `title`, `date`, `author`, `scope`, `question`, `stages`, `doctrine`, `evolution`, `type`, `asset`, `color`, `visibility`, `cost`, `team`, `build`, `buy`, `outsource`, `accelerating`, `stagnating`, `declining`, `co-evolution`, `red-queen`, `commoditization`, `network-effects`, `economies-of-scale`, `ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`, `explorer`, `settler`, `town-planner`, `pioneer`, `villager`, `tech`, `financial`, `human`, `relational`, `social`, `on`, `hygiene`, `context`, `excellence`

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
doctrine: context

stages: Genesis, Custom, Product, Commodity
legend

// Anchors
anchor Driver
anchor Local Authority

// Visible layer
Application : III.5
Displayed Route : III.2
Real-Time Traffic Alerts : II.3

// Core engine with qualified inertia and asset/cost
Route Calculation Engine : II.7 !!(tech,human) >> III.5 {
  type: build
  asset: tech
  color: #3498DB
  cost: "1.2M/year, 12 FTEs"
  note: "Key differentiator"
}

// Pipeline: the engine exists in multiple forms
pipeline Route Calculation Engine {
  Classic Dijkstra : III.5
  Predictive AI : II.3
  Quantum Algo : I.2
}

Cartographic Data Model : III.1 (buy)
B2G Partner API : II.1 {
  asset: relational
}

// Infrastructure
OSM Data : III.8 (buy)
Real-Time Sensor Feed : I.8 !(relational) >> II.5 {
  type: build
  asset: tech
  color: #E67E22
  cost: "300k/year"
  note: "Partnership in progress with Waze/TomTom"
}
Cloud Infrastructure : IV.3 (buy) {
  cost: "500k/year, rising 30%"
}
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

// Groups with team types
group Core Navigation Team {
  team: settler
  Route Calculation Engine
  Predictive AI
  Cartographic Data Model
}

group Platform Team {
  team: town-planner
  Cloud Infrastructure
  CDN
  Payment System
}

group Data Team {
  team: explorer
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

// Market signals and climatic patterns
signal accelerating on Predictive AI
signal co-evolution on Real-Time Sensor Feed
signal stagnating on Classic Dijkstra
signal commoditization on Cloud Infrastructure
signal declining on OSM Data
signal red-queen on Application

// Gameplays
gameplay strangler-fig "Replace Classic Dijkstra with Predictive AI" on Route Calculation Engine
gameplay open-source "Commoditize mapping data to reduce dependency" on OSM Data
gameplay ILC on B2G Partner API

// Focus
focus Route Calculation Engine
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

10. **Follow the canonical section order:** metadata, stages/legend, nodes, edges, groups, annotations.

11. **Annotate gameplays** when the map represents a competitive strategy. Ask: "What maneuver is being applied to which component?" Use `gameplay` to make strategic intent explicit.

12. **Classify asset types** for components where the strategic value is non-obvious. A database is `tech` capital; a brand is `relational` capital; a specialized team is `human` capital.

13. **Qualify inertia** beyond severity. When marking `!!`, ask: "Is this technical debt (`tech`), sunk cost (`financial`), skills gap (`human`), contractual lock-in (`relational`), or cultural resistance (`social`)?"

14. **Apply climatic patterns** to explain evolutionary forces. Components in phase III with many dependents should likely have `signal commoditization`. Use `signal co-evolution` when a practice and technology evolve together.

15. **Align teams to evolution phases.** When defining groups with `team:`, verify that explorer teams own Genesis-phase components and town-planner teams own Commodity-phase components. Highlight mismatches as strategic risks.

16. **Annotate cost** for components consuming significant budget. This enables run/change ratio analysis across the value chain. Use `cost:` in block config.

---

## Strategic Completeness Checklist

Before finalizing a map, verify:

1. Every anchor has at least one dependency chain to infrastructure
2. Components in Genesis (I) have movement (`>>`) or a signal
3. Inertia is qualified by kind when the nature of resistance is known
4. Groups carrying team types align with component evolution phases
5. At least one gameplay is identified if the map is strategic
6. Warnings exist for SPOF, vendor lock-in, or single-supplier risks
7. Cost annotations exist for high-budget components
8. Signals reflect observed market dynamics, not speculation

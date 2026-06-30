---
name: wtg2
description: "Generate Wardley Maps in the WTG2 domain-specific language and render them to SVG. ALWAYS use this skill when the user mentions Wardley Maps, strategic mapping, value chains, evolution axes, component evolution, or asks to visualize strategy — even if they don't say 'WTG2' explicitly. Also triggers on: 'wardley', 'strategy map', 'value chain map', 'strategic landscape'."
argument-hint: "[description of the map to generate]"
---

# WTG2 — Wardley Map Skill

You generate Wardley Maps in the **WTG2** domain-specific language, then render them to SVG. Your output must be a valid `.wtg2` file that can be parsed and rendered.

A Wardley Map visualizes a **value chain** (Y axis — user-visible at top, infrastructure at bottom) against **evolution** (X axis — four phases from left to right):

| Phase | Name | Character |
|-------|------|-----------|
| I | Genesis | Novel, uncertain, experiment |
| II | Custom-built | Understood but bespoke |
| III | Product/Rental | Standardized, available as products |
| IV | Commodity/Utility | Ubiquitous, pay-per-use |

**Anchors** (users/actors) sit at the top with a person icon. Components form the value chain below them, connected by dependency edges.

### Strategic Lens: Climat / Doctrine / Manoeuvre

Every strategic element on a map operates at one of three levels:

- **Climat** — External forces you cannot control: evolution toward commodity, network effects, co-evolution. Expressed as `signal` annotations. Read the climate; do not fight it.
- **Doctrine** — Universal principles you choose to follow: hygiene, context awareness, operational excellence. Expressed via the `doctrine:` metadata field.
- **Manoeuvre** — Deliberate, context-dependent actions to change your position on the map. Expressed as `gameplay` annotations. A manoeuvre repositions you on the map (changes the terrain); a *tactic* optimizes within the current position (does the same thing better). Only manoeuvres belong as `gameplay` — refactoring code is a tactic; replacing a legacy system via strangler-fig is a manoeuvre.

---

## Workflow

Follow this end-to-end process every time:

### 1. Analyze the request

Identify from the user's description:
- **Anchors** — who are the users/actors?
- **Components** — what capabilities, systems, or resources exist?
- **Dependencies** — what depends on what? Build the value chain top-down.
- **Evolution** — where does each component sit on the maturity axis?

### 2. Generate valid WTG2

Write a complete `.wtg2` file following the syntax reference below. Use the canonical section order: metadata → configuration → nodes → edges → groups → annotations.

### 3. Write the file

Save as `<descriptive-name>.wtg2` in the current directory (or wherever the user specifies). Use kebab-case filenames.

### 4. Render to SVG

Render with `wtg2svg`:

```bash
wtg2svg -auto < name.wtg2 > name.svg
```

If `wtg2svg` is not found, tell the user it can be installed from https://github.com/owulveryck/wardleyToGo with:
```
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
```

**Flags:**
- `-auto` — auto-compute viewBox dimensions (recommended default)
- `-static` — produce static SVG without CSS/JS interactivity
- `-legend` — force legend display even if not declared in source
- `-width N` / `-height N` — override viewBox dimensions (default: 1100×900)

### 5. Show the result

On macOS: `open name.svg`

### 6. Iterate

If the user wants changes, edit the `.wtg2` file and re-render.

---

## Bundled References

Read these files when you need deeper detail:

- **`references/syntax-reference.md`** (relative to this skill directory) — Full position mapping table, identifier rules, reserved keywords, semantic rules, block config fields. Read when you hit an edge-case syntax question.
- **`references/strategic-concepts.md`** (relative to this skill directory) — Gameplays (ILC, strangler-fig, etc.), Five Capitals, qualified inertia types, climatic patterns, EVT/PST team alignment. Read when the user asks for strategic depth, competitive analysis, or organizational mapping.
- **`wtg2/doc_en.md`** — Complete language documentation.
- **`wtg2/grammar.bnf`** — Formal EBNF grammar.

---

## Syntax Quick Reference

### Comments

```wtg2
// Single-line comment
/* Block comment — can span multiple lines */
```

### Document Structure

```
1. Metadata       (title, date, author, scope, question, doctrine)
2. Configuration   (stages, legend)
3. Nodes           (anchors, components, submaps, pipelines)
4. Value chain     (edges / dependencies)
5. Groups          (visual organization)
6. Annotations     (notes, warnings, signals, gameplays, focus)
```

### Metadata

```wtg2
title: My Wardley Map
date: 2026-01-15
author: Strategy Team
scope: B2C mobile platform, European market
question: "Where should we invest to differentiate?"
doctrine: context                              // hygiene | context | excellence | evolution
stages: Genesis, Custom, Product, Commodity    // exactly 4 labels
legend                                         // show auto-generated legend panel
```

Doctrine phases reflect organizational maturity — choose the one that best fits the organization being mapped:
- `hygiene` — Fixing basics: common language, user focus, eliminating NIH. Maps at this level often reveal semantic confusion and misplaced investments.
- `context` — Adapting method to problem type: Agile for Genesis, Lean for Product, Six Sigma for Commodity. Maps should show methodological awareness.
- `excellence` — Operational efficiency: automation, ownership ("you build it, you run it"), waste elimination. Maps focus on flow optimization.
- `evolution` — Self-reinvention: "No Core" principle, willingness to cannibalize yesterday's success. Maps at this level should challenge assumptions.

### Nodes

Three kinds: `anchor` (person icon, top of map), `component` (default), `submap`.

**Shorthand:**
```wtg2
anchor User                           // anchor, auto-positioned horizontally
anchor User : III.5                   // anchor with explicit horizontal position
Component Name : III.5                // component keyword is optional
Database : III.8 (buy)                // with sourcing type
Infrastructure : IV.3 (buy) @0.2     // with visibility override
submap Payment System : III.6         // encapsulated sub-map
```

**Block form:**
```wtg2
Route Engine : II.7 !!(tech,human) >> III.5 {
    type: build                       // build | buy | outsource
    asset: tech                       // tech | financial | human | relational | social
    color: #3498DB
    cost: "1.2M/year, 12 FTEs"
    note: "Key differentiator"
}
```

### Evolution Positioning

Format: `<roman>.<digit>` — roman numeral (I–IV) with decimal subdivision (0–9).

| Position | Meaning |
|----------|---------|
| `I.0`–`I.9` | Genesis (0%–22%) |
| `II.0`–`II.9` | Custom (25%–47%) |
| `III.0`–`III.9` | Product (50%–72%) |
| `IV.0`–`IV.9` | Commodity (75%–97%) |

**Phase characteristics for accurate positioning:**

| Property | I Genesis | II Custom | III Product | IV Commodity |
|----------|-----------|-----------|-------------|--------------|
| Market | Undefined | Forming | Growing, consolidating | Mature, stable |
| Failure | Tolerated | Disappointing | Not tolerated | Shocking |
| Management | Explore (Agile) | Expert craft (Lean Startup) | Industrialize (Lean) | Optimize (Six Sigma) |
| Profile | Explorer | Artisan | Product Manager | Industrialist |

**Movement** — component evolving: `Component : II.7 >> III.5`

**Inertia** — resistance to evolution (1–3 levels, optionally qualified):
```wtg2
Component : II.7 ! >> III.5                // moderate
Component : II.7 !!(tech,human) >> III.5   // strong, tech + human inertia
Component : II.7 !!! >> III.5              // blocking
```

Every component has a dual nature: *capability* (what it enables) and *carrying cost* (what it prevents from changing). Sourcing carries financial implications: `build` = CAPEX (sunk cost, high inertia), `buy`/`outsource` = OPEX (flexible, lower inertia).

**Visibility override** — manual Y position: `Component : III.5 @0.9` (0.0 = bottom, 1.0 = top)

### Edges (Value Chain)

```wtg2
A -> B                          // dependency
A <-> B                         // bidirectional
A -[label text]-> B             // annotated
A <-[label text]-> B            // annotated bidirectional
User -> App -> API -> DB        // chained (creates 3 edges)
Component -> Pipeline:Member    // target pipeline member
```

### Pipelines

Multiple implementations of a component at different evolution stages:

```wtg2
pipeline Route Engine {
    Classic Dijkstra : III.5
    Predictive AI : II.3
    Quantum Algo : I.2
}
```

The pipeline name must match a declared component. Members get their vertical position from the parent.

### Groups

Visual grouping with optional team alignment:

```wtg2
group Core Team {
    team: settler                // explorer | settler | town-planner | pioneer | villager
    color: #E74C3C
    Route Engine
    Cartographic Data
}
```

### Annotations

```wtg2
note "Description text" on Component Name
warning "Risk description" on Component Name

signal accelerating on Component         // market dynamics
signal stagnating on Component
signal declining on Component
signal co-evolution on Component          // climatic patterns
signal red-queen on Component
signal commoditization on Component
signal network-effects on Component
signal economies-of-scale on Component

gameplay ILC on Platform API              // strategic maneuvers
gameplay open-source "Rationale" on DB
gameplay strangler-fig on Legacy
gameplay due-diligence "Assess before acquisition" on Target Company
gameplay two-sided-market on Marketplace
// Also: land-grab, embrace-extend, tower-moat, FUD, signal-distortion

// focus is available but should NOT be used — see Generation Guidelines
```

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

```

---

## Generation Guidelines

1. **Start with the user need.** Identify anchors first — who is the user/customer?
2. **Build the value chain top-down.** "What does the user need?" → "What does *that* need?" → repeat until infrastructure.
3. **Position evolution realistically.** Phase I = research/novel. Phase IV = utilities (cloud, power, internet). Don't place things where you *want* them — place them where they *are*.
4. **Use `(buy)` and `(outsource)`** for consumed components. Leave untyped or `(build)` for in-house.
5. **Mark movement (`>>`)** only for components actively transitioning. Add inertia when resistance slows it.
6. **Use pipelines** when a component exists in multiple forms at different evolution stages.
7. **Add annotations sparingly.** `warning` for risks, `note` for observations, `signal` for market dynamics.
8. **Group by team or domain** to show ownership. Align team types to evolution phases.
9. **Keep names readable.** Natural language with spaces, not camelCase.
10. **Follow canonical section order:** metadata, stages/legend, nodes, edges, groups, annotations.
11. **Do not use `focus`.** The `focus` keyword dims everything except one subtree, which hides context and makes the map harder to read. Let the user decide what to highlight — never add `focus` to generated maps.
12. **Annotate gameplays** when the map is competitive strategy.
13. **Classify assets** for components where strategic value is non-obvious.
14. **Qualify inertia** — is it tech debt, sunk cost, skills gap, contractual, or cultural?
15. **Apply climatic patterns** — commoditization on phase III dependents, co-evolution on linked tech+practice.
16. **Verify team-evolution alignment** — explorers own Genesis, town-planners own Commodity. Mismatches are risks.
17. **Annotate cost** for high-budget components.
18. **Respect the Extremistan/Mediocristan divide.** Phases I–II operate in Extremistan (power laws, unpredictable ROI, winner-takes-all). Phases III–IV operate in Mediocristan (normal distributions, predictable metrics). Do not apply uniform KPIs across all phases — a startup metric does not apply to infrastructure, and an SLA does not apply to R&D.
19. **Flag doctrine violations.** Check for the anti-patterns below and annotate them with `warning` or `note`.

**Doctrine violations to flag:**

| Violation | Symptom on map |
|-----------|---------------|
| NIH | Component `build` in phase III–IV when market products exist |
| No differentiation | All components in III–IV, nothing in I–II |
| Dispersion | Many components in I without critical mass in any |
| Strategy theatre | No `question:`, no `gameplay`, no movement (`>>`) |
| Single method | Same approach applied uniformly regardless of evolution phase |

### Flow Anomalies

When tracing the value chain, flag these with `warning`:
- **One-way flow** — dependency where value only flows in one direction (data goes up, no value returns)
- **Bottleneck** — single component through which all value chains pass (structural SPOF)
- **Value leak** — value created internally is captured by an external `buy` component extracting rent

### Completeness Checklist

Before finalizing, verify:
- [ ] Every anchor has a dependency chain to infrastructure
- [ ] Genesis components have movement (`>>`) or a signal
- [ ] Inertia is qualified when the nature of resistance is known
- [ ] Groups with team types align with component evolution phases
- [ ] At least one gameplay if the map is strategic
- [ ] Warnings for SPOF, vendor lock-in, or single-supplier risks
- [ ] Cost annotations for high-budget components
- [ ] Signals reflect observed market dynamics, not speculation
- [ ] The `question:` metadata is defined — a map without a question is strategy theatre
- [ ] At least one component has movement (`>>`) — a static map is not strategic
- [ ] Components marked `build` in phase III–IV are justified (otherwise: NIH violation)
- [ ] Capital flows along edges are bidirectional — one-way flows deserve a `warning`
- [ ] The map contains at least one surprise — if nothing contradicts assumptions, you have not looked closely enough

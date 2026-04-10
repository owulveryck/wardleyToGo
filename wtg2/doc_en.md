# WTG2 -- Wardley Map Description Language

## Table of Contents

1. [Introduction](#introduction)
2. [Document Structure](#document-structure)
3. [Comments](#comments)
4. [Metadata](#metadata)
5. [Evolution Stage Labels](#evolution-stage-labels)
6. [Legend](#legend)
7. [Nodes](#nodes)
   - [Anchors](#anchors)
   - [Components](#components)
   - [Submaps](#submaps)
   - [Block Declaration](#block-declaration)
8. [Evolution Positioning](#evolution-positioning)
9. [Evolution Movement](#evolution-movement)
10. [Inertia](#inertia)
11. [Visibility](#visibility)
12. [Pipelines](#pipelines)
13. [Value Chain (Edges)](#value-chain-edges)
14. [Groups](#groups)
15. [Annotations](#annotations)
16. [Signals](#signals)
17. [Strategic Maneuvers (Gameplays)](#strategic-maneuvers-gameplays)
18. [Focus](#focus)
19. [Identifier Rules](#identifier-rules)
20. [Semantic Rules](#semantic-rules)
21. [Complete Example](#complete-example)

---

## Introduction

**WTG2** is a declarative language designed to describe [Wardley Maps](https://learnwardleymapping.com/). A Wardley Map is a strategic tool that represents a **value chain** (vertical axis) against the **evolution stage** of each component (horizontal axis).

### The Two Axes of a Wardley Map

**Vertical Axis -- Visibility (Value Chain):**
Components visible to the user sit at the top of the map. The further down you go, the closer you get to invisible infrastructure.

**Horizontal Axis -- Evolution:**
Components move from left to right as they mature, through four phases:

| Phase | Roman Numeral | Description |
|-------|--------------|-------------|
| Genesis | I | Novel, poorly understood, high uncertainty |
| Custom-Built | II | Understood but bespoke, requires expertise |
| Product | III | Increasingly standardized, available as products |
| Commodity | IV | Highly standardized, pay-per-use, invisible |

A `.wtg2` file is a plain text file that, once parsed, produces a map structure renderable to SVG.

---

## Document Structure

A WTG2 document is made up of **statements**, in any order. However, the recommended canonical order is:

```
1. Metadata        (title, date, author, scope, question, doctrine)
2. Configuration   (stages, legend)
3. Nodes           (anchor, component, submap, pipeline)
4. Value chain     (edges / dependencies)
5. Groups          (group)
6. Annotations     (note, warning, signal, gameplay)
```

All sections are optional. Comments and blank lines may appear anywhere.

---

## Comments

WTG2 supports two comment styles:

```
// Single-line comment

/* Multi-line
   comment */
```

Comments are ignored by the parser and serve to document the map.

---

## Metadata

Metadata fields describe the context of the map. Each field is optional.

```
title: Online Retail Platform -- Strategic Review 2026
date: 2026-04-09
author: Platform Strategy Team
scope: B2C e-commerce, EU market
question: "How do we reduce infrastructure costs while accelerating AI-driven personalisation?"
doctrine: context
```

| Field | Description |
|-------|-------------|
| `title` | Map title |
| `date` | Date in ISO format `YYYY-MM-DD` |
| `author` | Author or team |
| `scope` | Map scope / boundary |
| `question` | Strategic question the map addresses (quoted) |
| `doctrine` | Organizational maturity phase: `hygiene`, `context`, `excellence`, or `evolution` |

---

## Evolution Stage Labels

By default, phases are labeled with roman numerals I, II, III, IV. You can customize them:

```
stages: Genesis, Custom-Built, Product, Commodity
```

Exactly four labels, comma-separated. They appear at the bottom of the map as zone labels.

---

## Legend

The `legend` keyword activates a legend panel to the right of the map. Content is auto-detected from elements present in the map.

```
legend
```

No configuration needed. The SVG viewBox is automatically widened to accommodate the legend without distorting the map.

---

## Nodes

A node represents an element on the map. There are three kinds of nodes.

### Anchors

An **anchor** represents a user or actor. It is always placed at the top of the map and does not require an evolution position.

```
anchor Customer
anchor Merchant
```

An anchor may optionally have a position (for precise horizontal placement):

```
anchor Customer : III.5
```

### Components

A **component** is the fundamental element of a Wardley Map. The `component` keyword is optional -- a bare name followed by a position is automatically treated as a component.

```
// These two forms are equivalent:
component Web Storefront : III.7
Web Storefront : III.7
```

A component can be qualified with a **sourcing type**:

| Type | Meaning |
|------|---------|
| `build` | Built in-house |
| `buy` | Purchased as a product |
| `outsource` | Outsourced to a third party |

```
Cloud Compute : IV.5 (buy)
CDN : IV.7 (outsource)
Recommendation Engine : II.5 (build)
```

### Submaps

A **submap** represents a component encapsulating a separate Wardley Map. It is rendered as a component with a specific visual marker.

```
submap Fulfilment and Logistics : III.8
```

### Block Declaration

To add extra properties, use the block declaration with braces:

```
Recommendation Engine : II.5 {
  type: build
  color: #8E44AD
  asset: tech
  cost: "900k EUR/yr, 8 engineers"
  note: "Key differentiator"
}
```

| Field | Values | Description |
|-------|--------|-------------|
| `type` | `build`, `buy`, `outsource` | Sourcing type |
| `asset` | `tech`, `financial`, `human`, `relational`, `social` | Capital nature (see below) |
| `color` | `#RRGGBB`, `#RGB`, CSS color name | Rendering color |
| `visibility` | `0.0` to `1.0` | Forced vertical position |
| `cost` | Free text | Cost annotation |
| `note` | Quoted text | Component description |

#### Asset Classification

Each component represents a type of organizational capital:

| Type | Description | Example |
|------|------------|---------|
| `tech` | Technological capital: code, infrastructure, patents | A routing engine, a data pipeline |
| `financial` | Financial capital: revenue models, pricing power | A billing system, a licensing model |
| `human` | Human capital: expertise, skills, tacit knowledge | An ML engineering team, domain experts |
| `relational` | Relational capital: partnerships, brand, contracts | A partner API, a brand, a patent portfolio |
| `social` | Social/environmental capital: community, regulatory | Open-source community, regulatory compliance |

---

## Evolution Positioning

The horizontal position uses **roman numeral notation** with an optional decimal subdivision:

```
<roman>.<digit>
```

Where `<roman>` is `I`, `II`, `III`, or `IV`, and `<digit>` is `0` through `9`.

Each phase spans 25% of the axis. The decimal subdivides within the phase (0 = start, 9 = end). Without a decimal, the center of the phase is used.

### Position Mapping Table

| Position | Coordinate (0-100) | Meaning |
|----------|-------------------|---------|
| `I.0` | 0 | Start of Genesis |
| `I.5` | 12 | Middle of Genesis |
| `I.9` | 22 | End of Genesis |
| `II.0` | 25 | Start of Custom-Built |
| `II.5` | 37 | Middle of Custom-Built |
| `II.9` | 47 | End of Custom-Built |
| `III.0` | 50 | Start of Product |
| `III.5` | 62 | Middle of Product |
| `III.9` | 72 | End of Product |
| `IV.0` | 75 | Start of Commodity |
| `IV.5` | 87 | Middle of Commodity |
| `IV.9` | 97 | End of Commodity |

**Formula:** `floor((base + digit/10 * 0.25) * 100)` where `base` is `I=0.00, II=0.25, III=0.50, IV=0.75`.

Without a decimal (e.g., bare `III`), the component is placed at the center of the phase (12, 37, 62, or 87).

---

## Evolution Movement

The `>>` operator indicates that a component is actively evolving from one position to another:

```
Component : II.7 >> III.5
```

This renders a dashed arrow from the current position (II.7) to the target position (III.5), with a "ghost" component at the destination.

---

## Inertia

Inertia represents resistance to change. It is marked with exclamation marks (`!`) between the current position and the `>>` movement operator.

### Inertia Levels

```
Component : II.7 ! >> III.5       // moderate inertia
Component : II.7 !! >> III.5      // strong inertia
Component : II.7 !!! >> III.5     // blocking inertia
```

### Qualified Inertia

Inertia can be qualified by its nature using a comma-separated list of types in parentheses:

```
Component : II.7 !!(tech,human) >> III.5      // technical + human inertia
Component : II.7 !(financial) >> III.5         // financial inertia only
Component : II.7 !!!(tech,financial,human) >> III.5  // three types of resistance
```

| Type | Meaning | Typical Symptom |
|------|---------|----------------|
| `tech` | Technology lock-in, infrastructure debt | "We've always used Java" |
| `financial` | Sunk costs, established revenue models | "We've invested 5M in this" |
| `human` | Skills gap, identity threat, expertise obsolescence | "Our team doesn't know cloud-native" |
| `relational` | Contractual obligations, partner dependencies | "We have a 3-year vendor contract" |
| `social` | Cultural resistance, regulatory inertia | "That's not how we do things here" |

Qualified inertia is rendered visually with colored bars side by side on the evolution arrow.

---

## Visibility

By default, a component's vertical position is **computed automatically** from the dependency graph: anchors are at the top, infrastructure components at the bottom.

You can override the vertical position with the `@` operator:

```
Component : III.5 @0.9    // near top of map
Component : III.5 @0.1    // near bottom of map
```

Values range from `0.0` (bottom) to `1.0` (top). This override is rarely needed -- the automatic layout usually produces a good result.

Visibility can also be set in a block:

```
Component : III.5 {
  visibility: 0.8
}
```

---

## Pipelines

A **pipeline** models a strategic component existing in **multiple simultaneous forms** at different evolution stages. It is one of the most powerful concepts in Wardley Mapping.

### Concept

A pipeline represents, for example, a recommendation engine that exists in three implementations:
- Collaborative filtering (mature product)
- Deep learning recommender (custom-built, in development)
- Contextual bandit (research, genesis)

All three implementations serve the same function but sit at different evolution levels.

### Syntax

```
// First, declare the parent component
Recommendation Engine : II.5

// Then, declare the pipeline
pipeline Recommendation Engine {
  Collaborative Filtering : III.3
  Deep Learning Recommender : II.3
  Contextual Bandit : I.5
}
```

### Rules

- The pipeline name **must match** an already-declared component.
- Members are positioned on the evolution axis only; their vertical position is derived from the parent component.
- The pipeline's horizontal span covers from its leftmost to rightmost member.
- A pipeline **cannot contain** another pipeline.
- Each member is declared with a name and a position: `Name : Position`.

### Visual Rendering

The pipeline is displayed as a horizontal rectangle containing its members, positioned on the evolution axis. The parent component serves as the anchor point for value chain edges.

### Referencing a Pipeline Member

In edges, you can target a specific member of a pipeline using the `:` syntax:

```
Real-Time Personalisation -> Recommendation Engine:Deep Learning Recommender
```

---

## Value Chain (Edges)

Edges define dependencies between components and form the value chain.

### Basic Syntax

```
A -> B                              // A depends on B
A <-> B                             // bidirectional relationship
A -[label text]-> B                 // annotated dependency
A <-[label text]-> B                // annotated bidirectional
```

### Chains

Edges can be chained on a single line:

```
Customer -> Web Storefront -> Search and Discovery -> Recommendation Engine
```

This creates three separate edges: Customer->Web Storefront, Web Storefront->Search and Discovery, Search and Discovery->Recommendation Engine.

### Labels

Labels document the nature of the relationship:

```
Merchant -> Merchant Analytics -[sales data feed]-> Customer Data Platform
```

The text between brackets is displayed along the edge on the map.

### Pipeline Member Reference

```
Real-Time Personalisation -> Recommendation Engine:Deep Learning Recommender
```

The `Pipeline:Member` syntax creates an edge to a specific member of a pipeline rather than to the parent component.

---

## Groups

Groups visually cluster components, typically by team or domain. They are **purely visual** and do not create a namespace.

### Basic Syntax

```
group AI and Personalisation {
  Recommendation Engine
  Deep Learning Recommender
  Contextual Bandit
  Real-Time Personalisation
  Customer Data Platform
}
```

### Group Directives

| Directive | Values | Description |
|-----------|--------|-------------|
| `color:` | `#RRGGBB`, `#RGB`, CSS color name | Group envelope color |
| `team:` | `explorer`, `settler`, `town-planner`, `pioneer`, `villager` | EVT/PST team type |

```
group R&D Team {
  team: explorer
  color: #E74C3C
  Contextual Bandit
  Experimental Cache
}
```

### EVT/PST Alignment

The Explorer-Villager-Town-Planner model aligns team types to evolution phases:

| Team Type | Evolution Phase | Mindset |
|-----------|----------------|---------|
| `explorer` / `pioneer` | Genesis (I) | Discovery, intuition, high failure tolerance |
| `settler` / `villager` | Custom/Product (II-III) | Productization, standards, analysis |
| `town-planner` | Commodity (IV) | Industrialization, cost optimization |

A mismatch between team type and the evolution phase of owned components is a strategic signal worth investigating.

---

## Annotations

Annotations attach textual information to a component.

### Notes

```
note "Evaluate Stripe vs Adyen migration Q3 2026" on Payment Gateway
note "Target: real-time inventory sync by Q2 2026" on Inventory Service
```

### Warnings

```
warning "Single vendor lock-in -- AWS costs rising 25% YoY" on Cloud Compute
warning "GDPR compliance gap in cross-border data flows" on Customer Data Platform
```

Notes are rendered with a neutral style (blue) while warnings use an alert style (orange/red).

---

## Signals

Signals mark market dynamics and climatic patterns on a component.

```
signal accelerating on Deep Learning Recommender
signal stagnating on Collaborative Filtering
signal declining on Legacy Module
```

### Signal Types

| Signal | Meaning |
|--------|---------|
| `accelerating` | Rapidly moving toward commodity |
| `stagnating` | Evolution has plateaued |
| `declining` | Regression in relevance or usage |
| `co-evolution` | Technology and practice evolving together |
| `red-queen` | Must evolve constantly just to maintain position |
| `commoditization` | Gravitational pull toward utility |
| `network-effects` | Value grows with number of users/participants |
| `economies-of-scale` | Cost advantage from volume |

```
signal co-evolution on DevOps Practices
signal commoditization on Container Orchestration
signal network-effects on Merchant Portal
signal economies-of-scale on Cloud Compute
signal red-queen on Web Storefront
```

---

## Strategic Maneuvers (Gameplays)

Gameplays annotate deliberate strategic maneuvers applied to a component.

### Syntax

```
gameplay <type> on <Component>
gameplay <type> "optional description" on <Component>
```

### Gameplay Types

| Gameplay | Description | Typical Context |
|----------|-------------|-----------------|
| `ILC` | Innovate-Leverage-Commoditize: provide infrastructure, observe, absorb | Platform with ecosystem |
| `open-source` | Commoditize a layer to capture value in an adjacent layer | Competitor with proprietary rent |
| `land-grab` | Sacrifice profitability for rapid market share | New market with network effects |
| `embrace-extend` | Adopt a standard, add proprietary extensions | Standard you want to control |
| `tower-moat` | Erect barriers: patents, lock-in, closed protocols | Protecting an existing rent |
| `FUD` | Spread fear/uncertainty/doubt to slow competitor adoption | Competitor gaining traction |
| `strangler-fig` | Progressively replace a legacy system | Legacy system blocking evolution |
| `signal-distortion` | Mislead competitors about strategic intent | Competitive misdirection |

### Examples

```
gameplay ILC on Platform API
gameplay open-source "Commoditize compute to capture AI layer" on Cloud Infrastructure
gameplay strangler-fig on Legacy CRM
gameplay land-grab on New Market
gameplay FUD "Slow down competitor adoption" on Rival Solution
```

The quoted description is optional and provides additional strategic context.

---

## Focus

The `focus` keyword highlights a component and all its descendants in the value chain. Elements outside the focus are rendered with reduced opacity.

```
focus Recommendation Engine
```

Multiple focus declarations can be combined -- their subtrees are merged:

```
focus Recommendation Engine
focus Real-Time Personalisation
```

Focus is useful for presenting a specific part of the map during a discussion or presentation.

---

## Identifier Rules

Identifiers (component names, group names, etc.):

- Start with a letter or digit
- May contain letters, digits, `.`, `-`, `'`, `_`, and **spaces**
- Spaces are allowed within identifiers (e.g., `Web Storefront`)
- Cannot be a reserved keyword used alone

### Reserved Keywords

`anchor`, `component`, `submap`, `pipeline`, `group`, `note`, `warning`, `signal`, `gameplay`, `legend`, `focus`, `title`, `date`, `author`, `scope`, `question`, `stages`, `doctrine`, `evolution`, `type`, `asset`, `color`, `visibility`, `cost`, `team`, `build`, `buy`, `outsource`, `accelerating`, `stagnating`, `declining`, `co-evolution`, `red-queen`, `commoditization`, `network-effects`, `economies-of-scale`, `ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`, `explorer`, `settler`, `town-planner`, `pioneer`, `villager`, `tech`, `financial`, `human`, `relational`, `social`, `on`, `hygiene`, `context`, `excellence`, `legend`

---

## Semantic Rules

1. Every node referenced in an edge, annotation, signal, or pipeline **must be declared** with a position somewhere in the document.
2. A pipeline name **must match** a declared component.
3. Pipelines **cannot be nested**.
4. Groups **do not create namespaces** -- components remain global.
5. Anchors **do not require** an evolution position (they are placed at the top automatically).
6. The dependency graph **must be acyclic** (no circular dependencies).
7. An orphan node (not referenced in any edge) is placed at the bottom of the map by default.

---

## Complete Example

```wtg2
// Wardley Map -- Online Retail Platform

title: Online Retail Platform -- Strategic Review 2026
date: 2026-04-09
author: Platform Strategy Team
scope: B2C e-commerce, EU market
question: "How do we reduce infrastructure costs while accelerating AI-driven personalisation?"
doctrine: context

stages: Genesis, Custom-Built, Product, Commodity
legend

// Anchors
anchor Customer
anchor Merchant

// User-facing layer
Web Storefront : III.7
Mobile App : III.4
Search and Discovery : III.0

// Core business logic with qualified inertia
Recommendation Engine : II.5 !!(tech,human) >> III.3 {
  type: build
  asset: tech
  color: #8E44AD
  cost: "900k EUR/yr, 8 engineers"
  note: "Key differentiator"
}

// Pipeline: the engine exists in multiple forms
pipeline Recommendation Engine {
  Collaborative Filtering : III.3
  Deep Learning Recommender : II.3
  Contextual Bandit : I.5
}

Order Management : III.5 (buy)
Inventory Service : III.2 (buy)
Pricing Engine : II.8 ! >> III.4 {
  type: build
  color: #2980B9
}

// Merchant-facing
Merchant Portal : II.8
Merchant Analytics : II.2 {
  type: build
  asset: tech
  color: #27AE60
}

// Payments
Payment Gateway : IV.2 (outsource)

// Data layer
Customer Data Platform : II.6 {
  type: build
  asset: tech
  color: #E74C3C
}
Product Catalogue : III.6 (buy)

// Infrastructure
Cloud Compute : IV.5 (buy) {
  cost: "400k EUR/yr, rising 25%"
}
CDN : IV.7 (outsource)
Container Orchestration : IV.1 (buy)
Observability Stack : III.5 (buy)
Message Queue : IV.3 (buy)

// Emerging
Real-Time Personalisation : I.7 !(tech) >> II.4 {
  type: build
  asset: tech
  color: #E67E22
  cost: "150k EUR/yr, 3 engineers"
  note: "R&D initiative -- 6-month horizon"
}

// Sub-map
submap Fulfilment and Logistics : III.8

// Value chain -- Customer flow
Customer -> Web Storefront -> Search and Discovery -> Recommendation Engine
Customer -> Mobile App -> Search and Discovery
Web Storefront -> Order Management -> Inventory Service
Mobile App -> Order Management
Order Management -> Payment Gateway
Order Management -> Fulfilment and Logistics

// Recommendation pipeline
Recommendation Engine -> Customer Data Platform
Recommendation Engine -> Product Catalogue
Search and Discovery -> Product Catalogue

// Pricing
Web Storefront -> Pricing Engine -> Product Catalogue
Pricing Engine -> Customer Data Platform

// Merchant flow
Merchant -> Merchant Portal -> Inventory Service
Merchant -> Merchant Analytics -[sales data feed]-> Customer Data Platform
Merchant Portal -> Product Catalogue

// Real-time personalisation
Search and Discovery -> Real-Time Personalisation -> Customer Data Platform
Real-Time Personalisation -> Recommendation Engine:Deep Learning Recommender

// Infrastructure dependencies
Order Management -> Message Queue -> Cloud Compute
Customer Data Platform -> Cloud Compute
Cloud Compute -> Container Orchestration
Web Storefront -> CDN -> Cloud Compute
Mobile App -> CDN
Observability Stack -> Cloud Compute

// Groups with team types
group AI and Personalisation {
  team: explorer
  Recommendation Engine
  Deep Learning Recommender
  Contextual Bandit
  Real-Time Personalisation
  Customer Data Platform
}

group Commerce Core {
  team: town-planner
  Order Management
  Inventory Service
  Pricing Engine
  Payment Gateway
  Product Catalogue
}

group Platform Infrastructure {
  color: #95A5A6
  Cloud Compute
  CDN
  Container Orchestration
  Message Queue
  Observability Stack
}

group Merchant Experience {
  team: settler
  Merchant Portal
  Merchant Analytics
}

// Warnings
warning "Single vendor lock-in -- AWS costs rising 25% YoY" on Cloud Compute
warning "GDPR compliance gap in cross-border data flows" on Customer Data Platform
warning "Legacy monolith -- hard to scale independently" on Order Management

// Strategic notes
note "Evaluate Stripe vs Adyen migration Q3 2026" on Payment Gateway
note "Target: real-time inventory sync by Q2 2026" on Inventory Service
note "Patent filed for contextual bandit approach" on Contextual Bandit

// Market signals
signal accelerating on Deep Learning Recommender
signal accelerating on Real-Time Personalisation
signal stagnating on Collaborative Filtering
signal commoditization on Container Orchestration
signal network-effects on Merchant Portal
signal economies-of-scale on Cloud Compute

// Gameplays
gameplay strangler-fig "Replace Collaborative Filtering with Deep Learning" on Recommendation Engine
gameplay ILC on Merchant Portal
gameplay open-source "Commoditize observability to reduce vendor lock-in" on Observability Stack

// Focus
focus Recommendation Engine
```

---

## Quick Syntax Reference

| Element | Syntax | Example |
|---------|--------|---------|
| Metadata | `key: value` | `title: My Map` |
| Anchor | `anchor Name` | `anchor Customer` |
| Component | `Name : Position` | `API : III.5` |
| Typed component | `Name : Position (type)` | `DB : III.8 (buy)` |
| Submap | `submap Name : Position` | `submap Payment : III.6` |
| Block | `Name : Position { ... }` | see above |
| Evolution | `Name : Pos >> Pos` | `API : II.7 >> III.5` |
| Inertia | `! !! !!!` | `API : II.7 !! >> III.5` |
| Qualified inertia | `!!(types)` | `API : II.7 !!(tech,human) >> III.5` |
| Visibility | `@0.0-1.0` | `API : III.5 @0.8` |
| Pipeline | `pipeline Name { ... }` | see above |
| Edge | `A -> B` | `Customer -> App` |
| Annotated edge | `A -[text]-> B` | `A -[data feed]-> B` |
| Bidirectional | `A <-> B` | `A <-> B` |
| Chained edge | `A -> B -> C` | `Customer -> App -> DB` |
| Pipeline ref | `Pipeline:Member` | `Engine:Deep Learning` |
| Group | `group Name { ... }` | see above |
| Note | `note "text" on Name` | `note "Important" on API` |
| Warning | `warning "text" on Name` | `warning "Risk" on DB` |
| Signal | `signal type on Name` | `signal accelerating on AI` |
| Gameplay | `gameplay type on Name` | `gameplay ILC on API` |
| Focus | `focus Name` | `focus Engine` |
| Legend | `legend` | `legend` |
| Stages | `stages: A, B, C, D` | `stages: Genesis, Custom, Product, Commodity` |

# WTG2 Strategic Concepts

Strategic theory from Simon Wardley's framework. Read this when the user asks for competitive analysis, organizational mapping, or strategic depth in their Wardley Maps.

---

## Gameplays (Strategic Maneuvers)

A gameplay is a deliberate action to modify your position on the map. Annotate gameplays on the component that is the *target* of the maneuver.

| Gameplay | Description | Typical Context |
|----------|-------------|-----------------|
| `ILC` | Innovate-Leverage-Commoditize: provide infrastructure, observe what thrives, absorb | Platform with ecosystem |
| `open-source` | Commoditize a layer to capture value in an adjacent layer | Competitor with proprietary rent |
| `land-grab` | Sacrifice profitability for rapid market share | New market with strong network effects |
| `embrace-extend` | Adopt a standard, add proprietary extensions, close ecosystem | Standard you want to control |
| `tower-moat` | Erect barriers: patents, lock-in, closed protocols | Protecting an existing rent |
| `FUD` | Spread fear/uncertainty/doubt to slow competitor adoption | Competitor gaining traction |
| `strangler-fig` | Progressively replace a legacy system by commoditizing non-differentiating layers | Legacy system blocking evolution |
| `signal-distortion` | Mislead competitors about strategic intent | Competitive misdirection |
| `due-diligence` | Strategic due diligence: map M&A target's value chain; detect anomalies (Custom where should be Commodity), assess EVT alignment, reveal real vs. paper synergies | Merger/acquisition, asset evaluation |
| `two-sided-market` | Create obligatory passage point via cross-sided network effects: more producers attract consumers and vice versa; virtuous circle creates entry barriers at critical mass | Platform connecting producers and consumers |

```wtg2
gameplay ILC on Platform API
gameplay open-source "Commoditize compute to capture AI layer" on Cloud Infra
gameplay strangler-fig on Legacy CRM
gameplay due-diligence "Assess strategic coherence before acquisition" on Target Company
gameplay two-sided-market on Marketplace
```

### When to Use Which Gameplay

| Situation | Symptoms | Recommended Gameplay |
|-----------|----------|---------------------|
| Legacy system blocking innovation | Time-to-market > 6 months, massive debt, frustrated teams | `strangler-fig` |
| Competitor extracting rent from proprietary component | High margins, customer lock-in, no alternatives | `open-source` |
| New market with network effects emerging | Exponential growth, short window, race for leadership | `land-grab` |
| Established platform with ecosystem | Many partners, usage data, intermediary position | `ILC` |
| M&A under consideration | Financial due diligence OK, need strategic value assessment | `due-diligence` |
| Mature product threatened, margins to preserve | External innovation emerging, need time to pivot | `tower-moat` |
| Platform connecting producers and consumers | Cross-sided network effects, critical mass dynamics | `two-sided-market` |

---

## Five Capitals (Asset Classification)

Components represent different types of organizational capital. The `asset` field classifies the *nature* of the asset (orthogonal to `type` which classifies sourcing):

| Asset Type | Description | Example |
|------------|-------------|---------|
| `tech` | Technological capital: code, infrastructure, patents | A routing engine, a data pipeline |
| `financial` | Financial capital: revenue models, pricing power | A billing system, a licensing model |
| `human` | Human capital: expertise, skills, tacit knowledge | An ML engineering team, domain experts |
| `relational` | Relational capital: partnerships, brand, contracts | A partner API, a brand, a patent portfolio |
| `social` | Social/environmental capital: community, regulatory | Open-source community, regulatory compliance |

Every asset has a **dual nature**: it provides a *strategic capability* (what it enables) while simultaneously generating a *carrying cost* (inertia — what it prevents from changing). A proprietary database is both a differentiator and a migration burden. The `asset` field captures the capability side; `cost:` and inertia (`!`) capture the carrying cost side.

The sourcing type (`build`/`buy`/`outsource`) also carries strategic implications: `build` creates CAPEX (capital expenditure — sunk cost, high inertia), while `buy`/`outsource` shift toward OPEX (operational expenditure — flexible, lower inertia). Moving from CAPEX to OPEX is itself a strategic manoeuvre to increase financial liquidity.

```wtg2
AI Team : II.3 {
    type: build
    asset: human
    note: "12 ML engineers, hard to replace"
}
```

---

## Qualified Inertia

Inertia is not a defect — it is a mechanical consequence of past success. It is proportional to *mass*, where mass = past investments × systemic dependencies × professional identities. The more successful a component has been, the harder it is to change.

Phase transitions (e.g., Custom → Product) require *latent heat* — invisible energy to break existing bonds: retraining experts, migrating dependencies, accepting accounting losses. This cost is invisible on a balance sheet but real on the map.

Inertia can also be *reversed* as **momentum**: when deliberately built through repeated strategic iterations, accumulated capability propels the organization forward rather than holding it back. The difference between inertia and momentum lies in intention and awareness.

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

### Assessing Inertia Severity

Total inertia is proportional to: **past investments × systemic dependencies × built identities**. To assess severity, ask:

| Dimension | Questions | Higher = stronger inertia |
|-----------|-----------|--------------------------|
| Investments | How much invested? What's the residual book value? | Large sunk cost makes psychological abandonment harder |
| Dependencies | How many systems/processes rely on this? | More dependencies = more complex untangling |
| Identities | How many people built expertise, power, or career identity on this? | Stronger identities = stronger emotional resistance |

### Phase Transition Latent Heat

Each phase transition requires different types of energy to break the bonds of the previous state:

| Transition | Barrier type | Why it's hard |
|------------|-------------|---------------|
| I → II | Political & psychological | Breaking mental models, overcoming institutional risk aversion |
| II → III | Human capital & identity crisis | Experts lose unique status as knowledge standardizes |
| III → IV | Financial & relational | Revenue model destruction, margin collapse, sales structure obsolescence |

---

## Climatic Patterns (Extended Signals)

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

---

## EVT/PST Team Alignment

The Explorer-Villager-Town-planner model aligns team types to evolution phases:

| Team Type | Evolution Phase | Mindset |
|-----------|----------------|---------|
| `explorer` / `pioneer` | Genesis (I) | Discovery, intuition, high failure tolerance |
| `settler` / `villager` | Custom–Product (II–III) | Productization, standards, analysis |
| `town-planner` | Commodity (IV) | Industrialization, cost optimization |

A mismatch between team type and component evolution phase is a strategic signal worth highlighting.

```wtg2
group R&D Team {
    team: explorer
    Quantum Algo
}
```

### Organization-Phase Alignment Grid

| Evolution phase | Team profile | Organizational role |
|----------------|-------------|-------------------|
| Genesis (I) | Explorers / Commandos | R&D, skunkworks, innovation labs |
| Custom (II) | Artisans / Villagers | Engineering teams, consultants |
| Product (III) | Product managers / Settlers | Product teams, solution architects |
| Commodity (IV) | Ops/SRE / Town Planners | Platform teams, infrastructure |

**Formal transfer principle:** Explorers create, Villagers productize, Town Planners industrialize. The fatal error is assigning the wrong profile to a phase — e.g., asking Town Planners to innovate (they optimize) or Explorers to run production (they break things). Transitions between profiles generate friction; transfers must be orchestrated, not assumed.

---

## Carrying Cost

The book emphasizes that 70–80% of IT budgets go to maintenance. Use `cost:` to annotate financial context and enable run/change ratio analysis:

```wtg2
Legacy CRM : III.2 {
    type: buy
    cost: "850k/year, 80% maintenance"
}
```

---

## Strategic Cycle

Strategy is not an event — it is a continuous discipline following the **OODA loop** (Observe → Orient → Decide → Act). Each iteration reduces friction, transforming the laborious effort of a beginner into the fluidity of an experienced practitioner.

The **Value Flywheel Effect** structures this cycle in 4 phases:
1. **Clarity of Purpose** — understand the user need (anchors, value chain)
2. **Challenge and Landscape** — map the terrain (components, evolution, climate)
3. **Next Best Actions** — decide and act (gameplays, manoeuvres)
4. **Long Term Value** — harvest value and reinvest (momentum builds)

Each iteration feeds the next: capabilities developed in phase 1 facilitate phase 2; patterns recognized in phase 2 accelerate phase 3. The *velocity of adaptation* (how fast you complete the cycle) becomes a competitive advantage in itself.

---

## Doctrine Violations

Recurrent anti-patterns detectable on a map. Flag these with `warning` or `note` annotations:

| Violation | Map symptom | What to flag |
|-----------|------------|--------------|
| NIH (Not Invented Here) | Component marked `build` in phase III–IV when market alternatives exist | `warning "NIH — standard solutions exist at this evolution stage"` |
| No differentiation | All components in phase III–IV, zero investment in I–II | `warning "No differentiation — entire value chain is commodity"` |
| Dispersion | Many components in phase I without critical mass in any | `warning "Dispersion — too many bets, insufficient focus"` |
| Single method | Same approach applied uniformly regardless of evolution phase | `note "Agile/Six Sigma may not suit all phases equally"` |
| Strategy theatre | Map has no `question:`, no `gameplay`, no movement (`>>`) | A map without a question is a ritual, not strategy |

---

## Flow Analysis

A map is not just a snapshot — value flows through it. Two types of flow reveal invisible dynamics:

1. **Evolutionary flow** — The gravitational pull from left to right (Genesis → Commodity). Every component tends toward commoditization. Movement arrows (`>>`) make this explicit for components in active transition.
2. **Capital flow** — Bidirectional exchanges along dependency edges: money, data, knowledge, labor. Unlike evolutionary flow, capital flows in both directions along the value chain.

### Flow Anomalies

| Anomaly | Symptom | WTG2 annotation |
|---------|---------|-----------------|
| One-way flow | A dependency where value only flows in one direction (e.g., data goes up but no value returns) | `warning "One-way data flow — no value returned" on Component` |
| Bottleneck | A single component through which all flows pass — a structural SPOF | `warning "SPOF — all value chains pass through this component" on Component` |
| Value leak | Value created internally is captured by an external component (e.g., a `buy` component in phase III extracting rent) | `warning "Value leak — vendor captures margin on our differentiation" on Component` |

---

## Extremistan / Mediocristan Divide

Phases I–II (Genesis, Custom) operate in **Extremistan**: power laws, unpredictable ROI, winner-takes-all dynamics. Phases III–IV (Product, Commodity) operate in **Mediocristan**: normal distributions, predictable metrics, Six Sigma applicable.

Do not apply the same KPIs or risk models uniformly across all phases — a startup metric (burn rate) does not apply to infrastructure, and an SLA does not apply to R&D. When generating maps, flag components where the management approach does not match the evolution phase.

---

## Evolution Cheat Sheet

Detailed characteristics per evolution phase — use this to position components accurately and to detect phase mismatches.

| Property | I Genesis | II Custom | III Product | IV Commodity |
|----------|-----------|-----------|-------------|--------------|
| Market | Undefined | Forming, competing forms | Growing, consolidation | Mature, stabilized |
| User perception | Exciting, disorienting | Avant-garde, emergent | Common, disappointing if absent | Standard, expected |
| Ecosystem perception | Chaotic, "domain of madmen" | Domain of experts | Domain of professionals | Ordered, trivial |
| Value focus | High future value, investment | Profit & ROI search | Strong unit profitability | Volume, margin reduction |
| Failure tolerance | High, tolerated | Moderate, disappointing | Not tolerated, triggers improvement | Shocking, surprising |
| Decision mode | Heritage, intuition, bet | Analysis and synthesis | Analysis, competing proofs | Past experience |
| Knowledge | Emerging | Usage-based learning | Training available, metrics exist | Industry standards |
| Efficiency focus | — | — | Reduce waste through learning | Minimize variance |
| Dominant profile | Explorer / Inventor | Artisan / Expert | Product Manager | Industrialist |
| Management | Explore (Agile, Design Thinking) | Expert craft (Lean Startup) | Industrialize (Lean, Kanban) | Optimize (Six Sigma, SRE) |

### Phase Transition Signals

How to detect that a component is shifting phase — these often translate to `signal` annotations:

| Signal | Interpretation |
|--------|---------------|
| Competitors copy your solution | Transition to Phase III underway |
| Customers ask "why so expensive?" | Transition to Phase IV imminent |
| Margins erode despite innovation | Phase IV already in progress |
| Experts become "temple guardians" | Inertia II→III to anticipate |
| Newcomers understand faster than veterans | Standardization happening |

---

## Doctrine Phases

The `doctrine:` metadata field declares the organization's maturity. Each phase has distinct practices and anti-patterns — use this to calibrate the depth and focus of map annotations.

### Phase 1: Hygiene — Stop self-harm

**Objective:** Eliminate internal friction before pursuing strategy.

Key practices:
- **Common language** — enforce precise definitions across all actors (not just glossaries, but aligned understanding)
- **User focus** — privilege effectiveness (doing the right thing) over efficiency (doing the right thing fast)
- **Situational awareness** — require visual representation of terrain before decisions; no SWOT-only strategy

Anti-pattern: **Glossary illusion** — defining terms without changing incentive structures.

### Phase 2: Context — Become methodologically polyglot

**Objective:** Adapt method to problem type.

Key practices:
- Agile for Genesis (uncertain, exploratory), Lean for Product (optimize flow), Six Sigma for Commodity (minimize variance)
- **F.I.R.E. principle** (Fast, Inexpensive, Restrained, Elegant) — keep initiatives small and disposable
- Manage three types of inertia: skill, investment, and cultural

Anti-pattern: **Agile dogmatism** — forcing exploration methods onto stable production systems.

### Phase 3: Excellence — Do more with less

**Objective:** Ruthless efficiency, but only after validating effectiveness (Phase 1) and context-fit (Phase 2).

Key practices:
- Eliminate waste (Muda); reduce lead-time-to-touch-time ratio toward 1:1
- "You build it, you run it" — creators feel operational pain, quality rises
- Autonomy + Purpose + Mastery drive elite talent

Anti-pattern: **Premature optimization** — making the wrong thing efficient.

### Phase 4: Evolution — Self-reinvention

**Objective:** Maintain adaptability; no permanent core business.

Key practices:
- **"No Core" principle** — everything is transitional; willingness to cannibalize yesterday's success
- Landscape manipulation — don't just play the game, change the rules
- Competitive sensing: detect weak signals before competitors

Anti-pattern: **Premature cannibalization** — destroying the cash-cow before building a replacement.

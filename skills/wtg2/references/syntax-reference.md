# WTG2 Syntax Reference

Detailed syntax rules for the WTG2 domain-specific language. Read this when you need precise syntax details beyond the quick reference in SKILL.md.

---

## Comments

```wtg2
// Single-line comment
/* Block comment — can span multiple lines */
```

---

## Evolution Positioning

Format: `<roman>.<digit>` where roman is I/II/III/IV and digit is 0–9.

Each phase spans 25% of the axis. The digit subdivides within the phase (0 = start, 9 = end). Without a digit, the center of the phase is used.

### Full Position Mapping

| Position | Coordinate | Meaning |
|----------|------------|---------|
| `I.0` | 0 | Start of Genesis |
| `I.5` | 12 | Middle of Genesis |
| `I.9` | 22 | End of Genesis |
| `II.0` | 25 | Start of Custom |
| `II.5` | 37 | Middle of Custom |
| `II.9` | 47 | End of Custom |
| `III.0` | 50 | Start of Product |
| `III.5` | 62 | Middle of Product |
| `III.9` | 72 | End of Product |
| `IV.0` | 75 | Start of Commodity |
| `IV.5` | 87 | Middle of Commodity |
| `IV.9` | 97 | End of Commodity |

**Formula:** `floor((base + digit/10 * 0.25) * 100)` where base is I=0.00, II=0.25, III=0.50, IV=0.75.

---

## Node Declaration Details

### Shorthand Syntax

```
[kind] <name> : <evolution> [(<type>)] [@<visibility>]
```

The `component` keyword is optional — a bare name with a position is treated as a component.

### Block Config Fields

| Field | Values | Description |
|-------|--------|-------------|
| `type` | `build`, `buy`, `outsource` | Sourcing strategy |
| `asset` | `tech`, `financial`, `human`, `relational`, `social` | Capital classification |
| `evolution` | Evolution expression | Alternative to inline position |
| `color` | `#RRGGBB` or `#RGB` | Custom color |
| `visibility` | `0.0` to `1.0` | Vertical position override (0=bottom, 1=top) |
| `cost` | Free text | Cost annotation |
| `note` | Quoted text | Description |

### Inertia Syntax

Inertia appears between the current position and the `>>` movement operator:

```wtg2
Component : II.7 ! >> III.5               // moderate, unqualified
Component : II.7 !! >> III.5              // strong, unqualified
Component : II.7 !!! >> III.5             // blocking, unqualified
Component : II.7 !!(tech,human) >> III.5  // strong, qualified by kind
Component : II.7 !(financial) >> III.5    // moderate, qualified by kind
```

Qualification kinds: `tech`, `financial`, `human`, `relational`, `social`.

---

## Edge Syntax Details

Four forms of edges:

```wtg2
A -> B                          // A depends on B
A <-> B                         // bidirectional relationship
A -[label text]-> B             // annotated dependency
A <-[label text]-> B            // annotated bidirectional
```

Edges can be chained — `A -> B -> C -> D` creates three separate edges.

### Pipeline Member References

Target a specific member within a pipeline:

```wtg2
Component -> Pipeline:Member
```

---

## Pipeline Rules

```wtg2
pipeline <component-name> {
    Implementation A : III.5
    Implementation B : II.3
}
```

- The pipeline name **must** match an already-declared component.
- Members are positioned on the evolution axis only; vertical position comes from the parent.
- The pipeline's horizontal span covers from its leftmost to rightmost member.
- Pipelines cannot be nested.

---

## Group Directives

```wtg2
group Team Name {
    team: explorer         // explorer | settler | town-planner | pioneer | villager
    color: #E74C3C         // optional color
    Component A
    Component B
}
```

- Members must reference existing component names.
- Groups do not create namespaces — components remain global.

---

## Visibility Override

Override automatic vertical positioning:

```wtg2
Component : III.5 @0.9    // near top of map
Component : III.5 @0.1    // near bottom of map
```

Values range from `0.0` (bottom) to `1.0` (top). By default, vertical position is computed from the dependency graph.

---

## Identifier Rules

Identifiers (component names, group names, etc.):
- Start with a letter or digit
- May contain letters, digits, `.`, `-`, `'`, `_`, and spaces
- Spaces are allowed inside identifiers (e.g., `Application Mobile`)
- Cannot be a reserved keyword used alone

### Reserved Keywords

`anchor`, `component`, `submap`, `pipeline`, `group`, `note`, `warning`, `signal`, `gameplay`, `legend`, `focus`, `title`, `date`, `author`, `scope`, `question`, `stages`, `doctrine`, `evolution`, `type`, `asset`, `color`, `visibility`, `cost`, `team`, `build`, `buy`, `outsource`, `accelerating`, `stagnating`, `declining`, `co-evolution`, `red-queen`, `commoditization`, `network-effects`, `economies-of-scale`, `ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`, `explorer`, `settler`, `town-planner`, `pioneer`, `villager`, `tech`, `financial`, `human`, `relational`, `social`, `on`, `hygiene`, `context`, `excellence`

---

## Semantic Rules

1. Every node referenced in an edge, annotation, signal, or pipeline must be declared with a position somewhere in the document.
2. Pipeline names must match a declared component.
3. Pipelines cannot be nested.
4. Groups do not create namespaces — components remain global.
5. Anchors do not need an evolution position (they are placed at the top automatically).
6. The dependency graph must be acyclic (no circular dependencies).

---

## Gameplay Extensions

The parser accepts any string as a gameplay type without validation. Beyond the 8 types with dedicated badge colors (`ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`), the following types are also supported and render with a default gray badge:

- `due-diligence` — Strategic due diligence for M&A evaluation
- `two-sided-market` — Cross-sided network effects platform strategy

These are not in the formal BNF grammar but are accepted by the parser and documented in the skill.

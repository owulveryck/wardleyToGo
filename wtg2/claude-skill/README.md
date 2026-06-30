# WTG2 Skill for Claude Code

A [Claude Code](https://docs.anthropic.com/en/docs/claude-code) skill that generates [Wardley Maps](https://wardleymaps.com) using the WTG2 domain-specific language and renders them to SVG.

## What it does

When you ask Claude Code to create a Wardley Map, this skill:

1. Generates a valid `.wtg2` file from your description
2. Renders it to SVG using `wtg2svg`
3. Opens the result for you to review
4. Iterates on your feedback

## Install

### Via APM (recommended)

Add to your project's `apm.yml`:

```yaml
skills:
  - git: github.com/owulveryck/wardleyToGo
    path: skills/wtg2
```

Then run:

```bash
apm install
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
```

### Via Claude Code Plugin

```
/plugin marketplace add owulveryck/wardleyToGo
/plugin install wtg2@wardleyToGo
```

Then install the renderer:

```bash
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
```

### Via install script

```bash
./install.sh
```

This copies the skill to `~/.claude/skills/wtg2/` (available across all projects) and installs the `wtg2svg` renderer.

### Manual install

```bash
# Copy the skill (from repo root)
mkdir -p ~/.claude/skills/wtg2/references
cp skills/wtg2/SKILL.md ~/.claude/skills/wtg2/
cp skills/wtg2/references/*.md ~/.claude/skills/wtg2/references/

# Install the renderer
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
```

## Usage

The skill triggers automatically when you mention Wardley Maps, strategic mapping, value chains, or evolution in Claude Code. You can also invoke it explicitly:

```
/wtg2 a map of our e-commerce platform
```

## Files

```
skills/wtg2/                      Skill source (plugin layout)
├── SKILL.md                      Main skill (loaded when triggered)
└── references/
    ├── syntax-reference.md       Detailed syntax rules (loaded on demand)
    └── strategic-concepts.md     Strategic theory: gameplays, inertia, etc.

.claude-plugin/                   Plugin metadata
├── plugin.json
└── marketplace.json

wtg2/claude-skill/                Install helpers
├── install.sh                    Manual install script
└── README.md                     This file
```

## Requirements

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)
- [Go](https://go.dev) (for installing `wtg2svg`)

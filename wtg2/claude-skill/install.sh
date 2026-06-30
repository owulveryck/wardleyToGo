#!/bin/sh
set -e

SKILL_NAME="wtg2"
TARGET_DIR="$HOME/.claude/skills/$SKILL_NAME"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
SKILL_SRC="$REPO_ROOT/skills/wtg2"

echo "Installing WTG2 skill for Claude Code..."

if [ ! -f "$SKILL_SRC/SKILL.md" ]; then
    echo "Error: SKILL.md not found at $SKILL_SRC"
    echo "Make sure you're running this from the wardleyToGo repository."
    exit 1
fi

# Create target directory
mkdir -p "$TARGET_DIR/references"

# Copy skill files
cp "$SKILL_SRC/SKILL.md" "$TARGET_DIR/SKILL.md"
cp "$SKILL_SRC/references/syntax-reference.md" "$TARGET_DIR/references/syntax-reference.md"
cp "$SKILL_SRC/references/strategic-concepts.md" "$TARGET_DIR/references/strategic-concepts.md"

echo "Skill installed to $TARGET_DIR"

# Install wtg2svg
if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go is required to install wtg2svg."
    echo "Install Go from https://go.dev and re-run this script."
    exit 1
fi

echo "Installing wtg2svg..."
go install github.com/owulveryck/wardleyToGo/cmd/wtg2svg@latest
echo "wtg2svg installed to $(go env GOPATH)/bin/wtg2svg"

echo ""
echo "Done. The skill is available in Claude Code as /wtg2"
echo "Try: claude '/wtg2 a map of a coffee shop supply chain'"

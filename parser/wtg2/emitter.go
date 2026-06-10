package wtg2

import (
	"fmt"
	"io"
	"strings"
)

// Emit writes a canonical WTG2 text representation of the Document to w.
// The output, when re-parsed, produces a semantically equivalent Document.
// Comments from the original source are not preserved (they are not in the AST).
func Emit(w io.Writer, doc *Document) error {
	e := &emitter{w: w}

	e.emitMeta(doc)
	e.emitStages(doc)
	e.emitNodes(doc)
	e.emitPipelines(doc)
	e.emitEdges(doc)
	e.emitGroups(doc)
	e.emitAnnotations(doc)
	e.emitSignals(doc)
	e.emitGameplays(doc)
	e.emitFocuses(doc)
	e.emitLegend(doc)

	return e.err
}

type emitter struct {
	w             io.Writer
	err           error
	wroteSection  bool
}

func (e *emitter) writef(format string, args ...any) {
	if e.err != nil {
		return
	}
	_, e.err = fmt.Fprintf(e.w, format, args...)
}

func (e *emitter) sectionBreak() {
	if e.wroteSection {
		e.writef("\n")
	}
	e.wroteSection = true
}

func (e *emitter) emitMeta(doc *Document) {
	wrote := false
	if doc.Title != "" {
		e.writef("title: %s\n", doc.Title)
		wrote = true
	}
	if doc.Date != "" {
		e.writef("date: %s\n", doc.Date)
		wrote = true
	}
	if doc.Author != "" {
		e.writef("author: %s\n", doc.Author)
		wrote = true
	}
	if doc.Scope != "" {
		e.writef("scope: %s\n", doc.Scope)
		wrote = true
	}
	if doc.Question != "" {
		e.writef("question: %s\n", doc.Question)
		wrote = true
	}
	if doc.Doctrine != "" {
		e.writef("doctrine: %s\n", doc.Doctrine)
		wrote = true
	}
	if wrote {
		e.wroteSection = true
	}
}

func (e *emitter) emitStages(doc *Document) {
	hasStages := false
	for _, s := range doc.Stages {
		if s != "" {
			hasStages = true
			break
		}
	}
	if !hasStages {
		return
	}
	e.sectionBreak()
	e.writef("stages: %s, %s, %s, %s\n", doc.Stages[0], doc.Stages[1], doc.Stages[2], doc.Stages[3])
}

func (e *emitter) emitNodes(doc *Document) {
	if len(doc.Nodes) == 0 {
		return
	}
	e.sectionBreak()
	for _, n := range doc.Nodes {
		e.emitNode(n)
	}
}

func (e *emitter) emitNode(n *NodeDecl) {
	keyword := nodeKindKeyword(n.Kind)
	hasEvolution := n.Evolution != ""
	needsBlock := n.Color != "" || n.Cost != "" || n.Note != "" || n.Asset != ""
	if !hasEvolution && (n.Type != "" || n.Visibility >= 0) {
		needsBlock = true
	}

	e.writef("%s %s", keyword, n.Name)

	if hasEvolution {
		e.writef(" : %s", formatEvolution(n))
		if !needsBlock {
			if n.Type != "" {
				e.writef(" (%s)", n.Type)
			}
			if n.Visibility >= 0 {
				e.writef(" @%s", formatVisibility(n.Visibility))
			}
		}
	}

	if needsBlock {
		e.writef(" {\n")
		if hasEvolution && n.Type != "" {
			e.writef("  type: %s\n", n.Type)
		} else if !hasEvolution && n.Type != "" {
			e.writef("  evolution: %s\n", formatEvolution(n))
			e.writef("  type: %s\n", n.Type)
		}
		if !hasEvolution && n.Type == "" && n.Evolution == "" {
			// no evolution to emit
		}
		if n.Asset != "" {
			e.writef("  asset: %s\n", n.Asset)
		}
		if n.Color != "" {
			e.writef("  color: %s\n", n.Color)
		}
		if n.Visibility >= 0 {
			e.writef("  visibility: %s\n", formatVisibility(n.Visibility))
		}
		if n.Cost != "" {
			e.writef("  cost: %s\n", n.Cost)
		}
		if n.Note != "" {
			e.writef("  note: %s\n", n.Note)
		}
		e.writef("}\n")
	} else {
		e.writef("\n")
	}
}

func nodeKindKeyword(k NodeKind) string {
	switch k {
	case KindAnchor:
		return "anchor"
	case KindSubmap:
		return "submap"
	default:
		return "component"
	}
}

func formatEvolution(n *NodeDecl) string {
	var sb strings.Builder
	sb.WriteString(n.Evolution)

	if n.Inertia > 0 {
		sb.WriteByte(' ')
		for i := 0; i < n.Inertia; i++ {
			sb.WriteByte('!')
		}
		if len(n.InertiaKinds) > 0 {
			sb.WriteByte('(')
			sb.WriteString(strings.Join(n.InertiaKinds, ","))
			sb.WriteByte(')')
		}
	}

	if n.EvolvedTo != "" {
		sb.WriteString(" >> ")
		sb.WriteString(n.EvolvedTo)
	}

	return sb.String()
}

func formatVisibility(v float64) string {
	s := fmt.Sprintf("%g", v)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

func (e *emitter) emitPipelines(doc *Document) {
	if len(doc.Pipelines) == 0 {
		return
	}
	e.sectionBreak()
	for _, p := range doc.Pipelines {
		e.writef("pipeline %s {\n", p.Name)
		for _, m := range p.Members {
			e.writef("  %s : %s\n", m.Name, m.Position)
		}
		e.writef("}\n")
	}
}

func (e *emitter) emitEdges(doc *Document) {
	if len(doc.Edges) == 0 {
		return
	}
	e.sectionBreak()
	for _, edge := range doc.Edges {
		e.writef("%s %s %s\n", edge.From, formatLink(edge), edge.To)
	}
}

func formatLink(edge *EdgeDecl) string {
	if edge.Label != "" {
		return fmt.Sprintf("-[%s]->", edge.Label)
	}
	if edge.Bidirectional {
		return "<->"
	}
	return "->"
}

func (e *emitter) emitGroups(doc *Document) {
	if len(doc.Groups) == 0 {
		return
	}
	e.sectionBreak()
	for _, g := range doc.Groups {
		e.writef("group %s {\n", g.Name)
		if g.Color != "" {
			e.writef("  color: %s\n", g.Color)
		}
		if g.Team != "" {
			e.writef("  team: %s\n", g.Team)
		}
		for _, m := range g.Members {
			e.writef("  %s\n", m)
		}
		e.writef("}\n")
	}
}

func (e *emitter) emitAnnotations(doc *Document) {
	if len(doc.Annotations) == 0 {
		return
	}
	e.sectionBreak()
	for _, a := range doc.Annotations {
		e.writef("%s \"%s\" on %s\n", a.Kind, a.Text, a.Target)
	}
}

func (e *emitter) emitSignals(doc *Document) {
	if len(doc.Signals) == 0 {
		return
	}
	e.sectionBreak()
	for _, s := range doc.Signals {
		e.writef("signal %s on %s\n", s.Type, s.Target)
	}
}

func (e *emitter) emitGameplays(doc *Document) {
	if len(doc.Gameplays) == 0 {
		return
	}
	e.sectionBreak()
	for _, g := range doc.Gameplays {
		if g.Text != "" {
			e.writef("gameplay %s \"%s\" on %s\n", g.Type, g.Text, g.Target)
		} else {
			e.writef("gameplay %s on %s\n", g.Type, g.Target)
		}
	}
}

func (e *emitter) emitFocuses(doc *Document) {
	if len(doc.Focuses) == 0 {
		return
	}
	e.sectionBreak()
	for _, f := range doc.Focuses {
		e.writef("focus %s\n", f.Target)
	}
}

func (e *emitter) emitLegend(doc *Document) {
	if !doc.Legend {
		return
	}
	e.sectionBreak()
	e.writef("legend\n")
}

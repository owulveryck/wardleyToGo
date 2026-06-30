//go:build js && wasm

package main

import (
	"fmt"
	"strings"

	"github.com/owulveryck/wardleyToGo/v2/parser/wtg2"
)

func docToWTG2(doc *wtg2.Document) string {
	var b strings.Builder

	if doc.Title != "" {
		fmt.Fprintf(&b, "title: %s\n", doc.Title)
	}
	if doc.Date != "" {
		fmt.Fprintf(&b, "date: %s\n", doc.Date)
	}
	if doc.Author != "" {
		fmt.Fprintf(&b, "author: %s\n", doc.Author)
	}
	if doc.Scope != "" {
		fmt.Fprintf(&b, "scope: %s\n", doc.Scope)
	}
	if doc.Question != "" {
		fmt.Fprintf(&b, "question: \"%s\"\n", doc.Question)
	}
	if doc.Doctrine != "" {
		fmt.Fprintf(&b, "doctrine: %s\n", doc.Doctrine)
	}

	hasStages := false
	for _, s := range doc.Stages {
		if s != "" {
			hasStages = true
			break
		}
	}
	if hasStages {
		fmt.Fprintf(&b, "stages: %s\n", strings.Join(doc.Stages[:], ", "))
	}

	if doc.Legend {
		b.WriteString("legend\n")
	}

	b.WriteString("\n")

	pipelineSet := make(map[string]*wtg2.PipelineDecl)
	for _, p := range doc.Pipelines {
		pipelineSet[p.Name] = p
	}

	for _, n := range doc.Nodes {
		writeNode(&b, n)
	}
	if len(doc.Nodes) > 0 {
		b.WriteString("\n")
	}

	for _, p := range doc.Pipelines {
		fmt.Fprintf(&b, "pipeline %s {\n", p.Name)
		for _, m := range p.Members {
			if m.Position != "" {
				fmt.Fprintf(&b, "  %s : %s\n", m.Name, m.Position)
			} else {
				fmt.Fprintf(&b, "  %s\n", m.Name)
			}
		}
		b.WriteString("}\n\n")
	}

	for _, e := range doc.Edges {
		writeEdge(&b, e)
	}
	if len(doc.Edges) > 0 {
		b.WriteString("\n")
	}

	for _, g := range doc.Groups {
		writeGroup(&b, g)
	}

	for _, a := range doc.Annotations {
		fmt.Fprintf(&b, "%s \"%s\" on %s\n", a.Kind, a.Text, a.Target)
	}

	for _, s := range doc.Signals {
		fmt.Fprintf(&b, "signal %s on %s\n", s.Type, s.Target)
	}

	for _, g := range doc.Gameplays {
		if g.Text != "" {
			fmt.Fprintf(&b, "gameplay %s \"%s\" on %s\n", g.Type, g.Text, g.Target)
		} else {
			fmt.Fprintf(&b, "gameplay %s on %s\n", g.Type, g.Target)
		}
	}

	for _, f := range doc.Focuses {
		fmt.Fprintf(&b, "focus %s\n", f.Target)
	}

	return b.String()
}

func writeNode(b *strings.Builder, n *wtg2.NodeDecl) {
	switch n.Kind {
	case wtg2.KindAnchor:
		b.WriteString("anchor ")
	case wtg2.KindSubmap:
		b.WriteString("submap ")
	default:
	}

	b.WriteString(n.Name)

	if n.Evolution != "" {
		fmt.Fprintf(b, " : %s", n.Evolution)
	}

	if n.Inertia > 0 {
		b.WriteString(" ")
		for i := 0; i < n.Inertia; i++ {
			b.WriteString("!")
		}
		if len(n.InertiaKinds) > 0 {
			fmt.Fprintf(b, "(%s)", strings.Join(n.InertiaKinds, ","))
		}
	}

	if n.EvolvedTo != "" {
		fmt.Fprintf(b, " >> %s", n.EvolvedTo)
	}

	hasBlock := n.Type != "" || n.Color != "" || n.Asset != "" ||
		n.Cost != "" || n.Note != "" || n.Visibility != -1

	if hasBlock {
		b.WriteString(" {\n")
		if n.Type != "" {
			fmt.Fprintf(b, "  type: %s\n", n.Type)
		}
		if n.Color != "" {
			fmt.Fprintf(b, "  color: %s\n", n.Color)
		}
		if n.Asset != "" {
			fmt.Fprintf(b, "  asset: %s\n", n.Asset)
		}
		if n.Visibility != -1 {
			fmt.Fprintf(b, "  visibility: %g\n", n.Visibility)
		}
		if n.Cost != "" {
			fmt.Fprintf(b, "  cost: \"%s\"\n", n.Cost)
		}
		if n.Note != "" {
			fmt.Fprintf(b, "  note: \"%s\"\n", n.Note)
		}
		b.WriteString("}")
	} else if n.Type == "" && n.Kind == wtg2.KindComponent {
		// No block needed, but check for shorthand type via parentheses
	}

	b.WriteString("\n")
}

func writeEdge(b *strings.Builder, e *wtg2.EdgeDecl) {
	b.WriteString(e.From)
	if e.Bidirectional {
		if e.Label != "" {
			fmt.Fprintf(b, " <-[%s]-> ", e.Label)
		} else {
			b.WriteString(" <-> ")
		}
	} else {
		if e.Label != "" {
			fmt.Fprintf(b, " -[%s]-> ", e.Label)
		} else {
			b.WriteString(" -> ")
		}
	}
	b.WriteString(e.To)
	b.WriteString("\n")
}

func writeGroup(b *strings.Builder, g *wtg2.GroupDecl) {
	fmt.Fprintf(b, "group %s {\n", g.Name)
	for _, m := range g.Members {
		fmt.Fprintf(b, "  %s\n", m)
	}
	if g.Color != "" {
		fmt.Fprintf(b, "  color: %s\n", g.Color)
	}
	if g.Team != "" {
		fmt.Fprintf(b, "  team: %s\n", g.Team)
	}
	b.WriteString("}\n\n")
}

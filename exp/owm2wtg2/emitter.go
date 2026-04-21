package owm2wtg2

import (
	"fmt"
	"io"
	"strings"
)

// Emit writes the OWMDocument as WTG2 format to w.
func Emit(doc *OWMDocument, w io.Writer) error {
	e := &emitter{w: w, evolveMap: buildEvolveMap(doc)}

	if err := e.emitMetadata(doc); err != nil {
		return err
	}
	if err := e.emitAnchors(doc); err != nil {
		return err
	}
	if err := e.emitComponents(doc); err != nil {
		return err
	}
	if err := e.emitSubmaps(doc); err != nil {
		return err
	}
	if err := e.emitPipelines(doc); err != nil {
		return err
	}
	if err := e.emitEdges(doc); err != nil {
		return err
	}
	if err := e.emitAreas(doc); err != nil {
		return err
	}
	if err := e.emitNotes(doc); err != nil {
		return err
	}
	if err := e.emitSignals(doc); err != nil {
		return err
	}
	if err := e.emitStyle(doc); err != nil {
		return err
	}

	return nil
}

type emitter struct {
	w          io.Writer
	evolveMap  map[string]*OWMEvolve
	wroteBlock bool
}

func buildEvolveMap(doc *OWMDocument) map[string]*OWMEvolve {
	m := make(map[string]*OWMEvolve, len(doc.Evolves))
	for _, ev := range doc.Evolves {
		m[ev.Name] = ev
	}
	return m
}

func (e *emitter) section() error {
	if e.wroteBlock {
		_, err := fmt.Fprintln(e.w)
		return err
	}
	return nil
}

func (e *emitter) emitMetadata(doc *OWMDocument) error {
	if doc.Title != "" {
		if _, err := fmt.Fprintf(e.w, "title: %s\n", doc.Title); err != nil {
			return err
		}
		e.wroteBlock = true
	}
	if doc.Evolution != [4]string{} {
		if _, err := fmt.Fprintf(e.w, "stages: %s, %s, %s, %s\n",
			doc.Evolution[0], doc.Evolution[1], doc.Evolution[2], doc.Evolution[3]); err != nil {
			return err
		}
		e.wroteBlock = true
	}
	return nil
}

func (e *emitter) emitAnchors(doc *OWMDocument) error {
	if len(doc.Anchors) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, a := range doc.Anchors {
		evo := MaturityToEvolution(a.Maturity)
		if _, err := fmt.Fprintf(e.w, "anchor %s : %s @%.2f\n", a.Name, evo, a.Visibility); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitComponents(doc *OWMDocument) error {
	if len(doc.Components) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, c := range doc.Components {
		if err := e.emitComponent(c); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitComponent(c *OWMComponent) error {
	var b strings.Builder
	b.WriteString(c.Name)
	b.WriteString(" : ")

	evo := MaturityToEvolution(c.Maturity)
	b.WriteString(evo)

	if c.Inertia {
		b.WriteString(" !!")
	}

	ev := e.evolveMap[c.Name]
	if ev != nil {
		b.WriteString(" >> ")
		b.WriteString(MaturityToEvolution(ev.Maturity))
	}

	compType := c.Type
	if ev != nil && ev.Type != "" {
		compType = ev.Type
	}
	if compType != "" && compType != "market" {
		b.WriteString(" (")
		b.WriteString(compType)
		b.WriteString(")")
	}

	fmt.Fprintf(&b, " @%.2f", c.Visibility)

	if _, err := fmt.Fprintln(e.w, b.String()); err != nil {
		return err
	}

	if ev != nil && ev.NewName != "" {
		if _, err := fmt.Fprintf(e.w, "// OWM: evolved name %q\n", ev.NewName); err != nil {
			return err
		}
	}

	if c.Type == "market" {
		if _, err := fmt.Fprintf(e.w, "// OWM: market component\n"); err != nil {
			return err
		}
	}

	return nil
}

func (e *emitter) emitSubmaps(doc *OWMDocument) error {
	if len(doc.Submaps) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, sm := range doc.Submaps {
		evo := MaturityToEvolution(sm.Maturity)
		if _, err := fmt.Fprintf(e.w, "submap %s : %s @%.2f\n", sm.Name, evo, sm.Visibility); err != nil {
			return err
		}
		if sm.URLName != "" {
			url := doc.URLs[sm.URLName]
			if url != "" {
				if _, err := fmt.Fprintf(e.w, "// OWM URL: %s\n", url); err != nil {
					return err
				}
			}
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitPipelines(doc *OWMDocument) error {
	if len(doc.Pipelines) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, p := range doc.Pipelines {
		members := e.findPipelineMembers(doc, p)
		if _, err := fmt.Fprintf(e.w, "pipeline %s {\n", p.Name); err != nil {
			return err
		}
		if len(members) > 0 {
			for _, m := range members {
				evo := MaturityToEvolution(m.Maturity)
				if _, err := fmt.Fprintf(e.w, "  %s : %s\n", m.Name, evo); err != nil {
					return err
				}
			}
		} else if p.HasRange {
			if _, err := fmt.Fprintf(e.w, "  // OWM pipeline range: [%.2f, %.2f]\n", p.StartMat, p.EndMat); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(e.w, "}"); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) findPipelineMembers(doc *OWMDocument, p *OWMPipeline) []*OWMComponent {
	if !p.HasRange {
		return nil
	}

	connected := make(map[string]bool)
	for _, edge := range doc.Edges {
		if edge.From == p.Name {
			connected[edge.To] = true
		}
		if edge.To == p.Name {
			connected[edge.From] = true
		}
	}

	var members []*OWMComponent
	for _, c := range doc.Components {
		if !connected[c.Name] {
			continue
		}
		if c.Maturity >= p.StartMat && c.Maturity <= p.EndMat {
			members = append(members, c)
		}
	}
	return members
}

func (e *emitter) emitEdges(doc *OWMDocument) error {
	if len(doc.Edges) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, edge := range doc.Edges {
		switch edge.FlowType {
		case "regular", "future":
			if _, err := fmt.Fprintf(e.w, "%s -> %s\n", edge.From, edge.To); err != nil {
				return err
			}
		case "bidirectional":
			if _, err := fmt.Fprintf(e.w, "%s <-> %s\n", edge.From, edge.To); err != nil {
				return err
			}
		case "labeled":
			if _, err := fmt.Fprintf(e.w, "%s -[%s]-> %s\n", edge.From, edge.Label, edge.To); err != nil {
				return err
			}
		case "past":
			if _, err := fmt.Fprintf(e.w, "%s -> %s\n", edge.From, edge.To); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(e.w, "// OWM: past flow\n"); err != nil {
				return err
			}
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitAreas(doc *OWMDocument) error {
	if len(doc.Areas) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	teamMap := map[string]string{
		"pioneers":     "pioneer",
		"settlers":     "settler",
		"townplanners": "town-planner",
	}
	for _, a := range doc.Areas {
		team := teamMap[a.Kind]
		groupName := strings.ToUpper(a.Kind[:1]) + a.Kind[1:]
		if _, err := fmt.Fprintf(e.w, "group %s {\n", groupName); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(e.w, "  // OWM area: [%.2f, %.2f, %.2f, %.2f]\n",
			a.Vis1, a.Mat1, a.Vis2, a.Mat2); err != nil {
			return err
		}
		if team != "" {
			if _, err := fmt.Fprintf(e.w, "  team: %s\n", team); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(e.w, "}"); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitNotes(doc *OWMDocument) error {
	if len(doc.Notes) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, n := range doc.Notes {
		if _, err := fmt.Fprintf(e.w, "// OWM note at [%.2f, %.2f]: %s\n",
			n.Visibility, n.Maturity, n.Text); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitSignals(doc *OWMDocument) error {
	if len(doc.Signals) == 0 {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	for _, s := range doc.Signals {
		signalType := "accelerating"
		if s.Kind == "deaccelerator" {
			signalType = "declining"
		}
		if _, err := fmt.Fprintf(e.w, "signal %s on %s\n", signalType, s.Name); err != nil {
			return err
		}
	}
	e.wroteBlock = true
	return nil
}

func (e *emitter) emitStyle(doc *OWMDocument) error {
	if doc.Style == "" {
		return nil
	}
	if err := e.section(); err != nil {
		return err
	}
	_, err := fmt.Fprintf(e.w, "// OWM style: %s\n", doc.Style)
	if err != nil {
		return err
	}
	e.wroteBlock = true
	return nil
}

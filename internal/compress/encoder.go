package compress

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

// Compress encodes a WTG2 Document into a compressed binary stream.
func Compress(doc *wtg2.Document, w io.Writer) error {
	// Write header (raw bytes, not arithmetic-coded)
	if _, err := w.Write([]byte{magic0, magic1, version, 0}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	e := &encoder{
		ac:  NewEncoder(w),
		sym: newSymTable(),
	}

	if err := e.encode(doc); err != nil {
		return err
	}

	return e.ac.Finish()
}

// CompressBytes encodes a WTG2 Document into compressed binary bytes.
func CompressBytes(doc *wtg2.Document) ([]byte, error) {
	var buf bytes.Buffer
	if err := Compress(doc, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type encoder struct {
	ac  *Encoder
	sym *symTable
}

func (e *encoder) encodeSym(sym int, dist *Distribution) error {
	return e.ac.Encode(sym, dist)
}

func (e *encoder) encode(doc *wtg2.Document) error {
	// Emit all statements in canonical order.
	// Each section: hasNext=true + stmtType + content.

	// Meta fields
	if err := e.encodeMeta(doc); err != nil {
		return err
	}

	// Stages
	if err := e.encodeStages(doc); err != nil {
		return err
	}

	// Nodes
	for _, n := range doc.Nodes {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtNode, distStmtType); err != nil {
			return err
		}
		if err := e.encodeNode(n); err != nil {
			return err
		}
	}

	// Pipelines
	for _, p := range doc.Pipelines {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtPipeline, distStmtType); err != nil {
			return err
		}
		if err := e.encodePipeline(p); err != nil {
			return err
		}
	}

	// Edges
	for _, edge := range doc.Edges {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtEdge, distStmtType); err != nil {
			return err
		}
		if err := e.encodeEdge(edge); err != nil {
			return err
		}
	}

	// Groups
	for _, g := range doc.Groups {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtGroup, distStmtType); err != nil {
			return err
		}
		if err := e.encodeGroup(g); err != nil {
			return err
		}
	}

	// Annotations
	for _, a := range doc.Annotations {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtAnnotation, distStmtType); err != nil {
			return err
		}
		if err := e.encodeAnnotation(a); err != nil {
			return err
		}
	}

	// Signals
	for _, s := range doc.Signals {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtSignal, distStmtType); err != nil {
			return err
		}
		if err := e.encodeSignal(s); err != nil {
			return err
		}
	}

	// Gameplays
	for _, g := range doc.Gameplays {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtGameplay, distStmtType); err != nil {
			return err
		}
		if err := e.encodeGameplay(g); err != nil {
			return err
		}
	}

	// Focuses
	for _, f := range doc.Focuses {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtFocus, distStmtType); err != nil {
			return err
		}
		if err := e.encodeIdentifier(f.Target); err != nil {
			return err
		}
	}

	// Legend
	if doc.Legend {
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtLegend, distStmtType); err != nil {
			return err
		}
	}

	// End of document
	return e.encodeSym(boolNo, distHasNext)
}

func (e *encoder) encodeMeta(doc *wtg2.Document) error {
	metas := []struct {
		field int
		value string
	}{
		{metaTitle, doc.Title},
		{metaDate, doc.Date},
		{metaAuthor, doc.Author},
		{metaScope, doc.Scope},
		{metaQuestion, doc.Question},
		{metaDoctrine, doc.Doctrine},
	}

	for _, m := range metas {
		if m.value == "" {
			continue
		}
		if err := e.encodeSym(boolYes, distHasNext); err != nil {
			return err
		}
		if err := e.encodeSym(stmtMeta, distStmtType); err != nil {
			return err
		}
		if err := e.encodeSym(m.field, distMetaField); err != nil {
			return err
		}
		if m.field == metaDoctrine {
			if err := e.encodeSym(doctrineValueIndex(m.value), distDoctrineValue); err != nil {
				return err
			}
		} else {
			if err := e.encodeText(m.value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *encoder) encodeStages(doc *wtg2.Document) error {
	hasStages := false
	for _, s := range doc.Stages {
		if s != "" {
			hasStages = true
			break
		}
	}
	if !hasStages {
		return nil
	}

	if err := e.encodeSym(boolYes, distHasNext); err != nil {
		return err
	}
	if err := e.encodeSym(stmtStages, distStmtType); err != nil {
		return err
	}
	for _, s := range doc.Stages {
		if err := e.encodeText(s); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeNode(n *wtg2.NodeDecl) error {
	// Node kind
	kind := nodeComponent
	switch n.Kind {
	case wtg2.KindAnchor:
		kind = nodeAnchor
	case wtg2.KindSubmap:
		kind = nodeSubmap
	}
	if err := e.encodeSym(kind, distNodeKind); err != nil {
		return err
	}

	// Identifier
	if err := e.encodeIdentifier(n.Name); err != nil {
		return err
	}

	// Determine format
	hasEvolution := n.Evolution != ""
	needsBlock := n.Color != "" || n.Cost != "" || n.Note != "" || n.Asset != ""
	if !hasEvolution && (n.Type != "" || n.Visibility >= 0) {
		needsBlock = true
	}

	var fmt int
	switch {
	case needsBlock:
		fmt = fmtBlock
	case hasEvolution:
		fmt = fmtShorthand
	default:
		fmt = fmtBare
	}
	if err := e.encodeSym(fmt, distNodeFormat); err != nil {
		return err
	}

	if fmt == fmtBare {
		return nil
	}

	if fmt == fmtShorthand {
		return e.encodeShorthand(n)
	}

	// Block format
	return e.encodeBlock(n)
}

func (e *encoder) encodeShorthand(n *wtg2.NodeDecl) error {
	// Determine variant
	hasType := n.Type != ""
	hasVis := n.Visibility >= 0
	var variant int
	switch {
	case hasType && hasVis:
		variant = shEvoTypeVis
	case hasType:
		variant = shEvoType
	case hasVis:
		variant = shEvoVis
	default:
		variant = shEvoOnly
	}
	if err := e.encodeSym(variant, distShorthandVar); err != nil {
		return err
	}

	if err := e.encodeEvolution(n); err != nil {
		return err
	}

	if hasType {
		if err := e.encodeSym(typeValueIndex(n.Type), distTypeValue); err != nil {
			return err
		}
	}

	if hasVis {
		if err := e.encodeVisibility(n.Visibility); err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) encodeBlock(n *wtg2.NodeDecl) error {
	// Evolution (in shorthand part before block)
	hasEvolution := n.Evolution != ""
	if err := e.encodeSym(boolVal(hasEvolution), distBool); err != nil {
		return err
	}
	if hasEvolution {
		if err := e.encodeEvolution(n); err != nil {
			return err
		}
	}

	// Block fields: each has a presence bit followed by value
	fields := []struct {
		present bool
		encode  func() error
	}{
		{n.Type != "", func() error { return e.encodeSym(typeValueIndex(n.Type), distTypeValue) }},
		{n.Asset != "", func() error { return e.encodeSym(assetValueIndex(n.Asset), distAssetValue) }},
		{n.Color != "", func() error { return e.encodeText(n.Color) }},
		{n.Visibility >= 0, func() error { return e.encodeVisibility(n.Visibility) }},
		{n.Cost != "", func() error { return e.encodeText(n.Cost) }},
		{n.Note != "", func() error { return e.encodeText(n.Note) }},
	}

	for _, f := range fields {
		if err := e.encodeSym(boolVal(f.present), distBool); err != nil {
			return err
		}
		if f.present {
			if err := f.encode(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *encoder) encodeEvolution(n *wtg2.NodeDecl) error {
	// Determine form
	hasMove := n.EvolvedTo != ""
	hasInertia := n.Inertia > 0
	var form int
	switch {
	case hasInertia && hasMove:
		form = evoInertiaMove
	case hasMove:
		form = evoMove
	default:
		form = evoPositionOnly
	}
	if err := e.encodeSym(form, distEvoForm); err != nil {
		return err
	}

	// Encode start position
	if err := e.encodePosition(n.Evolution); err != nil {
		return err
	}

	// Encode inertia if present
	if hasInertia {
		if err := e.encodeSym(n.Inertia-1, distInertiaLevel); err != nil {
			return err
		}
		// Has qualified kinds?
		hasKinds := len(n.InertiaKinds) > 0
		if err := e.encodeSym(boolVal(hasKinds), distBool); err != nil {
			return err
		}
		if hasKinds {
			if err := e.encodeLength(len(n.InertiaKinds)); err != nil {
				return err
			}
			for _, k := range n.InertiaKinds {
				if err := e.encodeSym(inertiaKindIndex(k), distInertiaKind); err != nil {
					return err
				}
			}
		}
	}

	// Encode target position if move
	if hasMove {
		if err := e.encodePosition(n.EvolvedTo); err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) encodePosition(pos string) error {
	// Parse roman numeral and optional decimal
	roman, decimal, hasDecimal := parsePositionParts(pos)

	if err := e.encodeSym(roman, distRoman); err != nil {
		return err
	}
	if err := e.encodeSym(boolVal(hasDecimal), distHasDecimal); err != nil {
		return err
	}
	if hasDecimal {
		if err := e.encodeSym(decimal, distDecimalDigit); err != nil {
			return err
		}
	}
	return nil
}

// parsePositionParts parses "II.7" → (romanII, 7, true) or "III" → (romanIII, 0, false).
func parsePositionParts(pos string) (roman, decimal int, hasDecimal bool) {
	dotIdx := strings.Index(pos, ".")
	var romanStr string
	if dotIdx >= 0 {
		romanStr = pos[:dotIdx]
		hasDecimal = true
		if dotIdx+1 < len(pos) {
			decimal = int(pos[dotIdx+1] - '0')
		}
	} else {
		romanStr = pos
	}

	switch strings.TrimSpace(romanStr) {
	case "I":
		roman = romanI
	case "II":
		roman = romanII
	case "III":
		roman = romanIII
	case "IV":
		roman = romanIV
	}
	return
}

func (e *encoder) encodeVisibility(v float64) error {
	// Encode as a text string for simplicity
	s := formatVisFloat(v)
	return e.encodeText(s)
}

func formatVisFloat(v float64) string {
	s := fmt.Sprintf("%g", v)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

func (e *encoder) encodePipeline(p *wtg2.PipelineDecl) error {
	if err := e.encodeIdentifier(p.Name); err != nil {
		return err
	}
	if err := e.encodeLength(len(p.Members)); err != nil {
		return err
	}
	for _, m := range p.Members {
		if err := e.encodeIdentifier(m.Name); err != nil {
			return err
		}
		if err := e.encodePosition(m.Position); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeEdge(edge *wtg2.EdgeDecl) error {
	if err := e.encodeIdentifier(edge.From); err != nil {
		return err
	}

	// Link type
	var lt int
	switch {
	case edge.Label != "":
		lt = linkLabeled
	case edge.Bidirectional:
		lt = linkBidir
	default:
		lt = linkArrow
	}
	if err := e.encodeSym(lt, distLinkType); err != nil {
		return err
	}

	if edge.Label != "" {
		if err := e.encodeText(edge.Label); err != nil {
			return err
		}
	}

	return e.encodeIdentifier(edge.To)
}

func (e *encoder) encodeGroup(g *wtg2.GroupDecl) error {
	if err := e.encodeIdentifier(g.Name); err != nil {
		return err
	}

	// Color
	hasColor := g.Color != ""
	if err := e.encodeSym(boolVal(hasColor), distBool); err != nil {
		return err
	}
	if hasColor {
		if err := e.encodeText(g.Color); err != nil {
			return err
		}
	}

	// Team
	hasTeam := g.Team != ""
	if err := e.encodeSym(boolVal(hasTeam), distBool); err != nil {
		return err
	}
	if hasTeam {
		if err := e.encodeSym(teamTypeIndex(g.Team), distTeamType); err != nil {
			return err
		}
	}

	// Members
	if err := e.encodeLength(len(g.Members)); err != nil {
		return err
	}
	for _, m := range g.Members {
		if err := e.encodeIdentifier(m); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeAnnotation(a *wtg2.AnnotationDecl) error {
	kind := annoNote
	if a.Kind == "warning" {
		kind = annoWarning
	}
	if err := e.encodeSym(kind, distAnnotationKind); err != nil {
		return err
	}
	if err := e.encodeText(a.Text); err != nil {
		return err
	}
	return e.encodeIdentifier(a.Target)
}

func (e *encoder) encodeSignal(s *wtg2.SignalDecl) error {
	if err := e.encodeSym(signalTypeIndex(s.Type), distSignalType); err != nil {
		return err
	}
	return e.encodeIdentifier(s.Target)
}

func (e *encoder) encodeGameplay(g *wtg2.GameplayDecl) error {
	if err := e.encodeSym(gameplayTypeIndex(g.Type), distGameplayType); err != nil {
		return err
	}

	hasText := g.Text != ""
	if err := e.encodeSym(boolVal(hasText), distBool); err != nil {
		return err
	}
	if hasText {
		if err := e.encodeText(g.Text); err != nil {
			return err
		}
	}

	return e.encodeIdentifier(g.Target)
}

// encodeIdentifier encodes an identifier with the symbol table.
func (e *encoder) encodeIdentifier(name string) error {
	if idx, ok := e.sym.lookup(name); ok {
		// Known identifier
		if err := e.encodeSym(1, distIdentKnown); err != nil { // 1 = known
			return err
		}
		return e.encodeSym(idx, Uniform(e.sym.len()))
	}

	// New identifier
	if err := e.encodeSym(0, distIdentKnown); err != nil { // 0 = new
		return err
	}
	e.sym.add(name)
	return e.encodeIDString(name)
}

// encodeIDString encodes identifier text byte-by-byte.
func (e *encoder) encodeIDString(s string) error {
	return e.encodeBytes(s)
}

// encodeText encodes free text byte-by-byte.
func (e *encoder) encodeText(s string) error {
	return e.encodeBytes(s)
}

// encodeBytes encodes a string byte-by-byte using the byte distribution.
func (e *encoder) encodeBytes(s string) error {
	if err := e.encodeLength(len(s)); err != nil {
		return err
	}
	for i := range len(s) {
		if err := e.encodeSym(int(s[i]), distByte); err != nil {
			return err
		}
	}
	return nil
}

// encodeLength encodes a non-negative integer using a unary/binary hybrid.
// Lengths 0-15: 1 bit (0) + 4 bits value
// Lengths 16-271: 1 bit (1) + 8 bits value (offset by 16)
func (e *encoder) encodeLength(n int) error {
	if n < 16 {
		if err := e.encodeSym(0, distBool); err != nil { // short
			return err
		}
		return e.encodeSym(n, Uniform(16))
	}
	if err := e.encodeSym(1, distBool); err != nil { // long
		return err
	}
	return e.encodeSym(n-16, Uniform(256))
}

func boolVal(b bool) int {
	if b {
		return boolYes
	}
	return boolNo
}

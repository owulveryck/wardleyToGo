package compress

import (
	"bytes"
	"fmt"
	"io"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

// Decompress reads a compressed binary stream and returns the WTG2 Document.
func Decompress(r io.ByteReader) (*wtg2.Document, error) {
	// Read header
	var hdr [4]byte
	for i := range 4 {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("reading header: %w", err)
		}
		hdr[i] = b
	}
	if hdr[0] != magic0 || hdr[1] != magic1 {
		return nil, fmt.Errorf("invalid magic bytes: %x %x", hdr[0], hdr[1])
	}
	if hdr[2] != version {
		return nil, fmt.Errorf("unsupported version: %d", hdr[2])
	}

	d := &decoderState{
		ac:  NewDecoder(r),
		sym: newSymTable(),
	}

	return d.decode()
}

// DecompressBytes reads compressed binary bytes and returns the WTG2 Document.
func DecompressBytes(data []byte) (*wtg2.Document, error) {
	return Decompress(bytes.NewReader(data))
}

type decoderState struct {
	ac  *Decoder
	sym *symTable
}

func (d *decoderState) decodeSym(dist *Distribution) (int, error) {
	return d.ac.Decode(dist)
}

func (d *decoderState) decode() (*wtg2.Document, error) {
	doc := &wtg2.Document{
		Stages: [4]string{"", "", "", ""},
	}

	for {
		hasNext, err := d.decodeSym(distHasNext)
		if err != nil {
			return nil, err
		}
		if hasNext == boolNo {
			break
		}

		stmtType, err := d.decodeSym(distStmtType)
		if err != nil {
			return nil, err
		}

		switch stmtType {
		case stmtMeta:
			if err := d.decodeMeta(doc); err != nil {
				return nil, err
			}
		case stmtStages:
			if err := d.decodeStages(doc); err != nil {
				return nil, err
			}
		case stmtNode:
			n, err := d.decodeNode()
			if err != nil {
				return nil, err
			}
			doc.Nodes = append(doc.Nodes, n)
		case stmtPipeline:
			p, err := d.decodePipeline()
			if err != nil {
				return nil, err
			}
			doc.Pipelines = append(doc.Pipelines, p)
		case stmtEdge:
			edge, err := d.decodeEdge()
			if err != nil {
				return nil, err
			}
			doc.Edges = append(doc.Edges, edge)
		case stmtGroup:
			g, err := d.decodeGroup()
			if err != nil {
				return nil, err
			}
			doc.Groups = append(doc.Groups, g)
		case stmtAnnotation:
			a, err := d.decodeAnnotation()
			if err != nil {
				return nil, err
			}
			doc.Annotations = append(doc.Annotations, a)
		case stmtSignal:
			s, err := d.decodeSignal()
			if err != nil {
				return nil, err
			}
			doc.Signals = append(doc.Signals, s)
		case stmtGameplay:
			g, err := d.decodeGameplay()
			if err != nil {
				return nil, err
			}
			doc.Gameplays = append(doc.Gameplays, g)
		case stmtFocus:
			target, err := d.decodeIdentifier()
			if err != nil {
				return nil, err
			}
			doc.Focuses = append(doc.Focuses, &wtg2.FocusDecl{Target: target})
		case stmtLegend:
			doc.Legend = true
		default:
			return nil, fmt.Errorf("unknown statement type: %d", stmtType)
		}
	}

	return doc, nil
}

func (d *decoderState) decodeMeta(doc *wtg2.Document) error {
	field, err := d.decodeSym(distMetaField)
	if err != nil {
		return err
	}

	if field == metaDoctrine {
		idx, err := d.decodeSym(distDoctrineValue)
		if err != nil {
			return err
		}
		doc.Doctrine = doctrineValues[idx]
		return nil
	}

	text, err := d.decodeText()
	if err != nil {
		return err
	}

	switch field {
	case metaTitle:
		doc.Title = text
	case metaDate:
		doc.Date = text
	case metaAuthor:
		doc.Author = text
	case metaScope:
		doc.Scope = text
	case metaQuestion:
		doc.Question = text
	}
	return nil
}

func (d *decoderState) decodeStages(doc *wtg2.Document) error {
	for i := range 4 {
		s, err := d.decodeText()
		if err != nil {
			return err
		}
		doc.Stages[i] = s
	}
	return nil
}

func (d *decoderState) decodeNode() (*wtg2.NodeDecl, error) {
	n := &wtg2.NodeDecl{Visibility: -1}

	// Kind
	kind, err := d.decodeSym(distNodeKind)
	if err != nil {
		return nil, err
	}
	switch kind {
	case nodeComponent:
		n.Kind = wtg2.KindComponent
	case nodeAnchor:
		n.Kind = wtg2.KindAnchor
	case nodeSubmap:
		n.Kind = wtg2.KindSubmap
	}

	// Name
	name, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}
	n.Name = name

	// Format
	format, err := d.decodeSym(distNodeFormat)
	if err != nil {
		return nil, err
	}

	switch format {
	case fmtBare:
		return n, nil
	case fmtShorthand:
		if err := d.decodeShorthand(n); err != nil {
			return nil, err
		}
	case fmtBlock:
		if err := d.decodeBlock(n); err != nil {
			return nil, err
		}
	}

	return n, nil
}

func (d *decoderState) decodeShorthand(n *wtg2.NodeDecl) error {
	variant, err := d.decodeSym(distShorthandVar)
	if err != nil {
		return err
	}

	if err := d.decodeEvolution(n); err != nil {
		return err
	}

	if variant == shEvoType || variant == shEvoTypeVis {
		idx, err := d.decodeSym(distTypeValue)
		if err != nil {
			return err
		}
		n.Type = typeValues[idx]
	}

	if variant == shEvoVis || variant == shEvoTypeVis {
		vis, err := d.decodeVisibility()
		if err != nil {
			return err
		}
		n.Visibility = vis
	}

	return nil
}

func (d *decoderState) decodeBlock(n *wtg2.NodeDecl) error {
	// Has evolution?
	hasEvo, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasEvo == boolYes {
		if err := d.decodeEvolution(n); err != nil {
			return err
		}
	}

	// Type
	hasType, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasType == boolYes {
		idx, err := d.decodeSym(distTypeValue)
		if err != nil {
			return err
		}
		n.Type = typeValues[idx]
	}

	// Asset
	hasAsset, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasAsset == boolYes {
		idx, err := d.decodeSym(distAssetValue)
		if err != nil {
			return err
		}
		n.Asset = assetValues[idx]
	}

	// Color
	hasColor, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasColor == boolYes {
		color, err := d.decodeText()
		if err != nil {
			return err
		}
		n.Color = color
	}

	// Visibility
	hasVis, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasVis == boolYes {
		vis, err := d.decodeVisibility()
		if err != nil {
			return err
		}
		n.Visibility = vis
	}

	// Cost
	hasCost, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasCost == boolYes {
		cost, err := d.decodeText()
		if err != nil {
			return err
		}
		n.Cost = cost
	}

	// Note
	hasNote, err := d.decodeSym(distBool)
	if err != nil {
		return err
	}
	if hasNote == boolYes {
		note, err := d.decodeText()
		if err != nil {
			return err
		}
		n.Note = note
	}

	return nil
}

func (d *decoderState) decodeEvolution(n *wtg2.NodeDecl) error {
	form, err := d.decodeSym(distEvoForm)
	if err != nil {
		return err
	}

	pos, err := d.decodePosition()
	if err != nil {
		return err
	}
	n.Evolution = pos

	if form == evoInertiaMove {
		level, err := d.decodeSym(distInertiaLevel)
		if err != nil {
			return err
		}
		n.Inertia = level + 1

		hasKinds, err := d.decodeSym(distBool)
		if err != nil {
			return err
		}
		if hasKinds == boolYes {
			count, err := d.decodeLength()
			if err != nil {
				return err
			}
			n.InertiaKinds = make([]string, count)
			for i := range count {
				idx, err := d.decodeSym(distInertiaKind)
				if err != nil {
					return err
				}
				n.InertiaKinds[i] = inertiaKinds[idx]
			}
		}
	}

	if form == evoMove || form == evoInertiaMove {
		target, err := d.decodePosition()
		if err != nil {
			return err
		}
		n.EvolvedTo = target
	}

	return nil
}

func (d *decoderState) decodePosition() (string, error) {
	roman, err := d.decodeSym(distRoman)
	if err != nil {
		return "", err
	}

	hasDecimal, err := d.decodeSym(distHasDecimal)
	if err != nil {
		return "", err
	}

	romanStr := [romanCount]string{"I", "II", "III", "IV"}[roman]

	if hasDecimal == boolYes {
		digit, err := d.decodeSym(distDecimalDigit)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.%d", romanStr, digit), nil
	}

	return romanStr, nil
}

func (d *decoderState) decodeVisibility() (float64, error) {
	s, err := d.decodeText()
	if err != nil {
		return -1, err
	}
	var v float64
	if _, err := fmt.Sscanf(s, "%f", &v); err != nil {
		return -1, fmt.Errorf("parsing visibility %q: %w", s, err)
	}
	return v, nil
}

func (d *decoderState) decodePipeline() (*wtg2.PipelineDecl, error) {
	name, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}

	count, err := d.decodeLength()
	if err != nil {
		return nil, err
	}

	p := &wtg2.PipelineDecl{Name: name}
	p.Members = make([]*wtg2.PipelineMemberDecl, count)
	for i := range count {
		mName, err := d.decodeIdentifier()
		if err != nil {
			return nil, err
		}
		pos, err := d.decodePosition()
		if err != nil {
			return nil, err
		}
		p.Members[i] = &wtg2.PipelineMemberDecl{Name: mName, Position: pos}
	}
	return p, nil
}

func (d *decoderState) decodeEdge() (*wtg2.EdgeDecl, error) {
	from, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}

	lt, err := d.decodeSym(distLinkType)
	if err != nil {
		return nil, err
	}

	edge := &wtg2.EdgeDecl{From: from}

	if lt == linkLabeled {
		label, err := d.decodeText()
		if err != nil {
			return nil, err
		}
		edge.Label = label
	}
	if lt == linkBidir {
		edge.Bidirectional = true
	}

	to, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}
	edge.To = to

	return edge, nil
}

func (d *decoderState) decodeGroup() (*wtg2.GroupDecl, error) {
	name, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}

	g := &wtg2.GroupDecl{Name: name}

	hasColor, err := d.decodeSym(distBool)
	if err != nil {
		return nil, err
	}
	if hasColor == boolYes {
		color, err := d.decodeText()
		if err != nil {
			return nil, err
		}
		g.Color = color
	}

	hasTeam, err := d.decodeSym(distBool)
	if err != nil {
		return nil, err
	}
	if hasTeam == boolYes {
		idx, err := d.decodeSym(distTeamType)
		if err != nil {
			return nil, err
		}
		g.Team = teamTypes[idx]
	}

	count, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	g.Members = make([]string, count)
	for i := range count {
		m, err := d.decodeIdentifier()
		if err != nil {
			return nil, err
		}
		g.Members[i] = m
	}

	return g, nil
}

func (d *decoderState) decodeAnnotation() (*wtg2.AnnotationDecl, error) {
	kind, err := d.decodeSym(distAnnotationKind)
	if err != nil {
		return nil, err
	}

	text, err := d.decodeText()
	if err != nil {
		return nil, err
	}

	target, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}

	kindStr := "note"
	if kind == annoWarning {
		kindStr = "warning"
	}

	return &wtg2.AnnotationDecl{Kind: kindStr, Text: text, Target: target}, nil
}

func (d *decoderState) decodeSignal() (*wtg2.SignalDecl, error) {
	idx, err := d.decodeSym(distSignalType)
	if err != nil {
		return nil, err
	}

	target, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}

	return &wtg2.SignalDecl{Type: signalTypes[idx], Target: target}, nil
}

func (d *decoderState) decodeGameplay() (*wtg2.GameplayDecl, error) {
	idx, err := d.decodeSym(distGameplayType)
	if err != nil {
		return nil, err
	}

	hasText, err := d.decodeSym(distBool)
	if err != nil {
		return nil, err
	}

	gp := &wtg2.GameplayDecl{Type: gameplayTypes[idx]}

	if hasText == boolYes {
		text, err := d.decodeText()
		if err != nil {
			return nil, err
		}
		gp.Text = text
	}

	target, err := d.decodeIdentifier()
	if err != nil {
		return nil, err
	}
	gp.Target = target

	return gp, nil
}

func (d *decoderState) decodeIdentifier() (string, error) {
	known, err := d.decodeSym(distIdentKnown)
	if err != nil {
		return "", err
	}

	if known == 1 { // known
		idx, err := d.decodeSym(Uniform(d.sym.len()))
		if err != nil {
			return "", err
		}
		return d.sym.get(idx), nil
	}

	// New identifier
	name, err := d.decodeIDString()
	if err != nil {
		return "", err
	}
	d.sym.add(name)
	return name, nil
}

func (d *decoderState) decodeIDString() (string, error) {
	return d.decodeBytes()
}

func (d *decoderState) decodeText() (string, error) {
	return d.decodeBytes()
}

// decodeBytes decodes a string byte-by-byte using the byte distribution.
func (d *decoderState) decodeBytes() (string, error) {
	length, err := d.decodeLength()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	for i := range length {
		idx, err := d.decodeSym(distByte)
		if err != nil {
			return "", err
		}
		buf[i] = byte(idx)
	}
	return string(buf), nil
}

func (d *decoderState) decodeLength() (int, error) {
	isLong, err := d.decodeSym(distBool)
	if err != nil {
		return 0, err
	}
	if isLong == boolNo {
		n, err := d.decodeSym(Uniform(16))
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	n, err := d.decodeSym(Uniform(256))
	if err != nil {
		return 0, err
	}
	return n + 16, nil
}

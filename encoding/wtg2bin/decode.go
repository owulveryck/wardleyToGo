package wtg2bin

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/owulveryck/wardleyToGo/v2/parser/wtg2"
)

// Decode deserializes a WTG2 binary format payload back into a Document.
func Decode(data []byte) (*wtg2.Document, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: need at least 4 bytes, got %d", len(data))
	}
	if data[0] != magic0 || data[1] != magic1 {
		return nil, fmt.Errorf("invalid magic bytes: expected 0x%02X%02X, got 0x%02X%02X", magic0, magic1, data[0], data[1])
	}
	if data[2] != version {
		return nil, fmt.Errorf("unsupported version: %d (expected %d)", data[2], version)
	}
	flags := data[3]

	payload := data[4:]
	var r *bytes.Reader

	if flags&flagCompressed != 0 {
		fr := flate.NewReader(bytes.NewReader(payload))
		defer func() { _ = fr.Close() }()
		decompressed, err := io.ReadAll(fr)
		if err != nil {
			return nil, fmt.Errorf("decompressing payload: %w", err)
		}
		r = bytes.NewReader(decompressed)
	} else {
		r = bytes.NewReader(payload)
	}

	si, err := readStringIndex(r, r)
	if err != nil {
		return nil, fmt.Errorf("reading string table: %w", err)
	}

	doc := &wtg2.Document{
		Legend: flags&flagLegend != 0,
	}

	if err := decodeMetadata(r, si, doc); err != nil {
		return nil, err
	}

	doc.Nodes, err = decodeNodes(r, si)
	if err != nil {
		return nil, err
	}

	doc.Pipelines, err = decodePipelines(r, si)
	if err != nil {
		return nil, err
	}

	doc.Edges, err = decodeEdges(r, si)
	if err != nil {
		return nil, err
	}

	doc.Groups, err = decodeGroups(r, si)
	if err != nil {
		return nil, err
	}

	doc.Annotations, err = decodeAnnotations(r, si)
	if err != nil {
		return nil, err
	}

	doc.Signals, err = decodeSignals(r, si)
	if err != nil {
		return nil, err
	}

	doc.Gameplays, err = decodeGameplays(r, si)
	if err != nil {
		return nil, err
	}

	doc.Focuses, err = decodeFocuses(r, si)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func readEnum(r *bytes.Reader, si *stringIndex, codeToStr map[byte]string) (string, error) {
	b, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	if b == enumFallback {
		return readStringRef(r, si)
	}
	s, ok := codeToStr[b]
	if !ok {
		return "", fmt.Errorf("unknown enum code %d", b)
	}
	return s, nil
}

func readStringRef(r *bytes.Reader, si *stringIndex) (string, error) {
	id, err := binary.ReadUvarint(r)
	if err != nil {
		return "", fmt.Errorf("reading string ref: %w", err)
	}
	return si.get(id)
}

func readEvolution(r *bytes.Reader, si *stringIndex) (string, error) {
	b, err := r.ReadByte()
	if err != nil {
		return "", fmt.Errorf("reading evolution byte: %w", err)
	}
	if s, ok := decodeEvolution(b); ok {
		return s, nil
	}
	return readStringRef(r, si)
}

func decodeMetadata(r *bytes.Reader, si *stringIndex, doc *wtg2.Document) error {
	var err error
	doc.Title, err = readStringRef(r, si)
	if err != nil {
		return fmt.Errorf("decoding title: %w", err)
	}
	doc.Date, err = readStringRef(r, si)
	if err != nil {
		return fmt.Errorf("decoding date: %w", err)
	}
	doc.Author, err = readStringRef(r, si)
	if err != nil {
		return fmt.Errorf("decoding author: %w", err)
	}
	doc.Scope, err = readStringRef(r, si)
	if err != nil {
		return fmt.Errorf("decoding scope: %w", err)
	}
	doc.Question, err = readStringRef(r, si)
	if err != nil {
		return fmt.Errorf("decoding question: %w", err)
	}

	doc.Doctrine, err = readEnum(r, si, codeToDoctrine)
	if err != nil {
		return fmt.Errorf("decoding doctrine: %w", err)
	}

	for i := range doc.Stages {
		doc.Stages[i], err = readStringRef(r, si)
		if err != nil {
			return fmt.Errorf("decoding stage %d: %w", i, err)
		}
	}
	return nil
}

func decodeNodes(r *bytes.Reader, si *stringIndex) ([]*wtg2.NodeDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding node count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("node count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	nodes := make([]*wtg2.NodeDecl, count)
	for i := range nodes {
		n := &wtg2.NodeDecl{}

		n.Name, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("node %d name: %w", i, err)
		}

		kb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("node %d kind: %w", i, err)
		}
		n.Kind = wtg2.NodeKind(kb)

		n.Evolution, err = readEvolution(r, si)
		if err != nil {
			return nil, fmt.Errorf("node %d evolution: %w", i, err)
		}

		n.EvolvedTo, err = readEvolution(r, si)
		if err != nil {
			return nil, fmt.Errorf("node %d evolvedTo: %w", i, err)
		}

		ib, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("node %d inertia: %w", i, err)
		}
		n.Inertia = int(ib)

		ikb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("node %d inertia kinds: %w", i, err)
		}
		n.InertiaKinds = decodeInertiaKinds(ikb)

		n.Type, err = readEnum(r, si, codeToNodeType)
		if err != nil {
			return nil, fmt.Errorf("node %d type: %w", i, err)
		}

		n.Asset, err = readEnum(r, si, codeToAsset)
		if err != nil {
			return nil, fmt.Errorf("node %d asset: %w", i, err)
		}

		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("node %d flags: %w", i, err)
		}

		if fb&(1<<0) != 0 {
			n.Color, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("node %d color: %w", i, err)
			}
		}

		if fb&(1<<1) != 0 {
			var tmp [8]byte
			if _, err := io.ReadFull(r, tmp[:]); err != nil {
				return nil, fmt.Errorf("node %d visibility: %w", i, err)
			}
			n.Visibility = math.Float64frombits(binary.LittleEndian.Uint64(tmp[:]))
		} else {
			n.Visibility = -1
		}

		if fb&(1<<2) != 0 {
			n.Cost, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("node %d cost: %w", i, err)
			}
		}

		if fb&(1<<3) != 0 {
			n.Note, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("node %d note: %w", i, err)
			}
		}

		nodes[i] = n
	}
	return nodes, nil
}

func decodePipelines(r *bytes.Reader, si *stringIndex) ([]*wtg2.PipelineDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding pipeline count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("pipeline count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	pipelines := make([]*wtg2.PipelineDecl, count)
	for i := range pipelines {
		p := &wtg2.PipelineDecl{}
		p.Name, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("pipeline %d name: %w", i, err)
		}

		mc, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, fmt.Errorf("pipeline %d member count: %w", i, err)
		}
		if mc > maxSectionCount {
			return nil, fmt.Errorf("pipeline %d member count %d exceeds maximum", i, mc)
		}

		if mc > 0 {
			p.Members = make([]*wtg2.PipelineMemberDecl, mc)
			for j := range p.Members {
				m := &wtg2.PipelineMemberDecl{}
				m.Name, err = readStringRef(r, si)
				if err != nil {
					return nil, fmt.Errorf("pipeline %d member %d name: %w", i, j, err)
				}
				m.Position, err = readEvolution(r, si)
				if err != nil {
					return nil, fmt.Errorf("pipeline %d member %d position: %w", i, j, err)
				}
				p.Members[j] = m
			}
		}
		pipelines[i] = p
	}
	return pipelines, nil
}

func decodeEdges(r *bytes.Reader, si *stringIndex) ([]*wtg2.EdgeDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding edge count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("edge count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	edges := make([]*wtg2.EdgeDecl, count)
	for i := range edges {
		e := &wtg2.EdgeDecl{}
		e.From, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("edge %d from: %w", i, err)
		}
		e.To, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("edge %d to: %w", i, err)
		}

		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("edge %d flags: %w", i, err)
		}
		e.Bidirectional = fb&(1<<0) != 0

		if fb&(1<<1) != 0 {
			e.Label, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("edge %d label: %w", i, err)
			}
		}
		edges[i] = e
	}
	return edges, nil
}

func decodeGroups(r *bytes.Reader, si *stringIndex) ([]*wtg2.GroupDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding group count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("group count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	groups := make([]*wtg2.GroupDecl, count)
	for i := range groups {
		g := &wtg2.GroupDecl{}
		g.Name, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("group %d name: %w", i, err)
		}

		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("group %d flags: %w", i, err)
		}

		if fb&(1<<0) != 0 {
			g.Color, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("group %d color: %w", i, err)
			}
		}

		if fb&(1<<1) != 0 {
			g.Team, err = readEnum(r, si, codeToTeam)
			if err != nil {
				return nil, fmt.Errorf("group %d team: %w", i, err)
			}
		}

		mc, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, fmt.Errorf("group %d member count: %w", i, err)
		}
		if mc > maxSectionCount {
			return nil, fmt.Errorf("group %d member count %d exceeds maximum", i, mc)
		}

		if mc > 0 {
			g.Members = make([]string, mc)
			for j := range g.Members {
				g.Members[j], err = readStringRef(r, si)
				if err != nil {
					return nil, fmt.Errorf("group %d member %d: %w", i, j, err)
				}
			}
		}
		groups[i] = g
	}
	return groups, nil
}

func decodeAnnotations(r *bytes.Reader, si *stringIndex) ([]*wtg2.AnnotationDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding annotation count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("annotation count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	anns := make([]*wtg2.AnnotationDecl, count)
	for i := range anns {
		a := &wtg2.AnnotationDecl{}
		a.Kind, err = readEnum(r, si, codeToAnnotationKind)
		if err != nil {
			return nil, fmt.Errorf("annotation %d kind: %w", i, err)
		}

		a.Text, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("annotation %d text: %w", i, err)
		}
		a.Target, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("annotation %d target: %w", i, err)
		}
		anns[i] = a
	}
	return anns, nil
}

func decodeSignals(r *bytes.Reader, si *stringIndex) ([]*wtg2.SignalDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding signal count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("signal count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	signals := make([]*wtg2.SignalDecl, count)
	for i := range signals {
		s := &wtg2.SignalDecl{}
		s.Type, err = readEnum(r, si, codeToSignalType)
		if err != nil {
			return nil, fmt.Errorf("signal %d type: %w", i, err)
		}

		s.Target, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("signal %d target: %w", i, err)
		}
		signals[i] = s
	}
	return signals, nil
}

func decodeGameplays(r *bytes.Reader, si *stringIndex) ([]*wtg2.GameplayDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding gameplay count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("gameplay count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	gps := make([]*wtg2.GameplayDecl, count)
	for i := range gps {
		g := &wtg2.GameplayDecl{}
		g.Type, err = readEnum(r, si, codeToGameplayType)
		if err != nil {
			return nil, fmt.Errorf("gameplay %d type: %w", i, err)
		}

		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("gameplay %d flags: %w", i, err)
		}
		if fb&(1<<0) != 0 {
			g.Text, err = readStringRef(r, si)
			if err != nil {
				return nil, fmt.Errorf("gameplay %d text: %w", i, err)
			}
		}

		g.Target, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("gameplay %d target: %w", i, err)
		}
		gps[i] = g
	}
	return gps, nil
}

func decodeFocuses(r *bytes.Reader, si *stringIndex) ([]*wtg2.FocusDecl, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("decoding focus count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("focus count %d exceeds maximum %d", count, maxSectionCount)
	}
	if count == 0 {
		return nil, nil
	}

	focuses := make([]*wtg2.FocusDecl, count)
	for i := range focuses {
		f := &wtg2.FocusDecl{}
		var err error
		f.Target, err = readStringRef(r, si)
		if err != nil {
			return nil, fmt.Errorf("focus %d target: %w", i, err)
		}
		focuses[i] = f
	}
	return focuses, nil
}

package wtg2bin

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

// Options controls encoding behavior.
type Options struct {
	Compress bool // Apply deflate compression (default in Encode: true)
}

// Encode serializes a Document to the WTG2 binary format with deflate compression.
func Encode(doc *wtg2.Document) ([]byte, error) {
	return EncodeWithOptions(doc, Options{Compress: true})
}

// EncodeWithOptions serializes a Document with the specified options.
func EncodeWithOptions(doc *wtg2.Document, opts Options) ([]byte, error) {
	st := newStringTable()
	internAll(st, doc)

	var payload bytes.Buffer
	st.writeTo(&payload)
	encodeMetadata(&payload, st, doc)
	if err := encodeNodes(&payload, st, doc.Nodes); err != nil {
		return nil, fmt.Errorf("encoding nodes: %w", err)
	}
	encodePipelines(&payload, st, doc.Pipelines)
	encodeEdges(&payload, st, doc.Edges)
	encodeGroups(&payload, st, doc.Groups)
	encodeAnnotations(&payload, st, doc.Annotations)
	encodeSignals(&payload, st, doc.Signals)
	encodeGameplays(&payload, st, doc.Gameplays)
	encodeFocuses(&payload, st, doc.Focuses)

	var flags byte
	if doc.Legend {
		flags |= flagLegend
	}

	var out bytes.Buffer
	out.WriteByte(magic0)
	out.WriteByte(magic1)
	out.WriteByte(version)

	if opts.Compress {
		flags |= flagCompressed
		out.WriteByte(flags)
		w, err := flate.NewWriter(&out, flate.DefaultCompression)
		if err != nil {
			return nil, fmt.Errorf("creating deflate writer: %w", err)
		}
		if _, err := w.Write(payload.Bytes()); err != nil {
			return nil, fmt.Errorf("writing compressed payload: %w", err)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("closing deflate writer: %w", err)
		}
	} else {
		out.WriteByte(flags)
		out.Write(payload.Bytes())
	}

	return out.Bytes(), nil
}

// internAll walks the entire document and interns all strings into the table (pass 1).
func internAll(st *stringTable, doc *wtg2.Document) {
	st.intern(doc.Title)
	st.intern(doc.Date)
	st.intern(doc.Author)
	st.intern(doc.Scope)
	st.intern(doc.Question)
	internEnumFallback(st, doc.Doctrine, doctrineToCode)
	for _, s := range doc.Stages {
		st.intern(s)
	}
	for _, n := range doc.Nodes {
		st.intern(n.Name)
		internEvolution(st, n.Evolution)
		internEvolution(st, n.EvolvedTo)
		internEnumFallback(st, n.Type, nodeTypeToCode)
		internEnumFallback(st, n.Asset, assetToCode)
		st.intern(n.Color)
		st.intern(n.Cost)
		st.intern(n.Note)
	}
	for _, p := range doc.Pipelines {
		st.intern(p.Name)
		for _, m := range p.Members {
			st.intern(m.Name)
			internEvolution(st, m.Position)
		}
	}
	for _, e := range doc.Edges {
		st.intern(e.From)
		st.intern(e.To)
		st.intern(e.Label)
	}
	for _, g := range doc.Groups {
		st.intern(g.Name)
		st.intern(g.Color)
		internEnumFallback(st, g.Team, teamToCode)
		for _, m := range g.Members {
			st.intern(m)
		}
	}
	for _, a := range doc.Annotations {
		internEnumFallback(st, a.Kind, annotationKindToCode)
		st.intern(a.Text)
		st.intern(a.Target)
	}
	for _, s := range doc.Signals {
		internEnumFallback(st, s.Type, signalTypeToCode)
		st.intern(s.Target)
	}
	for _, g := range doc.Gameplays {
		internEnumFallback(st, g.Type, gameplayTypeToCode)
		st.intern(g.Text)
		st.intern(g.Target)
	}
	for _, f := range doc.Focuses {
		st.intern(f.Target)
	}
}

func internEvolution(st *stringTable, s string) {
	if _, ok := encodeEvolution(s); !ok {
		st.intern(s)
	}
}

func internEnumFallback(st *stringTable, s string, known map[string]byte) {
	if _, ok := known[s]; !ok && s != "" {
		st.intern(s)
	}
}

func encodeMetadata(buf *bytes.Buffer, st *stringTable, doc *wtg2.Document) {
	writeStringID(buf, st, doc.Title)
	writeStringID(buf, st, doc.Date)
	writeStringID(buf, st, doc.Author)
	writeStringID(buf, st, doc.Scope)
	writeStringID(buf, st, doc.Question)
	code, ok := doctrineToCode[doc.Doctrine]
	if ok {
		buf.WriteByte(code)
	} else {
		buf.WriteByte(enumFallback)
		writeStringID(buf, st, doc.Doctrine)
	}
	for _, s := range doc.Stages {
		writeStringID(buf, st, s)
	}
}

func encodeNodes(buf *bytes.Buffer, st *stringTable, nodes []*wtg2.NodeDecl) error {
	writeUvarint(buf, uint64(len(nodes)))
	for _, n := range nodes {
		writeStringID(buf, st, n.Name)
		buf.WriteByte(byte(n.Kind))
		writeEvolution(buf, st, n.Evolution)
		writeEvolution(buf, st, n.EvolvedTo)
		buf.WriteByte(byte(n.Inertia))

		ik, err := encodeInertiaKinds(n.InertiaKinds)
		if err != nil {
			return fmt.Errorf("node %q: %w", n.Name, err)
		}
		buf.WriteByte(ik)

		tc, ok := nodeTypeToCode[n.Type]
		if ok {
			buf.WriteByte(tc)
		} else {
			buf.WriteByte(enumFallback)
			writeStringID(buf, st, n.Type)
		}

		ac, ok := assetToCode[n.Asset]
		if ok {
			buf.WriteByte(ac)
		} else {
			buf.WriteByte(enumFallback)
			writeStringID(buf, st, n.Asset)
		}

		var flags byte
		if n.Color != "" {
			flags |= 1 << 0
		}
		if n.Visibility != -1 {
			flags |= 1 << 1
		}
		if n.Cost != "" {
			flags |= 1 << 2
		}
		if n.Note != "" {
			flags |= 1 << 3
		}
		buf.WriteByte(flags)

		if n.Color != "" {
			writeStringID(buf, st, n.Color)
		}
		if n.Visibility != -1 {
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], math.Float64bits(n.Visibility))
			buf.Write(tmp[:])
		}
		if n.Cost != "" {
			writeStringID(buf, st, n.Cost)
		}
		if n.Note != "" {
			writeStringID(buf, st, n.Note)
		}
	}
	return nil
}

func encodePipelines(buf *bytes.Buffer, st *stringTable, pipelines []*wtg2.PipelineDecl) {
	writeUvarint(buf, uint64(len(pipelines)))
	for _, p := range pipelines {
		writeStringID(buf, st, p.Name)
		writeUvarint(buf, uint64(len(p.Members)))
		for _, m := range p.Members {
			writeStringID(buf, st, m.Name)
			writeEvolution(buf, st, m.Position)
		}
	}
}

func encodeEdges(buf *bytes.Buffer, st *stringTable, edges []*wtg2.EdgeDecl) {
	writeUvarint(buf, uint64(len(edges)))
	for _, e := range edges {
		writeStringID(buf, st, e.From)
		writeStringID(buf, st, e.To)
		var flags byte
		if e.Bidirectional {
			flags |= 1 << 0
		}
		if e.Label != "" {
			flags |= 1 << 1
		}
		buf.WriteByte(flags)
		if e.Label != "" {
			writeStringID(buf, st, e.Label)
		}
	}
}

func encodeGroups(buf *bytes.Buffer, st *stringTable, groups []*wtg2.GroupDecl) {
	writeUvarint(buf, uint64(len(groups)))
	for _, g := range groups {
		writeStringID(buf, st, g.Name)
		var flags byte
		if g.Color != "" {
			flags |= 1 << 0
		}
		if g.Team != "" {
			flags |= 1 << 1
		}
		buf.WriteByte(flags)
		if g.Color != "" {
			writeStringID(buf, st, g.Color)
		}
		if g.Team != "" {
			tc, ok := teamToCode[g.Team]
			if ok {
				buf.WriteByte(tc)
			} else {
				buf.WriteByte(enumFallback)
				writeStringID(buf, st, g.Team)
			}
		}
		writeUvarint(buf, uint64(len(g.Members)))
		for _, m := range g.Members {
			writeStringID(buf, st, m)
		}
	}
}

func encodeAnnotations(buf *bytes.Buffer, st *stringTable, anns []*wtg2.AnnotationDecl) {
	writeUvarint(buf, uint64(len(anns)))
	for _, a := range anns {
		code, ok := annotationKindToCode[a.Kind]
		if ok {
			buf.WriteByte(code)
		} else {
			buf.WriteByte(enumFallback)
			writeStringID(buf, st, a.Kind)
		}
		writeStringID(buf, st, a.Text)
		writeStringID(buf, st, a.Target)
	}
}

func encodeSignals(buf *bytes.Buffer, st *stringTable, signals []*wtg2.SignalDecl) {
	writeUvarint(buf, uint64(len(signals)))
	for _, s := range signals {
		code, ok := signalTypeToCode[s.Type]
		if ok {
			buf.WriteByte(code)
		} else {
			buf.WriteByte(enumFallback)
			writeStringID(buf, st, s.Type)
		}
		writeStringID(buf, st, s.Target)
	}
}

func encodeGameplays(buf *bytes.Buffer, st *stringTable, gps []*wtg2.GameplayDecl) {
	writeUvarint(buf, uint64(len(gps)))
	for _, g := range gps {
		code, ok := gameplayTypeToCode[g.Type]
		if ok {
			buf.WriteByte(code)
		} else {
			buf.WriteByte(enumFallback)
			writeStringID(buf, st, g.Type)
		}
		var flags byte
		if g.Text != "" {
			flags |= 1 << 0
		}
		buf.WriteByte(flags)
		if g.Text != "" {
			writeStringID(buf, st, g.Text)
		}
		writeStringID(buf, st, g.Target)
	}
}

func encodeFocuses(buf *bytes.Buffer, st *stringTable, focuses []*wtg2.FocusDecl) {
	writeUvarint(buf, uint64(len(focuses)))
	for _, f := range focuses {
		writeStringID(buf, st, f.Target)
	}
}

func writeEvolution(buf *bytes.Buffer, st *stringTable, s string) {
	if b, ok := encodeEvolution(s); ok {
		buf.WriteByte(b)
		return
	}
	buf.WriteByte(evolFallback)
	writeStringID(buf, st, s)
}

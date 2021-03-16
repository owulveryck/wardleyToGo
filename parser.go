package main

import (
	"io"
	"strconv"
	"strings"
	"text/scanner"

	svg "github.com/ajstarks/svgo"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type parser struct {
	s        *scanner.Scanner
	currID   int
	g        *simple.DirectedGraph
	nodeDict map[string]graph.Node
	edges    []edge
}

func newParser(r io.Reader) *parser {
	var s scanner.Scanner
	s.Init(r)
	s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
	return &parser{
		s: &s,
	}
}

func (p *parser) parse() {
	p.nodeDict = make(map[string]graph.Node)
	p.edges = make([]edge, 0)
	p.g = simple.NewDirectedGraph()
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		switch p.s.TokenText() {
		case "component":
			e := p.parseComponent()
			e.id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.label] = e
		case "anchor":
			e := p.parseAnchor()
			e.id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.label] = e
		default:
			e := p.parseDefault(p.s.TokenText())
			switch e := e.(type) {
			case edge:
				p.edges = append(p.edges, e)
			}
		}
	}
	p.createEdges()
}

func (p *parser) createEdges() {
	for _, edge := range p.edges {
		edge.F = p.nodeDict[edge.fromLabel]
		edge.T = p.nodeDict[edge.toLabel]
		p.g.SetEdge(edge)
	}
}

func (p *parser) parseDefault(firstElement string) interface{} {
	var e edge
	var b strings.Builder
	b.WriteString(firstElement)
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Ident {
			b.WriteRune(' ')
			b.WriteString(p.s.TokenText())
		}
		if tok == '>' {
			e.fromLabel = strings.TrimLeft(b.String(), " ")
			b.Reset()
		}
	}
	if e.fromLabel != "" {
		e.toLabel = strings.TrimLeft(b.String(), " ")
		return e
	}
	return nil
}

func (p *parser) parseComponent() *component {
	c := &component{
		coords:     [2]int{-1, -1},
		labelCoord: [2]int{-1, -1},
	}
	var b strings.Builder
	inLabel := true
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == '[' {
			inLabel = false
		}
		if tok == scanner.Ident && inLabel {
			b.WriteString(p.s.TokenText())
			b.WriteRune(' ')
		}
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				panic(err)
			}
			if c.coords[0] == -1 {
				c.coords[0] = int(f * 100)
				continue
			}
			if c.coords[1] == -1 {
				c.coords[1] = int(f * 100)
				continue
			}
		}
		if tok == scanner.Int {
			i, err := strconv.Atoi(p.s.TokenText())
			if err != nil {
				panic(err)
			}
			if c.labelCoord[0] == -1 {
				c.labelCoord[0] = i
				continue
			}
			if c.labelCoord[1] == -1 {
				c.labelCoord[1] = i
				continue
			}
		}
	}
	c.label = strings.TrimRight(b.String(), " ")
	return c
}

func (p *parser) parseAnchor() *anchor {
	a := &anchor{
		coords: [2]int{-1, -1},
	}
	var b strings.Builder
	inLabel := true
	curLine := p.s.Pos().Line
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if curLine != p.s.Pos().Line {
			// emit the component
			break
		}
		if tok == '[' {
			inLabel = false
		}
		if tok == scanner.Ident && inLabel {
			b.WriteString(p.s.TokenText())
			b.WriteString(" ")
		}
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				panic(err)
			}
			if a.coords[0] == -1 {
				a.coords[0] = int(f * 100)
				continue
			}
			if a.coords[1] == -1 {
				a.coords[1] = int(f * 100)
				continue
			}
		}
	}
	a.label = strings.TrimRight(b.String(), " ")
	return a
}

type edge struct {
	toLabel   string
	fromLabel string
	T         graph.Node
	F         graph.Node
	edgeLabel string
}

func (e edge) From() graph.Node {
	return e.F
}

func (e edge) ReversedEdge() graph.Edge {
	return edge{
		F:         e.T,
		T:         e.F,
		toLabel:   e.fromLabel,
		fromLabel: e.toLabel,
		edgeLabel: e.edgeLabel,
	}
}

func (e edge) To() graph.Node {
	return e.T
}

func (e edge) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	fromCoord := e.F.(element).GetCoordinates()
	toCoord := e.T.(element).GetCoordinates()
	s.Line(fromCoord[1]*(width-padLeft)/100+padLeft,
		u(height-padLeft)-fromCoord[0]*(height-padLeft)/100,
		toCoord[1]*(width-padLeft)/100+padLeft,
		(height-padLeft)-toCoord[0]*(height-padLeft)/100,
		`stroke="grey"`, `stroke-width="1"`)
}

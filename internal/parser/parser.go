package parser

import (
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Parser struct {
	s        *scanner.Scanner
	currID   int
	G        *simple.DirectedGraph
	nodeDict map[string]graph.Node
	edges    []edge
}

func NewParser(r io.Reader) *Parser {
	var s scanner.Scanner
	s.Init(r)
	s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
	return &Parser{
		s: &s,
	}
}

func (p *Parser) Parse() {
	p.nodeDict = make(map[string]graph.Node)
	p.edges = make([]edge, 0)
	p.G = simple.NewDirectedGraph()
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		switch p.s.TokenText() {
		case "component":
			e := p.parseComponent()
			e.Id = int64(p.currID)
			p.currID++
			p.G.AddNode(e)
			p.nodeDict[e.Label] = e
		case "anchor":
			e := p.parseAnchor()
			e.Id = int64(p.currID)
			p.currID++
			p.G.AddNode(e)
			p.nodeDict[e.Label] = e
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

func (p *Parser) createEdges() {
	for _, edge := range p.edges {
		edge.F = p.nodeDict[edge.fromLabel]
		edge.T = p.nodeDict[edge.toLabel]
		p.G.SetEdge(edge)
	}
}

func (p *Parser) parseDefault(firstElement string) interface{} {
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

func (p *Parser) parseComponent() *wardley.Component {
	c := &wardley.Component{
		Coords:      [2]int{-1, -1},
		LabelCoords: [2]int{-1, -1},
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
			if c.Coords[0] == -1 {
				c.Coords[0] = int(f * 100)
				continue
			}
			if c.Coords[1] == -1 {
				c.Coords[1] = int(f * 100)
				continue
			}
		}
		if tok == scanner.Int {
			i, err := strconv.Atoi(p.s.TokenText())
			if err != nil {
				panic(err)
			}
			if c.LabelCoords[0] == -1 {
				c.LabelCoords[0] = i
				continue
			}
			if c.LabelCoords[1] == -1 {
				c.LabelCoords[1] = i
				continue
			}
		}
	}
	c.Label = strings.TrimRight(b.String(), " ")
	return c
}

func (p *Parser) parseAnchor() *wardley.Anchor {
	a := &wardley.Anchor{
		Coords: [2]int{-1, -1},
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
			if a.Coords[0] == -1 {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == -1 {
				a.Coords[1] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a
}

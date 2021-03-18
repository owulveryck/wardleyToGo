package parser

import (
	"fmt"
	"io"
	"log"
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
	g        *simple.DirectedGraph
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

func (p *Parser) Parse() (*wardley.Map, error) {
	p.nodeDict = make(map[string]graph.Node)
	p.edges = make([]edge, 0)
	p.g = simple.NewDirectedGraph()
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		switch p.s.TokenText() {
		case "component":
			e, err := p.parseComponent()
			if err != nil {
				return nil, err
			}
			e.Id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.Label] = e
		case "anchor":
			e, err := p.parseAnchor()
			if err != nil {
				return nil, err
			}
			e.Id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.Label] = e
		case "streamalignedteam":
			e, err := p.parseStreamAligned()
			if err != nil {
				return nil, err
			}
			e.Id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.Label] = e
		case "enablingteam":
			e, err := p.parseEnabling()
			if err != nil {
				return nil, err
			}
			e.Id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeDict[e.Label] = e
		default:
			e, err := p.parseDefault(p.s.TokenText())
			if err != nil {
				log.Println("Warning", err)
			}
			switch e := e.(type) {
			case edge:
				p.edges = append(p.edges, e)
			}
		}
	}
	err := p.createEdges()
	return &wardley.Map{
		DirectedGraph: p.g,
	}, err
}

func (p *Parser) createEdges() error {
	var ok bool
	for _, edge := range p.edges {
		edge.F, ok = p.nodeDict[edge.fromLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		edge.T, ok = p.nodeDict[edge.toLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		p.g.SetEdge(edge)
	}
	return nil
}

func (p *Parser) parseDefault(firstElement string) (interface{}, error) {
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
		return e, nil
	}
	return nil, fmt.Errorf("unhandled element at line %v: %v", p.s.Line, b.String())
}

func (p *Parser) parseComponent() (*wardley.Component, error) {
	c := &wardley.Component{
		Coords:      [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
		LabelCoords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
	}
	var b strings.Builder
	inLabel := true
	var prevTok rune
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
				return nil, err
			}
			if c.Coords[0] == wardley.UndefinedCoord {
				c.Coords[0] = int(f * 100)
				continue
			}
			if c.Coords[1] == wardley.UndefinedCoord {
				c.Coords[1] = int(f * 100)
				continue
			}
		}
		if tok == scanner.Int {
			sign := ""
			if prevTok == '-' {
				sign = "-"
			}
			i, err := strconv.Atoi(sign + p.s.TokenText())
			if err != nil {
				return nil, err
			}
			if c.LabelCoords[0] == wardley.UndefinedCoord {
				c.LabelCoords[0] = i
				continue
			}
			if c.LabelCoords[1] == wardley.UndefinedCoord {
				c.LabelCoords[1] = i
				continue
			}
		}
		prevTok = tok
	}
	c.Label = strings.TrimRight(b.String(), " ")
	return c, nil
}

func (p *Parser) parseAnchor() (*wardley.Anchor, error) {
	a := &wardley.Anchor{
		Coords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
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
				return nil, err
			}
			if a.Coords[0] == wardley.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == wardley.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

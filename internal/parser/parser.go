package parser

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Parser struct {
	s              *scanner.Scanner
	currID         int
	g              *simple.DirectedGraph
	nodeDict       map[string]graph.Node
	nodeEvolveDict map[string]graph.Node
	edges          []edge
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
	p.nodeEvolveDict = make(map[string]graph.Node)
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
		case "evolve":
			e, err := p.parseEvolve()
			if err != nil {
				return nil, err
			}
			e.Id = int64(p.currID)
			p.currID++
			p.g.AddNode(e)
			p.nodeEvolveDict[e.Label] = e
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
		case "platformteam":
			e, err := p.parsePlatform()
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
	err := p.completeEvolve()
	if err != nil {
		return nil, err
	}
	err = p.createEdges()
	if err != nil {
		return nil, err
	}
	return &wardley.Map{
		DirectedGraph: p.g,
	}, nil
}

func (p *Parser) parseDefault(firstElement string) (interface{}, error) {
	var e edge
	e.edgeType = wardley.RegularEdge
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

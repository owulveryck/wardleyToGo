package parser

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Parser struct {
	s                    *scanner.Scanner
	title                string
	g                    *simple.DirectedGraph
	nodeDict             map[string]graph.Node
	nodeEvolveDict       map[string]graph.Node
	edges                []plan.Edge
	annotations          []*plan.Annotation
	annotationsPlacement [2]int
}

func NewParser(r io.Reader) *Parser {
	var s scanner.Scanner
	s.Init(r)
	s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
	return &Parser{
		s:              &s,
		nodeDict:       make(map[string]graph.Node),
		nodeEvolveDict: make(map[string]graph.Node),
		edges:          make([]plan.Edge, 0),
		annotations:    make([]*plan.Annotation, 0),
		g:              simple.NewDirectedGraph(),
	}
}

func (p *Parser) Parse() (*plan.Map, error) {
	parsers := map[string]func() error{
		"title":                    p.parseTitle,
		"component":                p.parseComponent,
		"evolve":                   p.parseEvolve,
		"anchor":                   p.parseAnchor,
		"streamAlignedTeam":        p.parseStreamAligned,
		"enablingTeam":             p.parseEnabling,
		"platformTeam":             p.parsePlatform,
		"complicatedSubsystemTeam": p.parseComplicatedSubsystem,
		"annotation":               p.parseAnnotation,
		"annotations":              p.parseAnnotations,
	}
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		if tok == '\n' {
			continue
		}
		if parser, ok := parsers[p.s.TokenText()]; ok {
			err := parser()
			if err != nil {
				return nil, err
			}
			continue
		}
		e, err := p.parseDefault(p.s.TokenText())
		if err != nil {
			log.Println("Warning", err)
		}
		switch e := e.(type) {
		case plan.Edge:
			p.edges = append(p.edges, e)
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
	return &plan.Map{
		Title:                p.title,
		DirectedGraph:        p.g,
		Annotations:          p.annotations,
		AnnotationsPlacement: p.annotationsPlacement,
	}, nil
}

func (p *Parser) parseDefault(firstElement string) (interface{}, error) {
	var e plan.Edge
	e.EdgeType = plan.RegularEdge
	var b strings.Builder
	b.WriteString(firstElement)
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Ident {
			b.WriteRune(' ')
			b.WriteString(p.s.TokenText())
		}
		if tok == '>' {
			e.FromLabel = strings.TrimLeft(b.String(), " ")
			b.Reset()
		}
	}
	if e.FromLabel != "" {
		e.ToLabel = strings.TrimLeft(b.String(), " ")
		return e, nil
	}
	return nil, fmt.Errorf("unhandled element at line %v: %v", p.s.Line, b.String())
}

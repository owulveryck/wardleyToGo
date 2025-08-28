package owm

import (
	"fmt"
	"image"
	"io"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo"
	tt "github.com/owulveryck/wardleyToGo/components/teamtopologies"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Parser struct {
	s                    *scanner.Scanner
	title                string
	g                    *simple.DirectedGraph
	nodeDict             map[string]graph.Node
	nodeEvolveDict       map[string]graph.Node
	edges                []edge
	annotations          []*wardleyToGo.Annotation
	annotationsPlacement image.Point
	warnings             []error
}

func NewParser(r io.Reader) *Parser {
	var s scanner.Scanner
	s.Init(r)
	s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
	/*
		s.IsIdentRune = func(ch rune, i int) bool {
			fmt.Printf("%s", string(ch))
			if ch == '\'' && i > 0 {
				return true
			}
			if unicode.IsGraphic(ch) && !unicode.IsSpace(ch) {
				return true
			}
			fmt.Println("")
			return false
		}
	*/
	return &Parser{
		s:              &s,
		nodeDict:       make(map[string]graph.Node),
		nodeEvolveDict: make(map[string]graph.Node),
		edges:          make([]edge, 0),
		annotations:    make([]*wardleyToGo.Annotation, 0),
		g:              simple.NewDirectedGraph(),
		warnings:       make([]error, 0),
	}
}

func (p *Parser) Parse() (*wardleyToGo.Map, error) {
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
			p.warnings = append(p.warnings, err)
		}
		switch e := e.(type) {
		case edge:
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
	wmap := &wardleyToGo.Map{
		Title:                p.title,
		DirectedGraph:        p.g,
		Annotations:          p.annotations,
		AnnotationsPlacement: p.annotationsPlacement,
	}

	// Attach warnings to the map if needed
	// Note: wardleyToGo.Map doesn't have a warnings field, so we could add this as a method
	// or consider returning warnings separately

	return wmap, nil
}

// GetWarnings returns any warnings encountered during parsing
func (p *Parser) GetWarnings() []error {
	return p.warnings
}

func (p *Parser) parseDefault(firstElement string) (interface{}, error) {
	var e edge
	//e.EdgeType = wardleyToGo.RegularEdge
	var b strings.Builder
	b.WriteString(firstElement)
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Ident {
			b.WriteRune(' ')
			b.WriteString(p.s.TokenText())
		}
		if tok == '-' {
			e.FromLabel = strings.TrimLeft(b.String(), " ")
			b.Reset()
		}
		if tok == '>' {
			switch strings.TrimLeft(b.String(), " ") {
			case "collaboration":
				e.EdgeType = tt.CollaborationEdge
			case "facilitating":
				e.EdgeType = tt.FacilitatingEdge
			case "xAsAService":
				e.EdgeType = tt.XAsAServiceEdge
			}
			b.Reset()
		}
	}
	if e.FromLabel != "" {
		e.ToLabel = strings.TrimLeft(b.String(), " ")
		return e, nil
	}
	return nil, fmt.Errorf("unhandled element at line %v: %v", p.s.Line, b.String())
}

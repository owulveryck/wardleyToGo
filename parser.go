package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

type parser struct {
	s *scanner.Scanner
}

func newParser(r io.Reader) *parser {
	var s scanner.Scanner
	s.Init(r)
	s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
	return &parser{
		s: &s,
	}
}

func (p *parser) Parse() {
	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		switch p.s.TokenText() {
		case "component":
			fmt.Println(p.parseComponent())
		case "anchor":
			p.parseAnchor()
		default:
			p.parseDefault()
		}
	}
}

func (p *parser) parseDefault() {

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

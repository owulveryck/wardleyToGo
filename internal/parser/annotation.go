package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseAnnotation() error {
	var a *plan.Annotation
	coords := make([]int, 0)
	var b strings.Builder
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Int && a == nil {
			i, err := strconv.Atoi(p.s.TokenText())
			if err != nil {
				return err
			}
			a = plan.NewAnnotation(i)
		}
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				return err
			}
			c := int(f * 100)
			coords = append(coords, c)
		}
		if tok == scanner.Ident {
			b.WriteString(p.s.TokenText())
			b.WriteRune(' ')
		}
	}
	if a == nil {
		return errors.New("bad coordinates")
	}
	if len(coords)%2 != 0 || len(coords) == 0 {
		return fmt.Errorf("incomplete coordinates: %v", coords)
	}
	for i := 0; i < len(coords); i += 2 {
		a.Coords = append(a.Coords, [2]int{coords[i], coords[i+1]})
	}
	a.Label = strings.TrimRight(b.String(), " ")
	p.annotations = append(p.annotations, a)
	return nil
}

func (p *Parser) parseAnnotations() error {
	coords := make([]int, 0)
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				return err
			}
			c := int(f * 100)
			coords = append(coords, c)
		}
	}
	p.annotationsPlacement = [2]int{coords[0], coords[1]}
	return nil
}

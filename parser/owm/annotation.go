package owm

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo"
)

func (p *Parser) parseAnnotation() error {
	var a *wardleyToGo.Annotation
	coords := make([]int, 0)
	var b strings.Builder
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == scanner.Int && a == nil {
			i, err := strconv.Atoi(p.s.TokenText())
			if err != nil {
				return err
			}
			a = wardleyToGo.NewAnnotation(i)
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
		a.Placements = append(a.Placements, image.Point{100 - coords[i+1], coords[i]})
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
	p.annotationsPlacement = image.Point{100 - coords[1], coords[0]}
	return nil
}

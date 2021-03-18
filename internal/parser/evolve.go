package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

func (p *Parser) parseEvolve() (*wardley.EvolvedComponent, error) {
	c := &wardley.EvolvedComponent{
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
			inLabel = false
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				return nil, err
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

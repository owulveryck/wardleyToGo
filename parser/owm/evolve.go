package parser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo"
)

func (p *Parser) parseEvolve() error {
	c, err := scanEvolve(p.s, p.g.NewNode().ID())
	if err != nil {
		return err
	}
	p.g.AddNode(c)
	p.nodeEvolveDict[c.Label] = c
	return nil
}

func scanEvolve(s *scanner.Scanner, id int64) (*wardleyToGo.EvolvedComponent, error) {
	c := wardleyToGo.NewEvolvedComponent(id)
	var b strings.Builder
	var prevTok rune
	labelize := func(c *wardleyToGo.EvolvedComponent, b *strings.Builder) {
		if c.Label == "" {
			c.Label = strings.TrimRight(b.String(), " ")
		}
		b.Reset()
	}
	for tok := s.Scan(); tok != '\n' && tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '[':
			labelize(c, &b)
		case scanner.Ident:
			b.WriteString(s.TokenText())
			b.WriteRune(' ')
		case scanner.Float:
			c.Label = strings.TrimRight(b.String(), " ")
			b.Reset()
			f, err := strconv.ParseFloat(s.TokenText(), 64)
			if err != nil {
				return nil, err
			}
			if c.Coords[1] == wardleyToGo.UndefinedCoord {
				c.Coords[1] = int(f * 100)
				continue
			}
		case '(':
			labelize(c, &b)
		case ')':
			switch strings.TrimRight(b.String(), " ") {
			case "build":
				c.Type = wardleyToGo.BuildComponent
			case "buy":
				c.Type = wardleyToGo.BuyComponent
			case "outsource":
				c.Type = wardleyToGo.OutsourceComponent
			case "dataProduct":
				c.Type = wardleyToGo.DataProductComponent
			default:
				return nil, fmt.Errorf("unhandled type %v", strings.TrimRight(b.String(), " "))
			}
		case scanner.Int:
			sign := ""
			if prevTok == '-' {
				sign = "-"
			}
			i, err := strconv.Atoi(sign + s.TokenText())
			if err != nil {
				return nil, err
			}
			if c.LabelCoords[0] == wardleyToGo.UndefinedCoord {
				c.LabelCoords[0] = i
				continue
			}
			if c.LabelCoords[1] == wardleyToGo.UndefinedCoord {
				c.LabelCoords[1] = i
				continue
			}
		}
		prevTok = tok
	}
	labelize(c, &b)
	return c, nil
}

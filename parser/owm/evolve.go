package owm

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/components/wardley"
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

func scanEvolve(s *scanner.Scanner, id int64) (*wardley.EvolvedComponent, error) {
	c := wardley.NewEvolvedComponent(id)
	var b strings.Builder
	var prevTok rune
	labelize := func(c *wardley.EvolvedComponent, b *strings.Builder) {
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
			if c.Placement.X == components.UndefinedCoord {
				c.Placement.X = int(f * 100)
				continue
			}
		case '(':
			labelize(c, &b)
		case ')':
			switch strings.TrimRight(b.String(), " ") {
			case "build":
				c.Type = wardley.BuildComponent
			case "buy":
				c.Type = wardley.BuyComponent
			case "outsource":
				c.Type = wardley.OutsourceComponent
			case "dataProduct":
				c.Type = wardley.DataProductComponent
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
			if c.LabelPlacement.X == components.UndefinedCoord {
				c.LabelPlacement.X = i
				continue
			}
			if c.LabelPlacement.Y == components.UndefinedCoord {
				c.LabelPlacement.Y = i
				continue
			}
		}
		prevTok = tok
	}
	labelize(c, &b)
	return c, nil
}

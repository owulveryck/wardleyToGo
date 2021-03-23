package parser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseComponent() error {
	c, err := scanComponent(p.s, p.g.NewNode().ID())
	if err != nil {
		return err
	}
	p.g.AddNode(c)
	p.nodeDict[c.Label] = c
	return nil
}

func scanComponent(s *scanner.Scanner, id int64) (*plan.Component, error) {
	c := plan.NewComponent(id)
	labelize := func(c *plan.Component, b *strings.Builder) {
		if c.Label == "" {
			c.Label = strings.TrimRight(b.String(), " ")
		}
		b.Reset()
	}
	var b strings.Builder
	var prevTok rune
	for tok := s.Scan(); tok != '\n' && tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '[':
			labelize(c, &b)
		case scanner.Ident:
			b.WriteString(s.TokenText())
			b.WriteRune(' ')
		case scanner.Float:
			f, err := strconv.ParseFloat(s.TokenText(), 64)
			if err != nil {
				return nil, err
			}
			if c.Coords[0] == plan.UndefinedCoord {
				c.Coords[0] = int(f * 100)
				continue
			}
			if c.Coords[1] == plan.UndefinedCoord {
				c.Coords[1] = int(f * 100)
				continue
			}
		case '(':
			labelize(c, &b)
			//b.Reset()
		case ')':
			switch strings.TrimRight(b.String(), " ") {
			case "build":
				c.Type = plan.BuildComponent
			case "buy":
				c.Type = plan.BuyComponent
			case "outsource":
				c.Type = plan.OutsourceComponent
			case "dataProduct":
				c.Type = plan.DataProductComponent
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
			if c.LabelCoords[0] == plan.UndefinedCoord {
				c.LabelCoords[0] = i
				continue
			}
			if c.LabelCoords[1] == plan.UndefinedCoord {
				c.LabelCoords[1] = i
				continue
			}
		}
		prevTok = tok
	}
	labelize(c, &b)
	return c, nil
}

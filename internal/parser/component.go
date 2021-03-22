package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseComponent() error {
	c := plan.NewComponent(p.g.NewNode().ID())
	var b strings.Builder
	var prevTok rune
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if tok == '[' && c.Label == "" {
			c.Label = strings.TrimRight(b.String(), " ")
			b.Reset()
		}
		if tok == scanner.Ident {
			b.WriteString(p.s.TokenText())
			b.WriteRune(' ')
		}
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(p.s.TokenText(), 64)
			if err != nil {
				return err
			}
			if c.Coords[0] == plan.UndefinedCoord {
				c.Coords[0] = int(f * 100)
				continue
			}
			if c.Coords[1] == plan.UndefinedCoord {
				c.Coords[1] = int(f * 100)
				continue
			}
		}
		if tok == '(' {
			b.Reset()
		}
		if tok == ')' {
			switch strings.TrimRight(b.String(), " ") {
			case "build":
				c.Type = plan.BuildComponent
			case "buy":
				c.Type = plan.BuyComponent
			case "outsource":
				c.Type = plan.OutsourceComponent
			}
		}
		if tok == scanner.Int {
			sign := ""
			if prevTok == '-' {
				sign = "-"
			}
			i, err := strconv.Atoi(sign + p.s.TokenText())
			if err != nil {
				return err
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

	p.g.AddNode(c)
	p.nodeDict[c.Label] = c
	return nil
}

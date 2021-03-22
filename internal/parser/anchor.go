package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseAnchor() error {
	a := plan.NewAnchor(p.g.NewNode().ID())
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
				return err
			}
			if a.Coords[0] == plan.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == plan.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	p.g.AddNode(a)
	p.nodeDict[a.Label] = a
	return nil
}

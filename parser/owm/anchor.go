package owm

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func (p *Parser) parseAnchor() error {
	a, err := scanAnchor(p.s, p.g.NewNode().ID())
	if err != nil {
		return err
	}
	p.g.AddNode(a)
	p.nodeDict[a.Label] = a
	return nil
}

func scanAnchor(s *scanner.Scanner, id int64) (*wardley.Anchor, error) {
	a := wardley.NewAnchor(id)
	var b strings.Builder
	inLabel := true
	curLine := s.Pos().Line
	for tok := s.Scan(); tok != '\n' && tok != scanner.EOF; tok = s.Scan() {
		if curLine != s.Pos().Line {
			// emit the component
			break
		}
		if tok == '[' {
			inLabel = false
		}
		if tok == scanner.Ident && inLabel {
			b.WriteString(s.TokenText())
			b.WriteString(" ")
		}
		if tok == scanner.Float {
			f, err := strconv.ParseFloat(s.TokenText(), 64)
			if err != nil {
				return nil, err
			}
			if a.Placement.X == components.UndefinedCoord {
				a.Placement.X = int(f * 100)
				continue
			}
			if a.Placement.Y == components.UndefinedCoord {
				a.Placement.Y = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

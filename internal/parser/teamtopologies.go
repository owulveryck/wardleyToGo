package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

func (p *Parser) parseStreamAligned() (*wardley.StreamAlignedTeam, error) {
	a := &wardley.StreamAlignedTeam{
		Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
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
				return nil, err
			}
			if a.Coords[0] == wardley.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == wardley.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
			if a.Coords[2] == wardley.UndefinedCoord {
				a.Coords[2] = int(f * 100)
				continue
			}
			if a.Coords[3] == wardley.UndefinedCoord {
				a.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseStreamAligned() (*plan.StreamAlignedTeam, error) {
	a := &plan.StreamAlignedTeam{
		Coords: [4]int{plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord},
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
			if a.Coords[0] == plan.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == plan.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
			if a.Coords[2] == plan.UndefinedCoord {
				a.Coords[2] = int(f * 100)
				continue
			}
			if a.Coords[3] == plan.UndefinedCoord {
				a.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

func (p *Parser) parseEnabling() (*plan.EnablingTeam, error) {
	a := &plan.EnablingTeam{
		Coords: [4]int{plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord},
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
			if a.Coords[0] == plan.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == plan.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
			if a.Coords[2] == plan.UndefinedCoord {
				a.Coords[2] = int(f * 100)
				continue
			}
			if a.Coords[3] == plan.UndefinedCoord {
				a.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

func (p *Parser) parsePlatform() (*plan.PlatformTeam, error) {
	a := &plan.PlatformTeam{
		Coords: [4]int{plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord, plan.UndefinedCoord},
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
			if a.Coords[0] == plan.UndefinedCoord {
				a.Coords[0] = int(f * 100)
				continue
			}
			if a.Coords[1] == plan.UndefinedCoord {
				a.Coords[1] = int(f * 100)
				continue
			}
			if a.Coords[2] == plan.UndefinedCoord {
				a.Coords[2] = int(f * 100)
				continue
			}
			if a.Coords[3] == plan.UndefinedCoord {
				a.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	a.Label = strings.TrimRight(b.String(), " ")
	return a, nil
}

package owm

import (
	"image"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/components"
	tt "github.com/owulveryck/wardleyToGo/components/teamtopologies"
)

func (p *Parser) parseTeam() (*tt.Team, error) {
	coords := []int{components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord}
	team := tt.NewTeam(p.g.NewNode().ID())
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
			if coords[0] == components.UndefinedCoord {
				coords[0] = int(f * 100)
				continue
			}
			if coords[1] == components.UndefinedCoord {
				coords[1] = int(f * 100)
				continue
			}
			if coords[2] == components.UndefinedCoord {
				coords[2] = int(f * 100)
				continue
			}
			if coords[3] == components.UndefinedCoord {
				coords[3] = int(f * 100)
				continue
			}
		}
	}
	team.Area = image.Rect(coords[1], 100-coords[0], coords[3], 100-coords[2])
	team.Label = strings.TrimRight(b.String(), " ")
	return team, nil

}

func (p *Parser) parseStreamAligned() error {
	t, err := p.parseTeam()
	if err != nil {
		return err
	}
	s := &tt.StreamAlignedTeam{Team: t}
	p.g.AddNode(s)
	p.nodeDict[s.Label] = s
	return nil
}

func (p *Parser) parseEnabling() error {
	t, err := p.parseTeam()
	if err != nil {
		return err
	}
	s := &tt.EnablingTeam{Team: t}
	p.g.AddNode(s)
	p.nodeDict[s.Label] = s
	return nil
}

func (p *Parser) parsePlatform() error {
	t, err := p.parseTeam()
	if err != nil {
		return err
	}
	s := &tt.PlatformTeam{Team: t}
	p.g.AddNode(s)
	p.nodeDict[s.Label] = s
	return nil
}

func (p *Parser) parseComplicatedSubsystem() error {
	t, err := p.parseTeam()
	if err != nil {
		return err
	}
	s := &tt.ComplicatedSubsystemTeam{Team: t}
	p.g.AddNode(s)
	p.nodeDict[s.Label] = s
	return nil
}

package parser

import (
	"strconv"
	"strings"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
)

func (p *Parser) parseStreamAligned() error {
	streamAligned := plan.NewStreamAlignedTeam(p.g.NewNode().ID())
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
			if streamAligned.Coords[0] == plan.UndefinedCoord {
				streamAligned.Coords[0] = int(f * 100)
				continue
			}
			if streamAligned.Coords[1] == plan.UndefinedCoord {
				streamAligned.Coords[1] = int(f * 100)
				continue
			}
			if streamAligned.Coords[2] == plan.UndefinedCoord {
				streamAligned.Coords[2] = int(f * 100)
				continue
			}
			if streamAligned.Coords[3] == plan.UndefinedCoord {
				streamAligned.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	streamAligned.Label = strings.TrimRight(b.String(), " ")
	p.g.AddNode(streamAligned)
	p.nodeDict[streamAligned.Label] = streamAligned
	return nil
}

func (p *Parser) parseEnabling() error {
	enabling := plan.NewEnablingTeam(p.g.NewNode().ID())
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
			if enabling.Coords[0] == plan.UndefinedCoord {
				enabling.Coords[0] = int(f * 100)
				continue
			}
			if enabling.Coords[1] == plan.UndefinedCoord {
				enabling.Coords[1] = int(f * 100)
				continue
			}
			if enabling.Coords[2] == plan.UndefinedCoord {
				enabling.Coords[2] = int(f * 100)
				continue
			}
			if enabling.Coords[3] == plan.UndefinedCoord {
				enabling.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	enabling.Label = strings.TrimRight(b.String(), " ")
	p.g.AddNode(enabling)
	p.nodeDict[enabling.Label] = enabling
	return nil
}

func (p *Parser) parsePlatform() error {
	platform := plan.NewPlatformTeam(p.g.NewNode().ID())
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
			if platform.Coords[0] == plan.UndefinedCoord {
				platform.Coords[0] = int(f * 100)
				continue
			}
			if platform.Coords[1] == plan.UndefinedCoord {
				platform.Coords[1] = int(f * 100)
				continue
			}
			if platform.Coords[2] == plan.UndefinedCoord {
				platform.Coords[2] = int(f * 100)
				continue
			}
			if platform.Coords[3] == plan.UndefinedCoord {
				platform.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	platform.Label = strings.TrimRight(b.String(), " ")
	p.g.AddNode(platform)
	p.nodeDict[platform.Label] = platform
	return nil
}

package parser

/*
func (p *Parser) parseStreamAligned() error {
	streamAligned := wardleyToGo.NewStreamAlignedTeam(p.g.NewNode().ID())
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
			if streamAligned.Coords[0] == wardleyToGo.UndefinedCoord {
				streamAligned.Coords[0] = int(f * 100)
				continue
			}
			if streamAligned.Coords[1] == wardleyToGo.UndefinedCoord {
				streamAligned.Coords[1] = int(f * 100)
				continue
			}
			if streamAligned.Coords[2] == wardleyToGo.UndefinedCoord {
				streamAligned.Coords[2] = int(f * 100)
				continue
			}
			if streamAligned.Coords[3] == wardleyToGo.UndefinedCoord {
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
	enabling := wardleyToGo.NewEnablingTeam(p.g.NewNode().ID())
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
			if enabling.Coords[0] == wardleyToGo.UndefinedCoord {
				enabling.Coords[0] = int(f * 100)
				continue
			}
			if enabling.Coords[1] == wardleyToGo.UndefinedCoord {
				enabling.Coords[1] = int(f * 100)
				continue
			}
			if enabling.Coords[2] == wardleyToGo.UndefinedCoord {
				enabling.Coords[2] = int(f * 100)
				continue
			}
			if enabling.Coords[3] == wardleyToGo.UndefinedCoord {
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
	platform := wardleyToGo.NewPlatformTeam(p.g.NewNode().ID())
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
			if platform.Coords[0] == wardleyToGo.UndefinedCoord {
				platform.Coords[0] = int(f * 100)
				continue
			}
			if platform.Coords[1] == wardleyToGo.UndefinedCoord {
				platform.Coords[1] = int(f * 100)
				continue
			}
			if platform.Coords[2] == wardleyToGo.UndefinedCoord {
				platform.Coords[2] = int(f * 100)
				continue
			}
			if platform.Coords[3] == wardleyToGo.UndefinedCoord {
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

func (p *Parser) parseComplicatedSubsystem() error {
	complicatedSubsystem := wardleyToGo.NewComplicatedSubsystemTeam(p.g.NewNode().ID())
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
			if complicatedSubsystem.Coords[0] == wardleyToGo.UndefinedCoord {
				complicatedSubsystem.Coords[0] = int(f * 100)
				continue
			}
			if complicatedSubsystem.Coords[1] == wardleyToGo.UndefinedCoord {
				complicatedSubsystem.Coords[1] = int(f * 100)
				continue
			}
			if complicatedSubsystem.Coords[2] == wardleyToGo.UndefinedCoord {
				complicatedSubsystem.Coords[2] = int(f * 100)
				continue
			}
			if complicatedSubsystem.Coords[3] == wardleyToGo.UndefinedCoord {
				complicatedSubsystem.Coords[3] = int(f * 100)
				continue
			}
		}
	}
	complicatedSubsystem.Label = strings.TrimRight(b.String(), " ")
	p.g.AddNode(complicatedSubsystem)
	p.nodeDict[complicatedSubsystem.Label] = complicatedSubsystem
	return nil
}
*/

package parser

import (
	"strings"
	"text/scanner"
)

func (p *Parser) parseTitle() error {
	var b strings.Builder
	curLine := p.s.Pos().Line
	for tok := p.s.Scan(); tok != '\n' && tok != scanner.EOF; tok = p.s.Scan() {
		if curLine != p.s.Pos().Line {
			break
		}
		b.WriteString(p.s.TokenText())
		b.WriteString(" ")
	}
	p.title = strings.TrimRight(b.String(), " ")
	return nil
}

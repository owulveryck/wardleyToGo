package wtg2

import (
	"strings"
)

// TokenType classifies a lexed statement.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenBlank
	TokenComment
	TokenMeta
	TokenStages
	TokenAnchor
	TokenComponent
	TokenSubmap
	TokenPipeline
	TokenEdge
	TokenGroup
	TokenNote
	TokenWarning
	TokenSignal
	TokenGameplay
)

// Token represents a classified statement from the source.
type Token struct {
	Type  TokenType
	Line  int      // 1-based line number
	Text  string   // The main line text (trimmed)
	Block []string // Inner lines for block constructs ({...})
}

// Lexer performs line-oriented tokenization of WTG2 source.
type Lexer struct {
	lines []string
	pos   int
}

// NewLexer creates a lexer from the full source text.
func NewLexer(input string) *Lexer {
	return &Lexer{
		lines: strings.Split(input, "\n"),
		pos:   0,
	}
}

// Next returns the next token. Returns TokenEOF when done.
func (l *Lexer) Next() Token {
	for l.pos < len(l.lines) {
		lineNum := l.pos + 1
		raw := l.lines[l.pos]
		l.pos++

		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			continue
		}

		// Comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// Check for block: if line contains '{' but not '}', accumulate
		text, block := l.extractBlock(trimmed)

		tt := classifyLine(text)
		return Token{Type: tt, Line: lineNum, Text: text, Block: block}
	}
	return Token{Type: TokenEOF}
}

// extractBlock handles lines with '{'. If the line has '{' without matching '}',
// it accumulates subsequent lines until '}' is found.
func (l *Lexer) extractBlock(line string) (string, []string) {
	openIdx := strings.Index(line, "{")
	if openIdx < 0 {
		return line, nil
	}

	// Check if closing brace is on the same line
	closeIdx := strings.LastIndex(line, "}")
	if closeIdx > openIdx {
		header := strings.TrimSpace(line[:openIdx])
		inner := strings.TrimSpace(line[openIdx+1 : closeIdx])
		var block []string
		if inner != "" {
			block = []string{inner}
		}
		return header, block
	}

	// Multi-line block: accumulate until '}'
	header := strings.TrimSpace(line[:openIdx])
	// Content after '{' on the same line
	afterBrace := strings.TrimSpace(line[openIdx+1:])
	var block []string
	if afterBrace != "" {
		block = append(block, afterBrace)
	}

	depth := 1
	for l.pos < len(l.lines) && depth > 0 {
		bline := l.lines[l.pos]
		l.pos++
		trimmed := strings.TrimSpace(bline)

		if trimmed == "}" {
			depth--
			if depth == 0 {
				break
			}
		}
		if strings.Contains(trimmed, "{") {
			depth++
		}
		if strings.Contains(trimmed, "}") && !strings.HasPrefix(trimmed, "}") {
			depth--
			if depth == 0 {
				// Remove trailing '}'
				trimmed = strings.TrimSpace(strings.TrimSuffix(trimmed, "}"))
				if trimmed != "" {
					block = append(block, trimmed)
				}
				break
			}
		}
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			block = append(block, trimmed)
		}
	}

	return header, block
}

// classifyLine determines the token type from the first word(s) of a line.
func classifyLine(line string) TokenType {
	lower := strings.ToLower(line)

	// Keyword-prefixed statements
	switch {
	case hasKeywordPrefix(lower, "title:") ||
		hasKeywordPrefix(lower, "date:") ||
		hasKeywordPrefix(lower, "author:") ||
		hasKeywordPrefix(lower, "scope:") ||
		hasKeywordPrefix(lower, "question:"):
		return TokenMeta
	case hasKeywordPrefix(lower, "stages:"):
		return TokenStages
	case hasKeywordPrefix(lower, "anchor "):
		return TokenAnchor
	case hasKeywordPrefix(lower, "component "):
		return TokenComponent
	case hasKeywordPrefix(lower, "submap "):
		return TokenSubmap
	case hasKeywordPrefix(lower, "pipeline "):
		return TokenPipeline
	case hasKeywordPrefix(lower, "group "):
		return TokenGroup
	case hasKeywordPrefix(lower, "note "):
		return TokenNote
	case hasKeywordPrefix(lower, "warning "):
		return TokenWarning
	case hasKeywordPrefix(lower, "signal "):
		return TokenSignal
	case hasKeywordPrefix(lower, "gameplay "):
		return TokenGameplay
	case hasKeywordPrefix(lower, "doctrine:") || hasKeywordPrefix(lower, "doctrine :"):
		return TokenMeta
	case strings.Contains(line, " -> ") || strings.Contains(line, " <-> ") || strings.Contains(line, "-["):
		return TokenEdge
	case strings.Contains(line, " : "):
		// Bare component with shorthand position
		return TokenComponent
	}

	return TokenBlank
}

func hasKeywordPrefix(line, prefix string) bool {
	return strings.HasPrefix(line, prefix)
}

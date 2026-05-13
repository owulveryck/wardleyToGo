package wtg2

import (
	"strings"
	"testing"
)

func TestLexerClassification(t *testing.T) {
	input := `// comment
title: My Map
stages: A, B, C, D
anchor User
component App : III.5
submap Payment : IV.2
pipeline Engine {
  Algo A : III.5
  Algo B : II.3
}
User -> App -> Engine
group Team {
  App
  Engine
}
note "important" on App
warning "risk" on Engine
signal accelerating on App
legend
focus App
`

	lex := NewLexer(input)
	expected := []TokenType{
		TokenMeta,      // title
		TokenStages,    // stages
		TokenAnchor,    // anchor
		TokenComponent, // component
		TokenSubmap,    // submap
		TokenPipeline,  // pipeline
		TokenEdge,      // User -> App -> Engine
		TokenGroup,     // group
		TokenNote,      // note
		TokenWarning,   // warning
		TokenSignal,    // signal
		TokenLegend,    // legend
		TokenFocus,     // focus
	}

	for i, exp := range expected {
		tok := lex.Next()
		if tok.Type != exp {
			t.Errorf("token %d: got type %d, want %d (text: %q)", i, tok.Type, exp, tok.Text)
		}
	}

	tok := lex.Next()
	if tok.Type != TokenEOF {
		t.Errorf("expected EOF, got type %d (text: %q)", tok.Type, tok.Text)
	}
}

func TestLexerBlockAccumulation(t *testing.T) {
	input := `pipeline Engine {
  Algo A : III.5
  Algo B : II.3
}`

	lex := NewLexer(input)
	tok := lex.Next()

	if tok.Type != TokenPipeline {
		t.Fatalf("expected TokenPipeline, got %d", tok.Type)
	}
	if tok.Text != "pipeline Engine" {
		t.Errorf("text = %q, want %q", tok.Text, "pipeline Engine")
	}
	if len(tok.Block) != 2 {
		t.Fatalf("block has %d lines, want 2", len(tok.Block))
	}
	if tok.Block[0] != "Algo A : III.5" {
		t.Errorf("block[0] = %q, want %q", tok.Block[0], "Algo A : III.5")
	}
}

func TestLexerBareComponent(t *testing.T) {
	input := `Application Mobile : III.5`
	lex := NewLexer(input)
	tok := lex.Next()
	if tok.Type != TokenComponent {
		t.Errorf("expected TokenComponent for bare component, got %d", tok.Type)
	}
}

func TestLexerNestingDepthLimit(t *testing.T) {
	// Build input with nesting deeper than maxNestingDepth
	var sb strings.Builder
	sb.WriteString("component Root {\n")
	for i := 0; i < maxNestingDepth+5; i++ {
		sb.WriteString("nested {\n")
	}
	for i := 0; i < maxNestingDepth+5; i++ {
		sb.WriteString("}\n")
	}
	sb.WriteString("}\n")

	lex := NewLexer(sb.String())
	for {
		tok := lex.Next()
		if tok.Type == TokenEOF {
			break
		}
	}
	if lex.Err() == nil {
		t.Fatal("expected nesting depth error, got nil")
	}
	if !strings.Contains(lex.Err().Error(), "nesting depth") {
		t.Fatalf("unexpected error message: %v", lex.Err())
	}
}

func TestLexerNestingWithinLimit(t *testing.T) {
	// 3 levels of nesting — well within the limit
	input := `component A {
  evolution: III.5
  sub {
    inner {
      deep: value
    }
  }
}`
	lex := NewLexer(input)
	for {
		tok := lex.Next()
		if tok.Type == TokenEOF {
			break
		}
	}
	if lex.Err() != nil {
		t.Fatalf("unexpected error for valid nesting: %v", lex.Err())
	}
}

func TestParserNestingDepthError(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("component Root {\n")
	for i := 0; i < maxNestingDepth+5; i++ {
		sb.WriteString("nested {\n")
	}
	for i := 0; i < maxNestingDepth+5; i++ {
		sb.WriteString("}\n")
	}
	sb.WriteString("}\n")

	p, err := NewParser(strings.NewReader(sb.String()))
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Parse()
	if err == nil {
		t.Fatal("expected parser error for excessive nesting depth")
	}
}

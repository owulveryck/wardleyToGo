package wtg

import (
	"testing"
)

func TestLexer(t *testing.T) {
	src := `
	  
first identifier    second identifier
third-identifiér --   forth identifier  
third-identifiér ---   forth identifier  
blablabla

evolution: blabla
this is an evolution: |...|...|...|...| 
this is an incomplete evolution: |...|...|...|...
blabla: this is another word
block: {
	fdsfds: bdsfd
}
	`
	expectedTokens := []tokenType{
		identifierToken,      // first identifier
		identifierToken,      // second identifier
		identifierToken,      // third-identifiér
		visibilityToken,      // --
		identifierToken,      // forth identifier
		identifierToken,      // third-identifiér
		visibilityToken,      // ---
		identifierToken,      // forth identifier
		identifierToken,      // blablabla
		evolutionToken,       // evolution
		colonToken,           // :
		identifierToken,      // vlabla
		identifierToken,      // this is an evolution
		colonToken,           // :
		evolutionStringToken, // |...|...|...|...|
		identifierToken,      // this is an incomplete evolution:
		colonToken,           // :
		unkonwnToken,         // |...|...|...|...
		identifierToken,      // bkabka
		colonToken,           // :
		identifierToken,      // this is another word
		identifierToken,      // block
		colonToken,           // :
		startBlockToken,      // {
		identifierToken,      // fdsfds
		colonToken,           // :
		identifierToken,      // bdsfd
		endBlockToken,        // }
		eofToken,
	}
	l := newLexer(src, startState)
	l.Start()
	i := 0
	for tok := range l.tokens {
		if i >= len(expectedTokens) {
			t.Fatalf("bad number of test cases - missins |%v|%v|", tok.Type, tok.Value)

		}
		if tok.Type != expectedTokens[i] {
			t.Fatalf("on iteration %v, expected %v, got |%v|%v|", i, expectedTokens[i], tok.Type, tok.Value)

		}
		i++
	}
}

package wtg

import (
	"testing"
)

func TestLexer(t *testing.T) {
	src := `-
	|...|...|...|...|
--
first identifier    second identifier
third-identifiér --   forth identifier  
third-identifiér ---   forth identifier  
blablabla
identifier1
identifier1.1

type: mytype

evolution: blabla
this is an evolution: |...|...|...|...| 
this is an incomplete evolution: |...|...|...|...
blabla: this is another word
block: {
	fdsfds: bdsfd
}
test // comment on a line
// comment on another line
/*
*/
	`
	expectedTokens := []struct {
		t tokenType
		v string
	}{
		{t: visibilityToken, v: "-"},
		{t: evolutionItem, v: "|...|...|...|...|"},
		{t: visibilityToken, v: "--"},
		{t: identifierToken, v: "first identifier"},
		{t: identifierToken, v: "second identifier"},
		{t: identifierToken, v: "third-identifiér"},
		{t: visibilityToken, v: "--"},
		{t: identifierToken, v: "forth identifier"},
		{t: identifierToken, v: "third-identifiér"},
		{t: visibilityToken, v: "---"},
		{t: identifierToken, v: "forth identifier"},
		{t: identifierToken, v: "blablabla"},
		{t: identifierToken, v: "identifier1"},
		{t: identifierToken, v: "identifier1.1"},
		{t: typeToken, v: "type"},
		{t: typeItem, v: "mytype"},
		{t: evolutionToken, v: "evolution"},
		{t: colonToken, v: ":"},
		{t: identifierToken, v: "blabla"},
		{t: identifierToken, v: "this is an evolution"},
		{t: colonToken, v: ":"},
		{t: evolutionItem, v: "|...|...|...|...|"},
		{t: identifierToken, v: "this is an incomplete evolution"},
		{t: colonToken, v: ":"},
		{t: unkonwnToken, v: "|...|...|...|..."},
		{t: identifierToken, v: "blabla"},
		{t: colonToken, v: ":"},
		{t: identifierToken, v: "this is another word"},
		{t: identifierToken, v: "block"},
		{t: colonToken, v: ":"},
		{t: startBlockToken, v: "{"},
		{t: identifierToken, v: "fdsfds"},
		{t: colonToken, v: ":"},
		{t: identifierToken, v: "bdsfd"},
		{t: endBlockToken, v: "}"},
		{t: identifierToken, v: "test"},
		{t: commentToken, v: "// comment on a line"},
		{t: commentToken, v: "// comment on another line"},
		{t: startBlockCommentToken, v: "/*"},
		{t: endBlockCommentToken, v: "*/"},
		{t: eofToken, v: ""},
	}
	l := newLexer(src, startState)
	l.Start()
	i := 0
	for tok := range l.tokens {
		if i >= len(expectedTokens) {
			t.Fatalf("bad number of test cases - missing |%v|%v|", tok.Type, tok.Value)

		}
		if tok.Type != expectedTokens[i].t || tok.Value != expectedTokens[i].v {
			t.Fatalf("on iteration %v, expected %v, got %v", i, expectedTokens[i], tok)

		}
		i++
	}
}

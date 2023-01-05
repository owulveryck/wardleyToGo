package wtg

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("types", testLexerType)
	t.Run("full", testLexerFull)
	t.Run("bad content", testLexer0)
}
func testLexer0(t *testing.T) {
	t.SkipNow()
	src := `0 \x8a 0`
	expectedTokens := []struct {
		t tokenType
		v string
	}{
		{t: identifierToken, v: "0"},
		{t: unkonwnToken, v: "\x8a"},
		{t: identifierToken, v: "0"},
	}
	l := newLexer(src, startState)
	l.Start()
	i := 0
	var tok token
	for tok = range l.tokens {
		if i >= len(expectedTokens) {
			t.Fatalf("bad number of test cases - missing |%v|%v|", tok.Type, tok.Value)

		}
		if tok.Type != expectedTokens[i].t || tok.Value != expectedTokens[i].v {
			t.Fatalf("on iteration %v, expected (%v,%v), got (%v,%v)", i, expectedTokens[i].t, []byte(expectedTokens[i].v), tok.Type, []byte(tok.Value))

		}
		i++
	}
	if i < len(expectedTokens) {
		t.Errorf("bad number of elements, expected %v got %v (last token is %v)", len(expectedTokens), i, tok)
	}
}
func testLexerType(t *testing.T) {
	src := `
build: {
	type: build
}
buy: {
	type: buy
}
outsource: {
	type: outsource
}
	`
	expectedTokens := []struct {
		t tokenType
		v string
	}{
		{t: identifierToken, v: "build"},
		{t: colonToken, v: ":"},
		{t: startBlockToken, v: "{"},
		{t: typeToken, v: "type"},
		{t: colonToken, v: ":"},
		{t: typeItem, v: "build"},
		{t: endBlockToken, v: "}"},
		{t: identifierToken, v: "buy"},
		{t: colonToken, v: ":"},
		{t: startBlockToken, v: "{"},
		{t: typeToken, v: "type"},
		{t: colonToken, v: ":"},
		{t: typeItem, v: "buy"},
		{t: endBlockToken, v: "}"},
		{t: identifierToken, v: "outsource"},
		{t: colonToken, v: ":"},
		{t: startBlockToken, v: "{"},
		{t: typeToken, v: "type"},
		{t: colonToken, v: ":"},
		{t: typeItem, v: "outsource"},
		{t: endBlockToken, v: "}"},
		{t: eofToken, v: ""},
	}
	l := newLexer(src, startState)
	l.Start()
	i := 0
	var tok token
	for tok = range l.tokens {
		if i >= len(expectedTokens) {
			t.Fatalf("bad number of test cases - missing |%v|%v|", tok.Type, tok.Value)

		}
		if tok.Type != expectedTokens[i].t || tok.Value != expectedTokens[i].v {
			t.Fatalf("on iteration %v, expected %v, got %v", i, expectedTokens[i], tok)

		}
		i++
	}
	if i < len(expectedTokens) {
		t.Errorf("bad number of elements, expected %v got %v (last token is %v)", len(expectedTokens), i, tok)
	}
}
func testLexerFull(t *testing.T) {
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
   type: mytype
	fdsfds: bdsfd
	color: red
}
test // comment on a line
// comment on another line
/* comment
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
		{t: colonToken, v: ":"},
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
		{t: typeToken, v: "type"},
		{t: colonToken, v: ":"},
		{t: typeItem, v: "mytype"},
		{t: identifierToken, v: "fdsfds"},
		{t: colonToken, v: ":"},
		{t: identifierToken, v: "bdsfd"},
		{t: colorToken, v: "color"},
		{t: colonToken, v: ":"},
		{t: colorItem, v: "red"},
		{t: endBlockToken, v: "}"},
		{t: identifierToken, v: "test"},
		{t: singleLineCommentSeparator, v: "//"},
		{t: commentToken, v: " comment on a line"},
		{t: singleLineCommentSeparator, v: "//"},
		{t: commentToken, v: " comment on another line"},
		{t: startBlockCommentToken, v: "/*"},
		{t: commentToken, v: " comment\n"},
		{t: endBlockCommentToken, v: "*/"},
		{t: eofToken, v: ""},
	}
	l := newLexer(src, startState)
	l.Start()
	i := 0
	var tok token
	for tok = range l.tokens {
		if i >= len(expectedTokens) {
			t.Fatalf("bad number of test cases - missing |%v|%v|", tok.Type, tok.Value)

		}
		if tok.Type != expectedTokens[i].t || tok.Value != expectedTokens[i].v {
			t.Errorf("on iteration %v, expected %v, got %v", i, expectedTokens[i].t, tok.Type)
			t.Errorf("on iteration %v, expected %v, got %v", i, []byte(expectedTokens[i].v), []byte(tok.Value))

		}
		i++
	}
	if i < len(expectedTokens) {
		t.Errorf("bad number of elements, expected %v got %v (last token is %v)", len(expectedTokens), i, tok)
	}

}
func getFunctionName(i interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

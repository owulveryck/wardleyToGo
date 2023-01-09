package wtg

import (
	"sync"
	"testing"
)

func TestFirstRuneAfterSpaceState(t *testing.T) {
	t.Run("star", func(t *testing.T) {
		tst := `*aaa/ blabla`
		l := newLexer(tst, firstRuneAfterSpaceState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != unknownToken {
			t.Fatalf("expected %v, got %v", endBlockCommentToken, tok)
		}
		if tok.Value != "*" {
			t.Fatalf("expected *, got %v", tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})
	t.Run("/rubbish ", func(t *testing.T) {
		tst := `/rubbish`
		l := newLexer(tst, firstRuneAfterSpaceState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != unknownToken {
			t.Fatalf("expected %v, got %v", startBlockCommentToken, tok)
		}
		if tok.Value != "/" {
			t.Fatalf("expected /, got %v", tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})
}
func TestCommentBlockState(t *testing.T) {
	t.Run("comment ok single line ", func(t *testing.T) {
		tst := `comment */`
		l := newLexer(tst, commentBlockState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != commentToken {
			t.Fatalf("expected %v, got %v", commentToken, tok)
		}
		if tok.Value != "comment" {
			t.Fatalf("expected comment , got %v", tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})
	t.Run("comment ok multiple lines ", func(t *testing.T) {
		comment := `comment
		dsdsa`
		tst := comment + ` */`
		l := newLexer(tst, commentBlockState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != commentToken {
			t.Fatalf("expected %v, got %v", commentToken, tok)
		}
		if tok.Value != comment {
			t.Fatalf("expected %v , got %v", comment, tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})
	t.Run("comment ok multiple lines no space ", func(t *testing.T) {
		comment := `comment
		dsdsa`
		tst := comment + `*/`
		l := newLexer(tst, commentBlockState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != commentToken {
			t.Fatalf("expected %v, got %v", commentToken, tok)
		}
		if tok.Value != comment {
			t.Fatalf("expected %v , got %v", comment, tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})
	t.Run("unfinished comment", func(t *testing.T) {
		comment := `comment
		dsdsa`
		tst := comment
		l := newLexer(tst, commentBlockState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "" {
			t.Fatalf("expected nil func, got %v", retName)

		}
		if tok.Type != commentToken {
			t.Fatalf("expected %v, got %v", commentToken, tok)
		}
		if tok.Value != comment {
			t.Fatalf("expected %v , got %v", comment, tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})

}

func TestOneLineCommentState(t *testing.T) {
	t.Run("single line comment", func(t *testing.T) {
		tst := ` blabla
		ahah`
		l := newLexer(tst, oneLineCommentState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != commentToken {
			t.Fatalf("expected %v, got %v", commentToken, tok)
		}
		if tok.Value != " blabla" {
			t.Fatalf("expected blabla got %v", tok)
		}
		if l.Current() != "" {
			t.Errorf("expected buffer to be empty")
		}
	})

}

func TestEvolutionState(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		tst := `|...|...|...|...| `
		l := newLexer(tst, evolutionState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != evolutionItem {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "|...|...|...|...|" {
			t.Fatalf("expected %v, got %v", "|...|...|...|...|", tok.Value)
		}
		if l.Current() != "" {
			t.Fatalf("buffer shoud be empty but contains %v", l.Current())
		}
	})
	t.Run("complete 2", func(t *testing.T) {
		tst := `|x|||| `
		l := newLexer(tst, evolutionState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != evolutionItem {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "|x||||" {
			t.Fatalf("expected %v, got %v", "|x||||", tok.Value)
		}
		if l.Current() != "" {
			t.Fatalf("buffer shoud be empty but contains %v", l.Current())
		}
	})
	t.Run("complete 3", func(t *testing.T) {
		tst := `|x|>||| `
		l := newLexer(tst, evolutionState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != evolutionItem {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "|x|>|||" {
			t.Fatalf("expected %v, got %v", "|x||||", tok.Value)
		}
		if l.Current() != "" {
			t.Fatalf("buffer shoud be empty but contains %v", l.Current())
		}
	})
	t.Run("bad", func(t *testing.T) {
		tst := `|x|x||| `
		l := newLexer(tst, evolutionState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != unknownToken {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "|x|x|||" {
			t.Fatalf("expected %v, got %v", "|x||||", tok.Value)
		}
		if l.Current() != "" {
			t.Fatalf("buffer shoud be empty but contains %v", l.Current())
		}
	})
	t.Run("bad 2", func(t *testing.T) {
		tst := `|x||||| `
		l := newLexer(tst, evolutionState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != unknownToken {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "|x|||||" {
			t.Fatalf("expected %v, got %v", "|x||||", tok.Value)
		}
		if l.Current() != "" {
			t.Fatalf("buffer shoud be empty but contains %v", l.Current())
		}
	})

}

func TestVisibilityState(t *testing.T) {
	t.Run("two dashes", func(t *testing.T) {
		tst := `--

		`
		l := newLexer(tst, visibilityState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != visibilityToken {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "--" {
			t.Fatalf("expected --, got %v", tok.Value)
		}
		if l.Current() != "" {
			t.Fatal("buffer should be empty")
		}
	})
	t.Run("two dashesa", func(t *testing.T) {
		tst := `--a

		`
		l := newLexer(tst, visibilityState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "startState" {
			t.Fatalf("expected startState func, got %v", retName)

		}
		if tok.Type != visibilityToken {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
		if tok.Value != "--" {
			t.Fatalf("expected --, got %v", tok.Value)
		}
		if l.Current() != "" {
			t.Fatal("buffer should be empty")
		}
	})

}

func TestStartState(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		tst := ``
		l := newLexer(tst, startState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "" {
			t.Fatalf("expected nil func, got %v", retName)
		}
		if tok.Type != eofToken {
			t.Fatalf("expected %v, got %v", eofToken, tok)
		}
	})
	t.Run("space only file", func(t *testing.T) {
		tst := `       `
		l := newLexer(tst, startState)
		l.tokens = make(chan token)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "" {
			t.Fatalf("expected nil func, got %v", retName)

		}
		if tok.Type != eofToken {
			t.Fatalf("expected %v, got %v", eofToken, tok)
		}
	})
	t.Run("several empty lines", func(t *testing.T) {
		tst := `       



		`
		l := newLexer(tst, startState)
		l.tokens = make(chan token)
		defer close(l.tokens)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			tok := <-l.tokens
			if tok.Type != eofToken {
				t.Errorf("expected %v, got %v", eofToken, tok)
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "" {
			t.Fatalf("expected nil func, got %v", retName)

		}
		wg.Wait()
	})
	t.Run("one non space", func(t *testing.T) {
		tst := `       

    é

		`
		l := newLexer(tst, startState)
		l.tokens = make(chan token)
		defer close(l.tokens)
		var tok token
		go func() {
			for tok = range l.tokens {
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != "firstRuneAfterSpaceState" {
			t.Fatalf("expected firstRuneAfterSpaceState func, got %v", retName)

		}
		if tok.Type != 0 {
			t.Fatalf("expected %v, got %v", 0, tok)
		}
	})

}

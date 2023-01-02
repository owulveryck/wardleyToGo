package wtg

import (
	"testing"
	"unicode"
)

func startState(l *lexer) stateFunc {
	l.Next() // eat starting "
	l.Ignore()
	if l.Peek() == EOFRune {
		l.Emit(unkonwnToken)
		return nil
	}
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
	if unicode.IsLetter(l.Peek()) {
		return identifierState
	}
	return startState
}

func identifierState(l *lexer) stateFunc {
	l.Next() // eat starting "
	for {
		switch {
		case l.Peek() == '-':
			if unicode.IsSpace(l.CurrentRune()) {
				l.Rewind()
				l.Emit(identifierToken)
				return visibilityAction
			}
			l.Next()
		case unicode.IsSpace(l.Peek()):
			if unicode.IsSpace(l.CurrentRune()) {
				l.Rewind()
				l.Emit(identifierToken)
				return startState
			}
			l.Next()
		case l.Peek() == EOFRune:
			l.Emit(identifierToken)
			return nil
		default:
			l.Next()
		}
	}
}

func visibilityAction(l *lexer) stateFunc {
	return startState
}

func TestLexer(t *testing.T) {
	src := `
first identifier    second identifier
third-identifiér --   forth identifier  
blablabla
	`
	l := newLexer(src, startState)
	l.Start()
	for tok := range l.tokens {
		//t.Log(tok.Type)
		t.Logf("|%v|", tok.Value)
	}
}

package wtg

import "unicode"

func startState(l *lexer) stateFunc {
	l.Next() // eat starting "
	l.Ignore()
	if l.Peek() == eofRune {
		l.Emit(eofToken)
		return nil
	}
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
	l.Ignore()
	if unicode.IsLetter(l.Peek()) {
		return wordState
	}
	if l.Peek() == '|' {
		return evolutionState
	}
	if l.Peek() == '{' {
		l.Next()
		l.Emit(startBlockToken)
	}
	if l.Peek() == '}' {
		l.Next()
		l.Emit(endBlockToken)
	}
	return startState
}

func wordState(l *lexer) stateFunc {
	for unicode.IsLetter(l.Peek()) {
		l.Next()
	}
	return postWordState
}

func postWordState(l *lexer) stateFunc {
	switch l.Peek() {
	case '-':
		l.Next()
		return wordState
	case ' ':
		if unicode.IsLetter(l.PeekPeek()) {
			l.Next()
			return wordState
		}
		if l.PeekPeek() == '-' {
			l.Emit(identifierToken)
			l.Next()
			l.Ignore()
			return visibilityState
		}
	case ':':
		return assignationState

	}
	l.Emit(identifierToken)
	return startState
}

func assignationState(l *lexer) stateFunc {
	switch l.Current() {
	case "evolution":
		l.Emit(evolutionToken)
		// TODO: add mor keywords here
	default:
		l.Emit(identifierToken)
	}
	l.Ignore()
	l.Next()
	l.Emit(colonToken)
	l.Ignore()
	return startState
}

func visibilityState(l *lexer) stateFunc {
	for l.Peek() == '-' {
		l.Next()
	}
	l.Emit(visibilityToken)
	return startState
}

func evolutionState(l *lexer) stateFunc {
	l.Next()
	states := 0
	for states < 4 {
		p := l.Peek()
		switch {
		case p == '|':
			states++
		case unicode.IsSpace(p):
			l.Emit(unkonwnToken)
			return startState
		}
		l.Next()
	}
	l.Emit(evolutionStringToken)
	return startState
}

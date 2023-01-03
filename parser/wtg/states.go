package wtg

import "unicode"

func startState(l *lexer) stateFunc {
	//l.Next() // eat starting "
	//l.Ignore()
	if l.Peek() == eofRune {
		l.Emit(eofToken)
		return nil
	}
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
	l.Ignore()
	if isAllowedCharacterForIdentifier(l.Peek()) {
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
	if l.Peek() == '/' && l.PeekPeek() == '*' {
		l.Next()
		l.Next()
		l.Emit(startBlockCommentToken)
	}
	if l.Peek() == '*' && l.PeekPeek() == '/' {
		l.Next()
		l.Next()
		l.Emit(endBlockCommentToken)
	}
	if l.Peek() == '/' && l.PeekPeek() == '/' {
		for l.CurrentRune() != '\n' {
			l.Next()
		}
		l.Rewind()
		l.Emit(commentToken)
	}
	if l.Peek() == '-' {
		return visibilityState
	}
	return startState
}

func wordState(l *lexer) stateFunc {
	for isAllowedCharacterForIdentifier(l.Peek()) {
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
		if isAllowedCharacterForIdentifier(l.PeekPeek()) {
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
	case "type":
		l.Emit(typeToken)
		// TODO: add mor keywords here
		return typeState
	default:
		l.Emit(identifierToken)
	}
	l.Ignore()
	l.Next()
	l.Emit(colonToken)
	l.Ignore()
	return startState
}

func typeState(l *lexer) stateFunc {
	if l.Peek() != ':' {
		return startState
	}
	l.Next()
	for l.Peek() == ' ' {
		l.Next()
	}
	l.Ignore()
	for unicode.IsLetter(l.Peek()) {
		l.Next()
	}
	l.Emit(typeItem)
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
	l.Emit(evolutionItem)
	return startState
}

func isAllowedCharacterForIdentifier(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' {
		return true
	}
	return false
}

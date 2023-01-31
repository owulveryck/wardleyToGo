package wtg

import (
	"fmt"
	"strings"
	"unicode"
)

func startState(l *lexer) stateFunc {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
	l.Ignore()
	if l.Peek() == eofRune {
		l.Emit(eofToken)
		return nil
	}
	return firstRuneAfterSpaceState
}

func commentBlockState(l *lexer) stateFunc {
	for l.Peek() != '*' || l.Peek() != eofRune {
		if l.Peek() == eofRune {
			l.Emit(commentToken)
			return nil
		}
		l.Next()
		if l.PeekPeek() == '/' {
			if l.CurrentRune() == ' ' {
				l.Rewind()
			}
			l.Emit(commentToken)
			l.Ignore()
			return startState
		}
	}
	return startState
}
func oneLineCommentState(l *lexer) stateFunc {
	for l.Peek() != '\n' && l.Peek() != eofRune {
		l.Next()
	}
	l.Emit(commentToken)
	l.Ignore()
	return startState
}

func firstRuneAfterSpaceState(l *lexer) stateFunc {
	l.Next()
	switch l.CurrentRune() {
	case '-':
		return visibilityState
	case '|':
		return evolutionState
	case ':':
		l.Emit(colonToken)
		l.Ignore()
		return startState
	case '{':
		l.Emit(startBlockToken)
		l.Ignore()
		return startState
	case '}':
		l.Emit(endBlockToken)
		l.Ignore()
		return startState
	case '/':
		if l.Peek() == '*' {
			l.Next()
			l.Emit(startBlockCommentToken)
			l.Ignore()
			return commentBlockState
		}
		if l.Peek() == '/' {
			l.Next()
			l.Emit(singleLineCommentSeparator)
			l.Ignore()
			return oneLineCommentState
		}
		l.Emit(unknownToken)
		if l.Current() == "" {
			l.Next()
		}
		l.Ignore()
		return startState
	case '*':
		if l.Peek() == '/' {
			l.Next()
			l.Emit(endBlockCommentToken)
			l.Ignore()
			return startState
		}
		l.Emit(unknownToken)
		if l.Current() == "" {
			l.Next()
		}
		l.Ignore()
		return startState
	case eofRune:
		return nil
	default:
		return wordState
	}
}

func wordState(l *lexer) stateFunc {
	for isAllowedCharacterForIdentifier(l.Peek()) {
		if l.Peek() == ' ' && l.PeekPeek() == ' ' ||
			l.Peek() == ' ' && l.PeekPeek() == '-' ||
			l.Peek() == ' ' && !isAllowedCharacterForIdentifier(l.PeekPeek()) {
			break
		}
		l.Next()
		if l.Peek() == eofRune {
			break
		}
	}
	switch l.Current() {
	case "stage1":
		return stageState
	case "stage2":
		return stageState
	case "stage3":
		return stageState
	case "stage4":
		return stageState
	case "":
		// this is probably a control character
		l.Emit(unknownToken)
		if l.Current() == "" {
			l.Next()
		}
		l.Next()
	case "color":
		return colorState
	case "label":
		return labelState
	case "title":
		return titleState
	case "type":
		return typeState
	case "evolution":
		l.Emit(evolutionToken)
	default:
		l.Emit(identifierToken)
	}
	l.Ignore()
	return startState
}

func stageState(l *lexer) stateFunc {
	stage := 0
	if l.Peek() == ':' {
		switch l.Current() {
		case "stage1":
			stage = 1
			l.Emit(stage1Token)
		case "stage2":
			stage = 2
			l.Emit(stage2Token)
		case "stage3":
			stage = 3
			l.Emit(stage3Token)
		case "stage4":
			stage = 4
			l.Emit(stage4Token)
		}
		l.Ignore()
		l.Next()
		l.Emit(colonToken)
		// discard the leading space
		for unicode.IsSpace(l.Peek()) {
			l.Next()
		}
		l.Ignore()
		for l.Peek() != '\n' && !(l.Peek() == '/' && l.PeekPeek() == '/') && unicode.IsGraphic(l.Peek()) {
			l.Next()
		}
		switch stage {
		case 1:
			l.Emit(stage1Item)
		case 2:
			l.Emit(stage2Item)
		case 3:
			l.Emit(stage3Item)
		case 4:
			l.Emit(stage4Item)
		}
		l.Ignore()
		return startState
	}
	l.Emit(identifierToken)
	l.Ignore()

	return startState
}

func typeState(l *lexer) stateFunc {
	if l.Peek() == ':' {
		l.Emit(typeToken)
		l.Ignore()
		l.Next()
		l.Emit(colonToken)
		// discard the leading space
		for unicode.IsSpace(l.Peek()) {
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
	l.Emit(identifierToken)
	l.Ignore()

	return startState
}

func titleState(l *lexer) stateFunc {
	if l.Peek() == ':' {
		l.Emit(titleToken)
		l.Ignore()
		l.Next()
		l.Emit(colonToken)
		// discard the leading space
		for unicode.IsSpace(l.Peek()) {
			l.Next()
		}
		l.Ignore()
		for l.Peek() != '\n' && !(l.Peek() == '/' && l.PeekPeek() == '/') && l.Peek() != eofRune {
			l.Next()
		}
		l.Emit(titleItem)
		l.Ignore()
		return startState
	}
	l.Emit(identifierToken)
	l.Ignore()

	return startState
}

func labelState(l *lexer) stateFunc {
	if l.Peek() == ':' {
		l.Emit(labelToken)
		l.Ignore()
		l.Next()
		l.Emit(colonToken)
		// discard the leading space
		for unicode.IsSpace(l.Peek()) {
			l.Next()
		}
		l.Ignore()
		for unicode.IsLetter(l.Peek()) || unicode.IsDigit(l.Peek()) {
			l.Next()
		}
		l.Emit(labelItem)
		l.Ignore()
		return startState
	}
	l.Emit(identifierToken)
	l.Ignore()

	return startState
}

func colorState(l *lexer) stateFunc {
	if l.Peek() == ':' {
		l.Emit(colorToken)
		l.Ignore()
		l.Next()
		l.Emit(colonToken)
		// discard the leading space
		for unicode.IsSpace(l.Peek()) {
			l.Next()
		}
		l.Ignore()
		for unicode.IsLetter(l.Peek()) || unicode.IsDigit(l.Peek()) {
			l.Next()
		}
		l.Emit(colorItem)
		l.Ignore()
		return startState
	}
	l.Emit(identifierToken)
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
	l.Take("|.x>")
	if strings.Count(l.Current(), "|") != 5 ||
		strings.Count(l.Current(), "x") > 1 ||
		strings.Count(l.Current(), ">") > 1 {
		l.Emit(unknownToken)
		if l.Current() == "" {
			l.Next()
		}
		l.Ignore()
		return startState
	}
	if l.CurrentRune() != '|' {
		l.Err = fmt.Errorf("bad evolution %v", l.Current())
		return nil
	}
	l.Emit(evolutionItem)
	l.Ignore()
	return startState
}

func isAllowedCharacterForIdentifier(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == ' ' || r == '-' {
		return true
	}
	return false
}

package wtg

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type stateFunc func(*lexer) stateFunc

type tokenType int

const (
	EOFRune    rune      = -1
	emptyToken tokenType = iota
	identifierToken
	visibilityToken // multiple dashes token
	evolutionStringToken
	evolutionToken // evolution keywork
	unkonwnToken
)

type token struct {
	Type  tokenType
	Value string
}

type lexer struct {
	source          string
	start, position int
	startState      stateFunc
	Err             error
	tokens          chan token
	ErrorHandler    func(e string)
	rewind          runeStack
}

// newLexer creates a returns a lexer ready to parse the given source code.
func newLexer(src string, start stateFunc) *lexer {
	return &lexer{
		source:     src,
		startState: start,
		start:      0,
		position:   0,
		rewind:     newRuneStack(),
	}
}

// Start begins executing the Lexer in an asynchronous manner (using a goroutine).
func (l *lexer) Start() {
	// Take half the string length as a buffer size.
	buffSize := len(l.source) / 2
	if buffSize <= 0 {
		buffSize = 1
	}
	l.tokens = make(chan token, buffSize)
	go l.run()
}

func (l *lexer) StartSync() {
	// Take half the string length as a buffer size.
	buffSize := len(l.source) / 2
	if buffSize <= 0 {
		buffSize = 1
	}
	l.tokens = make(chan token, buffSize)
	l.run()
}

// Current returns the value being being analyzed at this moment.
func (l *lexer) Current() string {
	return l.source[l.start:l.position]
}

func (l *lexer) CurrentRune() rune {
	r, _ := utf8.DecodeLastRuneInString(l.Current())
	return r
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (l *lexer) Emit(t tokenType) {
	tok := token{
		Type:  t,
		Value: l.Current(),
	}
	l.tokens <- tok
	l.start = l.position
	l.rewind.clear()
}

// Ignore clears the rewind stack and then sets the current beginning position
// to the current position in the source which effectively ignores the section
// of the source being analyzed.
func (l *lexer) Ignore() {
	l.rewind.clear()
	l.start = l.position
}

// Peek performs a Next operation immediately followed by a Rewind returning the
// peeked rune.
func (l *lexer) Peek() rune {
	r := l.Next()
	l.Rewind()

	return r
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can
// occur more than once per call to Next but you can never rewind past the
// last point a token was emitted.
func (l *lexer) Rewind() {
	r := l.rewind.pop()
	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.position -= size
		if l.position < l.start {
			l.position = l.start
		}
	}
}

// Next pulls the next rune from the Lexer and returns it, moving the position
// forward in the source.
func (l *lexer) Next() rune {
	var (
		r rune
		s int
	)
	str := l.source[l.position:]
	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}
	l.position += s
	l.rewind.push(r)

	return r
}

// Take receives a string containing all acceptable strings and will contine
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *lexer) Take(chars string) {
	r := l.Next()
	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}
	l.Rewind() // last next wasn't a match
}

// NextToken returns the next token from the lexer and a value to denote whether
// or not the token is finished.
func (l *lexer) NextToken() (*token, bool) {
	if tok, ok := <-l.tokens; ok {
		return &tok, false
	} else {
		return nil, true
	}
}

// Partial yyLexer implementation

func (l *lexer) Error(e string) {
	if l.ErrorHandler != nil {
		l.Err = errors.New(e)
		l.ErrorHandler(e)
	} else {
		panic(e)
	}
}

// Private methods

func (l *lexer) run() {
	state := l.startState
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
}
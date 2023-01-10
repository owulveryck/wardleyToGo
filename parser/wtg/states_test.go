package wtg

import (
	"reflect"
	"sync"
	"testing"
)

func TestWordState(t *testing.T) {
	evaluatedState := wordState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "complete",
			corpus:         `abcde`,
			expectedTokens: []token{{identifierToken, "abcde"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestEvolutionState(t *testing.T) {
	evaluatedState := evolutionState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "complete",
			corpus:         `|...|.x.|...|...| `,
			expectedTokens: []token{{evolutionItem, "|...|.x.|...|...|"}},
			expectedState:  "startState",
		},
		{
			title:          "complete 2",
			corpus:         `|x|||| `,
			expectedTokens: []token{{evolutionItem, "|x||||"}},
			expectedState:  "startState",
		},
		{
			title:          "complete 3",
			corpus:         `|x|>||| `,
			expectedTokens: []token{{evolutionItem, "|x|>|||"}},
			expectedState:  "startState",
		},
		{
			title:          "bad",
			corpus:         `|x|x||| `,
			expectedTokens: []token{{unknownToken, "|x|x|||"}},
			expectedState:  "startState",
		},
		{
			title:          "bad 2",
			corpus:         `|x||||| `,
			expectedTokens: []token{{unknownToken, "|x|||||"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestStartState(t *testing.T) {
	evaluatedState := startState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "empty file",
			corpus:         ``,
			expectedTokens: []token{{eofToken, ""}},
			expectedState:  "",
		},
		{
			title:          "space only file",
			corpus:         `       `,
			expectedTokens: []token{{eofToken, ""}},
			expectedState:  "",
		},
		{
			title:          "several empty lines",
			corpus:         "\n\n\n\n",
			expectedTokens: []token{{eofToken, ""}},
			expectedState:  "",
		},
		{
			title:          "one non space",
			corpus:         "\n\n\né\n",
			expectedTokens: []token{},
			expectedState:  "firstRuneAfterSpaceState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestVisibilityState(t *testing.T) {
	evaluatedState := visibilityState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "one dash",
			corpus:         `-`,
			expectedTokens: []token{{visibilityToken, "-"}},
			expectedState:  "startState",
		},
		{
			title:          "two dashes",
			corpus:         `--`,
			expectedTokens: []token{{visibilityToken, "--"}},
			expectedState:  "startState",
		},
		{
			title:          "two dashesa",
			corpus:         `--a`,
			expectedTokens: []token{{visibilityToken, "--"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestOneLineCommentState(t *testing.T) {
	evaluatedState := oneLineCommentState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "single line comment space",
			corpus:         ` blabla`,
			expectedTokens: []token{{commentToken, " blabla"}},
			expectedState:  "startState",
		},
		{
			title:          "single line comment no space",
			corpus:         `blabla`,
			expectedTokens: []token{{commentToken, "blabla"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestCommentBlockState(t *testing.T) {
	evaluatedState := commentBlockState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "comment ok on single line",
			corpus:         `comment */`,
			expectedTokens: []token{{commentToken, "comment"}},
			expectedState:  "startState",
		},
		{
			title:          "comment ok on multiple line",
			corpus:         string("comment\ncomment */"),
			expectedTokens: []token{{commentToken, "comment\ncomment"}},
			expectedState:  "startState",
		},
		{
			title:          "comment ok on multiple line no space",
			corpus:         string("comment\ncomment*/"),
			expectedTokens: []token{{commentToken, "comment\ncomment"}},
			expectedState:  "startState",
		},
		{
			title:          "unfinished comment",
			corpus:         string("comment\ncomment"),
			expectedTokens: []token{{commentToken, "comment\ncomment"}},
			expectedState:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func TestFirstRuneAfterSpace(t *testing.T) {
	evaluatedState := firstRuneAfterSpaceState
	type testcase struct {
		title          string
		corpus         string
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title:          "single word",
			corpus:         `abcde`,
			expectedTokens: []token{},
			expectedState:  "wordState",
		},
		{
			title: "one dash",
			corpus: `-
		`,
			expectedTokens: []token{},
			expectedState:  "visibilityState",
		},
		{
			title: "two dashes",
			corpus: `--
		`,
			expectedTokens: []token{},
			expectedState:  "visibilityState",
		},
		{
			title:          "pipe",
			corpus:         `|-`,
			expectedTokens: []token{},
			expectedState:  "evolutionState",
		},
		{
			title:          "startBlock",
			corpus:         `{-`,
			expectedTokens: []token{{startBlockToken, "{"}},
			expectedState:  "startState",
		},
		{
			title:          "endBlock",
			corpus:         `}-`,
			expectedTokens: []token{{endBlockToken, "}"}},
			expectedState:  "startState",
		},
		{
			title:          "single line comment",
			corpus:         `// blabla`,
			expectedTokens: []token{{singleLineCommentSeparator, "//"}},
			expectedState:  "oneLineCommentState",
		},
		{
			title:          "block comment",
			corpus:         `/* blabla`,
			expectedTokens: []token{{startBlockCommentToken, "/*"}},
			expectedState:  "commentBlockState",
		},
		{
			title:          "end block comment",
			corpus:         `*/ blabla`,
			expectedTokens: []token{{endBlockCommentToken, "*/"}},
			expectedState:  "startState",
		},
		{
			title:          "star",
			corpus:         `*aaa/ blabla`,
			expectedTokens: []token{{unknownToken, "*"}},
			expectedState:  "startState",
		},
		{
			title:          "slashRubbish",
			corpus:         `/rubbish`,
			expectedTokens: []token{{unknownToken, "/"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, stateEvaluation(test.corpus, evaluatedState, test.expectedTokens, test.expectedState))
	}
}

func stateEvaluation(corpus string, evaluatedState stateFunc, expectedTokens []token, expectedState string) func(t *testing.T) {
	return func(t *testing.T) {
		l := newLexer(corpus, evaluatedState)
		l.tokens = make(chan token, len(expectedTokens))
		defer close(l.tokens)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			for i := 0; i < len(expectedTokens); i++ {
				tok := <-l.tokens
				t.Logf("received: %v", tok)
				if !reflect.DeepEqual(tok, expectedTokens[i]) {
					t.Errorf("expected %v, got %v", expectedTokens[i], tok)
				}
			}
		}()
		ret := l.startState(l)
		retName := getFunctionName(ret)
		if retName != expectedState {
			t.Errorf("expected %v func, got %v", expectedState, retName)

		}
		wg.Wait()

	}
}

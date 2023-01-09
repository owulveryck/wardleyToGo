package wtg

import (
	"reflect"
	"sync"
	"testing"
)

func TestSpaces(t *testing.T) {
	type testcase struct {
		title          string
		corpus         string
		evaluatedState stateFunc
		expectedState  string
		expectedTokens []token
	}
	tests := []testcase{
		{
			title: "firstRuneAfterSpace/one dash",
			corpus: `-
		`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{},
			expectedState:  "visibilityState",
		},
		{
			title: "firstRuneAfterSpace/two dashes",
			corpus: `--
		`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{},
			expectedState:  "visibilityState",
		},
		{
			title:          "firstRuneAfterSpace/pipe",
			corpus:         `|-`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{},
			expectedState:  "evolutionState",
		},
		{
			title:          "firstRuneAfterSpace/startBlock",
			corpus:         `{-`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{{startBlockToken, "{"}},
			expectedState:  "startState",
		},
		{
			title:          "firstRuneAfterSpace/endBlock",
			corpus:         `}-`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{{endBlockToken, "}"}},
			expectedState:  "startState",
		},
		{
			title:          "firstRuneAfterSpace/single line comment",
			corpus:         `// blabla`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{{singleLineCommentSeparator, "//"}},
			expectedState:  "oneLineCommentState",
		},
		{
			title:          "firstRuneAfterSpace/block comment",
			corpus:         `/* blabla`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{{startBlockCommentToken, "/*"}},
			expectedState:  "commentBlockState",
		},
		{
			title:          "firstRuneAfterSpace/end block comment",
			corpus:         `*/ blabla`,
			evaluatedState: firstRuneAfterSpaceState,
			expectedTokens: []token{{endBlockCommentToken, "*/"}},
			expectedState:  "startState",
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			l := newLexer(test.corpus, test.evaluatedState)
			l.tokens = make(chan token, len(test.expectedTokens))
			defer close(l.tokens)
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				for i := 0; i < len(test.expectedTokens); i++ {
					tok := <-l.tokens
					t.Logf("received: %v", tok)
					if !reflect.DeepEqual(tok, test.expectedTokens[i]) {
						t.Errorf("expected %v, got %v", test.expectedTokens[i], tok)
					}
				}
			}()
			ret := l.startState(l)
			retName := getFunctionName(ret)
			if retName != test.expectedState {
				t.Errorf("expected %v func, got %v", test.expectedState, retName)

			}
			wg.Wait()

		})
	}
}

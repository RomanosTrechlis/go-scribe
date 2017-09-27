package streamer

import "testing"

func TestLexicalToNumber(t *testing.T) {
	type testCase struct {
		lexical string
		number  int64
		err     bool
	}
	var tc = []testCase{
		{"100KB", 100000, false},
		{"5MB", 5000000, false},
		{"0.01GB", 10000000, false},
		{"0.01TB", 10000000000, false},
		{"0AB", 0, true},
		{"A0.1MB", 0, true},
		{"-1", -1, false},
	}
	for _, s := range tc {
		i, err := LexicalToNumber(s.lexical)
		if s.err && err == nil {
			t.Errorf("For %s expected an error, but got %d", s.lexical, i)
		}
		if !s.err && err != nil {
			t.Errorf("For %s expected %d, but got err '%v'", s.lexical, s.number, err)
		}
		if i != s.number {
			t.Errorf("For %s expected %d, but got %d", s.lexical, s.number, i)
		}
	}
}

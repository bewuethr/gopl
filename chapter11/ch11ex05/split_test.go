package ch11ex05

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	var tests = []struct {
		s, sep string
		want   int
	}{
		{"a:b:c", ":", 3},
		{"a:b:", ":", 3},
		{":a:b", ":", 3},
		{"abc", "", 3},
		{":", ":", 2},
		{"", "", 0},
		{"a::b:c", "::", 2},
		{"a:b:c", ",", 1},
		{"", ":", 1},
		{"abc", "abcd", 1},
	}

	for _, test := range tests {
		words := strings.Split(test.s, test.sep)
		if got, want := len(words), test.want; got != want {
			t.Errorf("Split(%q, %q) returned %d words, want %d", test.s, test.sep, got, want)
		}
	}
}

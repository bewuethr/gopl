package ch11ex01

import (
	"strings"
	"testing"
)

func TestCharCount(t *testing.T) {
	var tests = []struct {
		inStr   string
		lengths [5]int
		invalid int
		counts  map[rune]int
	}{
		{
			"abc", [...]int{0, 3, 0, 0, 0}, 0,
			map[rune]int{'a': 1, 'b': 1, 'c': 1},
		},
		{
			"", [...]int{0, 0, 0, 0, 0}, 0,
			map[rune]int{},
		},
		{
			"\xbd\xb2=\xbc ⌘", [...]int{0, 2, 0, 1, 0}, 3,
			map[rune]int{'=': 1, ' ': 1, '⌘': 1},
		},
		{
			"\xbd\xb2\xbc", [...]int{0, 0, 0, 0, 0}, 3,
			map[rune]int{},
		},
		{
			"xü語😀", [...]int{0, 1, 1, 1, 1}, 0,
			map[rune]int{'x': 1, 'ü': 1, '語': 1, '😀': 1},
		},
	}

	for _, test := range tests {
		input := strings.NewReader(test.inStr)
		want := Count{
			test.counts,
			test.lengths,
			test.invalid,
		}
		if got, _ := CharCount(input); !isEqual(got, want) {
			t.Errorf("CharCount(%q) = %+v\nexpected %+v", test.inStr, got, want)
		}
	}
}

func isEqual(c1, c2 Count) bool {
	if len(c1.Counts) != len(c2.Counts) {
		return false
	}
	for k, v1 := range c1.Counts {
		v2, ok := c2.Counts[k]
		if !ok || v1 != v2 {
			return false
		}
	}

	return c1.UTFLen == c2.UTFLen && c1.Invalid == c2.Invalid
}

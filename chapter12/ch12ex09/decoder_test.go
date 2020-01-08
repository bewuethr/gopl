package ch12ex09

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestDecoder(t *testing.T) {
	var tests = []string{
		"abc",
		`"abc"`,
		"123",
		"(1 2 3)",
		`("abc" "def" "ghi")`,
		"(a b c)",
		`(key "value")`,
		`((key1 "value1") (key2 "value2"))`,
	}

	for _, test := range tests {
		s := bytes.NewBufferString(test)
		d := NewDecoder(s)
		tok, err := d.Token()
		for err == nil {
			fmt.Printf("%T: %[1]v\n", tok)
			tok, err = d.Token()
		}
		if err != io.EOF {
			t.Errorf("error: %v", err)
		}
		fmt.Println()
	}
}

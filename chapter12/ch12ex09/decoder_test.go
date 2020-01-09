package ch12ex09

import (
	"bytes"
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
			t.Logf("%T: %[1]v\n", tok)
			tok, err = d.Token()
		}
		if err != io.EOF {
			t.Errorf("error: %v", err)
		}
	}
}

func TestDecoderErrors(t *testing.T) {
	var tests = []string{
		")",
		"(()",
		"1.23",
		"(1",
		`("abc"))`,
	}

	for _, test := range tests {
		s := bytes.NewBufferString(test)
		d := NewDecoder(s)
		tok, err := d.Token()
		for err == nil {
			t.Logf("%T: %[1]v\n", tok)
			tok, err = d.Token()
		}
		if err == io.EOF || err == nil {
			t.Error("expected error, but didn't get one")
		} else {
			t.Logf("expected error: %v", err)
		}
	}
}

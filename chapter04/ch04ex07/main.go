// ch04ex06 reverses a byte slice representing a UTF-8 string in-place.
package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func main() {
	// Tests for reverse
	testInputs := [][]byte{
		[]byte("abc"),
		[]byte("Hello, 世界"),
		[]byte("世ab界"),
		[]byte("a世b界"),
	}

	expected := [][]byte{
		[]byte("cba"),
		[]byte("界世 ,olleH"),
		[]byte("界ba世"),
		[]byte("界b世a"),
	}

	for i, v := range testInputs {
		var orig []byte
		copy(orig, v)
		reverse(v)
		if bytes.Equal(v, expected[i]) {
			fmt.Printf("%s: \u2713\n", v)
		} else {
			fmt.Printf("%s: expected <%s>, got <%s>\n", orig, expected[i], v)
		}
	}
}

// reverse reverses a byte slice of UTF-8 characters in place.
func reverse(s []byte) {
	i, j := 0, len(s)
	for {
		first, szFirst := utf8.DecodeRune(s[i:j])
		last, szLast := utf8.DecodeLastRune(s[i:j])

		if szFirst != szLast {
			// Shift middle part right by this much
			shift := szLast - szFirst
			copy(s[i+szFirst+shift:j-szLast+shift], s[i+szFirst:j-szLast])
		}
		utf8.EncodeRune(s[i:], last)
		utf8.EncodeRune(s[j-szFirst:], first)

		i += szLast
		j -= szFirst
		if i >= j-1 {
			return
		}
	}
}

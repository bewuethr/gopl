// Ch04ex06 squashes adjacent Unicode spaces in a UTF-8-encoded byte slice.
package main

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	testInputs := [][]byte{
		[]byte("abc"),
		[]byte("ab c"),
		[]byte("ab  c"),
		[]byte(" abc"),
		[]byte("  abc"),
		[]byte("abc "),
		[]byte("abc  "),
		[]byte("ab   c"),
		[]byte("ab\u00a0c"),
		[]byte("ab\u00a0\u1680\u2000\u3000c"),
	}

	expected := [][]byte{
		[]byte("abc"),
		[]byte("ab c"),
		[]byte("ab c"),
		[]byte(" abc"),
		[]byte(" abc"),
		[]byte("abc "),
		[]byte("abc "),
		[]byte("ab c"),
		[]byte("ab c"),
		[]byte("ab c"),
	}

	for i, v := range testInputs {
		v = squeeze(v)
		if bytes.Equal(v, expected[i]) {
			fmt.Printf("%s: \u2713\n", v)
		} else {
			fmt.Printf("%s: expected <%s>, got <%s>\n", v, expected[i], v)
		}
	}
}

// squeeze squashes adjacent Unicode spaces in the provided byte slice.
func squeeze(b []byte) []byte {
	s := b[:0]
	for idx := 0; idx < len(b); {
		r, size := utf8.DecodeRune(b[idx:])
		if unicode.IsSpace(r) {
			if len(s) == 0 || s[len(s)-1] != ' ' {
				s = append(s, ' ')
			}
		} else {
			s = append(s, string(r)...)
			// Was previously written as:
			// buf := make([]byte, utf8.RuneLen(r))
			// utf8.EncodeRune(buf, r)
			// s = append(s, buf...)
		}
		idx += size
	}
	return b[:len(s)]
}

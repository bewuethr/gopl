// Package ch11ex01 provides the Charcount function from Section 4.3 as a
// library function to be tested.
package ch11ex01

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

// Count contains character count data.
type Count struct {
	Counts  map[rune]int         // counts of Unicode characters
	UTFLen  [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	Invalid int                  // count of invalid UTF-8 characters
}

// CharCount reports the count of each rune, the byte lengths of the encodings
// and the number of invalid characters in input r.
func CharCount(r io.Reader) (Count, error) {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	in := bufio.NewReader(r)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			return Count{}, err
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}

	return Count{
		Counts:  counts,
		UTFLen:  utflen,
		Invalid: invalid,
	}, nil
}

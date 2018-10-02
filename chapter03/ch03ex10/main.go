// ch03ex10 is a non-recursive version of comma using bytes.Buffer instead of
// string concatenation.
package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	var buf bytes.Buffer
	revIdx := utf8.RuneCountInString(s)

	for i, r := range s {
		if revIdx%3 == 0 && i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteRune(r)
		revIdx--
	}

	return buf.String()
}

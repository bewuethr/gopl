// ch04ex08 adds Unicode category counts to the charcount program.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)      // counts of Unicode characters
	catCounts := make(map[string]int) // counts of characters by Unicode category
	var utflen [utf8.UTFMax + 1]int   // count of lengths of UTF-8 encodings
	invalid := 0                      // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++

		// Category counts: digits, lower, upper, letter, space
		if unicode.IsDigit(r) {
			catCounts["digit"]++
		}
		if unicode.IsLower(r) {
			catCounts["lower"]++
		}
		if unicode.IsUpper(r) {
			catCounts["upper"]++
		}
		if unicode.IsLetter(r) {
			catCounts["letter"]++
		}
		if unicode.IsSpace(r) {
			catCounts["space"]++
		}
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	fmt.Print("\ncat\tcount\n")
	for cat, n := range catCounts {
		fmt.Printf("%s\t%d\n", cat, n)
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

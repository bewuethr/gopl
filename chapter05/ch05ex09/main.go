// ch05ex09 expands each substring $foo to f(foo).
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	s1 := "This is $my $test $string."
	s2 := "This is $my$test$string."
	fmt.Println(expand(s1, strings.ToUpper))
	fmt.Println(expand(s2, strings.ToUpper))
	fmt.Println(expand(s1, toLen))
	fmt.Println(expand(s2, toLen))
	fmt.Println(expand(s1, strings.Title))
	fmt.Println(expand(s2, strings.Title))
}

func expand(s string, f func(string) string) string {
	var res, varName bytes.Buffer
	var inVar bool
	for _, r := range s {
		if r == '$' {
			if inVar {
				// We have a $foo$bar situation and we're at the "$" of "$bar"
				res.WriteString(f(varName.String()))
				varName.Reset()
			} else {
				inVar = true
			}
		} else if unicode.IsLetter(r) {
			if inVar {
				varName.WriteRune(r)
			} else {
				res.WriteRune(r)
			}
		} else {
			if inVar {
				// We've reached the end of a variable name
				inVar = false
				res.WriteString(f(varName.String()))
				res.WriteRune(r)
				varName.Reset()
			} else {
				res.WriteRune(r)
			}
		}
	}

	return res.String()
}

// toLen returns the rune count of the supplied string s as a string.
func toLen(s string) string {
	return strconv.Itoa(utf8.RuneCountInString(s))
}

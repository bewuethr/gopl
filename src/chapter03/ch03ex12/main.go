// ch03ex12 tests whether two strings are anagrams of each other.
package main

import (
	"fmt"
	"os"
)

func main() {
	if areAnagrams(os.Args[1], os.Args[2]) {
		fmt.Println("Anagrams")
	} else {
		fmt.Println("Not anagrams")
	}
}

func areAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	charCount1 := make(map[rune]int)
	charCount2 := make(map[rune]int)

	for _, r := range s1 {
		charCount1[r]++
	}
	for _, r := range s2 {
		charCount2[r]++
	}

	if len(charCount1) != len(charCount2) {
		return false
	}

	for r1, c1 := range charCount1 {
		c2, ok := charCount2[r1]
		if !ok || c1 != c2 {
			return false
		}
	}

	return true
}

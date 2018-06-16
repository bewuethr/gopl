// ch04ex05 implements an in-place function that eliminates adjacent
// duplicates in a []string slice.
package main

import "fmt"

func main() {
	s0 := []string{"abc", "abc", "def", "ghi"}
	s0 = uniq(s0)
	fmt.Printf("%v\n", s0) // [abc def ghi]

	s1 := []string{"abc", "abc", "abc", "def", "ghi"}
	s1 = uniq(s1)
	fmt.Printf("%v\n", s1) // [abc def ghi]

	s2 := []string{"abc", "abc", "abc", "def", "ghi", "ghi"}
	s2 = uniq(s2)
	fmt.Printf("%v\n", s2) // [abc def ghi]

	s3 := []string{}
	s3 = uniq(s3)
	fmt.Printf("%v\n", s3) // []

	s4 := []string{"abc"}
	s4 = uniq(s4)
	fmt.Printf("%v\n", s4) // [abc]

	s5 := []string{"abc", "abc"}
	s5 = uniq(s5)
	fmt.Printf("%v\n", s5) // [abc]

	s6 := []string{"abc", "def"}
	s6 = uniq(s6)
	fmt.Printf("%v\n", s6) // [abc def]
}

// uniq removes adjacent duplicates from s by appending to a slice that
// references the same array.
func uniq(s []string) []string {
	u := s[:0]
	for i, v := range s {
		if i == 0 || v != s[i-1] {
			u = append(u, v)
		}
	}
	return s[:len(u)]
}

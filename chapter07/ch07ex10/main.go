// Ch07ex10 implements a palindrome checker taking advantage of sort.Interface.
package main

import (
	"fmt"
	"sort"
)

func isPalindrome(s sort.Interface) bool {
	j := s.Len() - 1
	for i := 0; i < s.Len()/2; i++ {
		if !(!s.Less(i, j) && !s.Less(j, i)) {
			return false
		}
		j--
	}
	return true
}

func main() {
	intsYes := []int{1, 2, 3, 2, 1}
	intsNo := []int{1, 2, 3, 4, 5}
	fmt.Printf("%v: %v\n", intsYes, isPalindrome(sort.IntSlice(intsYes)))
	fmt.Printf("%v: %v\n", intsNo, isPalindrome(sort.IntSlice(intsYes)))

	stringsYes := []string{"one", "two", "three", "two", "one"}
	stringsNo := []string{"one", "two", "three", "four", "five"}
	fmt.Printf("%v: %v\n", stringsYes, isPalindrome(sort.StringSlice(stringsYes)))
	fmt.Printf("%v: %v\n", stringsNo, isPalindrome(sort.StringSlice(stringsNo)))
}

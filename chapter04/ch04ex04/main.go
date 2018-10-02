// ch04ex04 implements a single pass rotate function.
package main

import (
	"fmt"
)

func main() {
	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	rotate(s, 2)
	fmt.Println(s) // "[2 3 4 5 0 1]"

	// Rotate s right by three positions.
	rotate(s, -3)
	fmt.Println(s) // "[4 5 0 1 2 3]"
}

// rotate rotates a slice of ints, where the rotation is left for positive
// amounts and right for negative amounts.
func rotate(s []int, dist int) {
	var tmp []int
	if dist > 0 {
		tmp = append(s[dist:], s[:dist]...)
	} else {
		tmp = append(s[-dist:], s[:-dist]...)
	}
	copy(s, tmp)
}

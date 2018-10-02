// Package ch02ex03 with modified popcount using loop instead of single
// expression.
package ch02ex03

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	sum := 0
	for i := uint(0); i < 8; i++ {
		sum += int(pc[byte(x>>i*8)])
	}
	return sum
}

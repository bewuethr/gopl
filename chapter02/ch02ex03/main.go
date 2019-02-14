// Package ch02ex03 modifies popcount using a loop instead of a single
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

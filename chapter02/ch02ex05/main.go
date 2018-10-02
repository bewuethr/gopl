// Package ch02ex05 with modified popcount using x&(x-1)
package ch02ex05

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	count := 0
	for x > 0 {
		count++
		x &= x - 1
	}
	return count
}

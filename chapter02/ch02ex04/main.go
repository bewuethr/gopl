// Package ch02ex04 with modified popcount using shifting
package ch02ex04

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	count := 0
	for i := uint64(0); i < 64; i++ {
		count += int(x >> i & uint64(1))
	}
	return count
}

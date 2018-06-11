// ch04ex01 counts the number of bits that are different in two SHA256 hashes.
package main

import (
	"crypto/sha256"
	"fmt"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func main() {
	str1, str2 := "x", "X"
	c1 := sha256.Sum256([]byte(str1))
	c2 := sha256.Sum256([]byte(str2))
	fmt.Printf("%s vs. %s: %d\n", str1, str2, bitDiff(c1, c2))
}

func bitDiff(hashA, hashB [32]byte) int {
	res := 0
	for i := range hashA {
		res += int(pc[hashA[i]^hashB[i]])
	}
	return res
}

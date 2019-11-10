// Package ch09ex02 provides the PopCount function from Chapter 2 using lazy
// initialization.
package ch09ex02

import (
	"fmt"
	"sync"
)

// pc[i] is the population count of i.
var pc [256]byte

var initCountOnce sync.Once

func initCount() {
	fmt.Println("initializing pc...")
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	initCountOnce.Do(initCount)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

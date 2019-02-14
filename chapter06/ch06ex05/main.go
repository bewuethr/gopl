// Ch06ex06 uses a platform-specific word sizes for the bitset implementation.
package main

import (
	"fmt"

	"github.com/bewuethr/gopl/chapter06/ch06ex05/intset"
)

func main() {
	var x, y intset.IntSet
	x.Add(1)
	x.Add(155)
	fmt.Printf("x (expect {1 155}): %v\n", x.String())

	y.Add(2)
	y.Add(3)
	fmt.Printf("y (expect {2 3}): %v\n", y.String())

	fmt.Printf("x has 2 (expect false): %v\n", x.Has(2))
	x.UnionWith(&y)
	fmt.Printf("x union y (expect {1 2 3 155}): %v\n", x.String())
	fmt.Printf("x has 2 (expect true): %v\n", x.Has(2))

	// Exercise 6.1
	fmt.Println()
	fmt.Printf("x.Len() (expect 4): %d\n", x.Len())
	fmt.Printf("y.Len() (expect 2): %d\n", y.Len())

	// Exercise 6.1
	x.Remove(3)
	fmt.Printf("x after 3 was removed (expect {1 2 155}): %v\n", x.String())

	// Exercise 6.1
	y.Clear()
	fmt.Printf("y after it was cleared (expect {}): %v\n", y.String())

	// Exercise 6.1
	var z intset.IntSet
	z = *x.Copy()
	z.Add(100)
	z.Remove(155)
	fmt.Printf("z (copy of x) after adding 100 and removing 155 (expect {1 2 100}): %v\n", z.String())
	fmt.Printf("and x is still (expect {1 2 155}): %v\n", x.String())

	// Exercise 6.2
	fmt.Println()
	z.AddAll(200, 300, 400)
	fmt.Printf("z after adding 200, 300, 400 (expect {1 2 100 200 300 400}): %v\n", z.String())

	// Exercise 6.3
	fmt.Println()
	x.IntersectWith(&z)
	fmt.Printf("x after intersecting with z (expect {1 2}): %v\n", x.String())

	z.DifferenceWith(&x)
	fmt.Printf("z after subtracting x (expect {100 200 300 400}): %v\n", z.String())

	x.Clear()
	x.AddAll(1, 2, 3)
	y.AddAll(2, 3, 4)
	fmt.Printf("x: %v\n", x.String())
	fmt.Printf("y: %v\n", y.String())
	x.SymmetricDifference(&y)
	fmt.Printf("x after symmetric difference with y (expect {1 4}): %v\n", x.String())

	// Exercise 6.4
	fmt.Println()
	e := y.Elems()
	fmt.Printf("Elements of y (expect [2 3 4]): %v\n", e)
}

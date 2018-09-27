package main

import (
	"fmt"

	"chapter06/intset"
)

func main() {
	var x, y intset.IntSet
	x.Add(1)
	x.Add(155)
	fmt.Printf("x: %v\n", x.String())

	y.Add(2)
	y.Add(3)
	fmt.Printf("y: %v\n", y.String())

	fmt.Printf("x has 2: %v\n", x.Has(2))
	x.UnionWith(&y)
	fmt.Printf("x union y: %v\n", x.String())
	fmt.Printf("x has 2: %v\n", x.Has(2))

	// Exercise 5.1
	fmt.Printf("x.Len() = %d\n", x.Len())
	fmt.Printf("y.Len() = %d\n", y.Len())

	// Exercise 5.1
	x.Remove(3)
	fmt.Printf("x after 3 was removed: %v\n", x.String())

	// Exercise 5.1
	y.Clear()
	fmt.Printf("y after it was cleared: %v\n", y.String())

	// Exercise 5.1
	z := x.Copy()
	z.Add(100)
	z.Remove(155)
	fmt.Printf("z (copy of x) after adding 100 and removing 155: %v\n", z.String())
	fmt.Printf("and x is still: %v\n", x.String())
}

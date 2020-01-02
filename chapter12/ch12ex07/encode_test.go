package ch12ex07

import (
	"bytes"
	"fmt"
)

func ExampleEncoder_Encode() {
	anInt := 8
	var anInterface interface{} = []int{1, 2, 3}

	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Encode(8)                        // int
	enc.Encode(-8)                       // negative int
	enc.Encode(uint(8))                  // uint
	enc.Encode("abc")                    // string
	enc.Encode(&anInt)                   // pointer to int
	enc.Encode([]int{0, 1, 2})           // slice
	enc.Encode([...]int{0, 1, 2})        // array
	enc.Encode([]int{})                  // empty slice
	enc.Encode(struct{ m, n int }{5, 6}) // struct
	enc.Encode(true)                     // true
	enc.Encode(false)                    // false
	enc.Encode(3.14159)                  // floating point
	enc.Encode(1.1 + 2.1i)               // complex number
	enc.Encode(&anInterface)             // interface

	fmt.Print(b.String())
	// Output:
	// 8
	// -8
	// 8
	// "abc"
	// 8
	// (0 1 2)
	// (0 1 2)
	// ()
	// ((m 5) (n 6))
	// t
	// nil
	// 3.14159
	// #C(1.1 2.1)
	// ("[]int" (1 2 3))
}

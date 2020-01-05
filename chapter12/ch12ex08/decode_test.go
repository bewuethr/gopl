package ch12ex08

import (
	"fmt"
	"log"
)

/*
	8
	(0 1 2)
	(0 1 2)
	()
	((m 5) (n 6))
*/

func ExampleDecoder_Decode() {
	b := []byte("8")
	var anInt int
	if err := Unmarshal(b, &anInt); err != nil {
		log.Fatal("int: " + err.Error())
	}
	fmt.Printf("%+v\n", anInt)

	b = []byte(`"abc"`)
	var aStr string
	if err := Unmarshal(b, &aStr); err != nil {
		log.Fatal("string: " + err.Error())
	}
	fmt.Printf("%+v\n", aStr)

	b = []byte("(0 1 2)")
	var intSlice []int
	if err := Unmarshal(b, &intSlice); err != nil {
		log.Fatal("int slice: " + err.Error())
	}
	fmt.Printf("%+v\n", intSlice)

	b = []byte("(3 4 5)")
	var intArray [3]int
	if err := Unmarshal(b, &intArray); err != nil {
		log.Fatal("int array: " + err.Error())
	}
	fmt.Printf("%+v\n", intArray)

	b = []byte("((M 5) (N 6))")
	var aStruct struct{ M, N int }
	if err := Unmarshal(b, &aStruct); err != nil {
		log.Fatal("struct: " + err.Error())
	}
	fmt.Printf("%+v\n", aStruct)

	b = []byte(`(("key1" 1) ("key2" 2))`)
	var aMap map[string]int
	if err := Unmarshal(b, &aMap); err != nil {
		log.Fatal("map: " + err.Error())
	}
	fmt.Printf("%+v\n", aMap)

	// Output:
	// 8
	// abc
	// [0 1 2]
	// [3 4 5]
	// {M:5 N:6}
	// map[key1:1 key2:2]
}

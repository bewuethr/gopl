package ch12ex01

func ExampleDisplay_structs() {
	type keyStruct struct {
		name string
		n    int
	}

	m := map[keyStruct]string{
		{name: "a", n: 1}: "first",
		{name: "b", n: 2}: "second",
	}

	Display("m", m)
	// Unordered output:
	// Display m (map[ch12ex01.keyStruct]string):
	// m[{name:"a" n:1}] = "first"
	// m[{name:"b" n:2}] = "second"
}

func ExampleDisplay_nestedStructs() {
	type innerStruct struct {
		n int
	}

	type keyStruct struct {
		name innerStruct
	}

	m := map[keyStruct]string{
		{name: innerStruct{n: 1}}: "first",
		{name: innerStruct{n: 2}}: "second",
	}

	Display("m", m)
	// Unordered output:
	// Display m (map[ch12ex01.keyStruct]string):
	// m[{name:{n:1}}] = "first"
	// m[{name:{n:2}}] = "second"
}

func ExampleDisplay_arrays() {
	m := map[[2]int]string{
		[...]int{0, 1}: "first",
		[...]int{2, 3}: "second",
	}

	Display("m", m)
	// Unordered output:
	// Display m (map[[2]int]string):
	// m[[0 1]] = "first"
	// m[[2 3]] = "second"
}

func ExampleDisplay_structArrays() {
	type aStruct struct {
		name string
		n    int
	}

	m := map[[2]aStruct]string{
		[...]aStruct{{name: "abc", n: 1}, {name: "def", n: 2}}: "first",
		[...]aStruct{{name: "ghi", n: 3}, {name: "jkl", n: 4}}: "second",
	}

	Display("m", m)
	// Unordered output:
	// Display m (map[[2]ch12ex01.aStruct]string):
	// m[[{name:"abc" n:1} {name:"def" n:2}]] = "first"
	// m[[{name:"ghi" n:3} {name:"jkl" n:4}]] = "second"

}

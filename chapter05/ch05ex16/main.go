// Ch05ex16 implements strings.Join as a variadic function.
package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(join("X", ""))
	fmt.Println(join("X", "abc"))
	fmt.Println(join("X", "abc", "def"))
}

// join concatenates the elements of a to create a single string. The separator
// string sep is placed between elements in the resulting string.
func join(sep string, a ...string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	}

	res := bytes.NewBufferString(a[0])
	for _, s := range a[1:] {
		res.WriteString(sep + s)
	}

	return res.String()
}

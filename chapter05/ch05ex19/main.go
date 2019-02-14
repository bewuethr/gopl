// Ch05ex19 implements a function without return statement that returns a
// non-zero value.
package main

import "fmt"

func main() {
	fmt.Println(magicFunc())
}

func magicFunc() (val int) {
	defer func() {
		recover()
		val = 1
	}()
	panic("panic!")
}

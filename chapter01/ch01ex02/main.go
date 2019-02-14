// Ch01ex02 modifies Echo to print the index and value of each argument, one
// per line.
package main

import (
	"fmt"
	"os"
)

func main() {
	for i, arg := range os.Args[1:] {
		fmt.Printf("%d:\t%s\n", i+1, arg)
	}
}

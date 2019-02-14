// Ch01ex01 modifies Echo to print the name of the command that invoked it and
// all its arguments.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args, " "))
}

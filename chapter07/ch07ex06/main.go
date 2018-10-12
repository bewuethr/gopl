// ch07ex06 prints the value of its -temp (temperature) flag, which can be in
// Celsius, Fahrenheit or Kelvin.
package main

import (
	"flag"
	"fmt"

	"gopl/chapter07/ch07ex06/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

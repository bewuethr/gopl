// Ch07ex06 prints the value of its -temp (temperature) flag, which can be in
// Celsius, Fahrenheit or Kelvin.
package main

import (
	"flag"
	"fmt"

	"github.com/bewuethr/gopl/chapter07/ch07ex06/tempconv"
)

// Exercise 7.7: The help message will say "20°C" because it uses
// tempconv.Celsius.String(). It is triggered by flag.PrintDefaults().
var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

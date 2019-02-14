// Ch02ex02 is a general-purpose unit-conversion program.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, arg := range os.Args[1:] {
			printConvs(arg)
		}
	} else {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			printConvs(input.Text())
		}
		// NOTE: ignoring errors from input.Err()
	}
}

func printConvs(s string) {
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cf: %v\n", err)
		os.Exit(1)
	}
	f := Fahrenheit(t)
	c := Celsius(t)
	ft := Feet(t)
	m := Meters(t)
	lbs := Pounds(t)
	kg := Kilograms(t)

	fmt.Printf("%s = %s, %s = %s, %s = %s, %s = %s, %s = %s, %s = %s\n",
		f, FToC(f), c, CToF(c), ft, FtToM(ft), m, MToFt(m), lbs, LbsToKg(lbs), kg, KgToLbs(kg))
}

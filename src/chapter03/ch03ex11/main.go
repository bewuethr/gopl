// ch03ex11 extends comma to work with floating-point numbers and an optional
// sign.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

// comma inserts commas in numbers, including floating point with optional
// sign.
func comma(s string) string {
	// Cut off sign if present
	var sign byte
	if s[0] == '+' || s[0] == '-' {
		sign = s[0]
		s = s[1:]
	}

	// Check if there is a decimal point
	i := strings.Index(s, ".")
	if i == -1 {
		return string(sign) + commafy(s)
	}

	intPart := s[:i]
	fracPart := s[i+1:]

	return string(sign) + commafy(intPart) + "." + commafyFrac(fracPart)
}

// commafy inserts commas for the integer part
func commafy(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return commafy(s[:n-3]) + "," + s[n-3:]
}

// commafyFrac inserts commas for the fractional part
func commafyFrac(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return s[:3] + "," + commafyFrac(s[3:])
}

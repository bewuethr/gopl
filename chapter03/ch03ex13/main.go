// ch03ex13 declares constants for KB, MB etc. as compactly as possible.
package main

import "fmt"

const (
	KB = 1e3
	MB = 1e6
	GB = 1e9
	TB = 1e12
	PB = 1e15
	EB = 1e18
	ZB = 1e21
	YB = 1e24
)

func main() {
	fmt.Println("KB:", KB)
	fmt.Println("MB:", MB)
	fmt.Println("GB:", GB)
	fmt.Println("TB:", TB)
	fmt.Println("PB:", PB)
	fmt.Println("EB:", EB)
	fmt.Println("ZB:", ZB)
	fmt.Println("YB:", YB)
}

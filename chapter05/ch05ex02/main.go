// Ch05ex02 builds a mapping from element names to the number of occurrences of
// that element in an HTML tree.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch05ex02: %v\n", err)
		os.Exit(1)
	}
	for name, count := range visit(make(map[string]int), doc) {
		fmt.Printf("%s:\t%d\n", name, count)
	}
}

// visit updates the mapping for the element count and returns the result.
func visit(elCount map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return elCount
	}

	if n.Type == html.ElementNode {
		elCount[n.Data]++
	}

	elCount = visit(elCount, n.FirstChild)
	elCount = visit(elCount, n.NextSibling)
	return elCount
}

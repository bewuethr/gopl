// Ch05ex03 traverses and HTML tree and prints the content of all text nodes
// (except contents of script and style elements).
package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch05ex03: %v\n", err)
		os.Exit(1)
	}
	for _, content := range visit(nil, doc) {
		fmt.Println(content)
	}
}

// visit appends to contents the content of each text node found in n and
// returns the result.
func visit(contents []string, n *html.Node) []string {
	if n == nil || n.Data == "script" || n.Data == "style" {
		return contents
	}

	if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
		contents = append(contents, n.Data)
	}

	contents = visit(contents, n.FirstChild)
	contents = visit(contents, n.NextSibling)
	return contents
}

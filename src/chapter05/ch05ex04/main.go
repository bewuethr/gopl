// ch05ex04 extracts image, script and style sheet links.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch05ex04: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode {
		if n.Data == "img" || n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					links = append(links, a.Val)
				}
			}
		}
		if n.Data == "link" {
			var rel, target string
			for _, a := range n.Attr {
				if a.Key == "href" {
					target = a.Val
				} else if a.Key == "rel" {
					rel = a.Val
				}
			}
			if rel == "stylesheet" {
				links = append(links, target)
			}
		}
	}

	links = visit(links, n.FirstChild)
	links = visit(links, n.NextSibling)
	return links
}

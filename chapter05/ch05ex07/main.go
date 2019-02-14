// Ch05ex07 is an HTML pretty printer.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	forEachNode(doc, startElement, endElement)

	return nil
}

// forEachNode calls the functions pre(x) and post(x) for each node x in the
// tree rooted at n. Both functions are optional.  pre is called before the
// children are visited (preorder) and post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

func startElement(n *html.Node) {
	switch n.Type {
	case html.DoctypeNode:
		fmt.Printf("<!DOCTYPE %s>\n", n.Data)
	case html.ElementNode:
		fmt.Printf("%*s<%s%v%s\n", depth*2, "", n.Data, printAttr(n.Attr), closeTag(n))
		depth++
	case html.TextNode:
		if strings.TrimSpace(n.Data) != "" {
			fmt.Printf("%*s%s\n", depth*2, "", n.Data)
		}
	case html.CommentNode:
		fmt.Printf("%*s<!--%s-->\n", depth*2, "", n.Data)
		depth++
	}
}

func endElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		depth--
		if n.FirstChild != nil {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	case html.CommentNode:
		depth--
	}
}

func printAttr(attrs []html.Attribute) string {
	if len(attrs) == 0 {
		return ""
	}

	var res string

	for _, a := range attrs {
		res += fmt.Sprintf(" %s=%q", a.Key, a.Val)
	}

	return res
}

func closeTag(n *html.Node) string {
	if n.FirstChild == nil {
		return " />"
	}
	return ">"
}

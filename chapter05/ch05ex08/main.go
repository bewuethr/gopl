// Ch05ex08 finds the first element with the given ID in the supplied HTML
// documents.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var id string

func init() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ch05ex08 ID URL [URL...]")
		os.Exit(1)
	}
	id = os.Args[1]
}

func main() {
	for _, url := range os.Args[2:] {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if element := elementByID(doc, id); element != nil {
			fmt.Printf("%+v\n", element)
		}
	}
}

func elementByID(doc *html.Node, id string) *html.Node {
	return forEachNode(doc, checkID, nil)
}

// forEachNode calls the functions pre(x) and post(x) for each node x in the
// tree rooted at n. Both functions are optional. pre is called before the
// children are visited (preorder) and post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) *html.Node {
	if pre != nil {
		if !pre(n) {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rv := forEachNode(c, pre, post)
		if rv != nil {
			return rv
		}
	}

	if post != nil {
		if !post(n) {
			return n
		}
	}

	return nil
}

func checkID(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return true
	}

	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return false
		}
	}

	return true
}

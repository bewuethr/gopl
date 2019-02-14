// Ch05ex17 implements a variadic function that extracts HTML elements by name
// from an HTML document.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "usage: ch05ex16 URL NAME...")
		os.Exit(1)
	}
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "GET response is not HTTP 200")
		os.Exit(1)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	nodes := elementsByTagName(doc, os.Args[2:]...)
	for _, n := range nodes {
		fmt.Printf("%+v\n", n.Attr)
	}
}

// elementsByTagName traverses doc and returns a slice of all nodes with a name
// matching any of the name strings.
func elementsByTagName(doc *html.Node, name ...string) []*html.Node {
	var res []*html.Node
	m := make(map[string]bool)
	for _, n := range name {
		m[n] = true
	}
	if doc.Type == html.ElementNode && m[doc.Data] {
		res = append(res, doc)
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		res = append(res, elementsByTagName(c, name...)...)
	}

	return res
}

// Ch07ex04 implements NewReader and a simple string reader to use with the
// HTML parser from exercise 5.2.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/net/html"
)

// MyReader implements an io.Reader.
type MyReader struct {
	p   []byte
	idx int
}

func (r *MyReader) Read(p []byte) (int, error) {
	read := 0
	for i := range p {
		if r.idx == len(r.p) {
			return read, io.EOF
		}
		p[i] = r.p[r.idx]
		read++
		r.idx++
	}
	return read, nil
}

// NewReader reuturns a pointer to an instance of MyReader.
func NewReader(s string) io.Reader {
	r := MyReader{
		p:   []byte(s),
		idx: 0,
	}
	return &r
}

func main() {
	inBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch07ex04: %v\n", err)
		os.Exit(1)
	}
	reader := NewReader(string(inBytes))
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch07ex04: %v\n", err)
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

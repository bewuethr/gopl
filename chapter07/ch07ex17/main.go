// Ch07ex17 extends xmlselect to select elements not just by name, but by their
// attributes, too. Usage:
//
//     ./ch07ex17 h5 id=IDAKYDS < infile.xml
//
// selects any element where the stack contains h5 and the most recent element
// has an attribute id="IDAKYDS".
package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	names, attrs, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	dec := xml.NewDecoder(os.Stdin)
	var (
		nameStack []string   // stack of element names
		topAttrs  []xml.Attr // attributes of most recent start element
	)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			nameStack = append(nameStack, tok.Name.Local) // push
			topAttrs = tok.Attr
		case xml.EndElement:
			nameStack = nameStack[:len(nameStack)-1] // pop
		case xml.CharData:
			if containsAll(nameStack, names) && matchesAttributes(topAttrs, attrs) {
				fmt.Printf("%s: %s\n", strings.Join(nameStack, " "), tok)
			}
		}
	}
}

// parseFlags reads the command line arguments into a slice for names and one
// for XML attributes. Parameters with an "=" sign are considered attributes.
func parseFlags() ([]string, []xml.Attr, error) {
	if len(os.Args) < 2 {
		return nil, nil, errors.New("not enough command line arguments")
	}
	var (
		names []string
		attrs []xml.Attr
	)

	for _, arg := range os.Args[1:] {
		if strings.ContainsRune(arg, '=') {
			strs := strings.Split(arg, "=")
			attrs = append(attrs, xml.Attr{Name: xml.Name{Local: strs[0]}, Value: strs[1]})
			continue
		}
		names = append(names, arg)
	}

	return names, attrs, nil
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

// matchesAttributes reports wether x provides all attributes in y.
func matchesAttributes(x, y []xml.Attr) bool {
	for _, a := range y {
		if !contains(x, a) {
			return false
		}
	}
	return true
}

// contains reports if x contains y by comparing name and value.
func contains(x []xml.Attr, y xml.Attr) bool {
	for _, a := range x {
		if a.Name.Local == y.Name.Local && a.Value == y.Value {
			return true
		}
	}
	return false
}

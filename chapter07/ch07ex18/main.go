// Ch07ex18 reads an arbirtrary XML document and constructs a tree of generic
// nodes that represent it.
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// A Node represents a text string (CharData) or named element node (Element).
type Node interface{}

// CharData is a text string Node.
type CharData string

// An Element represents a named element and its attributes; it holds a slice
// of its child nodes.
type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func main() {
	tree := createXMLTree(os.Stdin)
	prettyPrint(tree, 0)
}

func createXMLTree(r io.Reader) *Element {
	dec := xml.NewDecoder(r)
	root := Element{}
	elemStack := []*Element{&root}
	curElem := elemStack[0]
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "ch07ex18: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			node := Element{
				Type: tok.Name,
				Attr: tok.Attr,
			}
			curElem.Children = append(curElem.Children, &node)
			elemStack = append(elemStack, &node)
			curElem = &node
		case xml.EndElement:
			elemStack = elemStack[:len(elemStack)-1]
			curElem = elemStack[len(elemStack)-1]
		case xml.CharData:
			curElem.Children = append(curElem.Children, CharData(tok))
		}
	}
	return &root
}

func prettyPrint(n Node, indent int) {
	switch n := n.(type) {
	case *Element:
		fmt.Printf("%*s<%s", indent, "", n.Type.Local)
		for _, a := range n.Attr {
			fmt.Printf(" %s=\"%s\"", a.Name.Local, a.Value)
		}
		fmt.Println(">")
		for _, c := range n.Children {
			prettyPrint(c, indent+2)
		}
		fmt.Printf("%*s</%s>\n", indent, "", n.Type.Local)
	case CharData:
		fmt.Printf("%*s%s\n", indent, "", n)
	default:
		fmt.Fprintf(os.Stderr, "unexpected type: %T\n", n)
		os.Exit(1)
	}
}

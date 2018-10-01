// ch07ex01 implements writers that count words and lines.
package main

import (
	"fmt"
	"strings"
)

type WordCounter int
type LineCounter int

func (w *WordCounter) Write(p []byte) (int, error) {
	*w += WordCounter(len(strings.Fields(string(p))))
	return len(p), nil
}

func (l *LineCounter) Write(p []byte) (int, error) {
	*l += LineCounter(len(strings.Split(string(p), "\n")))
	return len(p), nil
}

func main() {
	var w WordCounter
	w.Write([]byte("hello"))
	fmt.Println(w) // "1"

	w = 0 // reset the counter
	var name = "Dolly"
	fmt.Fprintf(&w, "hello, %s, four words", name)
	fmt.Println(w) // "4"

	var l LineCounter
	l.Write([]byte("hello"))
	fmt.Println(l) // "1"

	l = 0 // reset the counter
	fmt.Fprintf(&l, "hello, %s,\nsecond line\nthird line\nfourth line\nfifth line", name)
	fmt.Println(l) // "5"
}

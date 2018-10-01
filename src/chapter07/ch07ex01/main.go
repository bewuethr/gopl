// ch07ex01 implements writers that count words and lines.
package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int
type LineCounter int

func (w *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)

	written := 0

	for scanner.Scan() {
		*w++
		written += len(scanner.Bytes())
	}

	return written, scanner.Err()
}

func (l *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))

	written := 0

	for scanner.Scan() {
		*l++
		written += len(scanner.Bytes())
	}

	return written, scanner.Err()
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

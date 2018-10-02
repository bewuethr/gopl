// ch07ex02 implements a function that returns a writer wrapping another writer
// and an int64 representing the number of bytes written to the writer at any
// time.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type CtWriter struct {
	writer    io.Writer
	byteCount *int64
}

func (w CtWriter) Write(p []byte) (int, error) {
	count, err := w.writer.Write(p)
	*w.byteCount += int64(count)
	return count, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := CtWriter{
		writer:    w,
		byteCount: new(int64),
	}
	return cw, cw.byteCount
}

func main() {
	w, c := CountingWriter(os.Stdout)
	fmt.Fprintf(w, "12345\n")
	fmt.Printf("Counter after adding '12345\\n': %d\n", *c)
	fmt.Fprintf(w, "12345\n")
	fmt.Printf("Counter after adding '12345\\n one more time': %d\n", *c)

	var b bytes.Buffer
	w, c = CountingWriter(&b)
	w.Write([]byte("abcdefghij"))
	fmt.Printf("Counter after writing 'abcdefghij' to wrapped bytes.Buffer: %d\n", *c)
}

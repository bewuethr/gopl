package ch13ex03

import (
	"bytes" // reader
	"io"
	"testing"
)

func TestBzip2(t *testing.T) {
	var compressed1, compressed2 bytes.Buffer
	w1 := NewWriter(&compressed1)
	w2 := NewWriter(&compressed2)

	// Write a repetitive message in a million pieces, to both writers
	// concurrently.
	tee := io.MultiWriter(w1, w2)
	for i := 0; i < 1000000; i++ {
		io.WriteString(tee, "hello")
	}
	if err := w1.Close(); err != nil {
		t.Fatal(err)
	}
	if err := w2.Close(); err != nil {
		t.Fatal(err)
	}

	// Check the size of the compressed streams.
	if got, want := compressed1.Len(), 255; got != want {
		t.Errorf("1 million hellos compressed to %d bytes, want %d", got, want)
	}
	if got, want := compressed2.Len(), 255; got != want {
		t.Errorf("1 million hellos compressed to %d bytes, want %d", got, want)
	}
}

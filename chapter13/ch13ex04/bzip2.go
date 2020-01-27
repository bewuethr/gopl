// Package ch13ex04 provides a writer that uses bzip2 compression (bzip.org) by
// wrapping the /bin/bzip2 executable.
package ch13ex04

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type writer struct {
	bzipProc *exec.Cmd
	stdin    io.WriteCloser
}

// NewWriter returns a writer for bzip2-compressed streams.
func NewWriter(out io.Writer) io.WriteCloser {
	cmd := exec.Command("bzip2")
	cmd.Stdout = out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get stdin pipe: %v\n", err)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to start bzip2 process: %v\n", err)
		os.Exit(1)
	}
	w := &writer{bzipProc: cmd, stdin: stdin}
	return w
}

// Write writes data to the underlying WriteCloser.
func (w *writer) Write(data []byte) (int, error) {
	return w.stdin.Write(data)
}

// Close flushes the compressed data and closes the stream. It does not close
// the underlying io.Writer.
func (w *writer) Close() error {
	if err := w.stdin.Close(); err != nil {
		return err
	}
	if err := w.bzipProc.Wait(); err != nil {
		return err
	}
	return nil
}

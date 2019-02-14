// Ch04ex09 counts the frequency of each word in an input text file.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "wordfreq: no file argument provided")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "wordfreq: %v\n", err)
		os.Exit(1)
	}
	input := bufio.NewScanner(f)
	input.Split(bufio.ScanWords)

	freq := make(map[string]int)

	for input.Scan() {
		freq[input.Text()]++
	}

	for w, c := range freq {
		fmt.Printf("%s\t%d\n", w, c)
	}
}

// Modify dup2 to print names of files in which duplicate occur
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	fnames := make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, fnames, "")
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, fnames, arg)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
			if len(fnames[line]) > 0 {
				fmt.Println(strings.Join(fnames[line], " "))
			}
		}
	}
}

func countLines(f *os.File, counts map[string]int, fnames map[string][]string, fname string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		text := input.Text()
		counts[text]++
		if fname != "" {
			fnames[text] = append(fnames[text], fname)
		}
	}
}

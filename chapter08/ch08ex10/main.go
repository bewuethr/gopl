// Ch08ex10 is a concurrent web crawler that supports cancellation.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	tokens = make(chan struct{}, 20) // limit concurrency
	done   = make(chan struct{})     // signal cancellation
)

func main() {
	flag.Parse()

	worklist := make(chan []string)

	var n int // number of pending sends to worklist
	n++
	go func() { worklist <- flag.Args() }()
	go cancel()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

func crawl(link string) []string {
	if cancelled() {
		return nil
	}
	fmt.Println(link)

	tokens <- struct{}{} // acquire a token
	list, err := extract(link)
	<-tokens //release the token

	if err != nil {
		log.Print(err)
	}

	return list
}

func cancel() {
	os.Stdin.Read(make([]byte, 1))
	close(done)
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

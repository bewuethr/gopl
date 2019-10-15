// Ch08ex06 is a concurrent crawler with depth-limiting.
package main

import (
	"flag"
	"fmt"
	"log"
)

type depthLink struct {
	url   string
	depth int
}

var tokens = make(chan struct{}, 20)

func main() {
	var maxDepth int
	flag.IntVar(&maxDepth, "depth", 0, "depth limit for crawler")
	flag.Parse()

	worklist := make(chan []depthLink)
	var n int // number of pending sends to worklist

	// Add command-line arguments to worklist.
	var argLinks []depthLink
	for _, arg := range flag.Args() {
		argLinks = append(argLinks, depthLink{
			url:   arg,
			depth: 0,
		})
	}
	n++
	go func() { worklist <- argLinks }()

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.url] {
				seen[link.url] = true
				if maxDepth == 0 || maxDepth > 0 && link.depth <= maxDepth {
					n++
					go func(link depthLink) {
						worklist <- crawl(link)
					}(link)
				}
			}
		}
	}
}

func crawl(link depthLink) []depthLink {
	fmt.Printf("%d - %s\n", link.depth, link.url)

	tokens <- struct{}{} // acquire a token
	list, err := extract(link.url)
	<-tokens //release the token

	if err != nil {
		log.Print(err)
	}
	var links []depthLink
	for _, el := range list {
		links = append(links, depthLink{
			url:   el,
			depth: link.depth + 1,
		})
	}

	return links
}

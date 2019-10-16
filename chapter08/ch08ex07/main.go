// Ch08ex07 makes local copies of the pages it crawls concurrently, creating
// directories as necessary, except for pages from different domains. URLs
// within mirrored pages are altered to refer to the mirrored page instead of
// the original.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bewuethr/gopl/chapter08/ch08ex07/links"
)

var tokens = make(chan struct{}, 20)

func main() {
	worklist := make(chan []string)
	go func() { worklist <- os.Args[1:] }()
	// Crawl the web breadth-first, starting from the command-line arguments.
	breadthFirst(crawl, worklist)
}

// breadthFirst applies function f to each element of each slice received from
// worklist, and feeds the output of f back into worklist.
func breadthFirst(f func(item string) []string, worklist chan []string) {
	n := 1 // Counter for slices left to process
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- f(link)
				}(link)
			}
		}
	}
}

func crawl(inURL string) []string {
	fmt.Println(inURL)
	tokens <- struct{}{} // acquire a token
	parsedURL, body, err := links.GetBody(inURL)
	<-tokens // release the token
	if err != nil {
		log.Printf("can't get %s: %v", inURL, err)
		return nil
	}

	// Check if that didn't redirect elsewhere
	parsedInURL, err := url.Parse(inURL)
	if parsedInURL.Host != parsedURL.Host {
		return nil
	}

	body = writeToFile(body, parsedURL)

	list, err := links.Extract(body, parsedURL)
	if err != nil {
		log.Print(err)
	}

	// Filter links from other domains
	filteredList := filter(list, os.Args[1:])
	return filteredList
}

// writeToFile stores body in the file corresponding to the path given in the
// URL. Because the body is read, the function returns a ReadCloser so the body
// can be still read from after this function was called.
func writeToFile(body io.ReadCloser, parsedURL *url.URL) io.ReadCloser {
	fPath := parsedURL.Host + parsedURL.Path

	// Set default filename
	if !strings.Contains(fPath, "/") {
		// For the domain index
		fPath += "/"
	}
	if strings.HasSuffix(fPath, "/") {
		// We assume it's a directory and want to avoid clashes with
		// subdirectories
		fPath += "index.html"
	}

	dirPath := filepath.Dir(fPath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Print(err)
		return body
	}

	contents, err := ioutil.ReadAll(body)
	body = ioutil.NopCloser(bytes.NewReader(contents))
	if err != nil {
		log.Print(err)
		return body
	}

	err = ioutil.WriteFile(fPath, contents, 0644)
	if err != nil {
		log.Print(err)
	}

	return body
}

// filter takes a slice of URLs and a slice includes of original pages; it
// returns a slice of all the URLs that are in the same domain as any of the
// include URLs, i.e., it filters URLs from different domains.
func filter(urls, includes []string) []string {
	// Get map of domains that are not filtered
	includeDomains := make(map[string]bool)
	for _, rawurl := range includes {
		u, err := url.Parse(rawurl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't parse %v\n", rawurl)
			os.Exit(1)
		}
		includeDomains[u.Host] = true
	}

	// Remove URLs from other domains
	var filtered []string
	for _, rawurl := range urls {
		u, err := url.Parse(rawurl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't parse %v\n", rawurl)
			os.Exit(1)
		}
		if includeDomains[u.Host] {
			// Remove references and query string to avoid downloading the same
			// page multiple times
			u.Fragment = ""
			u.RawQuery = ""
			// Remove trailing slash if present
			filtered = append(filtered, strings.TrimSuffix(u.String(), "/"))
		}
	}

	return filtered
}

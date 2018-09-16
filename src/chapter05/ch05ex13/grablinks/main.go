// grablinks makes local copies of the pages it crawls, creating directories as
// necessary, except for pages from different domains.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"chapter05/ch05ex13/links"
)

// breadthFirst calls f for each item in the worklist. Any items returned by f
// are added to the worklist. f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(inURL string) []string {
	fmt.Println(inURL)
	parsedURL, body, err := links.GetBody(inURL)
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
	if fPath[len(fPath)-1] == '/' {
		// We assume it's a directory and want to avoid clashes with
		// subdirectories
		fPath += "index.html"
	}

	dirPath := fPath[:strings.LastIndex(fPath, "/")]

	err := os.MkdirAll(dirPath, 0755)
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
			filtered = append(filtered, u.String())
		}
	}

	// Remove slashes from end of URLs
	for i, url := range filtered {
		if url[len(url)-1] == '/' {
			filtered[i] = url[:len(url)-1]
		}
	}

	return filtered
}

func main() {
	// Crawl the web breadth-first, starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

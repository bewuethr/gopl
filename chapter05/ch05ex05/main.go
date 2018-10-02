// ch05ex05 fetches the document from the URL specified as an argument and
// prints word and image count
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No URL provided")
		os.Exit(1)
	}

	words, images, err := CountWordsAndImages(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch05ex05: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Words:\t%d\nImages:\t%d\n", words, images)
}

// CountWordsAndImages does an HTTP GET request for the HTML document url and
// returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

// countWordsAndImages returns the number of words and images in the given HTML
// document.
func countWordsAndImages(n *html.Node) (words, images int) {
	words, images = visit(0, 0, n)
	return
}

// visit traverses the given HTML document recursively and returns the number
// of images and words from text nodes.
func visit(words, images int, n *html.Node) (int, int) {
	if n == nil {
		return words, images
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	} else if n.Type == html.TextNode {
		words += countWords(n.Data)
	}

	words, images = visit(words, images, n.FirstChild)
	words, images = visit(words, images, n.NextSibling)
	return words, images
}

// countWords returns the number of words s.
func countWords(s string) int {
	input := bufio.NewScanner(strings.NewReader(s))
	input.Split(bufio.ScanWords)

	var count int
	for input.Scan() {
		count++
	}

	return count
}

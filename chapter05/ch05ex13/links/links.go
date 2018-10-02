// Package links provides a link-extraction function.
package links

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// GetBody makes an HTTP GET request to the specified URL and returns the
// request URL and body of the response.
func GetBody(url string) (*url.URL, io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	return resp.Request.URL, resp.Body, nil
}

// Extract extracts the links in the HTML document contained in the supplied
// response body.
func Extract(body io.ReadCloser, url *url.URL) ([]string, error) {
	defer body.Close()
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("problems parsing HTML of %v: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := url.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

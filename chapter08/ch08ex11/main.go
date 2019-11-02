// Ch08ex11 fetches multiple URLs in parallel and cancels all other requests
// when the first one has finished.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var done = make(chan struct{}) // signal cancellation

func main() {
	flag.Parse()
	local, url, n, err := fetch(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ch08ex11: %s: %v\n", url, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
}

// fetch starts downloading all urls in parallel; as soon as the first one has
// finished, the other requests are cancelled.
func fetch(urls []string) (filename, url string, n int64, err error) {
	resp := parallelFetch(urls)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.Request.URL.String(), 0, err
	}
	close(done)

	local := path.Base(resp.Request.URL.Path)
	if local == "/" || local == "." {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", resp.Request.URL.String(), 0, err
	}
	n, err = io.Copy(f, bytes.NewReader(bodyBytes))

	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return local, resp.Request.URL.String(), n, err
}

func parallelFetch(urls []string) *http.Response {
	responses := make(chan *http.Response)
	for _, url := range urls {
		go func(url string) {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not create request for %v: %v\n", url, err)
				return
			}
			req.Cancel = done

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error fetching %v: %v\n", url, err)
				return
			}
			responses <- resp
		}(url)
	}
	return <-responses
}

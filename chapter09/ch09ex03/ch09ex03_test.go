package ch09ex03

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestSequential(t *testing.T) {
	m := New(httpGetBody)
	sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := New(httpGetBody)
	concurrent(t, m)
}

func httpGetBody(url string, done <-chan struct{}) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type M interface {
	Get(key string, done <-chan struct{}) (interface{}, error)
	Close()
}

func sequential(t *testing.T, m M) {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	for url := range incomingURLs() {
		var (
			start    = time.Now()
			done     = make(chan struct{})
			finished = make(chan struct{})
		)

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			value, err := m.Get(url, done)
			close(finished)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
		}(url)

		select {
		case <-finished:
			// request finished before random timeout
		case <-time.After(time.Duration(rand.Float32()*1500) * time.Millisecond):
			close(done) // cancel request
		}

		wg.Wait()
	}
	m.Close()
}

func concurrent(t *testing.T, m M) {
	var n sync.WaitGroup
	for url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			start := time.Now()
			done := make(chan struct{})
			value, err := m.Get(url, done)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	n.Wait()
	m.Close()
}

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

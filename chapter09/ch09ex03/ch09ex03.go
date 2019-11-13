// Package ch09ex03 provides a concurrency-safe memoization a function of a
// function.  Requests for different keys proceed in parallel. Concurrent
// requests for the same key block until the first completes. This
// implementation uses a Mutex and lets clients provide an optional channel
// through which they can cancel the operation.
package ch09ex03

import "fmt"

// Func is the type of the function to memoize.
type Func func(key string, done <-chan struct{}) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	done     <-chan struct{}
	response chan<- result // the client wants a single result
}

// Memo is a memoizing wrapper around a Func.
type Memo struct {
	requests  chan request
	deletions chan string // for keys to be deleted
}

// New returns a memoization of f. Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{
		requests:  make(chan request),
		deletions: make(chan string),
	}
	go memo.server(f)
	return memo
}

// Get is the memoizing version of Func.
func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, done, response}
	select {
	case res := <-response:
		return res.value, res.err
	case <-done:
		memo.deletions <- key
		return nil, fmt.Errorf("request for %v got cancelled", key)
	}
}

// Close closes Memo.
func (memo *Memo) Close() {
	close(memo.requests)
	close(memo.deletions)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	// First, process cancels
	for {
	OuterLoop:
		for {
			select {
			case key := <-memo.deletions:
				delete(cache, key)
			default:
				break OuterLoop
			}
		}
		// Then, process requests
		select {
		case req := <-memo.requests:
			e := cache[req.key]
			if e == nil {
				// This is the first request for this key.
				e = &entry{ready: make(chan struct{})}
				cache[req.key] = e
				go e.call(f, req.key, req.done, memo.deletions) // call f(key)
			}
			go e.deliver(req.response)
		default:
		}
	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}, deletions chan<- string) {
	e.res.value, e.res.err = f(key, done)
	// Broadcast the ready condition.
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

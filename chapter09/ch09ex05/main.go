// Ch09ex05 creates two goroutines passing messages back and forth on two
// channels. Each goroutine keeps track of how many times it receives a message;
// after five seconds, the total number of messages and the number of messages
// per second is printed.
package main

import (
	"fmt"
	"sync"
	"time"
)

var count int64

func main() {
	c1, c2 := make(chan struct{}), make(chan struct{})

	go worker(c1, c2)
	go worker(c2, c1)

	const dur = 5
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		time.Sleep(dur * time.Second)
		fmt.Printf("%.3gM total (%.3gM per second)\n", float64(count)/1e6, float64(count)/1e6/dur)
	}()

	// Kick off ping-pong
	c1 <- struct{}{}
	wg.Wait()
}

func worker(in <-chan struct{}, out chan<- struct{}) {
	for {
		// Separate read and write to guarantee sequential access to shared
		// variable count
		<-in
		count++
		out <- struct{}{}
	}
}

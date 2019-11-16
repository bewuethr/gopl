// Ch09ex04 takes an optional parameter to indicate the number of goroutines to
// create, then creates that number of channels and goroutines and finally
// measures how long a message takes to pass through all of them.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var nChan int // number of channels
	if len(os.Args) < 2 {
		nChan = 10
	} else {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		nChan = n
	}

	// Create nChan channels
	channels := make([]chan struct{}, nChan)
	for i := range channels {
		channels[i] = make(chan struct{})
	}

	fmt.Printf("Creating %d goroutines/channels\n", nChan)

	// Create nChan goroutines
	for i := 0; i < len(channels)-1; i++ {
		go func(in <-chan struct{}, out chan<- struct{}) {
			out <- <-in
		}(channels[i], channels[i+1])
	}

	first, last := make(chan struct{}), make(chan struct{})

	// First worker
	go func(in <-chan struct{}, out chan<- struct{}) {
		out <- <-in
	}(first, channels[0])

	// Last worker
	go func(in <-chan struct{}, out chan<- struct{}) {
		out <- <-in
	}(channels[len(channels)-1], last)

	start := time.Now()
	first <- struct{}{}
	<-last
	elapsed := time.Since(start)
	fmt.Printf("time: %v\n", elapsed)
}

// Ch08ex08 is a reverb server that disconnects any client that doesn't shout
// anything for 10 seconds.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	shout := make(chan string)
	input := bufio.NewScanner(c)
	go func() {
		for input.Scan() {
			shout <- input.Text()
		}
	}()
OuterLoop:
	for {
		select {
		case <-time.After(10 * time.Second):
			break OuterLoop
		case t := <-shout:
			echo(c, t, 1*time.Second)
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}

// Ch08ex13 is the chat server from Chapter 8 Section 10: it disconnects idle
// clients (10 seconds timeout, adjustable with -timeout SECONDS flag).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	timeout  time.Duration
)

func init() {
	var timeoutFlag = flag.Int("timeout", 10, "time in seconds after which idle clients are disconnected")
	flag.Parse()
	timeout = time.Duration(*timeoutFlag) * time.Second
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all clients' outgoing message
			// channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	listener := make(chan string) // for incoming client messages
	ch := make(chan string)       // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			listener <- who + ": " + input.Text()
		}
		// NOTE: ignoring potential errors from input.Err()
		close(listener) // signal client has disconnected
	}()

outerloop:
	for {
		select {
		case msg, ok := <-listener:
			if !ok {
				// client has disconnected
				break outerloop
			}
			messages <- msg
		case <-time.After(timeout):
			break outerloop
		}
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

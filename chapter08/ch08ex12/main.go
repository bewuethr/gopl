// Ch08ex12 is the chat program from Chapter 8.10, announcing the current
// clients to new arrivals.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

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

type clientInfo struct {
	clientChan chan<- string
	clientName string
}

type client chan<- clientInfo // an outgoing message channel

var (
	entering = make(chan clientInfo)
	leaving  = make(chan clientInfo)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[clientInfo]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all clients' outgoing message
			// channels.
			for cli := range clients {
				cli.clientChan <- msg
			}

		case cli := <-entering:
			if len(clients) > 0 {
				cli.clientChan <- "chat participants:"
				for c := range clients {
					cli.clientChan <- c.clientName
				}
			} else {
				cli.clientChan <- "you're the only chat participant"
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.clientChan)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- clientInfo{
		clientChan: ch,
		clientName: who,
	}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- clientInfo{
		clientChan: ch,
		clientName: who,
	}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

// Clockwall is a client of multiple clock servers at once; it reads the time
// from each one and displays the results in a table. It takes table headings
// and server address/port as parameters, for example:
//
//    ./clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type clock struct {
	name string
	host string
	port int
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: ./clockwall CITY=HOST:PORT [CITY=HOST:PORT ...]")
		os.Exit(1)
	}

	clocks := []clock{}
	maxLen := 0
	for _, arg := range os.Args[1:] {
		fields := strings.FieldsFunc(arg, func(c rune) bool {
			return strings.ContainsRune("=:", c)
		})
		p, err := strconv.Atoi(fields[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockwall: %v\n", err)
			os.Exit(1)
		}
		clocks = append(clocks, clock{
			name: fields[0],
			host: fields[1],
			port: p,
		})
		if len(fields[0]) > maxLen {
			maxLen = len(fields[0])
		}
	}

	for _, c := range clocks {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		go readClock(c, conn, maxLen)
	}

	for {
		time.Sleep(time.Minute)
	}
}

func readClock(c clock, r io.Reader, maxLen int) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		fmt.Printf("%-*s %v\n", maxLen+1, c.name+":", s.Text())
	}
	fmt.Println("done")
	if s.Err() != nil {
		fmt.Printf("lost %s %s\n", c.name, s.Err())
	}
}

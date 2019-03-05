// Ch08ex02 is a minimal implementation of a concurrent FTP server.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	ctrlPort = 21
	dataPort = 20
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", ctrlPort))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	clientAddr := c.RemoteAddr().(*net.TCPAddr)
	log.Printf("%v: new session\n", clientAddr)

	wd, err := os.Getwd()
	if err != nil {
		io.WriteString(c, "421 Service not feeling well today :/\r\n")
		log.Printf("%v: responded with 421, can't get working directory\n", clientAddr)
		return
	}

	io.WriteString(c, "220 Service ready for action\r\n")

	var (
		user     string
		dataAddr = &net.TCPAddr{
			IP:   clientAddr.IP,
			Port: clientAddr.Port,
		}
	)

	scanner := bufio.NewScanner(c)
	scanner.Split(scanTelnetCmd)

	for scanner.Scan() {
		log.Printf("%v: received command %v\n", clientAddr, scanner.Text())
		fields := strings.Fields(scanner.Text())

		switch fields[0] {
		case "CWD":
			newDir, err := changeWorkingDir(c, clientAddr, user, wd, fields)
			if err != nil {
				continue
			}
			wd = newDir

		case "LIST":
			list(user, c, clientAddr, dataAddr, wd)

		case "PORT":
			ip, port, err := port(user, fields, c, clientAddr)
			if err != nil {
				continue
			}
			dataAddr.IP = ip
			dataAddr.Port = port

		case "PWD":
			workingDir(c, clientAddr, wd)

		case "QUIT":
			quit(c, clientAddr)
			return

		case "RETR":
			retrieve(user, wd, clientAddr, dataAddr, c, fields)

		case "SYST":
			system(c, clientAddr)

		case "USER":
			user = userName(fields, c, clientAddr)

		case "TYPE", "MODE", "STRU", "STOR", "NOOP":
			notImplemented(c, clientAddr)

		default:
			notRecognized(c, clientAddr)
		}
	}
}

// scanTelnetCmd is a split function for scanner that splits on CRLF terminated
// lines such as Telnet commands.
func scanTelnetCmd(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\r\n")); i >= 0 {
		// We have a full CRLF-terminated line.
		return i + 2, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// changeWorkingDir implements the CWD command.
func changeWorkingDir(c net.Conn, addr *net.TCPAddr, user, wd string, fields []string) (string, error) {
	if user == "" {
		io.WriteString(c, "530 I can't do that for non-logged in users.\r\n")
		log.Printf("%v: responded with 530\n", addr)
		return "", errors.New("user not logged in")
	}

	path := fields[1]
	newDir := filepath.Clean(wd + "/" + path)
	fileInfo, err := os.Stat(newDir)
	if err != nil {
		io.WriteString(c, fmt.Sprintf("550 Oops: %v\r\n", err))
		log.Printf("%v: responded with 550, %v\n", addr, err)
		return "", errors.New("could not stat new directory")
	}

	if !fileInfo.IsDir() {
		io.WriteString(c, fmt.Sprintf("550 Oops: %v is no not a directory\r\n", path))
		log.Printf("%v: responded with 550, %v not a directory\n", addr, path)
		return "", errors.New("not a directory")
	}

	wd = newDir
	io.WriteString(c, "250 Requested file action fiiiine, completed.\r\n")
	log.Printf("%v: responded with 250, changed directory to %v\n", addr, newDir)

	return wd, nil
}

func list(user string, c net.Conn, clientAddr, dataAddr *net.TCPAddr, wd string) {
	if user == "" {
		io.WriteString(c, "530 I can't do that for non-logged in users.\r\n")
		log.Printf("%v: responded with 530\n", clientAddr)
		return
	}

	files, err := ioutil.ReadDir(wd)
	if err != nil {
		io.WriteString(c, "450 File not available.\r\n")
		log.Printf("%v: could not read working directory: %v\n", clientAddr, err)
		return
	}

	io.WriteString(c, "150 File status okay; about to open data connection.\r\n")
	log.Printf("%v: opening data connection to %+v\n", clientAddr, dataAddr)
	dataConn, err := net.Dial("tcp", dataAddr.String())
	if err != nil {
		io.WriteString(c, "425 Can't open the data connection >:(\r\n")
		log.Printf("%v: could not open data connection: %v\n", clientAddr, err)
		return
	}
	defer dataConn.Close()

	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t)
		io.WriteString(dataConn, fmt.Sprintf("%s %3d %d %d %12d %s %s\r\n",
			file.Mode().String(), stat.Nlink, stat.Uid, stat.Gid, file.Size(),
			file.ModTime().Format("2006-01-02 15:04:05"), file.Name()))
	}

	io.WriteString(c, "226 Directory listing complete\r\n")
	log.Printf("%v: listing complete\n", clientAddr)
}

// port implements the PORT command.
func port(user string, fields []string, c net.Conn, clientAddr *net.TCPAddr) (net.IP, int, error) {
	if user == "" {
		io.WriteString(c, "530 I can't do that for non-logged in users.\r\n")
		log.Printf("%v: responded with 530\n", clientAddr)
		return nil, 0, errors.New("user not logged in")
	}

	csv := strings.Split(fields[1], ",")
	ipStr := strings.Join(csv[:4], ".")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		io.WriteString(c, "501 That address with ports is fishy.\r\n")
		log.Printf("%v: could not parse %s as  IP address\n", clientAddr, ipStr)
		return nil, 0, errors.New("could not parse IP address")
	}

	portHighBit, err := strconv.Atoi(csv[4])
	if err != nil {
		io.WriteString(c, "501 That address with ports is fishy.\r\n")
		log.Printf("%v: could not convert high bit of port: %v\n", clientAddr, err)
		return nil, 0, errors.New("could not convert high bit of port")
	}

	portLowBit, err := strconv.Atoi(csv[5])
	if err != nil {
		io.WriteString(c, "501 That address with ports is fishy.\r\n")
		log.Printf("%v: could not convert low bit of port: %v\n", clientAddr, err)
		return nil, 0, errors.New("could not convert low bit of port")
	}

	newPort := 256*portHighBit + portLowBit

	io.WriteString(c, "200 Port and address updated!\r\n")
	log.Printf("%v: set client data transfer connection address to %s:%d\n", clientAddr, ip, newPort)

	return ip, newPort, nil
}

// workingDir implements the PWD command.
func workingDir(c net.Conn, addr *net.TCPAddr, wd string) {
	io.WriteString(c, fmt.Sprintf("257 \"%v\" is the current directory\r\n", wd))
	log.Printf("%v: responded with 257, current directory is %v\n", addr, wd)
}

// quit implements the QUIT commmand.
func quit(c net.Conn, addr *net.TCPAddr) {
	io.WriteString(c, "221 Service closing down. Forever.\r\n")
	log.Printf("%v: responded with 221, closing connection\n", addr)
}

// retrieve implements the RETR command.
func retrieve(user, wd string, clientAddr, dataAddr *net.TCPAddr, c net.Conn, fields []string) {
	if user == "" {
		io.WriteString(c, "530 I can't do that for non-logged in users.\r\n")
		log.Printf("%v: responded with 530\n", clientAddr)
		return
	}

	file, err := os.Open(filepath.Clean(wd + "/" + fields[1]))
	if err != nil {
		if os.IsTimeout(err) {
			io.WriteString(c, "450 File not available, timeout.\r\n")
		} else {
			io.WriteString(c, "550 File not available.\r\n")
		}
		log.Printf("%v: could not open requested file: %v\n", clientAddr, err)
		return
	}
	defer file.Close()

	io.WriteString(c, "150 File status okay; about to open data connection.\r\n")
	log.Printf("%v: opening data connection to %+v\n", clientAddr, dataAddr)
	dataConn, err := net.Dial("tcp", dataAddr.String())
	if err != nil {
		io.WriteString(c, "425 Can't open the data connection >:(\r\n")
		log.Printf("%v: could not open data connection: %v\n", clientAddr, err)
		return
	}
	defer dataConn.Close()

	fs := bufio.NewScanner(file)
	for fs.Scan() {
		io.WriteString(dataConn, fmt.Sprintf("%s\r\n", fs.Text()))
	}
	if err := fs.Err(); err != nil {
		io.WriteString(c, "426 Transfer aborted, connection closed\r\n")
		log.Printf("%v: could not copy file to data connection: %v\n", clientAddr, err)
		return
	}

	io.WriteString(c, "250 File retrieval successful\r\n")
	log.Printf("%v: download complete\n", clientAddr)
}

// system implements the SYST command.
func system(c net.Conn, addr *net.TCPAddr) {
	io.WriteString(c, "215 UNIX system type. (Actually Linux.)\r\n")
	log.Printf("%v: returned system type\n", addr)
}

// userName implements the USER command.
func userName(fields []string, c net.Conn, addr *net.TCPAddr) string {
	user := strings.Join(fields[1:], " ")
	io.WriteString(c, "230 User logged in, go on.\r\n")
	log.Printf("%v: user %v logged in\n", addr, user)
	return user
}

// notImplemented is run for recognized, but not implemented commands.
func notImplemented(c net.Conn, addr *net.TCPAddr) {
	io.WriteString(c, "502 Command not (yet) implemented.\r\n")
	log.Printf("%v: responded with 502\n", addr)
}

// notRecognized is for commands that aren't recognized.
func notRecognized(c net.Conn, addr *net.TCPAddr) {
	io.WriteString(c, "500 Command not recognized :'(\r\n")
	log.Printf("%v: responded with 500\n", addr)
}

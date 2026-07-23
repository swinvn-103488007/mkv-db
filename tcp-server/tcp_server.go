package main

import (
	"net"
	"fmt"
	"io"
	"os"
	"trildd/mkv-db/multithread-counter"
)

var tcp_port = ":1606"

func main() {
	counter := counter.NewCounter()
	listener, err := net.Listen("tcp", tcp_port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %v, %v\n", tcp_port, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Start listening on %v\n", tcp_port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		handleConnection(conn, counter)
	}
}

func handleConnection(conn net.Conn, counter *counter.Counter) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		numBytes, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading from %s: %v\n", conn.RemoteAddr(), err)
			} else {
				fmt.Printf("Connection closed by %s\n", conn.RemoteAddr())
			}
			return
		}
		counter.Increment()
		// Print the raw bytes received
		fmt.Printf("Received %d bytes from %s: %v\n. Request count: %v", numBytes, conn.RemoteAddr(), buffer[:numBytes], counter.Value())
		// Also print as string for readability (if printable)
		fmt.Printf("As string: %q\n", buffer[:numBytes])
	}
}
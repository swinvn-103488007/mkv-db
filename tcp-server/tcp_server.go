package main

import (
	"net"
	"fmt"
	"io"
	"os"
	"trildd/mkv-db/config"
)

const ACCEPT_MSG = "accept_message"

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.TCP_PORT))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %v, %v\n", config.TCP_PORT, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Start listening on %v\n", config.TCP_PORT)
	semaphore := make(chan bool, config.MAX_CLIENT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}
		semaphore <- true
		conn.Write([]byte(fmt.Sprintf("%v\n", ACCEPT_MSG)))
		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		go handleConnection(conn, semaphore)
	}
}

func handleConnection(conn net.Conn, semaphore chan bool) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		numBytes, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading from %s: %v\n", conn.RemoteAddr(), err)
			} else {
				fmt.Printf("Connection closed by %s\n", conn.RemoteAddr())
				<-semaphore
			}
			return
		}
		// Print the raw bytes received
		fmt.Printf("Received %d bytes from %s: %v.\n", numBytes, conn.RemoteAddr(), buffer[:numBytes])
		// Also print as string for readability (if printable)
		fmt.Printf("As string: %q\n", buffer[:numBytes])
	}
}
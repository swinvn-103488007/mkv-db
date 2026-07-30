// package tcpclient

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const ACCEPT_MSG = "accept_message"

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1606")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to 127.0.0.1:1606")
	fmt.Println("Wait until server is ready to take message...")
	waitReadyMsg(conn)
	clearStdin()
	talkingToServer(conn)
}

func waitReadyMsg(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while waiting for start signal: %v\n", err)
			conn.Close()
			os.Exit(1)
			return
		}

		msg := strings.TrimSpace(line)
		fmt.Println(msg)
		if msg == ACCEPT_MSG {
			fmt.Println("Server is ready to take message. Proceeding...")
			conn.SetReadDeadline(time.Time{})
			break
		}
	}
}

// helper to clear any buffered data caused by user during wait time
func clearStdin() {
	fmt.Println("Clearing coincident buffer")
	done := make(chan struct{})
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			// just discard the data
		}
		close(done)
	}()

	select {
	case <-done:
		// successfully drained
	case <-time.After(150 * time.Millisecond):
	// 	// timeout – we assume the buffer is clear enough
	}
	fmt.Println("Done clear coincident buffer")
}

func talkingToServer(conn net.Conn) {
	fmt.Println("Start talking to server")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}
		msg := strings.TrimSpace(line)
		if msg == "quit" || msg == "exit" {
			fmt.Println("Closing connection...")
			break
		}

		// Send the message with a newline so the server sees it as a complete line
		_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send: %v\n", err)
			break
		}
	}
	fmt.Println("Disconnected.")
}
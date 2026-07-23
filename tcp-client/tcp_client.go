// package tcpclient

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1606")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to 127.0.0.1:1606")
	fmt.Println("Type messages and press Enter to send. Type 'quit' or 'exit' to leave.")

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
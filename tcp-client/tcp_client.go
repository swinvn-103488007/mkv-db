// package tcpclient

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"syscall"
	"strings"
	"time"
)

const (
	TCIFLUSH = 0
	TCFLSH   = 0x540B // Linux; on macOS/BSD use syscall.TIOCFLUSH instead
	ACCEPT_MSG = "accept_message"
)

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
	DrainStdin()
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
// NOTE: this helper cannot truly clear/drain the Stdin buffer
// func clearStdin() {
// 	fmt.Println("Clearing coincident buffer")
// 	done := make(chan struct{})
// 	go func() {
// 		buf := make([]byte, 4096)
// 		for {
// 			n, _ := os.Stdin.Read(buf)
// 			fmt.Printf("Read buffer as string: %q\n", buf[:n])
// 			if n == 0 {
// 				fmt.Println("End of old buffered data")
// 				break
// 			}
// 		}
// 		close(done)
// 	}()

// 	<-done
// 	fmt.Println("Done clear coincident buffer")
// }

// DrainStdin discards any unread input currently buffered by the terminal.
func DrainStdin() error {
	fd := os.Stdin.Fd()
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(TCFLSH),
		uintptr(TCIFLUSH),
	)
	if errno != 0 {
		return errno
	}
	return nil
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
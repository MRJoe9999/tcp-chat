package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Channel to receive server messages
	go receiveMessages(conn)

	// Read from stdin and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/quit" {
			fmt.Println("Disconnecting...")
			return
		}
		sendTime := time.Now()
		fmt.Printf("Sent at: %s\n", sendTime.Format("15:04:05.000"))
		_, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

	}

	if scanner.Err() != nil {
		fmt.Println("Error reading from input:", scanner.Err())
	}
}

// receiveMessages reads incoming messages from server
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Diconnected from server.")
			os.Exit(0)
		}
		receiveTime := time.Now()
		fmt.Printf("Received at: %s → %s", receiveTime.Format("15:04:05.000"), message)
	}
}

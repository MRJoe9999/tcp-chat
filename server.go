package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]bool) // All connected clients
	clientsMu sync.Mutex                // Protects the clients map
	broadcast = make(chan string)       // Channel for broadcasting messages
)

func main() {
	// Start listening on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080")

	// Start the broadcaster goroutine
	go handleBroadcast()

	for {
		// Accept new client connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Add client to map
		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		fmt.Println("New client connected:", conn.RemoteAddr())

		// Handle client in a separate goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		// Clean up when client disconnects
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		fmt.Println("Client disconnected:", conn.RemoteAddr())
	}()
	conn.Write([]byte("Welcome to the chat server! Enter your name to start chatting.\n"))
	reader := bufio.NewReader(conn)
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]
	conn.Write([]byte("Hello " + username + "! You can start chatting now.\n"))

	broadcast <- fmt.Sprintf("%s has joined the chat.\n", username)
	for {
		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			break // Client disconnected or error
		}

		// Send message to broadcast channel
		broadcast <- fmt.Sprintf("%s %s %s", username, conn.RemoteAddr(), message)

	}
}

func handleBroadcast() {
	for {
		// Receive a message to broadcast
		msg := <-broadcast

		// Send the message to all clients
		clientsMu.Lock()
		for client := range clients {
			_, err := fmt.Fprint(client, msg)
			if err != nil {
				// Problem with this client; close it
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

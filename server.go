package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type BroadcastMessage struct {
	sender net.Conn
	text   string
}

var (
	clients   = make(map[net.Conn]bool) // All connected clients
	clientsMu sync.Mutex                // Protects the clients map
	broadcast = make(chan BroadcastMessage)
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
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		fmt.Println("Client disconnected:", conn.RemoteAddr())
	}()

	conn.Write([]byte("Welcome to the chat server! Enter your name to start chatting.\n"))
	reader := bufio.NewReader(conn)
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // Remove newline
	conn.Write([]byte("Hello " + username + "! You can start chatting now.\n"))

	broadcast <- BroadcastMessage{
		sender: conn,
		text:   fmt.Sprintf("%s has joined the chat.\n", username),
	}

	for {
		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			break // Client disconnected or error
		}

		// Send message to broadcast channel
		broadcast <- BroadcastMessage{
			sender: conn,
			text:   fmt.Sprintf("%s [%s]: %s", username, conn.RemoteAddr(), message),
		}
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast

		clientsMu.Lock()
		for client := range clients {
			if client != msg.sender {
				_, err := fmt.Fprint(client, msg.text)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
		clientsMu.Unlock()
	}
}

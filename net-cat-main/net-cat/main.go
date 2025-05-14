package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"net-cat/chat"
)



func main() {
	// Default port is 8989
	port := "8989"

	if len(os.Args) > 2 {
		// If the argument is not a valid port number, show usage and exit
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Check if a port argument is provided
	if len(os.Args) == 2 {
		// Check if the argument is a valid number
		if _, err := strconv.Atoi(os.Args[1]); err != nil {
			// If the argument is not a valid port number, show usage and exit
			fmt.Println("[USAGE]: ./TCPChat $port (the port must be a number)")
			return
		}
		port = os.Args[1]
	}

	// Create a new server
	s := chat.NewServer()
	go s.Run()

	// Listen on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	defer listener.Close()

	log.Printf("Listening on the port :%s", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}
		chat.Count++

		if chat.Count > 10 {
			conn.Close()
			chat.Count--
			continue
		}
		c := s.NewClient(conn)
		go c.ReadInput()
	}
}

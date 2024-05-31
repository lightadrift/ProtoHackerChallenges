package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send some data to the server
	_, err = conn.Write([]byte(`{"method": "test", "number": 3}`))
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, 1024)

	// Read the response from the server
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	// Convert the bytes read into a string
	responseStr := string(buffer[:n])

	// Print the response received from the server
	fmt.Println("Server response:", responseStr)

	// Close the connection
	conn.Close()
}

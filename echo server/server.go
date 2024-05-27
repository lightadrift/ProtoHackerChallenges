package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var logger *log.Logger

func main() {
	// Open log file
	logFile, err := os.OpenFile("echo_server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	logger = log.New(logFile, "", log.LstdFlags)

	port := fmt.Sprintf(":%s", os.Args[1])
	// Map of allowed IPs for faster lookup
	allowedIPsMap := map[string]bool{
		"206.189.113.124": true,
		"127.0.0.1":       true,
		"::1":             true,
		// Add more IPs as needed
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Server initialization failed, error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server initialized on port:", port)

	connections := make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection, error:", err)
				close(connections)
				return
			}
			ip, _, err := net.SplitHostPort(strings.TrimSpace(conn.RemoteAddr().String()))
			if err != nil {
				fmt.Println("Error splitting host and port:", err)
				conn.Close()
				continue
			}
			// Check if IP is in the allowed IPs map
			if !allowedIPsMap[ip] {
				fmt.Println("Connection from unauthorized IP:", ip)
				conn.Close()
				continue
			}
			connections <- conn
		}
	}()

	for conn := range connections {
		go HandleConnections(conn)
	}
}

func HandleConnections(conn net.Conn) {
	defer conn.Close()
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("request:")

}

package main

import (
	"bufio"
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
	prefix := os.Args[2]
	allowedIP := "206.189.113.124"

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
			if err != nil || ip != allowedIP {
				fmt.Println("Connection from unauthorized IP:", ip)
				conn.Close()
				continue
			}
			connections <- conn
		}
	}()

	for conn := range connections {
		go HandleConnections(conn, prefix)
	}
}

func HandleConnections(conn net.Conn, prefix string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read data, error:", err)
			}
			return
		}
		request := fmt.Sprintf("%s", bytes)
		response := request                                             // Assuming the response is simply echoing the request
		logger.Printf("Request: %s\nResponse: %s\n", request, response) // Log the request and response
		conn.Write([]byte(response))
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	port := fmt.Sprintf(":%s", os.Args[1])
	prefix := os.Args[2]
	allowedIP := "206.189.113.124"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("server initialization failed, error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("server initialized in port:", port)

	connections := make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("failed to accpet connection, error:", err)
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
				fmt.Println("failed to read data, error:", err)
			}
			return
		}
		fmt.Printf("request: %s", bytes)
		line := fmt.Sprintf("%s %s", prefix, bytes)
		fmt.Printf("response: %s", line)
		conn.Write([]byte(line))
	}

}

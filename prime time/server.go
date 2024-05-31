package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

type Request struct {
	Method string      `json:"method"`
	Number interface{} `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func checkIfPrime(f float64) bool {
	x := new(big.Float).SetFloat64(f)
	if x.IsInt() {
		intVal, _ := x.Int64()
		return isPrime(intVal)
	}
	return false
}

func process_request(req Request, res *Response) {
	switch n := req.Number.(type) {
	case int:
		if big.NewInt(int64(n)).ProbablyPrime(0) {
			res.Prime = true
		} else {
			res.Prime = false
		}
	case float64:
		if checkIfPrime(n) {
			res.Prime = true
		} else {
			res.Prime = false
		}
	default:
		res.Prime = false

	}
}

func main() {
	port := fmt.Sprintf(":%s", os.Args[1])
	allowedIPsMap := map[string]bool{
		"206.189.113.124": true,
		"127.0.0.1":       true,
		"::1":             true,
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
	buf := make([]byte, 1024)

	n, err := conn.Read(buf[:])

	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	str := string(buf[:n])

	var request Request
	var response = &Response{Method: "isPrime"}
	err = json.Unmarshal([]byte(str), &request)

	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
		return
	}

	fmt.Printf("Received string: %s", request.Method)

	process_request(request, response)

	res_json, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
	} else {
		fmt.Println(string(res_json))
		conn.Write([]byte(res_json))
	}

}

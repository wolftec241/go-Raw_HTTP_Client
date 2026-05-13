package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// TCP client - example
/*
func StartTCPClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send a message
	fmt.Fprintln(conn, "Hello Server!")
	fmt.Fprintln(conn, "aya")

	// Read response
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Server replied:", response)
}
*/

func StartClient() {
	listener, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	request := "GET / HTTP/1.1\r\n" +
		"Host: httpbin.org\r\n" +
		"Connection: keep-alive\r\n" +
		"\r\n"

	listener.Write([]byte(request))

	response := bufio.NewReader(listener)
	msg, _ := io.ReadAll(response)
	fmt.Print("Server replied:", string(msg))
}

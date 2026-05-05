package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// TCP server - example
/*
func handleTCPConnection(conn net.Conn) {
	defer conn.Close() // always close when done

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println("Received:", msg)
		conn.Write([]byte("Echo: " + msg + "\n"))
	}
}

func StartTCPServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept() // blocks until a client connects
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleTCPConnection(conn) // handle each client in a goroutine
	}
}
*/

func handleConnection(conn net.Conn) {
	defer conn.Close() // always close when done

	scanner := bufio.NewReader(conn)
	line, _ := scanner.ReadString('\n')
	if !validMethod(line) {
		conn.Write([]byte("Error 405" + "\n"))
		return
	}
	response := "HTTP/1.1 200\r\n" + "\r\n"

	for {
		line, _ := scanner.ReadString('\n')
		if line == "\r\n" || line == "\n" || line == "" {
			break
		}

		fmt.Println("Received: ", line)
		//response += line + "\r\n"
	}
	response += "aba" + "\r\n"
	conn.Write([]byte(response))
}

func validMethod(line string) bool {
	methods := map[string]struct{}{
		"GET":    {},
		"SET":    {},
		"POST":   {},
		"DELETE": {},
	}

	slices := strings.Split(line, " ")
	if len(slices) != 3 {
		return false
	}

	_, methodExist := methods[slices[0]]
	if (!methodExist) || slices[2] != "HTTP/1.1\r\n" {
		fmt.Println(slices[2])
		return false
	}
	return true
}

func StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	//defer listener.Close()
	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error", err)
			continue
		}

		go handleConnection(conn)

	}
}

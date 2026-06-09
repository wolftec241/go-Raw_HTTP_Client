package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
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

type request struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

func sortRequest(conn net.Conn) *request {
	req := request{Headers: make(map[string]string)}
	var haveContentLength bool = false
	var contentLength int
	scanner := bufio.NewReader(conn)

	// 1. Deal with first line(Method, Path, Protocol)
	line, err := scanner.ReadString('\n') // Read the first line
	if err != nil {
		if err == io.EOF {
			fmt.Println("Empty request")
			return nil
		}
		fmt.Println("Error with ReadString: ", err)
		return nil
	}

	headParts := strings.Split(line, " ")
	if len(headParts) != 3 {
		return nil
	}

	req.Method = headParts[0]
	req.Path = headParts[1]
	req.Protocol = headParts[2]

	// 2. Handle Headers
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println("Error with ReadString: ", err)
			return nil
		}

		if line == "\r\n" {
			break
		}

		head := strings.SplitN(line, ": ", 2)
		if len(head) == 2 {
			req.Headers[head[0]] = strings.TrimSpace(head[1])
			if head[0] == "Content-Length" {
				contentLength, err = strconv.Atoi(strings.TrimSpace(head[1]))
				if err != nil {
					fmt.Println("Error with strconv: ", err)
					return nil
				}
				haveContentLength = true
			}
		}
	}

	// 3. Handle Body(if exist)
	// Check if should be body but didn't get Content Length
	if !haveContentLength && req.Method == "POST" {
		fmt.Println("411 Length Required")
		return nil
	}

	if haveContentLength && contentLength > 0 {
		bodyBuff := make([]byte, contentLength)
		_, err := io.ReadFull(scanner, bodyBuff)
		if err != nil {
			fmt.Println("Problem with reading the body part")
			return nil
		}
		req.Body = string(bodyBuff)
	}

	return &req
}

func handleRequest(req *request) string {
	if req.Path != "/" {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	protocols := map[string]struct{}{
		"GET":    {},
		"SET":    {},
		"POST":   {},
		"DELETE": {},
	}

	_, protocolExist := protocols[req.Protocol]
	if (!protocolExist) || req.Method != "HTTP/1.1\r\n" {
		return "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
	}

	return "HTTP/1.1 200 OK\r\n\r\n"
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req := sortRequest(conn)
	if req == nil {
		fmt.Println("Problem with request")
		return
	}

	conn.Write([]byte(handleRequest(req)))
}

func StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
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

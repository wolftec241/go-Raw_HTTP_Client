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

/*
func sortRequest(conn net.Conn) *request {
	req := request{Headers: make(map[string]string)}

	// Separate Header and Body
	parts := strings.Split(str, "\r\n\r\n")
	headerSection := parts[0]
	if len(parts) > 1 {
		req.Body = parts[1]
	}

	// Split Header section into lines
	lines := strings.Split(headerSection, "\r\n")

	// Parse Request Line
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) >= 2 {
		req.Method = requestLine[0]
		req.Url = requestLine[1]
	} else {
		fmt.Printf("Not valid request")
	}

	// Parse Headers
	for i := 1; i < len(lines); i++ {
		headerPart := strings.SplitN(lines[i], ": ", 2)
		if len(headerPart) == 2 {
			req.Headers[headerPart[0]] = headerPart[1]
		}
	}

	return &req

}
*/

func sortRequest(conn net.Conn) *request {
	req := request{Headers: make(map[string]string)}
	var haveContentLength bool = false
	var contentLength int
	total := 0 // For future contentLength check
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
	if !validHead(line, &req) {
		fmt.Println("Error 405 Method Not Allowed")
		conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
		return nil
	}
	total += len(line)

	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println("Error with ReadString: ", err)
			return nil
		}
		total += len(line)

		if line == "\r\n" {
			break
		}

		head := strings.SplitN(line, ": ", 2)
		req.Headers[head[0]] = strings.TrimSpace(head[1])
		if head[0] == "Content-Length" {
			contentLength, err = strconv.Atoi(strings.TrimSpace(head[1]))
			if err != nil {
				fmt.Println("Error with strconv: ", err)
				return nil
			}
		}
	}

	if !haveContentLength {
		fmt.Println("411 Length Required")
		return nil
	}

	body, err := io.ReadFull(scanner, make([]byte, contentLength))
	if err != nil {
		fmt.Println("Problem with reading the body part")
		return nil
	}
	req.Body = string(body)

	return &req
}

/*
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewReader(conn)

	for {
		requestLine, err := scanner.ReadString('\n')
		if err != nil {
			return
		}

		if !validMethod(requestLine) {
			conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
			return
		}

		keepAlive := true
		for {
			line, err := scanner.ReadString('\n')
			if err != nil {
				return
			}
			if line == "\r\n" || line == "\n" || line == "" {
				break
			}
			fmt.Println("Received: ", line)
			if strings.EqualFold(strings.TrimSpace(line), "connection: close") {
				keepAlive = false
			}
		}

		response := "HTTP/1.1 200 OK\r\n"
		if keepAlive {
			response += "Connection: keep-alive\r\n"
		} else {
			response += "Connection: close\r\n"
		}
		response += "\r\naba\r\n"
		conn.Write([]byte(response))

		if !keepAlive {
			return
		}
	}
}
*/

// Check if head of the request is valid and put the relevant data to request struct
func validHead(line string, req *request) bool {
	methods := map[string]struct{}{
		"GET":    {},
		"SET":    {},
		"POST":   {},
		"DELETE": {},
	}

	headParts := strings.Split(line, " ")
	if len(headParts) != 3 {
		return false
	}

	_, methodExist := methods[headParts[0]]
	if (!methodExist) || headParts[2] != "HTTP/1.1\r\n" {
		fmt.Println(headParts[2])
		return false
	}

	req.Method = headParts[0]
	req.Path = headParts[1]
	req.Protocol = headParts[2]
	return true
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	//scanner := bufio.NewReader()
	//req := sortRequest(conn)

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

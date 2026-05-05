package server

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("connection test", func(t *testing.T) {
		client, server := net.Pipe()
		defer client.Close()

		// run handler on server side
		go handleConnection(server)

		//send raw HTTP request
		fmt.Fprint(client, "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")

		//read response
		buf := make([]byte, 1024)
		n, _ := client.Read(buf)
		response := string(buf[:n])

		if !strings.Contains(response, "HTTP/1.1 200") {
			t.Errorf("expected 200, got: %s", response)
		}

	})

	t.Run("response test", func(t *testing.T) {
		client, server := net.Pipe()
		defer client.Close()

		// run handler on server side
		go handleConnection(server)

		//send raw HTTP request
		fmt.Fprint(client, "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")

		//read response
		buf := make([]byte, 1024)
		n, _ := client.Read(buf)
		response := string(buf[:n])

		if !strings.Contains(response, "HTTP/1.1 200") {
			t.Errorf("expected 200, got: %s", response)
		}

	})
}

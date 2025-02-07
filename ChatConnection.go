package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// ChatConnection : A wrapper for net.TCPConn to provide chat convencience methods
type ChatConnection struct {
	net.Conn
	server *ChatServer
}

func (connection *ChatConnection) receive() *string {
	result, err := bufio.NewReader(connection).ReadString('\n')

	if err != nil {
		log.Printf("%s: Couldn't read from client connection: %v\n", app.Name, err)
		return nil
	}

	result = strings.TrimSpace(result)
	return &result
}

func (connection *ChatConnection) send(what string) {
	// Unhandled error here, add a few lines to log the error if one occurs during write
	connection.Write([]byte(what + "\n"))
}

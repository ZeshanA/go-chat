//revive:disable:dot-imports

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
)

func main() {
	godotenv.Load()

	host := os.Getenv("CHAT_HOST")
	port, _ := strconv.Atoi(os.Getenv("CHAT_PORT"))

	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	server := ChatServer{
		listener,
		[]*ChatChannel{},
		[]*ChatClient{},
	}

	if err != nil {
		log.Fatalf("%s: Couldn't listen on %s:%d:\n %v", app.Prefix, host, port, err)
	}

	log.Printf("%s: Server listening on: %v\n", app.Prefix, server.Addr())

	defer server.Close()

	for {
		connection, err := server.Accept()

		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}

		conn := ChatConnection{
			connection.(*net.TCPConn),
			&server,
		}

		user := connect(conn, server)
		go user.receive()

		server.users = append(server.users, user)
		conn.send("\nAvailable commands for this server are:")
		user.Help()
	}
}

func connect(connection ChatConnection, server ChatServer) *ChatClient {
	client := &ChatClient{connection: &connection}
	server.users = append(server.users, client)

	intro := fmt.Sprintf("🚀 %s v%s\n", Blue(app.Name), Bold(app.Version))
	connection.send(intro)

	connection.send(fmt.Sprintf("%v What shall we call you?", Magenta(">")))
	client.username = *connection.receive()
	connection.send(fmt.Sprintf("%v Ahoy, %s", Magenta(">"), client.username))

	return client
}

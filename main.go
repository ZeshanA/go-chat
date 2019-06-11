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
	// Don't disregard the error returned by Atoi here
	port, _ := strconv.Atoi(os.Getenv("CHAT_PORT"))

	// You've fetched your port above as a string using os.Getenv, converted it
	// to an integer and then used Itoa to convert it back to a string - you
	// could just use the host + ":" + os.Getenv("CHAT_PORT") and skip all
	// of the conversions
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

		// If you change the ChatConnection struct to embed a net.Conn
		// instead of specifically a net.TCPConn, you can get rid of the
		// type assertion here - I don't think you're using any methods
		// that are specific to TCPConn and not Conn
		conn := ChatConnection{
			connection.(*net.TCPConn),
			&server,
		}

		user := connect(conn, server)
		go user.receive()

		// You've added the client to the users list here and in
		// the second line of connect - I think you only need to do it once?
		// I'd say get rid of this line and leave the one in the connect function
		// so that everything that needs to be done to register a new user is nicely
		// contained in the connect function.
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

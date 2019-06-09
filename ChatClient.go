//revive:disable:dot-imports

package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

const usage = `  /help   Prints available commands and their description
  /join   Joins a channel. Usage: /join <channel>
  /list   Lists all available channels. Usage: /list
  /part   Leaves the current channel. Usage: /part
  /quit   Disconnects from the chat server
`

// ChatClient : A chat client/user
type ChatClient struct {
	channel    *ChatChannel
	connection *ChatConnection
	username   string
}

func logToFile(message string, username string) {
	logDir := os.Getenv("CHAT_LOG_DIR")
	fileName := fmt.Sprintf("%s/go-chat-%s.log", logDir, username)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(message + "\n"); err != nil {
		panic(err)
	}
}

func (user *ChatClient) message(payload ...interface{}) {
	message := payload[0].(string)
	fromUser := payload[1].(*ChatClient)

	if user.username != fromUser.username {
		now := time.Now().Format("15:04:05")
		message = fmt.Sprintf("[%s] %s", now, message)
		user.connection.send(message)
	}

	logToFile(message, user.username)
}

func (user *ChatClient) receive() {
	for {
		message := *user.connection.receive()
		parts := strings.Split(message, " ")
		methodName := strings.Title(parts[0][1:])
		method := reflect.ValueOf(user).MethodByName(methodName)
		zero := reflect.Value{}

		if method == zero {
			if user.channel != nil {
				user.channel.send(fmt.Sprintf("<%s> %s", user.username, message), user)
			} else {
				user.send(fmt.Sprintf("%v Unrecognized command. You must join a channel before you can send a message", Red("Error:")))
			}
			continue
		}

		args := make([]interface{}, len(parts)-1)
		for index, arg := range parts[1:] {
			args[index] = arg
		}
		call(method, args...)

		if message == "/quit" {
			break
		}
	}
}

// Help : Displays usage information
func (user *ChatClient) Help() {
	user.connection.send(usage)
}

// Join : Joins a channel
func (user *ChatClient) Join(name string) {
	channel := user.connection.server.joinChannel(name, user)
	channel.on("message", user.message)
	user.channel = channel
}

// List : Lists all available channels
func (user *ChatClient) List() {
	what := "Active Channels:\n"
	channels := user.connection.server.channels

	if len(channels) == 0 {
		what += "  Aww, there aren't any channels. Create one with /join"
	} else {
		for _, channel := range channels {
			what += "  " + channel.name + "\n"
		}
	}

	user.send(strings.TrimSpace(what))
}

// Part : Leaves the current channel
func (user *ChatClient) Part() {
	if user.channel != nil {
		user.channel.part(user)
		user.send(fmt.Sprintf("%v You've left #%s", Magenta("⊙"), user.channel.name))
		user.channel.off("message", user.message)
		user.channel = nil
	} else {
		user.send(fmt.Sprintf("%v You must join a channel before you can leave it", Red("Error:")))
	}
}

// Quit : Disconnects from the chat server
func (user *ChatClient) Quit() {
	if user.channel != (*ChatChannel)(nil) {
		user.Part()
	}
	user.send(fmt.Sprintf("%v Goodbye", Magenta("⊙")))
	user.connection.Close()
}

func (user *ChatClient) send(what string) {
	user.connection.send(what)
}

// shamelessly lifted from https://stackoverflow.com/a/19721562/416845
func call(method reflect.Value, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	method.Call(inputs)
}

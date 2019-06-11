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
		// Using log.Panic, log.Panicln, log.Panicf etc. instead of just panic(err)
		// is free and prints useful information like the date/time that the panic occurred
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

		// Try to avoid using reflection if possible, it's very, very slow and you lose
		// type safety. I think in this case you could use a map of method name to functions
		// instead:
		//
		// methods := map[string]func(string){
		//   "Join": user.Join,
		// }
		// methodName := "Join"
		// methods[methodName]("myChannel")
		//
		// You'd have to modify the methods of User to all have the same signature
		// but they're pretty close as-is. Would probably have to make them all take a string
		// parameter (because user.Join does), and then you could ignore the param in functions that don't
		// need it. Not a perfect solution either, I agree.
		//
		// I completely see the appeal of using reflect and interface{} coming from a JS
		// background; I don't think it should be the first you way you try to solve a problem in Go,
		// but of course, sometimes the benefit is worth the performance degradation/loss of type safety,
		// which is up to you to determine.
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
	// Maybe use a more descriptive variable name than 'what' ;)
	what := "Active Channels:\n"
	channels := user.connection.server.channels

	// Using += allocates a new string in memory everytime you use it,
	// consider using a StringBuilder: https://www.calhoun.io/concatenating-and-building-strings-in-go/
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
	// You never need to cast nil to a different pointer type before doing a comparison,
	// you can just do: if user.channel != nil
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

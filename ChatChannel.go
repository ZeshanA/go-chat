//revive:disable:dot-imports

package main

import (
	"fmt"

	events "github.com/kataras/go-events"
	. "github.com/logrusorgru/aurora"
)

// ChatChannel : A chat channel
type ChatChannel struct {
	emitter events.EventEmmiter
	name    string
	server  *ChatServer
	users   []*ChatClient
}

func (channel *ChatChannel) init() {
	channel.users = []*ChatClient{}
	channel.emitter = events.New()
}

func (channel *ChatChannel) on(eventName events.EventName, listener events.Listener) {
	channel.emitter.On(eventName, listener)
}

func (channel *ChatChannel) off(eventName events.EventName, listener events.Listener) {
	channel.emitter.RemoveListener(eventName, listener)
}

func (channel *ChatChannel) join(user *ChatClient) {
	channel.users = append(channel.users, user)
	channel.send(fmt.Sprintf("%v %s has joined the channel", Magenta("⊙"), user.username), user)
	user.send(fmt.Sprintf("\n%v Welcome to #%s", Magenta("⊙"), channel.name))
}

func (channel *ChatChannel) part(user *ChatClient) {
	for index := 0; index < len(channel.users); index++ {
		if channel.users[index] == user {
			channel.users = append(channel.users[:index], channel.users[index+1:]...)
			index--
		}
	}
	channel.send(fmt.Sprintf("%v %s has left the channel", Magenta("⊙"), user.username), user)
}

func (channel *ChatChannel) send(what string, user *ChatClient) {
	channel.emitter.Emit("message", what, user)
}

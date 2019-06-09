# go-chat

Go, chat.

## Overview

This was a rather in-depth introduction to the Go language. The code is more verbose and performs more actions than it probably has to, as this is my first foray into the language, but it does demonstrate many concepts of the language and gave me a good excuse to deep-dive some of the core features of Go.

It's worth noting that I prefer to be "thrown into the fire" as it were when it comes to learning new languages. The code in this repo was gleaned from comparison to knowledge in Node.js and searching for equivalent concepts. I did not participate in any bootcamps or courses to assemble the code in this repository. Overall, I really enjoyed getting to know the language.

Inspiration for server commands and UX were taken straight from using IRC back in the day.

## Go Concepts Covered

Though this list may not be complete, these are some concepts that stood out to me:

- Reflection
- Wrapping `struct`
- Stand-alone `struct`
- Struct methods
- Package Scope
- Third-Party library imports
- Pointers and references
- Slices
- Standard Library methods

## App Requirements

- [x] Spin up a server
- [x] Connect via telnet
- [x] Send messages to server from client
- [x] Relay messages to all client
- [x] Messages (from other users) should have timestamp and username prefix
- [x] Messages (from server and other users) logged to file
- [x] Configuration from file (using .env pattern)

## Usage

In a terminal window, run:

```console
$ go run ./pkg
```

In a separate (or multiple) terminal window, run:

```console
$ telnet 0.0.0.0 55555
````

You should see messages from the telnet client, and the server which look something like:

```console
Trying 0.0.0.0...
Connected to 0.0.0.0.
Escape character is '^]'.
🚀 Go Chat v0.1.0

> What shall we call you?
```

Type your desired username, and press return/enter. You'll then be shown the available commands for the server. Go bananas from there.

## Key Struct Notes

### `app`

- contains buuld metadata

Some habits are hard to let go of. I've long been a fan of `package.json` in Node.js and not including some kind of build/package metadata feels wrong.

### `ChatChannel`

- manages messaging the channel

### `ChatClient`

- implements the methods that correspond to available commands
- handles sending data to the client
- listens to the channel
- dispatches server functions

### `ChatConnection`

- handles I/O between server and client

### `ChatServer`

- wraps `net.Listener`
- handles new client connections
- manages channels

## Room for Improvement / Research

- Keep auxiliary Go files in separate directory from `main.go`
- Use reflection (if possible) to detect how many arguments a `ChatClient` method has, and throw an error if the user passes more (or less) arguments than the method requires
- The multiline string for `usage` isn't printing correctly. Must be a Go thing
- Handle an abrupt client disconnect, such as `^]` signal in telnet
- Prevent two users from connecting at the same username
- Learn more about the repercussions and performance issues surrounding pointers

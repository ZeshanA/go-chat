package main

// AppMeta : Metadata for this application
type AppMeta struct {
	// Name : Application name
	Name string
	// Prefix : Application log prefix
	Prefix string
	// Version : Application version
	Version string
}

var app = AppMeta{
	Name:    "Go Chat",
	Prefix:  "go-chat",
	Version: "0.1.0",
}

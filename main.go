package main

import (
	"flag"
	room "web-chat-room/models"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func main() {
	room.Start(*port)
}
package main

import (
	"flag"
	"log"

	"github.com/StreamMeBots/pkg/commands"
	"github.com/StreamMeBots/pkg/tcpclient"
)

// NOTE: for a higher abstraction use the bot package
func main() {
	// command line flags
	publicId := flag.String("publicId", "", "room you want to join to")
	key := flag.String("bot-key", "", "bot key")
	secret := flag.String("bot-secret", "", "bot secret")
	host := flag.String("chat-host", "www.stream.me:2020", "Chat server address")

	flag.Parse()

	// chat client.
	c := tcpclient.New(*host)

	// create room, gives us access to available chat room commands
	room := commands.NewRoom(*publicId)

	// authenticate
	if err := c.Write(room.Pass(*key, *secret), 0); err != nil {
		log.Fatal(err)
	}

	// join room
	if err := c.Write(room.Join(), 0); err != nil {
		log.Fatal(err)
	}

	// Hello, World!
	if err := c.Write(room.Say("Hello, World!"), 0); err != nil {
		log.Println("write error:", err)
	}

	// Leave the chat room
	if err := c.Write(room.Leave(), 0); err != nil {
		log.Println("write error:", err)
	}

	// close the client connection
	c.Close()
}

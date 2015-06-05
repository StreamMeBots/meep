package main

import (
	"flag"
	"log"

	"github.com/StreamMeBots/pkg/bot"
)

func main() {
	// command line flags
	publicId := flag.String("userPublicId", "", "User's public ID of the room you want to join")
	key := flag.String("bot-key", "", "bot key")
	secret := flag.String("bot-secret", "", "bot secret")
	host := flag.String("host", "www.stream.me:2020", "bot chat host")

	flag.Parse()

	// Create a bot
	b, err := bot.New(*host, *key, *secret, *publicId)
	if err != nil {
		log.Fatal(err)
	}
	if err := bot.JoinRoom(); err != nil {
		log.Fatal(err)
	}

	if err := b.Say("Hello, World!"); err != nil {
		log.Println("Say error:", err)
	}

	b.Leave()
}

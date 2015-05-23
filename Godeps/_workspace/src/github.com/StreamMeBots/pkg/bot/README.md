# StreamMe Bot Package

Provides a simple library for creating a StreamMe chat bot. 

## API Documentation

[![GoDoc](https://godoc.org/github.com/StreamMeBots/pkg/bot?status.svg)](https://godoc.org/github.com/StreamMeBots/pkg/bot)

### Hello World Example.

Join a user's chat room and say "Hello, World!"

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
	
		if err := b.Say("Hello, World!"); err != nil {
			log.Println("Say error:", err)
		}
	
		b.Leave()
	}

More examples can be found [here](https://github.com/StreamMeBots/pkg/tree/master/bot/examples).



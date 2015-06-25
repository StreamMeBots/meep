package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/StreamMeBots/pkg/bot"
	"github.com/StreamMeBots/pkg/commands"
)

func main() {
	// command line flags
	publicId := flag.String("publicId", "", "room you want to join to")
	key := flag.String("key", "", "key")
	secret := flag.String("secret", "", "secret")
	host := flag.String("host", "pds.dev.ifi.tv:2020", "bot chat host")

	flag.Parse()

	// Create a bot
	b, err := bot.New(*host, *key, *secret, *publicId)
	if err != nil {
		log.Fatal(err)
	}
	if err := b.JoinRoom(); err != nil {
		log.Fatal(err)
	}

	// handle LEAVE on server shutdown
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			// Wait for a command to come from the chat server
			cmd, err := b.Read()
			if err != nil {
				log.Println("read error:", err)
				continue
			}

			// handle say commands that are not from bots
			switch cmd.Name {
			case commands.LSay:
				if cmd.Get("bot") == "false" {
					if err := b.Say(cmd.Get("message")); err != nil {
						log.Println("write error:", err)
					}
				}
			}
		}
	}()

	<-done
	fmt.Println("Leaving chat room")
	b.Leave()
}

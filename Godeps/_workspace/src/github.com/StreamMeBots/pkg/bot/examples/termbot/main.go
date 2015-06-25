/*
* termbot demonstrates a trivial command line chat client that lets you chat as your bot.
 */
package main

import (
	"bufio"
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
	log.SetFlags(log.Lshortfile)

	// command line flags
	publicId := flag.String("publicId", "", "room you want to join to")
	key := flag.String("key", "", "Bot key")
	secret := flag.String("secret", "", "Bot secret")
	host := flag.String("host", "", "bot host server")

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

	// say loop
	go botSay(b)
	// read loop
	go botRead(b)

	<-done
	fmt.Println("Leaving chat room")
	b.Leave()
}

// botRead prints JOIN, LEAVE and SAY commands to the terminal from the chat server
func botRead(b *bot.Bot) chan *commands.Command {
	for {
		cmd, err := b.Read()
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		fmt.Println("command", cmd.Name, cmd.Args)
		switch cmd.Name {
		case commands.LSay:
			if cmd.Get("bot") == "true" {
				continue
			}
			fmt.Printf("%s: %s\n\n", cmd.Get("username"), cmd.Get("message"))
		case commands.LJoin:
			fmt.Printf("%s has joined chat\n\n", cmd.Get("username"))
		case commands.LLeave:
			fmt.Printf("%s has left chat\n\n", cmd.Get("username"))
		}
	}
}

// botSay sends SAY commands from text entered from the terminal.
func botSay(b *bot.Bot) {
	for {
		l, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println("Error reading stdin:", err)
		}
		b.Say(l[:len(l)-1])
	}
}

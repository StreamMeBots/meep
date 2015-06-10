package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/StreamMeBots/meep/pkg/buckets"
	"github.com/StreamMeBots/meep/pkg/command"
	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/boltdb/bolt"

	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/StreamMeBots/meep/pkg/greetings"
	"github.com/StreamMeBots/meep/pkg/stats"
	pkgBot "github.com/StreamMeBots/pkg/bot"
	"github.com/StreamMeBots/pkg/commands"
)

// Errors
var (
	ErrBotAlreadyStarted = errors.New("Bot is already running")
	ErrAuthNon200        = errors.New("Unable to authorize bot")
)

// NewBots is the constructor for Bots
func NewBots() Bots {
	return Bots{
		bots: map[string]Bot{},
	}
}

// Bots is used to safely control access to all running bots
type Bots struct {
	bots map[string]Bot
	sync.RWMutex
	commandTimers []*command.Command
}

// Start starts a user's bot
func (bs *Bots) Start(userPublicId string, client *http.Client) error {
	bs.RLock()
	if _, ok := bs.bots[userPublicId]; ok {
		return ErrBotAlreadyStarted
	}
	bs.RUnlock()

	bt, err := NewBot(userPublicId, client)
	if err != nil {
		return err
	}

	bs.Lock()
	defer bs.Unlock()
	bs.bots[userPublicId] = bt
	return nil
}

// TODO(james): finish
func (b *Bot) startCommandTimers() {
	/*
		var err error
		b.commandTimers, err = command.GetCommandsWithTimers(b.UserPublicId)
		if err != nil {

		}
	*/
}

func (bs *Bots) Info(userPublicId string) pkgBot.Info {
	bs.RLock()
	defer bs.RUnlock()

	b, ok := bs.bots[userPublicId]
	if !ok {
		return pkgBot.Info{State: "notStarted"}
	}

	return b.bot.GetInfo()
}

// NewBot is the constructor for Bot
//	- creates a bucket in the bolt db
//	- starts the bots read loop
//	- authorizes the bot
func NewBot(userPublicId string, client *http.Client) (Bot, error) {
	bt := Bot{
		UserPublicId: userPublicId,
		stop:         make(chan struct{}),
		client:       client,
	}

	conf := []pkgBot.Config{}
	if config.Conf.Debug {
		conf = append(conf, pkgBot.LogCommands)
	}

	var err error
	bt.bot, err = pkgBot.New(config.Conf.ChatHost, config.Conf.BotKey, config.Conf.BotSecret, userPublicId, conf...)
	if err != nil {
		return bt, err
	}

	// auth bot with user's chat room
	if err := bt.auth(); err != nil {
		return bt, err
	}

	if err := bt.bot.JoinRoom(); err != nil {
		return bt, err
	}

	go bt.read()

	return bt, nil
}

func (bs *Bots) Startup() {

}

// Close stops all bots and saves the user's public id so the bots can be restarted on startup
func (bs *Bots) Close() {
	bs.Lock()
	defer bs.Unlock()

	db.DB.Update(func(tx *bolt.Tx) error {
		bkt := buckets.RunningBots(tx)
		for id, b := range bs.bots {
			close(b.stop)
			delete(bs.bots, id)

			if err := bkt.Put([]byte(id), []byte(id)); err != nil {
				log.Println("msg='error-saving-running-bot-id', id='%s' error='%v'", id, err)
				continue
			}
		}
		return nil
	})
}

// Stop stops a user's bot
func (bs *Bots) Stop(userPublicId string) {
	bs.Lock()
	b, ok := bs.bots[userPublicId]
	delete(bs.bots, userPublicId)
	bs.Unlock()

	if ok {
		close(b.stop)
	}
}

// LogStream returns a channel that can be used to listen for events
func (bs *Bots) LogStream(userPublicId string) (chan interface{}, error) {
	bs.RLock()
	defer bs.RUnlock()

	b, ok := bs.bots[userPublicId]
	if !ok {
		return nil, fmt.Errorf("Bot is not running")
	}

	c := b.bot.Subscribe(userPublicId)
	return c, nil
}

func (bs *Bots) CloseLogStream(userPublicId string) {
	bs.RLock()
	defer bs.RUnlock()

	b, ok := bs.bots[userPublicId]
	if !ok {
		return
	}

	b.bot.Unsubscribe(userPublicId)
}

// Bot represents a bot that is associated to a stream.me user
type Bot struct {
	UserPublicId string
	bot          *pkgBot.Bot
	stop         chan struct{}
	client       *http.Client
}

func (b *Bot) bucketKey() []byte {
	return []byte(b.UserPublicId)
}

// read is responsible for reading commands from the chat room then routing the commands to a bot method
func (b *Bot) read() {
	for {
		// check if we need to close down
		select {
		case <-b.stop:
			b.bot.Leave()
			return
		default:
		}

		// read chat command
		cmd, err := b.bot.Read()
		if err != nil {
			continue
		}

		// route
		switch cmd.Name {
		case commands.LJoin:
			b.join(cmd)
		case commands.LSay:
			b.say(cmd)
		}
	}
}

// auth authorizes the bot with the user's chat room
func (b *Bot) auth() error {
	url := fmt.Sprintf(
		// /v1/rooms/:roomPublicId/authorized-bots/:botId
		config.Conf.Url+"/api-chat/v1/rooms/%s/authorized-bots/%s",
		b.bot.RoomId(),
		b.bot.Key,
	)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Printf("msg='error-creating-request', error='%v'\n", err)
		return err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		log.Printf("msg='request-error', error='%v'\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("msg='error-reading-response-body', error='%v'\n", err)
		return err
	}

	log.Printf("auth-error-body='%s' statusCode='%v'\n", string(bd), resp.StatusCode)

	return ErrAuthNon200
}

func (b *Bot) say(cmd *commands.Command) {
	isCommand := false
	defer func() {
		if !isCommand {
			stats.Line(b.bucketKey())
		}
	}()

	m := cmd.Get("message")
	if len(m) > 2 && m[0] == '!' {
		c, err := command.Get(b.bucketKey(), m)
		if err != nil {
			return
		}

		if stats.Command(b.bucketKey(), c) {
			if msg := c.Parse(cmd); len(msg) > 0 {
				b.bot.Say(msg)
				isCommand = true
			}
		}
	}
}

func (b *Bot) join(cmd *commands.Command) {
	// noop for bots
	if bot := cmd.Get("bot"); bot == "true" {
		return
	}
	// noop for owners
	/*
		if bot := cmd.Get("role") == "owner" {
			return
		}
	*/

	e := greetings.Join(b.bucketKey(), cmd)
	if len(e.Response) > 0 {
		if e.Private {
			// TODO: meep command only
		} else {
			if e.Type == "answeringMachine" {
				if stats.Command(b.bucketKey(), &command.Command{Name: "answeringMachine"}) {
					fmt.Println("response:", e.Response)
					b.bot.Say(e.Response)
				}
			} else {
				b.bot.Say(e.Response)
			}
		}
	}
}

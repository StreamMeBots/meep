package bot

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/StreamMeBots/meep/pkg/greetings"
	"github.com/StreamMeBots/meep/pkg/stats"
	"github.com/StreamMeBots/meep/pkg/user"
	pkgBot "github.com/StreamMeBots/pkg/bot"
	"github.com/StreamMeBots/pkg/commands"
	"github.com/boltdb/bolt"
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

// Bots is used to safely contorl access to all running bots
type Bots struct {
	bots map[string]Bot
	sync.RWMutex
}

// Start starts a user's bot
func (bs *Bots) Start(userPublicId string, client *http.Client) error {
	bs.RLock()
	if _, ok := bs.bots[userPublicId]; ok {
		return ErrBotAlreadyStarted
	}
	bs.RUnlock()

	bt, err := NewBot(userPublicId)
	if err != nil {
		return err
	}

	bs.Lock()
	defer bs.Unlock()
	bs.bots[userPublicId] = bt
	return nil
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
func NewBot(userPublicId string) (Bot, error) {
	var bt Bot
	b, err := pkgBot.New(config.Conf.ChatHost, config.Conf.BotKey, config.Conf.BotSecret, userPublicId)
	if err != nil {
		return bt, err
	}

	bt = Bot{
		UserPublicId: userPublicId,
		bot:          b,
		stop:         make(chan struct{}),
	}

	// auth bot with user's chat room
	/* TODO: blocked by node work
	if err := bt.auth(); err != nil {
		return bt, err
	}
	*/

	// create buckets for bot
	err = db.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(bt.bucketKey())
		if err != nil {
			return err
		}

		_, err = bkt.CreateBucketIfNotExists(greetings.GreetingsKeyName)
		if err != nil {
			return err
		}

		_, err = bkt.CreateBucketIfNotExists(stats.KeyName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return bt, err
	}

	go bt.read()

	return bt, nil
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
	return []byte(`bot:` + b.UserPublicId)
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

		b.stat(cmd)

		// route
		switch cmd.Name {
		case commands.LJoin:
			b.join(cmd)
		}
	}
}

// auth authorizes the bot with the user's chat room
func (b *Bot) auth() error {
	url := fmt.Sprintf(
		config.Conf.Url+"/api-chat/v1/users/%s/authorized-bots/%s",
		b.UserPublicId,
		b.bot.RoomId(),
	)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return ErrAuthNon200
}

func (b *Bot) stat(cmd *commands.Command) {

}

func (b *Bot) join(cmd *commands.Command) {
	// noop for bots
	if bot := cmd.Get("bot"); bot == "true" {
		return
	}

	msg := greetings.Join(user.BucketName(b.UserPublicId), b.bucketKey(), cmd)
	if len(msg) > 0 {
		b.bot.Say(msg)
	}
}

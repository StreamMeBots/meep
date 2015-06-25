/*
Package bot provides a simple library for creating a StreamMe chat bot
*/
package bot

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/StreamMeBots/pkg/commands"
	"github.com/StreamMeBots/pkg/tcpclient"
)

// State represents the state of the ot
type State string

// Connection States
const (
	Disconnected State = "Disconnected" // Bot is disconnected
	Connecting   State = "Connecting"   // Bot is connecting
	Connected    State = "Connected"    // Bot is connected
	Joined       State = "Joined"       // Bot joined a chat room
)

// Event types that can be sent down the subscribe channel
type (
	EventReadError   error
	EventReadCommand commands.Command
	EventWrite       string
	EventWriteError  error
	EventStateChange State
)

// Bot represents a bot user to a stream.me chat server.
type Bot struct {
	Room        commands.Room
	Key         string
	Secret      string
	tcp         *tcpclient.Client
	state       State
	mx          sync.RWMutex
	subs        map[string]chan interface{}
	started     time.Time
	logCommands bool
}

// RoomId returns the chat room ID that the bot uses to connect to the chat room.
func (b *Bot) RoomId() string {
	return string(b.Room)
}

// Info represents stats and state about a bot.
type Info struct {
	State   State     // Current state of the bot
	Started time.Time // When the bot was started
}

// Config is used to configure the bot
type Config func(*Bot)

// LogCommands is used to turn on tcpclient's command logging
var LogCommands = func(b *Bot) {
	b.logCommands = true
}

// New is the constructor for Bot. The Bot will connect to the chat server
func New(host, key, secret, userPublicId string, confs ...Config) (*Bot, error) {
	b := &Bot{
		Room:    commands.NewRoom(userPublicId),
		Key:     key,
		Secret:  secret,
		started: time.Now().UTC(),
		state:   Disconnected,
		subs:    make(map[string]chan interface{}),
	}
	for _, conf := range confs {
		conf(b)
	}
	if b.logCommands {
		b.tcp = tcpclient.New(host, tcpclient.LogCommands)
	} else {
		b.tcp = tcpclient.New(host)
	}

	if err := b.isOnline(); err != nil {
		return nil, err
	}

	return b, nil
}

// JoinRoom joins the room
func (b *Bot) JoinRoom() error {
	if err := b.Pass(); err != nil {
		return err
	}

	if err := b.Join(); err != nil {
		return err
	}

	return nil
}

// GetInfo gets state info about the Bot
func (b *Bot) GetInfo() Info {
	b.mx.Lock()
	defer b.mx.Unlock()
	return Info{
		State:   b.state,
		Started: b.started,
	}
}

// Subscribe returns a channel where all bot activity can be monitored
func (b *Bot) Subscribe(id string) chan interface{} {
	b.mx.Lock()
	defer b.mx.Unlock()
	c := make(chan interface{}, 10)
	b.subs[id] = c
	return c
}

// Unsubscribe removes the channel associated to the id from the bots subscriber list
func (b *Bot) Unsubscribe(id string) {
	b.mx.Lock()
	defer b.mx.Unlock()
	if c, ok := b.subs[id]; ok {
		close(c)
	}
	delete(b.subs, id)
}

// Read reads one command from the chat server
func (b *Bot) Read() (*commands.Command, error) {
	return b.ReadTimeout(0)
}

// ReadTimeout is the same as Read but adds a timeout.
func (b *Bot) ReadTimeout(d time.Duration) (*commands.Command, error) {
	cmd, err := b.tcp.Read(d)

	// emit
	if err != nil {
		b.emit(EventReadError(err))
	} else {
		b.emit(EventReadCommand(*cmd))
	}

	if err == io.EOF {
		b.setState(Connecting)
		// if we disconnected we need to re-auth and re-join the room
		if err := b.Pass(); err != nil {
			return nil, err
		}
		if err := b.Join(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// Pass authenticates the bot to the their chat room
func (b *Bot) Pass() error {
	cmd, err := b.retryTimeout(b.Room.Pass(b.Key, b.Secret), commands.LPass, 5, time.Second*5)
	if err != nil {
		return err
	}

	if cmd.Get("result") != "success" {
		return fmt.Errorf("Auth failure. Check your creds")
	}

	return nil
}

// Join joins the bot to their chat room
func (b *Bot) Join() error {
	_, err := b.retryTimeout(b.Room.Join(), commands.LJoin, 5, time.Second*5)
	if err != nil {
		return err
	}
	b.setState(Joined)
	return nil
}

// Say sends the SAY command to the chat room
func (b *Bot) Say(msg string) error {
	return b.write(b.Room.Say(msg))
}

// Kick is used to kick a user from chat
func (b *Bot) Kick(userPublicId string) error {
	return b.write(b.Room.Kick(userPublicId))
}

// Ban is used to Ban a user from chat
func (b *Bot) Ban(userPublicId string) error {
	return b.write(b.Room.Ban(userPublicId))
}

// Mod is used to mod a user in chat
func (b *Bot) Mod(userPublicId string) error {
	return b.write(b.Room.Mod(userPublicId))
}

// MuteGuest is used to mute a user with a role of "guest"
func (b *Bot) MuteGuest(userPublicId string) error {
	return b.write(b.Room.MuteGuest(userPublicId))
}

// Mute is used to mute a user
func (b *Bot) Mute(userPublicId string) error {
	return b.write(b.Room.Mute(userPublicId))
}

// Unmute unmutes a previously muted user
func (b *Bot) UnMute(userPublicId string) error {
	return b.write(b.Room.UnMute(userPublicId))
}

// Erase a message
func (b *Bot) Erase(messageId string) error {
	return b.write(b.Room.Erase(messageId))
}

// Leave sends the Leave command and disconnects the bot from the chat server
func (b *Bot) Leave() {
	b.write(b.Room.Leave())
	b.tcp.Close()
	b.setState(Disconnected)
}

func (b *Bot) setState(s State) {
	b.mx.Lock()
	b.state = s
	b.mx.Unlock()
	b.emit(EventStateChange(s))
}

func (b *Bot) write(msg string) error {
	b.emit(EventWrite(msg))
	err := b.tcp.Write(msg, time.Second*5)
	if err != nil {
		b.emit(EventWriteError(err))
	}
	return err
}

func (b *Bot) retryTimeout(command string, checkCommand string, retryCount int, timeout time.Duration) (*commands.Command, error) {
	// FIXME: update to retry however many times for a given timeout.
	// FIXME: pass in a callback instead of a checkCommand string, will allow for more fine grained control

	if err := b.tcp.Write(command, timeout); err != nil {
		return nil, err
	}

	cmd := &commands.Command{}
	var err error
	for i := 0; i < retryCount; i++ {
		cmd, err = b.ReadTimeout(timeout)
		if err != nil {
			log.Printf("%s read timeout error: %v. Trying again...\n", checkCommand, err)
			continue
		}
		if cmd.Name != checkCommand {
			log.Printf("%s: wrong command: %s Args: %+v. Trying again...\n", checkCommand, cmd.Name, cmd.Args)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("Never heard back from %s after %v retries", checkCommand, retryCount)
	}

	return cmd, nil
}

// isOnline is a helper function that blocks for up to 10 seconds waiting for the tcp connection to occur
func (b *Bot) isOnline() error {
	b.setState(Connecting)
	for i := 0; i < 10; i++ {
		if b.tcp.Stats().Online {
			b.setState(Connected)
			return nil
		}
		time.Sleep(time.Second)
	}
	b.tcp.Close()
	b.setState(Disconnected)
	return fmt.Errorf("Unable to connect to the chat server")
}

func (b *Bot) emit(i interface{}) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	for id, ch := range b.subs {
		select {
		case ch <- i:
		default:
			// cleanup subscribers who are not listening
			b.Unsubscribe(id)
		}
	}
}

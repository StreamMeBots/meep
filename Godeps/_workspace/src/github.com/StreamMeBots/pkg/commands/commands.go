/*
* Package commands is used to format TCP messages for stream.me chat
 */
package commands

import "fmt"

// Listen Commands. The commands that can come back from the stream.me chat server
const (
	LJoin  = "JOIN"
	LSay   = "SAY"
	LLeave = "LEAVE"
	LError = "ERROR"
	LPass  = "PASS"
)

// Command represents a command from a stream.me chat server
type Command struct {
	Name string
	Args map[string]string
}

// Get is a helper function to get an argument, if the key is not found an empty string is returned
func (c *Command) Get(key string) string {
	if c == nil {
		return ""
	}
	if a, ok := c.Args[key]; ok {
		return a
	}
	return ""
}

// Say formats the command to send a message to the chat room
func Say(msg string) string {
	return "SAY " + msg
}

// Room defines the chat actions that can be taken by the owner/moderator/admin of the chat room.
//
// NOTE: use NewRoom to initialize
type Room string

// NewRoom is the constructor for Room. NewRoom should always be used to create a new Room type.
func NewRoom(userPublicId string) Room {
	return Room(fmt.Sprintf("user:%s:web", userPublicId))
}

// Pass formats the command used to authenticate with stream.me
func (r Room) Pass(key string, secret string) string {
	//return fmt.Sprintf("PASS %s %s %s", key, secret, r)
	return fmt.Sprintf("PASS %s %s", key, secret)
}

// Say formats the command used to chat
func (r Room) Say(msg string) string {
	return Say(msg)
}

// Join formats the command used to join a stream.me chat room
func (r Room) Join() string {
	return "JOIN " + string(r)
}

// Kick formats the command a user with a role of moderator/owner/admin can send to kick a user from a chat room
func (r Room) Kick(userPublicId string) string {
	return fmt.Sprintf("KICK %s %s", r, userPublicId)
}

// Ban formats the command a a user with a role of moderator/owner/admin can send to ban a user from a chat room
func (r Room) Ban(userPublicId string) string {
	return r.ChangeRole(userPublicId, "banned")
}

// Unban is an alias for ChangeRole that changes a user's role to 'user'
func (r Room) Unban(userPublicId string) string {
	return r.ChangeRole(userPublicId, "user")
}

// Mod is an alias for ChangeRole that changes a user's role to 'moderator'
func (r Room) Mod(userPublicId string) string {
	return r.ChangeRole(userPublicId, "moderator")
}

// MuteGuest is an alias for ChangeRole that changes a user's role to 'mutedGuest'
func (r Room) MuteGuest(userPublicId string) string {
	return r.ChangeRole(userPublicId, "mutedGuest")
}

// Mute is an alias for ChangeRole that changes a user's role to 'mute'
func (r Room) Mute(userPublicId string) string {
	return r.ChangeRole(userPublicId, "mute")
}

// UnMute is an alias for ChangeRole that changes a user's role to 'user'
func (r Room) UnMute(userPublicId string) string {
	return r.ChangeRole(userPublicId, "user")
}

// Leave formats the LEAVE command. This command should be used when you want to leave a chat room
func (r Room) Leave() string {
	return fmt.Sprintf("LEAVE %s", r)
}

func (r Room) Erase(messageId string) string {
	return fmt.Sprintf("ERASE %s", messageId)
}

// ChangeRole formats the command a user with a role of moderator/owner/admin can send to change a user's role
func (r Room) ChangeRole(userPublicId string, role string) string {
	return fmt.Sprintf("CHANGEROLE %s %s %s", userPublicId, r, role)
}

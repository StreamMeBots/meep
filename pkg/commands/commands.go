package command

import "github.com/StreamMeBots/pkg/commands"

type Command struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

// Save saves the command
func (c *Command) Save(userBucket []byte) error {
	return nil
}

// Get gets a single command
func Get(userBucket []byte, name string) (*Command, error) {
	return nil, nil

}

// GetAll gets all of a user's commands
func GetAll(userBucket []byte) ([]*Command, error) {
	return nil, nil
}

// Say checks if the message is a command and if it is provies and answer to the command
func Say(userBucket []byte, cmd *commands.Command) string {
	return ""
}

package command

import (
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/StreamMeBots/pkg/commands"
	"github.com/boltdb/bolt"
)

// BucketName bucket name in bolt
var BucketName = []byte(`commands`)

// Command represents a command response template.
type Command struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

// Validate validates the Command
func (c *Command) Validate() error {
	if len(c.Name) == 0 || len(c.Name) > 100 {
		return fmt.Errorf("Command Name should be between 1 and 100 characters")
	}
	if len(c.Template) == 0 || len(c.Template) > 500 {
		return fmt.Errorf("Command Template should be between 1 and 500 characters")
	}

	if _, err := template.New("foo").Parse(c.Template); err != nil {
		return fmt.Errorf("Error parsing Template: %v", err)
	}

	return nil
}

// Save saves the command
func (c *Command) Save(userBucket []byte) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}

		bkt, err := ubkt.CreateBucketIfNotExists(BucketName)
		if err != nil {
			return err
		}

		b, err := json.Marshal(c)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(c.Name), b)
	})

	if err != nil {
		log.Println("msg='error-saving-command', error='%v', userBucket='%s'", err, string(userBucket))
		return err
	}

	return nil
}

// Get gets a single command
func Get(userBucket []byte, name string) (*Command, error) {
	var cmd *Command
	err := db.DB.View(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}

		bkt, err := ubkt.CreateBucketIfNotExists(BucketName)
		if err != nil {
			return err
		}

		b := bkt.Get([]byte(name))
		if b == nil {
			return nil
		}

		if err := json.Unmarshal(b, &cmd); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("msg='error-reading-command', error='%v', userBucket='%s'", err, string(userBucket))
		return nil, err
	}

	return cmd, nil
}

// GetAll gets all of a user's commands
func GetAll(userBucket []byte, name string) ([]*Command, error) {
	cmds := []*Command{}
	/*
		err := db.DB.View(func(tx *bolt.Tx) error {
			ubkt, err := tx.CreateBucketIfNotExists(userBucket)
			if err != nil {
				return err
			}

			bkt, err := ubkt.CreateBucketIfNotExists(BucketName)
			if err != nil {
				return err
			}

			b := bkt.Get(name)
			if b == nil {
				return nil
			}

			if err := json.Unmarshal(b, &cmd); err != nil {
				return nil, err
			}

			return nil
		})

		if err != nil {
			log.Println("msg='error-reading-command', error='%v', userBucket='%s'", err, string(userBucket))
			return nil, err
		}
	*/

	return cmds, nil
}

// Say checks if the message is a command and if it is provies and answer to the command
func Say(userBucket []byte, cmd *commands.Command) string {
	return ""
}

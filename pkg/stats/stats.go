package stats

import (
	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/boltdb/bolt"
)

var KeyName = []byte("stats")

func Line(botBucket []byte) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		return nil
	})

	if err != nil {

	}

	return nil
}

func LinesSinceLastCommand() int {
	return 0
}

func Command() {

}

/*
func StatSay(botBucket []byte, cmd *commands.Command) {
	ts := now.BeginningOfDay().Format(time.RFC3339)
}

func StatGreeting(botBucket []byte, cmd *commands.Command) {

}

func StatJoin(botBucket []byte, cmd *commands.Command) {

}*/

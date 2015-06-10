/*
* Package buckets defines the bolt buckets and keys used
 */
package buckets

import (
	"bytes"
	"log"
	"strconv"

	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/boltdb/bolt"
)

// buckets
var (
	userData              = []byte(`user.data`)
	userGreetingTemplates = []byte(`user.greetings.templates`)
	runningBots           = []byte(`bots.running`)

	// partial
	botGreetings            = []byte(`bot.greetings:`)
	botStatsLinesPerHour    = []byte(`bot.stats.lines.perhour:`)
	botStatsLinesPerDay     = []byte(`bot.stats.lines.perday:`)
	botStatsCommandsPerHour = []byte(`bot.stats.commands.perhour:`)
	botStatsCommandsPerDay  = []byte(`bot.stats.commands.perday:`)
	botStatsLastCommand     = []byte(`bot.stats.commands.last:`)

	userCommands = []byte(`user.commands:`)
)

// Bucket wraps the bolt bucket - future proofing
type Bucket struct {
	*bolt.Bucket
}

func Init() {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(userData); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(userGreetingTemplates); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(runningBots); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// LinesPerHour
func LinesPerHour(tx *bolt.Tx, botUserPublicId []byte) (Bucket, error) {
	return createBucket(tx, createKey(botStatsLinesPerHour, botUserPublicId))
}

// LinesPerDay
func LinesPerDay(tx *bolt.Tx, botUserPublicId []byte) (Bucket, error) {
	return createBucket(tx, createKey(botStatsLinesPerDay, botUserPublicId))
}

// CommandsPerDay
func CommandsPerDay(tx *bolt.Tx, botUserPublicId, command []byte) (Bucket, error) {
	return createBucket(tx, createKey(botStatsCommandsPerDay, botUserPublicId, command))
}

// CommandsPerHour
func CommandsPerHour(tx *bolt.Tx, botUserPublicId, command []byte) (Bucket, error) {
	return createBucket(tx, createKey(botStatsCommandsPerHour, botUserPublicId, command))
}

func LastCommand(tx *bolt.Tx, botUserPublicId []byte) (Bucket, error) {
	return createBucket(tx, createKey(botStatsLastCommand, botUserPublicId))
}

func BotGreetings(tx *bolt.Tx, botUserPublicId []byte) (Bucket, error) {
	return createBucket(tx, createKey(botGreetings, botUserPublicId))
}

func UserCommands(userBucket []byte, tx *bolt.Tx) (Bucket, error) {
	return createBucket(tx, createKey(userCommands, userBucket))
}

func UserData(tx *bolt.Tx) Bucket {
	return Bucket{tx.Bucket(userData)}
}

func UserGreetingTemplates(tx *bolt.Tx) Bucket {
	return Bucket{tx.Bucket(userGreetingTemplates)}
}

func RunningBots(tx *bolt.Tx) Bucket {
	return Bucket{tx.Bucket(runningBots)}
}

// createKey is a helper function to join multiple slices with ':'
func createKey(keys ...[]byte) []byte {
	return bytes.Join(keys, []byte(`:`))
}

// createBucket is a helper function for creating a Bucket
func createBucket(tx *bolt.Tx, key []byte) (Bucket, error) {
	bkt, err := tx.CreateBucketIfNotExists(key)
	if err != nil {
		return Bucket{}, err
	}
	return Bucket{bkt}, nil
}

// Decr decrements
func Decr(bkt *bolt.Bucket, key []byte) (int64, error) {
	c, err := GetInt64(bkt, key)
	if err != nil {
		return 0, err
	}

	v := c - 1
	if err := SetInt64(bkt, key, v); err != nil {
		return 0, nil
	}
	return v, nil
}

// Incr increments
func Incr(bkt *bolt.Bucket, key []byte) (int64, error) {
	c, err := GetInt64(bkt, key)
	if err != nil {
		return 0, err
	}

	v := c + 1
	if err := SetInt64(bkt, key, v); err != nil {
		return 0, err
	}
	return v, nil
}

// GetInt64 gets and int64
func GetInt64(bkt *bolt.Bucket, key []byte) (int64, error) {
	b := bkt.Get(key)
	if b == nil {
		return 0, nil
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		log.Println("msg='error-converting-bytes-to-int', error='%v', value='%s'", err, string(b))
		return 0, err
	}

	return i, nil
}

// SetInt64 sets and int64 value
func SetInt64(bkt *bolt.Bucket, key []byte, count int64) error {
	n := []byte(strconv.FormatInt(count, 10))
	if err := bkt.Put(key, n); err != nil {
		log.Println("msg='error-putting-count', error='%v', count='%v', bucket='%v'", err, count)
		return err
	}
	return nil
}

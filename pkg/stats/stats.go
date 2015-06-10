/*
* say commands per hour
* say commands per day
 */
package stats

import (
	"fmt"
	"log"
	"time"

	"github.com/StreamMeBots/meep/pkg/buckets"
	"github.com/StreamMeBots/meep/pkg/command"
	"github.com/StreamMeBots/meep/pkg/db"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/now"
)

// CommandThrottle number of lines between a command can be displayed
var CommandThrottle int64 = 10

// Line keep track of line stats
func x(tx *bolt.Tx, publicId []byte) error {
	return nil
}

// Line write line stats, aky SAY commands
func Line(userPublicId []byte) {
	db.DB.Update(func(tx *bolt.Tx) error {
		day := []byte(now.BeginningOfDay().Format(time.RFC3339))
		bkt, err := buckets.LinesPerDay(tx, userPublicId)
		if err != nil {
			return fmt.Errorf("msg='error-getting-lines-per-day-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
		}
		if _, err := buckets.Incr(bkt.Bucket, day); err != nil {
			return err
		}
		return nil
	})

	db.DB.Update(func(tx *bolt.Tx) error {
		hour := []byte(now.BeginningOfHour().Format(time.RFC3339))
		bkt, err := buckets.LinesPerHour(tx, userPublicId)
		if err != nil {
			return fmt.Errorf("msg='error-getting-lines-per-hour-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
		}
		if _, err := buckets.Incr(bkt.Bucket, hour); err != nil {
			return err
		}
		return nil
	})

	cmds, err := command.GetAll(userPublicId)
	if err != nil {
		log.Printf("msg='error-getting-commands', error='%v' userPublicId='%s'\n", err, userPublicId)
		return
	}

	cmds = append(cmds, &command.Command{
		Name: "answeringMachine", // this is a bit of a hack...
	})
	for _, cmd := range cmds {
		db.DB.Update(func(tx *bolt.Tx) error {
			bkt, err := buckets.LastCommand(tx, userPublicId)
			if err != nil {
				return fmt.Errorf("msg='error-getting-last-commands-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
			}
			if _, err := buckets.Incr(bkt.Bucket, []byte(cmd.Name)); err != nil {
				return err
			}
			return nil
		})
	}
}

// Command checks if a command should be written and writes command stats
func Command(userPublicId, command []byte) (ok bool) {
	db.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := buckets.LastCommand(tx, userPublicId)
		if err != nil {
			return err
		}
		c, err := buckets.GetInt64(bkt.Bucket, command)
		if err == buckets.ErrIntNotSet {
			// edge case
			c = CommandThrottle + 1
		} else if err != nil {
			return err
		}
		if c > CommandThrottle {
			ok = true
			// reset
			return buckets.SetInt64(bkt.Bucket, command, 0)
		}
		return nil
	})

	if !ok {
		return false
	}

	db.DB.Update(func(tx *bolt.Tx) error {
		day := []byte(now.BeginningOfDay().Format(time.RFC3339))
		bkt, err := buckets.CommandsPerDay(tx, userPublicId, command)
		if err != nil {
			log.Printf("msg='error-getting-commands-per-day-bucket', error='%v', userPublicId='%v', command='%s'\n", err, userPublicId, string(command))
			return err
		}
		_, err = buckets.Incr(bkt.Bucket, day)
		return err
	})
	db.DB.Update(func(tx *bolt.Tx) error {
		hour := []byte(now.BeginningOfHour().Format(time.RFC3339))
		bkt, err := buckets.CommandsPerHour(tx, userPublicId, command)
		if err != nil {
			log.Printf("msg='error-getting-commands-per-hour-bucket', error='%v', userPublicId='%v', command='%s'\n", err, userPublicId, string(command))
			return err
		}
		_, err = buckets.Incr(bkt.Bucket, hour)
		return err
	})

	return true
}

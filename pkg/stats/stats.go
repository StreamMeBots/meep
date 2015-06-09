/*
* say commands per hour
* say commands per day
 */
package stats

import (
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
		hour := []byte(now.BeginningOfHour().Format(time.RFC3339))

		if bkt, err := buckets.LinesPerDay(tx, userPublicId); err != nil {
			log.Printf("msg='error-getting-lines-per-day-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
		} else {
			buckets.Incr(bkt.Bucket, day)
		}

		if bkt, err := buckets.LinesPerHour(tx, userPublicId); err != nil {
			log.Printf("msg='error-getting-lines-per-hour-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
		} else {
			buckets.Incr(bkt.Bucket, hour)
		}

		if cmds, err := command.GetAll(userPublicId); err != nil {
			log.Printf("msg='error-getting-commands', error='%v' userPublicId='%s'\n", err, userPublicId)
		} else {
			cmds = append(cmds, &command.Command{
				Name: "answeringMachine", // this is a bit of a hack...
			})
			for _, cmd := range cmds {
				if bkt, err := buckets.LastCommand(tx, userPublicId); err != nil {
					log.Printf("msg='error-getting-last-commands-bucket', error='%v', userPublicId='%v'\n", err, userPublicId)
				} else {
					buckets.Incr(bkt.Bucket, []byte(cmd.Name))
				}
			}
		}
		return nil
	})
}

// Command checks if a command should be written and writes command stats
func Command(userPublicId, command []byte) (ok bool) {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := buckets.LastCommand(tx, userPublicId)
		if err != nil {
			return err
		}
		c, err := buckets.GetInt64(bkt.Bucket, command)
		if err != nil {
			return err
		}
		if c < CommandThrottle {
			return nil
		}

		day := []byte(now.BeginningOfDay().Format(time.RFC3339))
		hour := []byte(now.BeginningOfHour().Format(time.RFC3339))

		if bkt, err := buckets.CommandsPerDay(tx, userPublicId, command); err != nil {
			log.Printf("msg='error-getting-commands-per-day-bucket', error='%v', userPublicId='%v', command='%s'\n", err, userPublicId, string(command))
		} else {
			buckets.Incr(bkt.Bucket, day)
		}

		if bkt, err := buckets.CommandsPerHour(tx, userPublicId, command); err != nil {
			log.Printf("msg='error-getting-commands-per-hour-bucket', error='%v', userPublicId='%v', command='%s'\n", err, userPublicId, string(command))
		} else {
			buckets.Incr(bkt.Bucket, hour)
		}

		// reset
		return buckets.SetInt64(bkt.Bucket, command, 0)
	})
	if err != nil {
		log.Println("msg='stat-command-error', error='%v', command='%s'\n", err, string(command))
	}

	return ok
}

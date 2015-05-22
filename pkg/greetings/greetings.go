package greetings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/StreamMeBots/meep/pkg/db"
	"github.com/StreamMeBots/pkg/commands"
	"github.com/boltdb/bolt"
	"github.com/jinzhu/now"
)

var TemplateKeyName = []byte("grettingsTemplate")

var GreetingsKeyName = []byte("greetings")

// Max length of greeting
var MaxGreetingLen = 500

// One day
var DayDuration = time.Hour * 24

// Template represents the various greetings the bot can perform
type Template struct {
	NewUser         string `json:"newUser"`
	ReturningUser   string `json:"returningUser"`
	ConsecutiveUser string `json:"consecutiveUser"`
	GreetTrolls     bool   `json:"greetTrolls"`
}

// Validate validates the Template
func (t *Template) Validate() error {
	if len(t.NewUser) > 500 {
		return fmt.Errorf("newUser greeting cannot exceed 500 characters")
	} else if _, err := template.New("msg").Parse(t.NewUser); err != nil {
		return fmt.Errorf("newUser is not a valid template: error %v", err)
	}

	if len(t.ReturningUser) > 500 {
		return fmt.Errorf("returningUser greeting cannot exceed 500 characters")
	} else if _, err := template.New("msg").Parse(t.ReturningUser); err != nil {
		return fmt.Errorf("returningUser is not a valid template: error %v", err)
	}

	if len(t.ConsecutiveUser) > 500 {
		return fmt.Errorf("consecutiveUser greeting cannot exceed 500 characters")
	} else if _, err := template.New("msg").Parse(t.ConsecutiveUser); err != nil {
		return fmt.Errorf("consecutiveUser is not a valid template: error %v", err)
	}

	return nil
}

// Save saves a Template to a bucket
func (t *Template) Save(bucket []byte) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		b, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return bkt.Put(TemplateKeyName, b)
	})
}

// Get gets a Template from a bucket
func Get(bucket []byte) (*Template, error) {
	tmpl := &Template{}
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		b := bkt.Get(TemplateKeyName)
		if b == nil {
			return nil
		}

		return json.Unmarshal(b, &tmpl)
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Join handles if a user should be greeted and what type of greeting they should receive
func Join(userBucket, botBucket []byte, cmd *commands.Command) string {
	msg := ""
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(botBucket).Bucket(GreetingsKeyName)

		username := cmd.Get("username")
		if len(username) == 0 {
			return fmt.Errorf("chat command did not have a username")
		}

		// get greeting stat
		grt := greetingStat{
			Username: username,
			troll:    (cmd.Get("role") == "guest"),
		}
		b := bkt.Get([]byte(username))
		if b != nil {
			if err := json.Unmarshal(b, &grt); err != nil {
				return err
			}
		}

		// get greetings template for chat room owner
		tmpl := &Template{}
		userBkt := tx.Bucket(userBucket)
		if userBkt == nil {
			// user has not a greetting template yet
			return nil
		}
		b = userBkt.Get(TemplateKeyName)
		if b == nil {
			// no message if we don't have any templates
			return nil
		}
		if err := json.Unmarshal(b, &tmpl); err != nil {
			return err
		}
		grt.tmpl = tmpl

		// get greeting message
		msg = grt.message()
		if len(msg) == 0 {
			return nil
		}

		// save greeting stat
		b, err := json.Marshal(grt)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(username), b)
	})

	if err != nil {
		log.Println("msg='greetings-join-error', error='%s'", err)
	}

	return msg
}

type greetingStat struct {
	Time       time.Time
	DaysInARow int
	tmpl       *Template
	Username   string `json:"-"`
	troll      bool
}

func (g *greetingStat) message() (msg string) {
	defer func() {
		if len(msg) > 0 {
			g.Time = time.Now()
		}
	}()

	if !g.tmpl.GreetTrolls && g.troll {
		return ""
	}

	if g.Time.IsZero() {
		// greet user as a new user
		g.DaysInARow = 1
		return parseTemplate(g.tmpl.NewUser, g)
	}

	today := now.BeginningOfDay()
	lastVisit := now.New(g.Time).BeginningOfDay()
	if today.Equal(lastVisit.Add(time.Hour * 24)) {
		// greet as a consecutive user if the user returned the next day
		g.DaysInARow++
		return parseTemplate(g.tmpl.ConsecutiveUser, g)
	}

	if time.Now().After(g.Time.Add(DayDuration)) {
		// greet as a returning user if it's been more than day
		return parseTemplate(g.tmpl.ReturningUser, g)
	}

	// no message
	return ""
}

func parseTemplate(tmpl string, d interface{}) string {
	t, err := template.New("msg").Parse(tmpl)
	if err != nil {
		log.Println("msg='error parsing template', template='%s', error='%v'", tmpl, err)
		return ""
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, d); err != nil {
		log.Println("msg='error executing template', template='%s', data='%+v', error='%v'", tmpl, d, err)
		return ""
	}

	return buf.String()
}

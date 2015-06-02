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

// greeting types
var (
	newUser          = "newUser"
	returningUser    = "returningUser"
	consecutiveUser  = "consecutiveUser"
	answeringMachine = "answeringMachine"
)

type Event struct {
	Type       string    `json:"type"`
	Response   string    `json:"response"`
	Username   string    `json:"username"`
	DaysInARow int       `json:"daysInARow"`
	NewUser    bool      `json:"newUser"`
	LastVisit  time.Time `json:"lastVisit"`
	Time       time.Time `json:"time"`
	Private    bool      `json:"private"`

	troll bool
	tmpl  *Template
}

// Max length of greeting
var MaxGreetingLen = 500

// One day
var DayDuration = time.Hour * 24

// Template represents the various greetings the bot can perform
type Template struct {
	NewUser            string `json:"newUser"`
	ReturningUser      string `json:"returningUser"`
	ConsecutiveUser    string `json:"consecutiveUser"`
	Private            bool   `json:"private"`
	GreetTrolls        bool   `json:"greetTrolls"`
	AnsweringMachine   string `json:"answeringMachine"`
	AnsweringMachineOn bool   `json:"answeringMachineOn"`
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

	if len(t.AnsweringMachine) > 500 {
		return fmt.Errorf("answeringMachine greeting cannot exceed 500 characters")
	} else if _, err := template.New("msg").Parse(t.AnsweringMachine); err != nil {
		return fmt.Errorf("answeringMachine is not a valid template: error %v", err)
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
func Join(userBucket, botBucket []byte, cmd *commands.Command) Event {
	e := Event{}
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(botBucket)
		if bkt == nil {
			return fmt.Errorf("missing bot bucket: %s", string(botBucket))
		}
		bkt = bkt.Bucket(GreetingsKeyName)
		if bkt == nil {
			return fmt.Errorf("missing grettings bucket: %s", string(GreetingsKeyName))
		}

		e.troll = cmd.Get("role") == "guest"
		e.Username = cmd.Get("username")
		if len(e.Username) == 0 {
			return fmt.Errorf("chat command did not have a username")
		}

		// get chat user's info
		b := bkt.Get([]byte(e.Username))
		if b != nil {
			if err := json.Unmarshal(b, &e); err != nil {
				return err
			}
		}

		// get greetings template for chat room owner
		userBkt := tx.Bucket(userBucket)
		if userBkt == nil {
			// user does not have greetting template yet
			return nil
		}
		b = userBkt.Get(TemplateKeyName)
		if b == nil {
			// no message if we don't have any templates
			return nil
		}
		if err := json.Unmarshal(b, &e.tmpl); err != nil {
			return err
		}

		// clear out old response and type
		e.Response = ""
		e.Type = ""

		// populate response and type
		e.populate()
		if len(e.Response) == 0 {
			return nil
		}

		// save greeting stat
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(e.Username), b)
	})

	if err != nil {
		log.Printf("msg='greetings-join-error', error='%s'\n", err)
	}

	return e
}

func (e *Event) populate() {
	defer func() {
		if e.tmpl.AnsweringMachineOn {
			e.Private = false
			return
		}
		if len(e.Response) > 0 {
			e.LastVisit = e.Time
			e.Time = time.Now()
			e.Private = e.tmpl.Private
		}
	}()

	if !e.tmpl.GreetTrolls && e.troll {
		return
	}

	if e.tmpl.AnsweringMachineOn {
		e.Type = answeringMachine
		e.parseTemplate(e.tmpl.AnsweringMachine)
		return
	}

	if e.Time.IsZero() {
		// greet user as a new user
		e.DaysInARow = 1
		e.Type = newUser
		e.parseTemplate(e.tmpl.NewUser)
		return
	}

	today := now.BeginningOfDay()
	lastVisit := now.New(e.Time).BeginningOfDay()
	if today.Equal(lastVisit.Add(time.Hour * 24)) {
		// greet as a consecutive user if the user returned the next day
		e.DaysInARow++
		e.Type = consecutiveUser
		e.parseTemplate(e.tmpl.ConsecutiveUser)
		return
	}

	if time.Now().After(e.Time.Add(DayDuration)) {
		// greet as a returning user if it's been more than day
		e.Type = returningUser
		e.parseTemplate(e.tmpl.ReturningUser)
	}
}

func (e *Event) parseTemplate(tmpl string) {
	t, err := template.New("msg").Parse(tmpl)
	if err != nil {
		log.Println("msg='error parsing template', template='%s', error='%v'", tmpl, err)
		return
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, e); err != nil {
		log.Println("msg='error executing template', template='%s', data='%+v', error='%v'", tmpl, e, err)
		return
	}

	e.Response = buf.String()
}

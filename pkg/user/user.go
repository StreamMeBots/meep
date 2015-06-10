package user

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/StreamMeBots/meep/pkg/buckets"
	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/StreamMeBots/meep/pkg/db"

	"github.com/boltdb/bolt"
)

// Errors
var (
	ErrNotFound = errors.New("User not found")
)

// Links represents a user's links
type Links struct {
	Avatar struct {
		Href     string `json:"href"`
		Template string `json:"template,omitempty"`
	} `json:"avatar,omitempty"`

	FallbackAvatar struct {
		Href     string `json:"href"`
		Template string `json:"template,omitempty"`
	} `json:"fallbackAvatar,omitempty"`
}

func Users() ([]*User, error) {
	users := []*User{}
	err := db.DB.View(func(tx *bolt.Tx) error {
		crs := buckets.UserData(tx).Cursor()
		for k, v := crs.First(); k != nil; k, v = crs.Next() {
			u := &User{}
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			users = append(users, u)
		}

		return nil
	})
	if err != nil {
		log.Printf("msg='error-getting-users', error='%v'\n", err)
		return nil, err
	}

	return users, nil
}

// User represents the fields that belong to a stream.me user
type User struct {
	Name       string `json:"displayName"`
	Username   string `json:"username"`
	Slug       string `json:"slug"`
	PublicId   string `json:"publicId"`
	Email      string `json:"email"`
	ChatRoomId string `json:"chatRoomId"`
	SessId     string `json:"sessId"`
	Links      Links  `json:"_links"`
}

func (u *User) BucketKey() []byte {
	return []byte(u.PublicId)
}

func (u *User) Save() error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return buckets.UserData(tx).Put(u.BucketKey(), b)
	})
	if err != nil {
		log.Println("msg='error-saving-user-data' error='%v' user='%+v'\n", err, u)
		return err
	}

	return nil
}

// Get gets a user from stream.me using a pre-authorized http client
func GetByClient(client *http.Client, userIp string) (*User, error) {

	resp, err := client.Get(config.Conf.Url + "/api-user/v1/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrNotFound
	}

	u := &User{}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	if len(u.PublicId) == 0 {
		return nil, ErrNotFound
	}

	// hack until the api provies the user's chat room
	u.ChatRoomId = "user:" + u.PublicId + ":web"

	// create a session ID
	u.SessId = fmt.Sprintf("%x", sha1.Sum([]byte(u.PublicId+strconv.FormatInt(time.Now().UnixNano(), 10))))

	return u, nil
}

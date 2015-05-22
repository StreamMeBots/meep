package user

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Errors
var (
	ErrNotFound = errors.New("User not found")
)

func BucketName(userPublicId string) []byte {
	return []byte(`user:` + userPublicId)
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
	Ip         string `json:"ip"`
}

// Get gets a user from stream.me using a pre-authorized http client
func GetByClient(client *http.Client, userIp string) (*User, error) {
	resp, err := client.Get("http://pds.dev.ifi.tv/api-user/v1/me")
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
	u.Ip = userIp

	return u, nil
}

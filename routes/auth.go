package routes

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/StreamMeBots/meep/pkg/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// oauth config and user authed clients
var (
	auth        authConfig
	userClients = UserClients{clients: map[string]UserClient{}}
	userKey     = "user"
)

// UserClients is used to store a user's authorized http client
type UserClients struct {
	sync.RWMutex
	clients map[string]UserClient
}

// UserClient contains the user info and the authed client used to interact with stream.me
type UserClient struct {
	client *http.Client
	User   user.User
	Token  string
}

// Get a user's http client
func (uc *UserClients) Get(sessid string) (UserClient, bool) {
	uc.RLock()
	defer uc.RUnlock()
	c, ok := uc.clients[sessid]
	return c, ok
}

// Add a user's http client
func (uc *UserClients) Add(sessid string, u user.User, client *http.Client) {
	uc.Lock()
	defer uc.Unlock()
	uc.clients[sessid] = UserClient{
		User:   u,
		client: client,
	}
}

func newAuth(c config.Config) authConfig {
	conf := oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{"chat", "account"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURL,
			TokenURL: c.TokenURL,
		},
		RedirectURL: c.RedirectURL,
	}

	return authConfig{
		conf:    conf,
		clients: make(map[string]*http.Client),
	}
}

// authConfig holds the oauth2 config for creating new authorized clients
type authConfig struct {
	conf         oauth2.Config
	sync.RWMutex                         // gaurds clients
	clients      map[string]*http.Client // user's authorized clients
}

// loginHandler redirects to stream.me to start the oauth2 process
func (a *authConfig) loginHandler(ctx *gin.Context) {
	// check if we have an authed client
	if isUserAuthed(ctx) {
		// go home if we are already logged in
		ctx.Redirect(302, config.Conf.Host())
		return
	}

	url := a.conf.AuthCodeURL("fooBarBaz")
	Debugln("AuthCodeURL:", url)
	ctx.Redirect(302, url)
}

// redirectHandler handles the redirect from stream.me and attempts to get the user's token and user information.
func (a *authConfig) redirectHandler(ctx *gin.Context) {
	Debugln("Oauth Redirect URL:", ctx.Request.URL.String())

	if errMsg := ctx.Request.FormValue("error"); len(errMsg) > 0 {
		ctx.Redirect(302, fmt.Sprintf("%s/?error='%s'", config.Conf.Host(), ctx.Request.FormValue("error_description")))
		return
	}

	// get auth code
	code := ctx.Request.FormValue("code")

	// get auth token
	tok, err := a.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("exchange error:", err)
		ctx.Redirect(302, fmt.Sprintf("%s/?error='Unable to get user auth token from stream.me'", config.Conf.Host()))
		return
	}

	// create authorized client
	client := a.conf.Client(oauth2.NoContext, tok)

	// get user from stream.me using the authed client
	u, err := user.GetByClient(client, ctx.Request.RemoteAddr)
	if err != nil {
		ctx.Redirect(302, fmt.Sprintf("%s/?error='Unable to get user from stream.me'", config.Conf.Host()))
		return
	}

	// save user info
	if err := u.Save(); err != nil {
		ctx.JSON(500, map[string]string{
			"message": "Error saving user information",
		})
	}

	// save the user
	userClients.Add(u.SessId, *u, client)

	// write session cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "sessid",
		Value:   u.SessId,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 30),
	})

	ctx.Redirect(302, config.Conf.Host())
}

// isUserAuthed is a helper function to check if the user is already authenticated
func isUserAuthed(ctx *gin.Context) bool {
	_, ok := userClients.Get(getSessId(ctx))
	return ok
}

func getSessId(ctx *gin.Context) string {
	c, err := ctx.Request.Cookie("sessid")
	if err != nil {
		return ""
	}

	return c.Value
}

func checkAuth(ctx *gin.Context) {
	u, ok := userClients.Get(getSessId(ctx))
	if !ok {
		ctx.JSON(401, map[string]string{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}
	ctx.Set(userKey, u)
	ctx.Next()
}

func getAuthedUser(ctx *gin.Context) UserClient {
	u, ok := ctx.Get(userKey)
	if !ok {
		panic("missing checkAuth middleware") // our route is missing auth middleware if we hit here
	}
	return u.(UserClient)
}

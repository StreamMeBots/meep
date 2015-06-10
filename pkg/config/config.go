package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

// Conf is the loaded config
var Conf Config

// Config represents the configuration options
type Config struct {
	ConfigPath        string `json:"-"`
	BotKey            string `json:"botKey"`
	BotSecret         string `json:"botSecret"`
	ClientId          string `josn:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	ChatHost          string `json:"chatHost"`
	ServerPort        string `json:"serverPort"`
	ServerHost        string `json:"serverHost"`
	ServerBehindProxy bool   `json:"serverBehindProxy"`
	AuthURL           string `json:"authURL"`
	TokenURL          string `json:"tokenURL"`
	RedirectURL       string `json:"redirectURL"`
	Url               string `json:"URL"`
	Debug             bool   `json:"debug"`
}

func (c *Config) Host() string {
	if c.ServerBehindProxy {
		return c.ServerHost + ":" + c.ServerPort
	}
	return c.ServerHost
}

func init() {
	// load config via command line flags
	flag.StringVar(&Conf.ConfigPath, "config-path", "", "path to a JSON config file")
	flag.StringVar(&Conf.BotKey, "bot-key", "", "bot key from stream.me")
	flag.StringVar(&Conf.BotSecret, "bot-secret", "", "bot secret from stream.me")
	flag.StringVar(&Conf.ClientId, "client-id", "", "oauth2 client id from stream.me")
	flag.StringVar(&Conf.ClientSecret, "client-secret", "", "oauth2 client secret from stream.me")
	flag.StringVar(&Conf.ChatHost, "chat-host", "", "stream.me chat server address")
	flag.StringVar(&Conf.ServerPort, "server-port", "", "http server port")
	flag.StringVar(&Conf.ServerHost, "server-host", "", "http server host")
	flag.StringVar(&Conf.AuthURL, "auth-url", "", "oauth2 auth url")
	flag.StringVar(&Conf.TokenURL, "token-url", "", "oauth2 token url")
	flag.StringVar(&Conf.RedirectURL, "redirect-url", "http://localhost:8888/redirect-url", "oauth redirect url")
	flag.StringVar(&Conf.Url, "url", "", "stream.me address")
	flag.BoolVar(&Conf.ServerBehindProxy, "behind-proxy", false, "indicate if the server is behind a proxy")
	flag.BoolVar(&Conf.Debug, "debug", false, "enable debug logging")
}

// CheckConfigPath checks the config if the 'config-path' flag was set. If the flag was set the config
// is loaded from the specified config path
func CheckConfigPath() {
	if len(Conf.ConfigPath) == 0 {
		return
	}

	b, err := ioutil.ReadFile(Conf.ConfigPath)
	if err != nil {
		log.Fatal("Error reading config JSON config file:", err)
	}

	if err := json.Unmarshal(b, &Conf); err != nil {
		log.Fatal("Error trying to JSON Unmarshal the config file:", err)
	}
}

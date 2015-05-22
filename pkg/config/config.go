package config

import "flag"

var Conf Config

type Config struct {
	BotKey            string
	BotSecret         string
	ClientId          string
	ClientSecret      string
	ChatHost          string
	ServerPort        string
	ServerHost        string
	ServerBehindProxy bool
	AuthURL           string
	TokenURL          string
	Url               string
	Debug             bool
}

func init() {
	// load config via command line flags
	flag.StringVar(&Conf.BotKey, "bot-key", "", "bot key from stream.me")
	flag.StringVar(&Conf.BotSecret, "bot-secret", "", "bot secret from stream.me")
	flag.StringVar(&Conf.ClientId, "client-id", "", "oauth2 client id from stream.me")
	flag.StringVar(&Conf.ClientSecret, "client-secret", "", "oauth2 client secret from stream.me")
	flag.StringVar(&Conf.ChatHost, "chat-host", "", "stream.me chat server address")
	flag.StringVar(&Conf.ServerPort, "server-port", "", "http server port")
	flag.StringVar(&Conf.ServerHost, "server-host", "", "http server host")
	flag.StringVar(&Conf.AuthURL, "auth-url", "", "oauth2 auth url")
	flag.StringVar(&Conf.TokenURL, "token-url", "", "oauth2 token url")
	flag.StringVar(&Conf.Url, "url", "", "stream.me address")
	flag.BoolVar(&Conf.ServerBehindProxy, "behind-proxy", false, "indicate if the server is behind a proxy")
	flag.BoolVar(&Conf.Debug, "debug", false, "enable debug logging")
}

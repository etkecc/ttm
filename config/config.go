package config

import (
	"gitlab.com/etke.cc/go/env"
)

// Config struct
type Config struct {
	// Log option to send full log output to matrix
	Log bool
	// NoHTML option to skip HTML formatted message to matrix
	NoHTML bool
	// NoText option to skip plaintext message to matrix
	NoText bool
	// NoTime option to skip time stats in message
	NoTime bool
	// MsgType option to set message type, default is "m.text", but "m.notice" is ok, too
	MsgType string
	// NoticeFail option to explicitly send "m.notice" message on exit code > 0
	NoticeFail bool

	// Token is an matrix user access token, to use without login and password
	Token string
	// Login is a matrix user login, to use with password
	Login string
	// Room is a target room ID or alias where messages will be sent
	Room string
	// Password is a matrix user password, to use with login
	Password string
	// Homeserver is a target HS url (delegation not supported)
	Homeserver string
}

const prefix = "ttm"

// New config
func New() *Config {
	env.SetPrefix(prefix)
	room := env.String("roomid", "")
	if room == "" {
		room = env.String("room", "")
	}

	return &Config{
		// connection
		Homeserver: env.String("homeserver", ""),
		Password:   env.String("password", ""),
		Room:       room,
		Login:      env.String("login", ""),
		Token:      env.String("token", ""),

		// options
		Log:        env.Bool("log"),
		NoHTML:     env.Bool("nohtml"),
		NoText:     env.Bool("notext"),
		NoTime:     env.Bool("notime"),
		MsgType:    env.String("msgtype", ""),
		NoticeFail: env.Bool("noticefail"),
	}
}

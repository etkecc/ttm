package config

import (
	"os"
	"strings"
)

// Config ...
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

func env(shortkey string) string {
	key := strings.ToUpper(prefix + "_" + strings.ReplaceAll(shortkey, ".", "_"))
	return strings.TrimSpace(os.Getenv(key))
}

func envBool(shortkey string) bool {
	key := strings.ToUpper(prefix + "_" + strings.ReplaceAll(shortkey, ".", "_"))
	value := strings.TrimSpace(os.Getenv(key))
	return (value == "1" || value == "true" || value == "yes")
}

// New config
func New() *Config {
	room := env("roomid")
	if room == "" {
		room = env("room")
	}

	return &Config{
		// connection
		Homeserver: env("homeserver"),
		Password:   env("password"),
		Room:       room,
		Login:      env("login"),
		Token:      env("token"),

		// options
		Log:        envBool("log"),
		NoHTML:     envBool("nohtml"),
		NoText:     envBool("notext"),
		NoTime:     envBool("notime"),
		MsgType:    env("msgtype"),
		NoticeFail: envBool("noticefail"),
	}
}

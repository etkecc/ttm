package config

import (
	"os"
	"strings"
)

// Config ...
type Config struct {
	Login      string
	RoomID     string
	Password   string
	Homeserver string
}

const prefix = "ttm"

func readEnv(shortkey string) string {
	key := strings.ToUpper(prefix + "_" + strings.ReplaceAll(shortkey, ".", "_"))
	return strings.TrimSpace(os.Getenv(key))
}

// New config
func New() *Config {
	return &Config{
		Homeserver: readEnv("homeserver"),
		Password:   readEnv("password"),
		RoomID:     readEnv("roomid"),
		Login:      readEnv("login"),
	}
}

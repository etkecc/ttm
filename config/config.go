package config

import (
	"os"
	"strings"
)

// Config ...
type Config struct {
	Log    bool
	NoTime bool
	NoHTML bool

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

func readEnvBool(shortkey string) bool {
	key := strings.ToUpper(prefix + "_" + strings.ReplaceAll(shortkey, ".", "_"))
	value := strings.TrimSpace(os.Getenv(key))
	return (value == "1" || value == "true" || value == "yes")
}

// New config
func New() *Config {
	return &Config{
		// connection
		Homeserver: readEnv("homeserver"),
		Password:   readEnv("password"),
		RoomID:     readEnv("roomid"),
		Login:      readEnv("login"),

		// options
		NoTime: readEnvBool("notime"),
		NoHTML: readEnvBool("nohtml"),
		Log:    readEnvBool("log"),
	}
}

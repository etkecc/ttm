package config

import (
	"testing"
)

var env = map[string]string{
	"TTM_HOMESERVER": "https://matrix.example.com",
	"TTM_LOGIN":      "test",
	"TTM_PASSWORD":   "password",
	"TTM_TOKEN":      "",
	"TTM_ROOMID":     "!test:example.com",
	"TTM_NOTIME":     "1",
	"TTM_NOHTML":     "1",
	"TTM_LOG":        "1",
}

func TestNew(t *testing.T) {
	for key, value := range env {
		t.Setenv(key, value)
	}

	cfg := New()

	if cfg.Homeserver != "https://matrix.example.com" {
		t.Fail()
	}
	if cfg.Login != "test" {
		t.Fail()
	}
	if cfg.Password != "password" {
		t.Fail()
	}
	if cfg.Token != "" {
		t.Fail()
	}
	if !cfg.NoTime {
		t.Fail()
	}
	if !cfg.NoHTML {
		t.Fail()
	}
	if !cfg.Log {
		t.Fail()
	}
}

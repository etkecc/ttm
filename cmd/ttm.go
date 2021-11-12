package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gitlab.com/etke.cc/ttm/compose"
	"gitlab.com/etke.cc/ttm/config"
	"gitlab.com/etke.cc/ttm/matrix"
	"gitlab.com/etke.cc/ttm/term"
)

func main() {
	cfg := config.New()
	if cfg.NoText && cfg.NoHTML {
		panic(errors.New("you can't use both TTM_NOHTML and TTM_NOTEXT at the same time"))
	}

	command := compose.Command()
	sender := matrix.New(cfg.Homeserver, cfg.Login, cfg.Password, cfg.Token, cfg.RoomID, cfg.MsgType)
	// login (password auth only) in separate goroutine, to save some time
	go login(sender)
	process, err := term.RunCommand(command, cfg.NoTime, cfg.Log)
	if err != nil {
		panic(err)
	}
	// override msgtype if TTM_NOTICEFAIL and exit code != 0
	if process.Exit != 0 && cfg.NoticeFail {
		sender.SetMsgType("m.notice")
	}

	plaintext, html := compose.Message(process, cfg.NoTime, cfg.NoHTML, cfg.NoText)
	err = sender.SendMessage(plaintext, html)
	if err != nil {
		panic(err)
	}

	os.Exit(process.Exit)
}

func login(client *matrix.Client) {
	ctx := context.Background()
	err := client.Login(ctx)
	if err != nil {
		fmt.Println("TTM ERROR (matrix):", err)
	}
}

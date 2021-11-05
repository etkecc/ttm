package main

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/etke.cc/ttm/compose"
	"gitlab.com/etke.cc/ttm/config"
	"gitlab.com/etke.cc/ttm/matrix"
	"gitlab.com/etke.cc/ttm/term"
)

func main() {
	cfg := config.New()
	command := compose.Command()
	sender := matrix.New(cfg.Homeserver, cfg.Login, cfg.Password, cfg.RoomID)
	// login in separate goroutine, to save some time
	go login(sender)
	process, err := term.RunCommand(command, cfg.NoTime, cfg.Log)
	if err != nil {
		panic(err)
	}

	plaintext, html := compose.Message(process, cfg.NoTime, cfg.NoHTML)
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

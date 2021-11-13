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
	cfg := getConfig()
	command := getCommand()
	client := getClient(cfg.Homeserver, cfg.Login, cfg.Password, cfg.Token, cfg.Room, cfg.MsgType)
	process := runCommand(command, cfg.NoTime, cfg.Log)
	plaintext, html := compose.Message(process, cfg.NoTime, cfg.NoHTML, cfg.NoText)

	// override msgtype if TTM_NOTICEFAIL and exit code != 0
	if process.Exit != 0 && cfg.NoticeFail {
		client.MsgType = "m.notice"
	}

	sendMessage(client, plaintext, html)

	os.Exit(process.Exit)
}

func getConfig() *config.Config {
	cfg := config.New()
	if cfg.NoText && cfg.NoHTML {
		fmt.Println("TTM ERROR (config): you can't use both TTM_NOHTML and TTM_NOTEXT at the same time")
		fmt.Println(compose.Help())
		os.Exit(1)
	}

	return cfg
}

func getCommand() string {
	command := compose.Command()
	if command == "--help" || command == "-h" {
		fmt.Println(compose.Help())
		os.Exit(0)
	}

	return command
}

func getClient(homeserver, login, password, token, room, msgtype string) *matrix.Client {
	ctx := context.Background()
	client := matrix.New(homeserver, login, password, token, room, msgtype)
	go func(ctx context.Context, c *matrix.Client, room string) {
		roomID, err := c.ResolveRoom(ctx, room)
		if err != nil {
			fmt.Println("TTM ERROR (matrix):", err)
		} else {
			c.Room = roomID
		}

		err = c.Login(ctx)
		if err != nil {
			fmt.Println("TTM ERROR (matrix):", err)
		}
	}(ctx, client, room)

	return client
}

func runCommand(command string, notime bool, log bool) *term.Process {
	process, err := term.RunCommand(command, notime, log)
	if err != nil {
		fmt.Println("TTY ERROR:", err)
		fmt.Println(compose.Help())
		os.Exit(1)
	}

	return process
}

func sendMessage(client *matrix.Client, plaintext, html string) {
	err := client.SendMessage(plaintext, html)
	if err != nil {
		fmt.Println("TTM ERROR (matrix):", err)
		fmt.Println(compose.Help())
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.com/etke.cc/ttm/config"
	"gitlab.com/etke.cc/ttm/matrix"
	"gitlab.com/etke.cc/ttm/term"
)

func main() {
	cfg := config.New()
	sender := matrix.New(cfg.Homeserver, cfg.Login, cfg.Password, cfg.RoomID)
	// login in separate goroutine, to save some time
	go login(sender)
	process, err := term.RunCommand(getCommand())
	if err != nil {
		panic(err)
	}

	plaintext, html := composeMessage(process)
	err = sender.SendMessage(plaintext, html)
	if err != nil {
		panic(err)
	}

	os.Exit(process.Exit)
}

// getCommand from arguments
func getCommand() string {
	if len(os.Args) == 1 {
		return ""
	}

	return strings.Join(os.Args[1:], " ")
}

func composeMessage(process *term.Process) (string, string) {
	var text strings.Builder
	var html strings.Builder
	// plain text
	text.WriteString("ttm report\n")
	text.WriteString("" + process.Command + "\n\n")
	if len(process.Log) > 0 {
		for _, line := range process.Log {
			text.WriteString(line + "\n")
		}
	}
	text.WriteString("\n\n")
	text.WriteString("real\t" + process.Time.Real + "\n")
	text.WriteString("user\t" + process.Time.User + "\n")
	text.WriteString("sys\t" + process.Time.User + "\n")
	text.WriteString("\n\n")
	text.WriteString("Exit code: " + strconv.Itoa(process.Exit))

	// html
	html.WriteString("<b>ttm report</b>")
	html.WriteString("<pre>" + process.Command + "</pre><br>")
	if len(process.Log) > 0 {
		html.WriteString("<pre>\n")
		for _, line := range process.Log {
			html.WriteString(line + "\n")
		}
		html.WriteString("</pre>\n\n")
	}
	html.WriteString("<pre>\n")
	html.WriteString("real\t" + process.Time.Real + "\n")
	html.WriteString("user\t" + process.Time.User + "\n")
	html.WriteString("sys\t" + process.Time.User + "\n")
	html.WriteString("</pre>\n\n")
	html.WriteString("Exit code: <code>" + strconv.Itoa(process.Exit) + "</code>")

	return text.String(), html.String()
}

func login(client *matrix.Client) {
	ctx := context.Background()
	err := client.Login(ctx)
	if err != nil {
		fmt.Println("TTM ERROR (matrix):", err)
	}
}

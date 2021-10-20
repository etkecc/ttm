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
	process, err := term.RunCommand(getCommand(), cfg.NoTime, cfg.Log)
	if err != nil {
		panic(err)
	}

	plaintext, html := composeMessage(process, cfg.NoTime)
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

func composeMessage(process *term.Process, notime bool) (string, string) {
	var text strings.Builder
	var html strings.Builder
	var logsb strings.Builder
	for _, line := range process.Log {
		logsb.WriteString(line + "\n")
	}
	log := logsb.String()
	// Here we try to roughly calculate if log will exceed matrix payload limit and shirnk it a bit
	// note that we don't do precise calculations, because in that case we will need to generate payload multiple times,
	// so following solution will work in 99% cases and the last 1% will be passed as-is.
	maxLogSize := matrix.MaxPayloadSize - matrix.InfrastructurePayloadSize - matrix.SuggestedPayloadBuffer
	// log will be sent twice - in html and plaintext form, so we use log length * 2
	logSizeDiff := logsb.Len()*2 - maxLogSize
	if logSizeDiff > 0 {
		singleLogSizeDiff := (logSizeDiff / 2)
		log = "# the beginning of the log is omitted due to protocol limitations\n" + log[singleLogSizeDiff:]
	}

	// plain text
	text.WriteString("ttm report\n")
	text.WriteString("" + process.Command + "\n\n")
	if len(process.Log) > 0 {
		text.WriteString(log)
	}
	text.WriteString("\n\n")
	if !notime {
		text.WriteString("real\t" + process.Time.Real + "\n")
		text.WriteString("user\t" + process.Time.User + "\n")
		text.WriteString("sys\t" + process.Time.User + "\n")
		text.WriteString("\n\n")
	}
	text.WriteString("Exit code: " + strconv.Itoa(process.Exit))

	// html
	html.WriteString("<b>ttm report</b>")
	html.WriteString("<pre>" + process.Command + "</pre><br>")
	if len(process.Log) > 0 {
		html.WriteString("<pre>\n")
		html.WriteString(log)
		html.WriteString("</pre>")
	}
	html.WriteString("\n\n")
	if !notime {
		html.WriteString("<pre>\n")
		html.WriteString("real\t" + process.Time.Real + "\n")
		html.WriteString("user\t" + process.Time.User + "\n")
		html.WriteString("sys\t" + process.Time.User + "\n")
		html.WriteString("</pre>\n\n")
	}
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

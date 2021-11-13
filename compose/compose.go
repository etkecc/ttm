package compose

import (
	"os"
	"strconv"
	"strings"

	"gitlab.com/etke.cc/ttm/matrix"
	"gitlab.com/etke.cc/ttm/term"
)

// Command compose command to run from cli args
func Command() string {
	if len(os.Args) == 1 {
		return ""
	}

	return strings.Join(os.Args[1:], " ")
}

// Help returns usage help/docs
func Help() string {
	var help strings.Builder
	help.WriteString("\n")
	help.WriteString("Time To Matrix is a time-like command that will send end of an arbitrary command output and some other info (like exit status) to matrix room.\n\n")
	help.WriteString("Usage:\n\n")
	help.WriteString("\tttm <command>\n\n")
	help.WriteString("If you want to get current configuration, run the following command: env | grep TTM_\n")
	help.WriteString("Check the https://gitlab.com/etke.cc/ttm for list of available configuration params and examples\n")
	return help.String()
}

// Message compose plaintext and html message to send
func Message(process *term.Process, notime bool, nohtml bool, notext bool) (string, string) {
	var text string
	var html string
	log := getLog(process, nohtml, notext)

	// plain text
	if !notext {
		text = getText(process, log, notime)
	}

	// html
	if !nohtml {
		html = getHTML(process, log, notime)
	}

	return text, html
}

func getLog(process *term.Process, nohtml bool, notext bool) string {
	var logsb strings.Builder
	for _, line := range process.Log {
		logsb.WriteString(line + "\n")
	}
	log := logsb.String()
	llm := getLogLengthModifier(nohtml, notext)
	// Here we try to roughly calculate if log will exceed matrix payload limit and shirnk it a bit
	// note that we don't do precise calculations, because in that case we will need to generate payload multiple times,
	// so following solution will work in 99% cases and the last 1% will be passed as-is.
	maxLogSize := matrix.MaxPayloadSize - matrix.InfrastructurePayloadSize - matrix.SuggestedPayloadBuffer
	// log will be sent twice - in html and plaintext form, so we use log length * 2
	logSizeDiff := logsb.Len()*llm - maxLogSize
	if logSizeDiff > 0 {
		singleLogSizeDiff := (logSizeDiff / llm)
		log = "# the beginning of the log is omitted due to protocol limitations\n" + log[singleLogSizeDiff:]
	}

	return log
}

func getText(process *term.Process, log string, notime bool) string {
	var text strings.Builder
	text.WriteString("ttm report\n")
	text.WriteString("" + process.Command + "\n\n")
	if len(log) > 0 {
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

	return text.String()
}

func getHTML(process *term.Process, log string, notime bool) string {
	var html strings.Builder
	html.WriteString("<b>ttm report</b>")
	html.WriteString("<pre>" + process.Command + "</pre><br>")
	if len(log) > 0 {
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

	return html.String()
}

// getLogLengthModifier returns the modifier of the allowed log length, based on HTML and plaintext settings.
// by default, both plaintext and HTML formatted body will be sent, so max log size <=~31kb,
// because log output should be duplicated in both formats, but if you want to skip the HTML formatted body or plaintext message,
// the max log size <=~63kb
func getLogLengthModifier(nohtml bool, notext bool) int {
	mod := 2
	if nohtml || notext {
		mod = 1
	}

	return mod
}

package term

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
)

// from https://github.com/acarl005/stripansi
const colorsregex = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var colors = regexp.MustCompile(colorsregex)

// Process info
type Process struct {
	Command string
	Time    ProcessTime
	Log     []string
	Exit    int
}

// ProcessTime duration info
type ProcessTime struct {
	Real string
	User string
	Sys  string
}

// RunCommand and return process info
func RunCommand(command string) (*Process, error) {
	process := &Process{Command: command}
	err := process.Run()

	return process, err
}

// Run the command
func (p *Process) Run() error {
	args := strings.Split(p.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	size, err := pty.GetsizeFull(os.Stdout)
	if err != nil {
		return err
	}
	startedAt := time.Now()
	stdout, err := pty.StartWithSize(cmd, size)
	if err != nil {
		return err
	}
	p.log(stdout)

	err = cmd.Wait()
	endedAt := time.Now()
	p.Time = ProcessTime{
		Real: endedAt.Sub(startedAt).String(),
		User: cmd.ProcessState.UserTime().String(),
		Sys:  cmd.ProcessState.SystemTime().String(),
	}
	p.Exit = cmd.ProcessState.ExitCode()

	fmt.Println("")
	fmt.Println("real\t", p.Time.Real)
	fmt.Println("user\t", p.Time.User)
	fmt.Println("sys\t", p.Time.Sys)
	if err != nil && strings.HasPrefix(err.Error(), "exit status") {
		return nil
	}

	return err
}

func (p *Process) log(r io.Reader) {
	var shouldLog bool
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)

		if !shouldLog && strings.HasPrefix(text, "PLAY RECAP") {
			shouldLog = true
		}

		if shouldLog {
			p.Log = append(p.Log, colors.ReplaceAllString(text, ""))
		}
	}
}

package compose

import (
	"os"
	"strings"
	"testing"

	"gitlab.com/etke.cc/ttm/matrix"
	"gitlab.com/etke.cc/ttm/term"
)

func TestCommand(t *testing.T) {
	expected := "echo test"
	os.Args = []string{"ttm", "echo", "test"}

	actual := Command()

	if expected != actual {
		t.Fail()
	}
}

func TestCommand_Empty(t *testing.T) {
	os.Args = []string{"only1"}

	actual := Command()

	if actual != "" {
		t.Fail()
	}
}

func TestMessage(t *testing.T) {
	expectedText := `ttm report
echo test

line 1
line 2


real	1m
user	1s
sys	1s


Exit code: 0`
	expectedHTML := `<b>ttm report</b><pre>echo test</pre><br><pre>
line 1
line 2
</pre>

<pre>
real	1m
user	1s
sys	1s
</pre>

Exit code: <code>0</code>`
	process := mockProcess(t, []string{"line 1", "line 2"})

	actualText, actualHTML := Message(process, false, false)

	if expectedText != actualText {
		t.Errorf("expected plaintext message doesn't equal to actual\nexpected:\n%s\n\nactual:\n%s", expectedText, actualText)
	}
	if expectedHTML != actualHTML {
		t.Errorf("expected html message doesn't equal to actual\nexpected:\n%s\n\nactual:\n%s", expectedHTML, actualHTML)
	}
}

func TestMessage_NoHTML(t *testing.T) {
	expectedText := `ttm report
echo test

line 1
line 2


real	1m
user	1s
sys	1s


Exit code: 0`
	expectedHTML := ""
	process := mockProcess(t, []string{"line 1", "line 2"})

	actualText, actualHTML := Message(process, false, true)

	if expectedText != actualText {
		t.Errorf("expected plaintext message doesn't equal to actual\nexpected:\n%s\n\nactual:\n%s", expectedText, actualText)
	}
	if expectedHTML != actualHTML {
		t.Errorf("expected html message doesn't equal to actual\nexpected:\n%s\n\nactual:\n%s", expectedHTML, actualHTML)
	}
}

func TestMessage_Shrink(t *testing.T) {
	var expectedText strings.Builder
	expectedText.WriteString("ttm report\n")
	expectedText.WriteString("echo test\n\n")
	expectedText.WriteString("# the beginning of the log is omitted due to protocol limitations\n")
	expectedText.WriteString(strings.Repeat("t", 62535) + "\n\n\n")
	expectedText.WriteString("Exit code: 0")
	expectedHTML := ""
	process := mockProcess(t, []string{strings.Repeat("t", matrix.MaxPayloadSize)})

	actualText, actualHTML := Message(process, true, true)

	if expectedText.String() != actualText {
		t.Errorf("expected plaintext message doesn't equal to actual\nexpected:\n%d\n\nactual:\n%d", len(expectedText.String()), len(actualText))
	}
	if expectedHTML != actualHTML {
		t.Errorf("expected html message doesn't equal to actual\nexpected:\n%s\n\nactual:\n%s", expectedHTML, actualHTML)
	}
}

func mockProcess(t *testing.T, lines []string) *term.Process {
	t.Helper()
	return &term.Process{
		Command: "echo test",
		Log:     lines,
		Time: term.ProcessTime{
			Real: "1m",
			User: "1s",
			Sys:  "1ms",
		},
	}
}

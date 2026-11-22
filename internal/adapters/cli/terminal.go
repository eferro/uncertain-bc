package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type lineReader interface {
	ReadLine() (string, error)
}

type promptedScanner struct {
	lines  *bufio.Scanner
	stdout io.Writer
}

func newPromptedScanner(in io.Reader, stdout io.Writer) promptedScanner {
	return promptedScanner{lines: bufio.NewScanner(in), stdout: stdout}
}

func (p promptedScanner) ReadLine() (string, error) {
	fmt.Fprint(p.stdout, prompt)
	if p.lines.Scan() {
		return p.lines.Text(), nil
	}
	if err := p.lines.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func newLineEditor(screen io.ReadWriter, history term.History) *term.Terminal {
	editor := term.NewTerminal(screen, prompt)
	editor.History = history
	return editor
}

// terminalRepl runs the session with line editing when stdin is a real
// terminal, and falls back to plain line reading otherwise.
func (a App) terminalRepl(in *os.File, stdout, stderr io.Writer, run func(lineReader, io.Writer, io.Writer) int) int {
	restore, err := enterRawMode(in, stdout)
	if err != nil {
		return run(newPromptedScanner(in, stdout), stdout, stderr)
	}
	defer restore()
	editor := newLineEditor(struct {
		io.Reader
		io.Writer
	}{in, stdout}, loadHistory(a.HistoryPath))
	return run(editor, editor, editor)
}

// enterRawMode prepares the console for the line editor: raw input
// (ADR 0023) and, on Windows, escape processing on the output (ADR 0037).
// restore puts both back as they were.
func enterRawMode(in *os.File, stdout io.Writer) (restore func(), err error) {
	fd := int(in.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	restoreOutput := func() {}
	if out, ok := stdout.(*os.File); ok {
		restoreOutput = enableVirtualTerminal(out)
	}
	return func() {
		restoreOutput()
		_ = term.Restore(fd, state)
	}, nil
}

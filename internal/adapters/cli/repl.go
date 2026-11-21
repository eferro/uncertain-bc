package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
)

const prompt = "ubc> "

type replSession struct {
	session *application.Session
	style   style
	stdout  io.Writer
	stderr  io.Writer
}

// repl evaluates one line at a time in a single session and keeps going after
// errors, until the input ends or :quit (SPEC §7).
func repl(session *application.Session, st style, lines lineReader, stdout, stderr io.Writer) int {
	r := replSession{session: session, style: st, stdout: stdout, stderr: stderr}
	for {
		line, err := lines.ReadLine()
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(stdout)
			return exitOK
		}
		if err != nil {
			fmt.Fprintf(stderr, "ubc: stdin: %v\n", err)
			return exitEvalError
		}
		if quit := r.handle(line); quit {
			return exitOK
		}
	}
}

func (r replSession) handle(line string) (quit bool) {
	if strings.HasPrefix(line, ":") {
		quit, err := r.command(line[1:])
		if err != nil {
			fmt.Fprintf(r.stderr, "ubc: %v\n", err)
		}
		return quit
	}
	results, err := r.session.Evaluate(line)
	for _, result := range results {
		fmt.Fprintln(r.stdout, display(result, r.style))
	}
	if err != nil {
		fmt.Fprintf(r.stderr, "ubc: %v\n", err)
	}
	return false
}

func (r replSession) printVariable(v application.Variable) {
	fmt.Fprintf(r.stdout, "%s = %s\n", v.Name, format(v.Value, r.style))
}

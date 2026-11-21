package cli

import (
	"bufio"
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
func repl(session *application.Session, st style, stdin io.Reader, stdout, stderr io.Writer) int {
	r := replSession{session: session, style: st, stdout: stdout, stderr: stderr}
	lines := bufio.NewScanner(stdin)
	for fmt.Fprint(stdout, prompt); lines.Scan(); fmt.Fprint(stdout, prompt) {
		if quit := r.handle(lines.Text()); quit {
			return exitOK
		}
	}
	fmt.Fprintln(stdout)
	if err := lines.Err(); err != nil {
		fmt.Fprintf(stderr, "ubc: stdin: %v\n", err)
		return exitEvalError
	}
	return exitOK
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
		fmt.Fprintln(r.stdout, format(result.Value, r.style))
	}
	if err != nil {
		fmt.Fprintf(r.stderr, "ubc: %v\n", err)
	}
	return false
}

func (r replSession) printVariable(v application.Variable) {
	fmt.Fprintf(r.stdout, "%s = %s\n", v.Name, format(v.Value, r.style))
}

package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

const prompt = "ubc> "

// repl evaluates one line at a time in a single session and keeps going after
// errors, until the input ends (SPEC §7).
func repl(session *application.Session, st style, stdin io.Reader, stdout, stderr io.Writer) int {
	lines := bufio.NewScanner(stdin)
	for fmt.Fprint(stdout, prompt); lines.Scan(); fmt.Fprint(stdout, prompt) {
		results, err := session.Evaluate(lines.Text())
		printEvery(results, st, stdout)
		if err != nil {
			fmt.Fprintf(stderr, "ubc: %v\n", err)
		}
	}
	fmt.Fprintln(stdout)
	if err := lines.Err(); err != nil {
		fmt.Fprintf(stderr, "ubc: stdin: %v\n", err)
		return exitEvalError
	}
	return exitOK
}

func printEvery(results []domain.Result, st style, stdout io.Writer) {
	for _, result := range results {
		fmt.Fprintln(stdout, format(result.Value, st))
	}
}

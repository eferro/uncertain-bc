// Package cli is the command-line adapter: it turns arguments into
// evaluations and results into text.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

const (
	exitOK         = 0
	exitEvalError  = 1
	exitUsageError = 2
	exitParseError = 3
)

func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	expr := fs.String("e", "", "evaluate EXPR")
	if err := fs.Parse(args); err != nil {
		return exitUsageError
	}
	value, err := application.NewSession(application.DefaultSettings()).Evaluate(*expr)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitCode(err)
	}
	fmt.Fprintln(stdout, format(value))
	return exitOK
}

func exitCode(err error) int {
	var syntaxErr *domain.SyntaxError
	if errors.As(err, &syntaxErr) {
		return exitParseError
	}
	return exitEvalError
}

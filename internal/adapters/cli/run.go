// Package cli is the command-line adapter: it turns arguments into
// evaluations and results into text.
package cli

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	expr := fs.String("e", "", "evaluate EXPR")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	value, err := domain.Evaluate(*expr)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return 3
	}
	fmt.Fprintln(stdout, strconv.FormatFloat(value, 'g', -1, 64))
	return 0
}

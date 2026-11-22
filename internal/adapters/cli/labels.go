package cli

import (
	"fmt"
	"io"
	"regexp"

	"github.com/eferro/uncertain-bc/internal/domain"
)

// printResults prints expressions, and assignments too with --echo. In a
// program that prints several expressions each one says what it is (ADR 0053).
func printResults(results []domain.Result, st style, stdout io.Writer) {
	label := !st.plain && countExpressions(results) > 1
	for _, result := range results {
		if line, shown := resultLine(result, label, st); shown {
			fmt.Fprintln(stdout, line)
		}
	}
}

func resultLine(result domain.Result, label bool, st style) (string, bool) {
	switch {
	case result.Name != "":
		return display(result, st), st.echo
	case !label:
		return display(result, st), true
	}
	if named, ok := labelled(result); ok {
		return display(named, st), true
	}
	if st.json {
		return display(result, st), true
	}
	return fmt.Sprintf("L%d: %s", result.Line, display(result, st)), true
}

func countExpressions(results []domain.Result) int {
	expressions := 0
	for _, result := range results {
		if result.Name == "" {
			expressions++
		}
	}
	return expressions
}

// maxLabelLength keeps a label short enough to read as a name; a longer
// expression is better pointed at by its line (ADR 0053).
const maxLabelLength = 30

var variableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// labelled names an expression after its source text, as an assignment is
// named after its variable, when that text reads as a name.
func labelled(result domain.Result) (domain.Result, bool) {
	if len(result.Source) > maxLabelLength && !variableName.MatchString(result.Source) {
		return result, false
	}
	result.Name = result.Source
	return result, true
}

package cli

import (
	"fmt"
	"io"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

// checkAssertions reports every failed --assert, so one CI run shows them
// all, but stops at the first that cannot be evaluated (ADR 0029).
func checkAssertions(session *application.Session, asserts []string, st style, stderr io.Writer) int {
	code := exitOK
	for _, expr := range asserts {
		check, err := session.Assert(expr)
		if err != nil {
			reportError(source{name: "--assert " + expr}, err, stderr)
			return exitCode(err)
		}
		if !check.Holds() {
			fmt.Fprintf(stderr, "ubc: --assert %s failed: %s\n", expr, actual(check, st))
			code = exitAssertion
		}
	}
	return code
}

// actual shows what was compared, and how likely it held when the operands
// are uncertain.
func actual(check domain.Check, st style) string {
	if check.Operator == "" {
		return operand(check.Value, st)
	}
	compared := fmt.Sprintf("%s %s %s", operand(check.Left, st), check.Operator, operand(check.Right, st))
	if check.Value == domain.Probability(0) {
		return compared
	}
	return fmt.Sprintf("%s (P = %s)", compared, operand(check.Value, st))
}

// operand is a value's summary without the sparkline, whatever the output
// format, because it goes into a one-line message on stderr.
func operand(value domain.Value, st style) string {
	text := style{approx: st.approx, ellipsis: st.ellipsis}
	samples, uncertain := value.(domain.SampleSet)
	if !uncertain || samples.Dropped() == len(samples) || samples.Quantile(0) == samples.Quantile(1) {
		return format(value, text)
	}
	return withDropped(summaryNumbers(samples, text), samples.Dropped())
}

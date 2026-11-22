package cli

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
)

const sensitivityCells = 20

// sensitivityReport is a tornado list: one bar per variable, as long as its
// rank correlation with expr is strong, so the drivers stand out (ADR 0033).
func sensitivityReport(expr string, drivers []application.Driver, st style) string {
	width := 0
	for _, d := range drivers {
		width = max(width, len(d.Name))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "sensitivity of %s (rank correlation):\n", expr)
	for _, d := range drivers {
		bar := strings.Repeat(st.bar, int(math.Round(math.Abs(d.Correlation)*sensitivityCells)))
		line := fmt.Sprintf("  %-*s  %+.2f  %s", width, d.Name, d.Correlation, bar)
		b.WriteString(strings.TrimRight(line, " ") + "\n")
	}
	return b.String()
}

func showSensitivity(r replSession, expr string) (bool, error) {
	if expr == "" {
		return false, errors.New(":sens expects an expression")
	}
	drivers, err := r.session.Sensitivity(expr)
	if errors.Is(err, application.ErrCertainTarget) || errors.Is(err, application.ErrNoDrivers) {
		return false, fmt.Errorf(":sens %w", err)
	}
	if err != nil {
		return false, err
	}
	fmt.Fprint(r.stdout, sensitivityReport(expr, drivers, r.style))
	return false, nil
}

// printSensitivity answers --sens after the program, before the assertions,
// so a failing check still shows what drives the result.
func printSensitivity(session *application.Session, expr string, st style, stdout, stderr io.Writer) int {
	if expr == "" {
		return exitOK
	}
	drivers, err := session.Sensitivity(expr)
	if err != nil {
		reportError(source{name: "--sens " + expr}, err, stderr)
		return exitCode(err)
	}
	fmt.Fprint(stdout, sensitivityReport(expr, drivers, st))
	return exitOK
}

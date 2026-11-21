package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

type override struct {
	name string
	expr string
}

// parseOverrides splits each --set at its first '=', so the expression may
// hold a comparison such as ok=a==b (ADR 0028).
func parseOverrides(values []string) ([]override, error) {
	overrides := make([]override, 0, len(values))
	for _, value := range values {
		name, expr, found := strings.Cut(value, "=")
		if !found {
			return nil, fmt.Errorf("--set expects name=expr, got '%s'", value)
		}
		name = strings.TrimSpace(name)
		if !domain.IsName(name) {
			return nil, fmt.Errorf("--set: '%s' is not a variable name", name)
		}
		overrides = append(overrides, override{name: name, expr: expr})
	}
	return overrides, nil
}

// startSession applies the --set overrides before any program runs.
func startSession(opts options, settings application.Settings, stderr io.Writer) (*application.Session, int) {
	session := application.NewSession(settings)
	for _, o := range opts.overrides {
		if err := session.Override(o.name, o.expr); err != nil {
			reportError(source{name: "--set " + o.name}, err, stderr)
			return nil, exitCode(err)
		}
	}
	return session, exitOK
}

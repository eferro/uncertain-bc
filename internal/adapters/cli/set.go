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
		name, expr, err := splitNamed("--set", "name=expr", value)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, override{name: name, expr: expr})
	}
	return overrides, nil
}

func splitNamed(flagName, form, value string) (name, rest string, err error) {
	name, rest, found := strings.Cut(value, "=")
	if !found {
		return "", "", fmt.Errorf("%s expects %s, got '%s'", flagName, form, value)
	}
	name = strings.TrimSpace(name)
	if !domain.IsName(name) {
		return "", "", fmt.Errorf("%s: '%s' is not a variable name", flagName, name)
	}
	return name, rest, nil
}

// startSession loads the --data files and then applies the --set overrides,
// so overrides can use the data, all before any program runs.
func startSession(opts options, settings application.Settings, stderr io.Writer) (*application.Session, int) {
	session := application.NewSession(settings)
	for _, d := range opts.datasets {
		if err := loadDataset(session, d); err != nil {
			reportError(source{name: "--data " + d.name}, err, stderr)
			return nil, exitEvalError
		}
	}
	for _, o := range opts.overrides {
		if err := session.Override(o.name, o.expr); err != nil {
			reportError(source{name: "--set " + o.name}, err, stderr)
			return nil, exitCode(err)
		}
	}
	return session, exitOK
}

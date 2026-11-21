package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

// A preset binds a variable from the command line before any program runs:
// --data and --set (ADR 0028, ADR 0031).
type preset interface {
	flag() string
	apply(session *application.Session) error
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

// startSession applies the presets in order, the --data files before the
// --set overrides so that overrides can use the data.
func startSession(opts options, settings application.Settings, stderr io.Writer) (*application.Session, int) {
	session := application.NewSession(settings)
	for _, p := range opts.presets {
		if err := p.apply(session); err != nil {
			reportError(source{name: p.flag()}, err, stderr)
			return nil, exitCode(err)
		}
	}
	return session, exitOK
}

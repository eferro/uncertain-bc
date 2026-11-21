package cli

import "github.com/eferro/uncertain-bc/internal/application"

type override struct {
	name string
	expr string
}

// parseOverrides splits each --set at its first '=', so the expression may
// hold a comparison such as ok=a==b (ADR 0028).
func parseOverrides(values []string) ([]preset, error) {
	overrides := make([]preset, 0, len(values))
	for _, value := range values {
		name, expr, err := splitNamed("--set", "name=expr", value)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, override{name: name, expr: expr})
	}
	return overrides, nil
}

func (o override) flag() string { return "--set " + o.name }

func (o override) apply(session *application.Session) error {
	return session.Override(o.name, o.expr)
}

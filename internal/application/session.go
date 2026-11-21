// Package application holds ubc's use cases: a Session evaluates statements
// and owns the state that lives across them.
package application

import "github.com/eferro/uncertain-bc/internal/domain"

type Settings struct {
	Seed    uint64
	Samples int
}

func DefaultSettings() Settings {
	return Settings{Seed: 1, Samples: 5000}
}

type Session struct {
	env domain.Env
}

func NewSession(settings Settings) *Session {
	return &Session{env: domain.NewEnv(settings.Seed, settings.Samples)}
}

func (s *Session) Evaluate(statement string) (domain.Value, error) {
	return domain.Evaluate(statement, s.env)
}

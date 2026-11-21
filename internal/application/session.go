// Package application holds ubc's use cases: a Session evaluates statements
// and owns the state that lives across them.
package application

import (
	"slices"

	"github.com/eferro/uncertain-bc/internal/domain"
)

type Settings struct {
	Seed    uint64
	Samples int
}

func DefaultSettings() Settings {
	return Settings{Seed: 1, Samples: 5000}
}

type Variable struct {
	Name  string
	Value domain.Value
}

type Session struct {
	env   domain.Env
	names []string
}

func NewSession(settings Settings) *Session {
	return &Session{env: domain.NewEnv(settings.Seed, settings.Samples)}
}

func (s *Session) Evaluate(program string) ([]domain.Result, error) {
	results, err := domain.Evaluate(program, s.env)
	for _, result := range results {
		if result.Name != "" && !slices.Contains(s.names, result.Name) {
			s.names = append(s.names, result.Name)
		}
	}
	return results, err
}

func (s *Session) Variables() []Variable {
	variables := make([]Variable, 0, len(s.names))
	for _, name := range s.names {
		value, _ := s.env.Variable(name)
		variables = append(variables, Variable{Name: name, Value: value})
	}
	return variables
}

func (s *Session) Reseed(seed uint64) {
	s.env.Reseed(seed)
}

func (s *Session) SetSamples(samples int) {
	s.env = s.env.WithSamples(samples)
	s.names = nil
}

// Package application holds ubc's use cases: a Session evaluates statements
// and owns the state that lives across them.
package application

import (
	"slices"
	"strings"

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
	env        domain.Env
	names      []string
	transcript []string
}

func NewSession(settings Settings) *Session {
	return &Session{env: domain.NewEnv(settings.Seed, settings.Samples)}
}

func (s *Session) Evaluate(program string) ([]domain.Result, error) {
	results, err := domain.Evaluate(program, s.env)
	for _, result := range results {
		if result.Name != "" {
			s.remember(result.Name)
		}
	}
	if err == nil && strings.TrimSpace(program) != "" {
		s.transcript = append(s.transcript, program)
	}
	return results, err
}

// Override fixes a variable from outside the program, as --set does: the
// programs that follow cannot reassign it (ADR 0028).
func (s *Session) Override(name, expr string) error {
	if err := domain.Override(name, expr, s.env); err != nil {
		return err
	}
	s.remember(name)
	return nil
}

// OverrideData fixes a variable to observed values, as --data does: they
// are resampled and the programs that follow cannot reassign it (ADR 0031).
func (s *Session) OverrideData(name string, values []float64) error {
	if err := domain.OverrideData(name, values, s.env); err != nil {
		return err
	}
	s.remember(name)
	return nil
}

// Assert checks a condition against the session without changing it, as
// --assert does after the program (ADR 0029).
func (s *Session) Assert(expr string) (domain.Check, error) {
	return domain.Assert(expr, s.env)
}

func (s *Session) remember(name string) {
	if !slices.Contains(s.names, name) {
		s.names = append(s.names, name)
	}
}

// Transcript is every program that evaluated without errors, in order, so a
// session can be saved and replayed.
func (s *Session) Transcript() []string {
	return s.transcript
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
	s.transcript = nil
}

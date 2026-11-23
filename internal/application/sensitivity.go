package application

import (
	"cmp"
	"errors"
	"math"
	"slices"
	"strings"

	"github.com/eferro/uncertain-bc/internal/domain"
)

var (
	ErrCertainTarget = errors.New("needs an uncertain value, got a plain number")
	ErrNoDrivers     = errors.New("needs other uncertain variables to compare with")
)

// A Driver is a variable and the rank correlation of its samples with the
// samples of a result: how much of that result's spread moves with it.
type Driver struct {
	Name        string
	Correlation float64
}

// Sensitivity ranks the uncertain variables by how strongly they move expr,
// strongest first, without changing the session (ADR 0033). Sample i of
// every variable belongs to the same scenario, so correlating them pairs
// each input with the result it produced.
func (s *Session) Sensitivity(expr string) ([]Driver, error) {
	value, err := domain.Inspect(expr, s.env)
	if err != nil {
		return nil, err
	}
	target, ok := domain.Samples(value)
	if !ok {
		return nil, ErrCertainTarget
	}
	drivers := s.driversOf(target, strings.TrimSpace(expr))
	if len(drivers) == 0 {
		return nil, ErrNoDrivers
	}
	slices.SortStableFunc(drivers, func(a, b Driver) int {
		return cmp.Compare(math.Abs(b.Correlation), math.Abs(a.Correlation))
	})
	return drivers, nil
}

func (s *Session) driversOf(target domain.SampleSet, self string) []Driver {
	var drivers []Driver
	for _, v := range s.Variables() {
		samples, ok := domain.Samples(v.Value)
		if !ok || v.Name == self {
			continue
		}
		if correlation, defined := domain.RankCorrelation(samples, target); defined {
			drivers = append(drivers, Driver{Name: v.Name, Correlation: correlation})
		}
	}
	return drivers
}

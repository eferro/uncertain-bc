package domain

import (
	"fmt"
	"math"
)

// paramCheck returns why the parameters of a distribution are invalid, or "" when they are fine.
type paramCheck func(p []float64) string

func positive(i int, name string) paramCheck {
	return func(p []float64) string {
		if p[i] <= 0 {
			return fmt.Sprintf("%s must be positive, got %g", name, p[i])
		}
		return ""
	}
}

func all(checks ...paramCheck) paramCheck {
	return func(p []float64) string {
		for _, check := range checks {
			if problem := check(p); problem != "" {
				return problem
			}
		}
		return ""
	}
}

func ordered(lo, hi int) paramCheck {
	return func(p []float64) string {
		if p[lo] >= p[hi] {
			return fmt.Sprintf("lo must be below hi, got %g and %g", p[lo], p[hi])
		}
		return ""
	}
}

var triangle = all(ordered(0, 2), modeInside)

func modeInside(p []float64) string {
	if p[1] < p[0] || p[1] > p[2] {
		return fmt.Sprintf("mode must be between lo and hi, got %g, %g and %g", p[0], p[1], p[2])
	}
	return ""
}

func finite(i int, name string) paramCheck {
	return func(p []float64) string {
		if math.IsInf(p[i], 0) || math.IsNaN(p[i]) {
			return fmt.Sprintf("%s must be finite, got %g", name, p[i])
		}
		return ""
	}
}

func positiveWhole(i int, name string) paramCheck {
	return func(p []float64) string {
		if p[i] < 1 || p[i] != math.Trunc(p[i]) || math.IsInf(p[i], 0) {
			return fmt.Sprintf("%s must be a positive whole number, got %g", name, p[i])
		}
		return ""
	}
}

func probability(i int, name string) paramCheck {
	return func(p []float64) string {
		if !(p[i] >= 0 && p[i] <= 1) {
			return fmt.Sprintf("%s must be between 0 and 1, got %g", name, p[i])
		}
		return ""
	}
}

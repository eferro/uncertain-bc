package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ScalesANumber_WhenItHasASuffix(t *testing.T) {
	cases := map[string]float64{
		"1k":       1e3,
		"2.5M":     2.5e6,
		"3B":       3e9,
		"1T":       1e12,
		"5%":       0.05,
		"1e3k":     1e6,
		"1k * 3":   3000,
		"5% * 200": 10,
		"-2k":      -2000,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(want), evaluate(t, program))
		})
	}
}

func TestEvaluate_AcceptsSuffixesInRanges(t *testing.T) {
	s := evaluateSamples(t, "1k to 5k")

	assertNear(t, 1000, s.Quantile(0.05), 4000)
	assertNear(t, 5000, s.Quantile(0.95), 4000)
}

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_AddsNumbers(t *testing.T) {
	cases := map[string]float64{
		"1 + 2":       3,
		"1+2":         3,
		"1 + 2 + 3.5": 6.5,
		"  7  ":       7,
		"1\t+\t2":     3,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(want), evaluate(t, program))
		})
	}
}

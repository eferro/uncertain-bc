package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_AppliesArithmeticWithPrecedenceAndAssociativity(t *testing.T) {
	cases := map[string]float64{
		"1 + 2":       3,
		"1 + 2 + 3.5": 6.5,
		"7 - 2":       5,
		"2 * 3":       6,
		"7 / 2":       3.5,
		"2 ^ 3":       8,
		"1 + 2 * 3":   7,
		"(1 + 2) * 3": 9,
		"10 - 4 - 3":  3,
		"8 / 4 / 2":   1,
		"2 ^ 3 ^ 2":   512,
		"-3":          -3,
		"--3":         3,
		"-2 ^ 2":      -4,
		"2 ^ -1":      0.5,
		"2 * -3":      -6,
		"((4))":       4,
		"1 / 0":       math.Inf(1),
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(want), evaluate(t, program))
		})
	}
}

func TestEvaluate_IgnoresSpacesAndTabsAroundTokens(t *testing.T) {
	for _, program := range []string{"1+2", "  1 + 2  ", "1\t+\t2"} {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(3), evaluate(t, program))
		})
	}
}

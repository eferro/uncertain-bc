package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSampleArithmetic_ScalesEverySample_WhenCombinedWithAScalar(t *testing.T) {
	for _, program := range []string{"(8 to 25) * 12", "12 * (8 to 25)"} {
		t.Run(program, func(t *testing.T) {
			s := evaluateSamples(t, program)

			assertNear(t, 8*12, s.Quantile(0.05), 17*12)
			assertNear(t, math.Sqrt(8*25)*12, s.Quantile(0.5), 17*12)
			assertNear(t, 25*12, s.Quantile(0.95), 17*12)
		})
	}
}

func TestSampleArithmetic_NegatesEverySample(t *testing.T) {
	s := evaluateSamples(t, "-(10 to 20)")

	assertNear(t, -math.Sqrt(10*20), s.Quantile(0.5), 10)
}

func TestSampleArithmetic_TreatsTwoRangesAsIndependentDraws(t *testing.T) {
	s := evaluateSamples(t, "(1 to 2) - (1 to 2)")

	assertNear(t, 0, s.Quantile(0.5), 1)
	assert.Greater(t, s.Quantile(0.95)-s.Quantile(0.05), 0.5)
}

func TestRange_Fails_WhenABoundIsUncertain(t *testing.T) {
	cases := map[string]int{
		"(1 to 2) to 3": 10,
		"1 to (2 to 3)": 3,
	}
	for program, col := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, domain.EvalError{Line: 1, Col: col, Msg: "range bounds must be plain numbers"}, *evalErr)
		})
	}
}

func TestEvaluate_ReportsTheInnermostError_WhenASubexpressionFails(t *testing.T) {
	cases := map[string]int{
		"-(10 to 5)":      6,
		"1 + (10 to 5)":   9,
		"(10 to 5) * 2":   5,
		"(10 to 5) to 20": 5,
		"1 to (10 to 5)":  10,
	}
	for program, col := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, col, evalErr.Col)
		})
	}
}

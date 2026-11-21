package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestMathFunctions_ApplyToAPlainNumber(t *testing.T) {
	cases := map[string]float64{
		"sqrt(16)":     4,
		"exp(0)":       1,
		"log(1)":       0,
		"log10(1000)":  3,
		"abs(-3)":      3,
		"round(2.5)":   3,
		"round(-2.5)":  -3,
		"round(2.4)":   2,
		"floor(2.7)":   2,
		"floor(-2.1)":  -3,
		"ceil(2.1)":    3,
		"exp(log(42))": 42,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.InDelta(t, want, float64(evaluate(t, program).(domain.Scalar)), 1e-9)
		})
	}
}

func TestMathFunctions_ApplyToEverySample(t *testing.T) {
	s := evaluateSamples(t, "x = 1 to 100; log10(x)")

	assertNear(t, 0, s.Quantile(0.05), 2)
	assertNear(t, 2, s.Quantile(0.95), 2)
}

func TestMathFunctions_TakeExactlyOneArgument(t *testing.T) {
	for program, message := range map[string]string{
		"sqrt()":    "sqrt expects 1 argument, got 0",
		"abs(1, 2)": "abs expects 1 argument, got 2",
	} {
		_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

		var evalErr *domain.EvalError
		require.ErrorAs(t, err, &evalErr)
		assert.Equal(t, message, evalErr.Msg)
	}
}

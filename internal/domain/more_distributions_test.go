package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestExponential_HasMedianLn2OverRate(t *testing.T) {
	s := evaluateSamples(t, "exponential(2)")

	assertNear(t, math.Ln2/2, s.Quantile(0.5), 1.5)
	assertNear(t, -math.Log(0.05)/2, s.Quantile(0.95), 1.5)
	assert.InDelta(t, 0.5, s.Mean(), 0.02)
}

func TestGamma_HasMeanShapeTimesScale(t *testing.T) {
	assert.InDelta(t, 6, evaluateSamples(t, "gamma(2, 3)").Mean(), 0.15)
}

func TestGamma_IsExponential_WhenTheShapeIsOne(t *testing.T) {
	assertNear(t, 2*math.Ln2, evaluateSamples(t, "gamma(1, 2)").Quantile(0.5), 6)
}

func TestPert_StaysWithinItsBoundsAndHasThePertMean(t *testing.T) {
	s := evaluateSamples(t, "pert(0, 1, 10)")

	assert.GreaterOrEqual(t, s.Quantile(0), 0.0)
	assert.LessOrEqual(t, s.Quantile(1), 10.0)
	assert.InDelta(t, (0+4*1+10)/6.0, s.Mean(), 0.05)
}

func TestPert_IsSymmetric_WhenTheModeIsCentred(t *testing.T) {
	s := evaluateSamples(t, "pert(100, 105, 110)")

	assertNear(t, 105, s.Quantile(0.5), 10)
	assertNear(t, 105-s.Quantile(0.1), s.Quantile(0.9)-105, 10)
}

func TestMoreDistributions_RejectInvalidParameters(t *testing.T) {
	cases := map[string]string{
		"exponential(0)": "exponential: rate must be positive, got 0",
		"gamma(0, 1)":    "gamma: shape must be positive, got 0",
		"gamma(1, -2)":   "gamma: scale must be positive, got -2",
		"pert(1, 0, 10)": "pert: mode must be between lo and hi, got 1, 0 and 10",
		"pert(5, 5, 5)":  "pert: lo must be below hi, got 5 and 5",
		"exponential()":  "exponential expects 1 argument, got 0",
	}
	for program, message := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, message, evalErr.Msg)
		})
	}
}

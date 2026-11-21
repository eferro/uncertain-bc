package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestBeta_IsUniform_WhenBothShapesAreOne(t *testing.T) {
	s := evaluateSamples(t, "beta(1, 1)")

	for _, q := range []float64{0.05, 0.25, 0.75, 0.95} {
		assert.InDelta(t, q, s.Quantile(q), 0.01, "q=%v", q)
	}
}

func TestBeta_FollowsTheArcsineLaw_WhenBothShapesAreOneHalf(t *testing.T) {
	s := evaluateSamples(t, "beta(0.5, 0.5)")

	for _, q := range []float64{0.05, 0.25, 0.5, 0.95} {
		assertNear(t, math.Pow(math.Sin(math.Pi*q/2), 2), s.Quantile(q), 1)
	}
}

func TestBeta_HasMeanAOverAPlusB(t *testing.T) {
	s := evaluateSamples(t, "beta(2, 8)")

	assert.InDelta(t, 0.2, sampleMean(s), 0.005)
	assert.Greater(t, s.Quantile(0), 0.0)
	assert.Less(t, s.Quantile(1), 1.0)
}

func TestBeta_RejectsNonPositiveShapes(t *testing.T) {
	cases := map[string]string{
		"beta(0, 1)":  "beta: a must be positive, got 0",
		"beta(1, -1)": "beta: b must be positive, got -1",
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

func sampleMean(s domain.SampleSet) float64 {
	total := 0.0
	for _, x := range s {
		total += x
	}
	return total / float64(len(s))
}

func TestBeta_MatchesClosedFormQuantilesClosely_WithManySamples(t *testing.T) {
	cases := []struct {
		program  string
		quantile func(q float64) float64
	}{
		{"beta(1, 1)", func(q float64) float64 { return q }},
		{"beta(0.5, 0.5)", func(q float64) float64 { return math.Pow(math.Sin(math.Pi*q/2), 2) }},
		{"beta(5, 1)", func(q float64) float64 { return math.Pow(q, 1.0/5) }},
		{"beta(1, 5)", func(q float64) float64 { return 1 - math.Pow(1-q, 1.0/5) }},
		{"beta(0.3, 1)", func(q float64) float64 { return math.Pow(q, 1/0.3) }},
	}
	for _, c := range cases {
		t.Run(c.program, func(t *testing.T) {
			results, err := domain.Evaluate(c.program, domain.NewEnv(1, 200_000))
			require.NoError(t, err)
			s := results[0].Value.(domain.SampleSet)

			for _, q := range []float64{0.02, 0.1, 0.5, 0.9, 0.98} {
				assert.InDelta(t, c.quantile(q), s.Quantile(q), 0.004, "q=%v", q)
			}
		})
	}
}

package domain_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestPoisson_MatchesTheExactDistribution_ForSmallAndLargeMeans(t *testing.T) {
	for _, lambda := range []float64{0.5, 4, 15.9, 16, 50, 1000} {
		t.Run(fmt.Sprint(lambda), func(t *testing.T) {
			s := manySamples(t, fmt.Sprintf("poisson(%g)", lambda))

			assertWholeNumbersFrom(t, 0, s)
			assertMatchesCDF(t, s, func(k float64) float64 { return poissonCDF(lambda, k) })
		})
	}
}

func TestBinomial_MatchesTheExactDistribution_ForFewAndManyTrials(t *testing.T) {
	cases := []struct{ n, p float64 }{{1, 0.5}, {10, 0.3}, {15, 0.9}, {16, 0.5}, {100, 0.05}, {1000, 0.3}, {5000, 0.99}}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%g,%g", c.n, c.p), func(t *testing.T) {
			s := manySamples(t, fmt.Sprintf("binomial(%g, %g)", c.n, c.p))

			assertWholeNumbersFrom(t, 0, s)
			assert.LessOrEqual(t, s.Quantile(1), c.n)
			assertMatchesCDF(t, s, func(k float64) float64 { return binomialCDF(c.n, c.p, k) })
		})
	}
}

func TestBinomial_IsCertain_WhenTheProbabilityIsZeroOrOne(t *testing.T) {
	assert.Equal(t, domain.Scalar(0), evaluate(t, "x = binomial(40, 0); max(x)"))
	assert.Equal(t, domain.Scalar(40), evaluate(t, "x = binomial(40, 1); min(x)"))
	assert.Equal(t, domain.Scalar(3), evaluate(t, "x = binomial(3, 1); min(x)"))
}

func TestDiscreteDistributions_RejectInvalidParameters(t *testing.T) {
	cases := map[string]string{
		"poisson(0)":           "poisson: lambda must be positive, got 0",
		"poisson(1 / 0)":       "poisson: lambda must be finite, got +Inf",
		"poisson(0 / 0)":       "poisson: lambda must be finite, got NaN",
		"binomial(0, 0.5)":     "binomial: n must be a positive whole number, got 0",
		"binomial(2.5, 0.5)":   "binomial: n must be a positive whole number, got 2.5",
		"binomial(1 / 0, 0.5)": "binomial: n must be a positive whole number, got +Inf",
		"binomial(10, -0.1)":   "binomial: p must be between 0 and 1, got -0.1",
		"binomial(10, 1.5)":    "binomial: p must be between 0 and 1, got 1.5",
		"binomial(10, 0 / 0)":  "binomial: p must be between 0 and 1, got NaN",
		"binomial(10)":         "binomial expects 2 arguments, got 1",
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

const cdfTolerance = 0.006

func manySamples(t *testing.T, program string) domain.SampleSet {
	t.Helper()
	results, err := domain.Evaluate(program, domain.NewEnv(1, 100_000))
	require.NoError(t, err)
	require.IsType(t, domain.SampleSet{}, results[0].Value)
	return results[0].Value.(domain.SampleSet)
}

func assertWholeNumbersFrom(t *testing.T, lowest float64, s domain.SampleSet) {
	t.Helper()
	for _, x := range s {
		require.Equal(t, math.Trunc(x), x)
	}
	assert.GreaterOrEqual(t, s.Quantile(0), lowest)
}

func assertMatchesCDF(t *testing.T, s domain.SampleSet, cdf func(k float64) float64) {
	t.Helper()
	for _, q := range []float64{0.05, 0.25, 0.5, 0.75, 0.95} {
		k := s.Quantile(q)
		assert.InDelta(t, cdf(k), fractionAtMost(s, k), cdfTolerance, "k=%v", k)
		assert.InDelta(t, cdf(k-1), fractionAtMost(s, k-1), cdfTolerance, "k=%v", k-1)
	}
}

func fractionAtMost(s domain.SampleSet, k float64) float64 {
	count := 0
	for _, x := range s {
		if x <= k {
			count++
		}
	}
	return float64(count) / float64(len(s))
}

func poissonCDF(lambda, k float64) float64 {
	total := 0.0
	for i := 0.0; i <= k; i++ {
		lgamma, _ := math.Lgamma(i + 1)
		total += math.Exp(i*math.Log(lambda) - lambda - lgamma)
	}
	return total
}

func binomialCDF(n, p, k float64) float64 {
	total := 0.0
	for i := 0.0; i <= k; i++ {
		total += math.Exp(logChoose(n, i) + i*math.Log(p) + (n-i)*math.Log1p(-p))
	}
	return total
}

func logChoose(n, k float64) float64 {
	a, _ := math.Lgamma(n + 1)
	b, _ := math.Lgamma(k + 1)
	c, _ := math.Lgamma(n - k + 1)
	return a - b - c
}

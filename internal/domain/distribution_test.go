package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestNormal_DrawsFromAMeanAndStandardDeviation(t *testing.T) {
	s := evaluateSamples(t, "normal(100, 15)")

	assertNear(t, 100-1.645*15, s.Quantile(0.05), 2*1.645*15)
	assertNear(t, 100, s.Quantile(0.5), 2*1.645*15)
	assertNear(t, 100+1.645*15, s.Quantile(0.95), 2*1.645*15)
}

func TestLognormal_DrawsFromTheMeanAndSigmaOfTheLogarithm(t *testing.T) {
	s := evaluateSamples(t, "lognormal(-2, 1)")

	assertNear(t, -2-1.645, math.Log(s.Quantile(0.05)), 2*1.645)
	assertNear(t, -2, math.Log(s.Quantile(0.5)), 2*1.645)
	assertNear(t, -2+1.645, math.Log(s.Quantile(0.95)), 2*1.645)
}

func TestUniform_DrawsEvenlyBetweenItsBounds(t *testing.T) {
	s := evaluateSamples(t, "uniform(1, 6)")

	assert.GreaterOrEqual(t, s.Quantile(0), 1.0)
	assert.LessOrEqual(t, s.Quantile(1), 6.0)
	assertNear(t, 1.25, s.Quantile(0.05), 5)
	assertNear(t, 3.5, s.Quantile(0.5), 5)
}

func TestTriangular_PeaksAtItsMode(t *testing.T) {
	s := evaluateSamples(t, "triangular(1, 3, 10)")

	assert.GreaterOrEqual(t, s.Quantile(0), 1.0)
	assert.LessOrEqual(t, s.Quantile(1), 10.0)
	assertNear(t, 10-math.Sqrt(0.5*9*7), s.Quantile(0.5), 9)
	assertNear(t, 1+math.Sqrt(0.1*9*2), s.Quantile(0.1), 9)
}

func TestTriangular_InvertsBothSidesOfTheMode_WhenItStartsFarFromZero(t *testing.T) {
	s := evaluateSamples(t, "triangular(100, 101, 110)")

	assertNear(t, 100+math.Sqrt(0.05*10*1), s.Quantile(0.05), 10)
	assertNear(t, 110-math.Sqrt(0.1*10*9), s.Quantile(0.9), 10)
}

func TestDistributions_DrawFreshSamplesOnEveryUse(t *testing.T) {
	s := evaluateSamples(t, "normal(0, 1) - normal(0, 1)")

	assert.Greater(t, s.Quantile(0.95)-s.Quantile(0.05), 1.0)
}

func TestDistributions_RejectInvalidParameters(t *testing.T) {
	cases := map[string]string{
		"normal(1)":                "normal expects 2 arguments, got 1",
		"normal(0, -1)":            "normal: sd must be positive, got -1",
		"normal(0, 0)":             "normal: sd must be positive, got 0",
		"lognormal(1, 0)":          "lognormal: sigma must be positive, got 0",
		"uniform(6, 1)":            "uniform: lo must be below hi, got 6 and 1",
		"uniform(1, 1)":            "uniform: lo must be below hi, got 1 and 1",
		"triangular(1, 11, 10)":    "triangular: mode must be between lo and hi, got 1, 11 and 10",
		"triangular(1, 0, 10)":     "triangular: mode must be between lo and hi, got 1, 0 and 10",
		"triangular(5, 5, 5)":      "triangular: lo must be below hi, got 5 and 5",
		"normal(0 to 1, 1)":        "normal expects plain numbers",
		"uniform(1, normal(5, 1))": "uniform expects plain numbers",
	}
	for program, message := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, domain.EvalError{Line: 1, Col: 1, Msg: message}, *evalErr)
		})
	}
}

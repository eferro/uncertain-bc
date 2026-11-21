package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSummaries_AreTheIdentity_OnAPlainNumber(t *testing.T) {
	cases := map[string]float64{
		"mean(5)":     5,
		"median(5)":   5,
		"stdev(5)":    0,
		"min(5)":      5,
		"max(5)":      5,
		"p(5, 0.9)":   5,
		"mean(2 + 3)": 5,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(want), evaluate(t, program))
		})
	}
}

func TestSummaries_DescribeASampleSetWithAPlainNumber(t *testing.T) {
	cases := []struct {
		summary string
		want    float64
	}{
		{"mean(x)", 5},
		{"median(x)", 5},
		{"stdev(x)", 10 / math.Sqrt(12)},
		{"min(x)", 0},
		{"max(x)", 10},
		{"p(x, 0.9)", 9},
		{"p(x, 0)", 0},
		{"p(x, 1)", 10},
	}
	for _, c := range cases {
		t.Run(c.summary, func(t *testing.T) {
			got := evaluate(t, "x = uniform(0, 10); "+c.summary)

			require.IsType(t, domain.Scalar(0), got)
			assertNear(t, c.want, float64(got.(domain.Scalar)), 10)
		})
	}
}

func TestMinAndMax_CompareEverySample_WhenGivenSeveralArguments(t *testing.T) {
	assert.Equal(t, domain.Scalar(1), evaluate(t, "min(3, 2, 1)"))
	assert.Equal(t, domain.Scalar(3), evaluate(t, "max(1, 2, 3)"))
	assert.Equal(t, 5.0, evaluateSamples(t, "x = 0 to 10; max(x, 5)").Quantile(0))
	assert.Equal(t, 5.0, evaluateSamples(t, "x = 0 to 10; min(x, 5)").Quantile(1))
}

func TestSummaries_RejectWrongArguments(t *testing.T) {
	cases := map[string]string{
		"mean()":       "mean expects 1 argument, got 0",
		"stdev(1, 2)":  "stdev expects 1 argument, got 2",
		"median(1, 2)": "median expects 1 argument, got 2",
		"min()":        "min expects at least 1 argument, got 0",
		"max()":        "max expects at least 1 argument, got 0",
		"p(1)":         "p expects 2 arguments, got 1",
		"p(1, 1.5)":    "p: q must be between 0 and 1, got 1.5",
		"p(1, -0.1)":   "p: q must be between 0 and 1, got -0.1",
		"p(1, 0 to 1)": "p: q must be a plain number",
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

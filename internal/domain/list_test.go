package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestList_ResamplesItsValuesWithReplacement(t *testing.T) {
	s := evaluateSamples(t, "[3, 5, 8, 13]")

	counts := map[float64]int{}
	for _, x := range s {
		counts[x]++
	}
	require.Len(t, counts, 4)
	for _, value := range []float64{3, 5, 8, 13} {
		assert.InDelta(t, samples/4, counts[value], 120, "value %v", value)
	}
}

func TestList_AcceptsExpressionsThatAreNumbers(t *testing.T) {
	s := evaluateSamples(t, "[1 + 1, -3]")

	assert.Equal(t, -3.0, s.Quantile(0))
	assert.Equal(t, 2.0, s.Quantile(1))
}

func TestList_BehavesLikeAnyOtherSampleSet(t *testing.T) {
	assert.Equal(t, domain.Scalar(0), evaluate(t, "x = [1, 2]; max(x - x)"))
	assert.InDelta(t, 2.5, float64(evaluate(t, "mean([1, 2, 3, 4])").(domain.Scalar)), 0.05)
	assert.Equal(t, 20.0, evaluateSamples(t, "[1, 2] * 10").Quantile(1))
}

func TestList_WithOneValueIsCertain(t *testing.T) {
	s := evaluateSamples(t, "[7]")

	assert.Equal(t, 7.0, s.Quantile(0))
	assert.Equal(t, 7.0, s.Quantile(1))
}

func TestList_RejectsWhatItCannotResample(t *testing.T) {
	cases := map[string]domain.EvalError{
		"[]":              {Line: 1, Col: 1, Msg: "an empty list has no values to resample"},
		"1 + [1, 2 to 3]": {Line: 1, Col: 5, Msg: "list elements must be plain numbers"},
		"[1, y]":          {Line: 1, Col: 5, Msg: "unknown variable 'y'"},
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, want, *evalErr)
		})
	}
}

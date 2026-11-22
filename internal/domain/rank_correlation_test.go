package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestRankCorrelation_IsOne_WhenOneGrowsWithTheOther(t *testing.T) {
	got, ok := domain.RankCorrelation(domain.SampleSet{1, 3, 2, 4}, domain.SampleSet{10, 1000, 100, 10000})

	require.True(t, ok)
	assert.InDelta(t, 1.0, got, 1e-12)
}

func TestRankCorrelation_IsMinusOne_WhenOneShrinksAsTheOtherGrows(t *testing.T) {
	got, ok := domain.RankCorrelation(domain.SampleSet{1, 2, 3}, domain.SampleSet{-1, -8, -27})

	require.True(t, ok)
	assert.InDelta(t, -1.0, got, 1e-12)
}

func TestRankCorrelation_GivesTiedSamplesTheirAverageRank(t *testing.T) {
	got, ok := domain.RankCorrelation(domain.SampleSet{1, 2, 3, 4, 5}, domain.SampleSet{5, 6, 7, 8, 7})

	require.True(t, ok)
	assert.InDelta(t, 8/math.Sqrt(95), got, 1e-12)
}

func TestRankCorrelation_IsNearZero_ForIndependentDraws(t *testing.T) {
	x := evaluateSamples(t, "normal(0, 1)")
	y := evaluateSamples(t, "a = normal(0, 1); normal(0, 1)")

	got, ok := domain.RankCorrelation(x, y)

	require.True(t, ok)
	assert.InDelta(t, 0, got, 0.05)
}

func TestRankCorrelation_SkipsSamplesThatAreNotFiniteOnEitherSide(t *testing.T) {
	x := domain.SampleSet{1, 2, math.NaN(), 3, 4}
	y := domain.SampleSet{1, 2, 0, math.Inf(1), 3}

	got, ok := domain.RankCorrelation(x, y)

	require.True(t, ok)
	assert.InDelta(t, 1.0, got, 1e-12)
}

func TestRankCorrelation_IsUndefined_WhenEitherSideDoesNotVary(t *testing.T) {
	cases := []struct{ x, y domain.SampleSet }{
		{domain.SampleSet{1, 2, 3}, domain.SampleSet{7, 7, 7}},
		{domain.SampleSet{7, 7, 7}, domain.SampleSet{1, 2, 3}},
		{domain.SampleSet{1, math.NaN(), math.NaN()}, domain.SampleSet{1, 2, 3}},
		{domain.SampleSet{}, domain.SampleSet{}},
	}
	for _, c := range cases {
		_, ok := domain.RankCorrelation(c.x, c.y)

		assert.False(t, ok, "%v with %v", c.x, c.y)
	}
}

func TestInspect_EvaluatesAnExpression_WithoutBindingLast(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("x = 2", env)
	require.NoError(t, err)

	got, err := domain.Inspect("x * 3", env)

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(6), got)
	_, bound := env.Variable("last")
	assert.False(t, bound)
}

func TestInspect_RejectsMoreThanOneExpression(t *testing.T) {
	_, err := domain.Inspect("1; 2", domain.NewEnv(1, samples))

	assert.Error(t, err)
}

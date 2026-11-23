package domain_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ReturnsAProbability_WhenATopLevelStatementIsAComparison(t *testing.T) {
	cases := map[string]float64{
		"2 > 1":      1,
		"2 < 1":      0,
		"1 <= 1":     1,
		"1 >= 2":     0,
		"1 == 1":     1,
		"1 != 1":     0,
		"1 < 1":      0,
		"1 > 1":      0,
		"1 >= 1":     1,
		"1 + 1 == 2": 1,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Probability(want), evaluate(t, program))
		})
	}
}

func TestEvaluate_ComparesEverySample(t *testing.T) {
	event := evaluateEvent(t, "x = 0 to 10; x > 5")

	assert.InDelta(t, 0.5, event.Probability(), 0.03)
}

func TestEvaluate_KeepsTheSamplesOfAComparison_WhenItIsAssigned(t *testing.T) {
	event := evaluateEvent(t, "x = 0 to 10; hit = x > 5")

	assert.Len(t, event, samples)
	assert.Equal(t, 0.0, slices.Min(event))
	assert.Equal(t, 1.0, slices.Max(event))
}

func TestEvaluate_ComputesWithAnAssignedComparison_SampleBySample(t *testing.T) {
	p := evaluate(t, "hit = uniform(0, 1) < 0.1; loss = hit * (50k to 500k); P(loss > 0)")

	assert.InDelta(t, 0.1, float64(p.(domain.Probability)), 0.01)
}

func TestEvaluate_KeepsAPlainComparisonAsAProbability_WhenItIsAssigned(t *testing.T) {
	assert.Equal(t, domain.Probability(1), evaluate(t, "ok = 3 > 2; ok"))
}

func TestEvaluate_ReturnsAnEvent_FromAFunctionWhoseBodyIsAComparison(t *testing.T) {
	event := evaluateEvent(t, "above(x, t) = x > t; above(0 to 10, 5)")

	assert.InDelta(t, 0.5, event.Probability(), 0.03)
}

func TestEvaluate_BindsTheEventToLast_WhenAComparisonIsAnExpression(t *testing.T) {
	s := evaluateSamples(t, "x = 0 to 10; x > 5; last * 2")

	assert.Equal(t, 0.0, s.Quantile(0))
	assert.Equal(t, 2.0, s.Quantile(1))
}

func TestP_ReturnsTheProbability_OfAnAssignedComparison(t *testing.T) {
	p := evaluate(t, "hit = 0 to 10 > 5; P(hit)")

	assert.InDelta(t, 0.5, float64(p.(domain.Probability)), 0.03)
}

func TestSummaries_UseTheSamples_OfAnEvent(t *testing.T) {
	assert.Equal(t, domain.Scalar(1), evaluate(t, "hit = 0 to 10 > 5; max(hit)"))
	assert.Equal(t, domain.Scalar(0), evaluate(t, "hit = 0 to 10 > 5; p(hit, 0.01)"))
}

func TestEvaluate_KeepsAComparisonAsZerosAndOnes_InsideArithmetic(t *testing.T) {
	s := evaluateSamples(t, "x = 0 to 10; (x > 5) * 1")

	assert.Equal(t, 0.0, s.Quantile(0))
	assert.Equal(t, 1.0, s.Quantile(1))
}

func TestP_ReturnsTheShareOfSamplesWhereTheConditionHolds(t *testing.T) {
	p := evaluate(t, "P(0 to 10 > 5)")

	require.IsType(t, domain.Probability(0), p)
	assert.InDelta(t, 0.5, float64(p.(domain.Probability)), 0.03)
}

func TestP_IsCertain_WhenTheConditionIsAPlainComparison(t *testing.T) {
	assert.Equal(t, domain.Probability(1), evaluate(t, "P(2 > 1)"))
}

func TestP_CanBeUsedInArithmetic(t *testing.T) {
	assert.Equal(t, domain.Scalar(100), evaluate(t, "P(2 > 1) * 100"))
}

func TestEvaluate_ReportsMisusedFunctions(t *testing.T) {
	cases := []struct {
		program string
		want    domain.EvalError
	}{
		{"P()", domain.EvalError{Line: 1, Col: 1, Msg: "P expects 1 argument, got 0"}},
		{"1 + P(2 > 1, 3)", domain.EvalError{Line: 1, Col: 5, Msg: "P expects 1 argument, got 2"}},
		{"P(5)", domain.EvalError{Line: 1, Col: 1, Msg: "P expects a comparison"}},
		{"P(0 to 10)", domain.EvalError{Line: 1, Col: 1, Msg: "P expects a comparison"}},
		{"x = 0 to 10; P((x > 5) * 2)", domain.EvalError{Line: 1, Col: 14, Msg: "P expects a comparison"}},
		{"x = 0 to 10; P((x <= 5) * 2)", domain.EvalError{Line: 1, Col: 14, Msg: "P expects a comparison"}},
		{"foo(1)", domain.EvalError{Line: 1, Col: 1, Msg: "unknown function 'foo'"}},
		{"P(y > 1)", domain.EvalError{Line: 1, Col: 3, Msg: "unknown variable 'y'"}},
		{"1 < y", domain.EvalError{Line: 1, Col: 5, Msg: "unknown variable 'y'"}},
	}
	for _, c := range cases {
		t.Run(c.program, func(t *testing.T) {
			_, err := domain.Evaluate(c.program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, c.want, *evalErr)
		})
	}
}

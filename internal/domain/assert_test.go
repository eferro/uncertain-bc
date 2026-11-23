package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func check(t *testing.T, env domain.Env, source string) domain.Check {
	t.Helper()
	result, err := domain.Assert(source, env)
	require.NoError(t, err)
	return result
}

func TestAssert_Holds_WhenTheConditionIsCertain(t *testing.T) {
	cases := map[string]bool{
		"1 < 2":                true,
		"2 < 1":                false,
		"1 <= 1":               true,
		"1 == 1":               true,
		"1 != 1":               false,
		"(1 < 2)":              true,
		"1":                    true,
		"0":                    false,
		"P(uniform(0, 1) > 0)": true,
	}
	for source, holds := range cases {
		t.Run(source, func(t *testing.T) {
			assert.Equal(t, holds, check(t, domain.NewEnv(1, samples), source).Holds())
		})
	}
}

func TestAssert_DoesNotHold_WhenTheConditionIsOnlyLikely(t *testing.T) {
	result := check(t, domain.NewEnv(1, samples), "normal(0, 1) < 3")

	assert.False(t, result.Holds())
	assert.InDelta(t, 0.99, float64(result.Value.(domain.Probability)), 0.01)
}

func TestAssert_ReportsTheOperands_OfAComparison(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "fixed", "150")

	result := check(t, env, "fixed >= 2 * 100")

	assert.Equal(t, domain.Check{
		Value: domain.Probability(0), Left: domain.Scalar(150), Operator: ">=", Right: domain.Scalar(200),
	}, result)
}

func TestAssert_ComputesTheProbability_FromTheReportedOperands(t *testing.T) {
	result := check(t, domain.NewEnv(1, samples), "uniform(0, 1) < 0.5")

	left := result.Left.(domain.SampleSet)
	below := 0
	for _, x := range left {
		if x < 0.5 {
			below++
		}
	}
	assert.Equal(t, domain.Probability(float64(below)/float64(len(left))), result.Value)
}

func TestAssert_HasNoOperands_WhenItIsNotAComparison(t *testing.T) {
	result := check(t, domain.NewEnv(1, samples), "P(uniform(0, 1) > 2)")

	assert.Equal(t, domain.Check{Value: domain.Probability(0)}, result)
}

func TestAssert_ChecksAVariableBoundToAComparison_AsItsProbability(t *testing.T) {
	cases := map[string]bool{"0 to 10 > -100": true, "0 to 10 > 5": false}
	for comparison, holds := range cases {
		t.Run(comparison, func(t *testing.T) {
			env := domain.NewEnv(1, samples)
			override(t, env, "hit", comparison)

			result := check(t, env, "hit")

			require.IsType(t, domain.Probability(0), result.Value)
			assert.Equal(t, holds, result.Holds())
		})
	}
}

func TestAssert_LeavesTheSessionUnchanged(t *testing.T) {
	env := domain.NewEnv(1, samples)

	check(t, env, "1 < 2")

	_, ok := env.Variable("last")
	assert.False(t, ok)
}

func TestAssert_Fails_WhenTheValueIsNotAConditionOrAProbability(t *testing.T) {
	for _, source := range []string{"2", "0.5", "uniform(0, 1)", "(uniform(0, 1) > 0.5) * 1", "-1"} {
		t.Run(source, func(t *testing.T) {
			_, err := domain.Assert(source, domain.NewEnv(1, samples))

			assert.Equal(t, &domain.EvalError{Line: 1, Col: 1, Msg: "an assertion must be a condition or a probability"}, err)
		})
	}
}

func TestAssert_Fails_WhenAnOperandFails(t *testing.T) {
	_, err := domain.Assert("1 < nobody", domain.NewEnv(1, samples))

	assert.Equal(t, &domain.EvalError{Line: 1, Col: 5, Msg: "unknown variable 'nobody'"}, err)
}

func TestAssert_Fails_WhenTheSourceIsNotOneExpression(t *testing.T) {
	_, err := domain.Assert("1 < 2; 3", domain.NewEnv(1, samples))

	var syntaxErr *domain.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr)
}

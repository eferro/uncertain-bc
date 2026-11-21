package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_KnowsPiAndE(t *testing.T) {
	assert.Equal(t, domain.Scalar(math.Pi), evaluate(t, "pi"))
	assert.Equal(t, domain.Scalar(math.E), evaluate(t, "e"))
}

func TestEvaluate_LetsPiBeRedefined(t *testing.T) {
	assert.Equal(t, domain.Scalar(3), evaluate(t, "pi = 3; pi"))
}

func TestLast_IsTheResultOfThePreviousExpression(t *testing.T) {
	assert.Equal(t, domain.Scalar(30), evaluate(t, "1 + 2; last * 10"))
}

func TestLast_IgnoresAssignments(t *testing.T) {
	assert.Equal(t, domain.Scalar(1), evaluate(t, "1; x = 7; last"))
}

func TestLast_IsUnknown_BeforeAnyExpression(t *testing.T) {
	_, err := domain.Evaluate("x = 5; last", domain.NewEnv(1, samples))

	var evalErr *domain.EvalError
	require.ErrorAs(t, err, &evalErr)
	assert.Equal(t, "unknown variable 'last'", evalErr.Msg)
}

func TestLast_SurvivesAcrossCallsWithTheSameEnv(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("41", env)
	require.NoError(t, err)

	results, err := domain.Evaluate("last + 1", env)

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(42), results[0].Value)
}

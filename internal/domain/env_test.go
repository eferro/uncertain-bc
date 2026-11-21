package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEnv_Variable_ReturnsWhatAStatementBound(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("x = 4", env)
	require.NoError(t, err)

	value, ok := env.Variable("x")

	assert.True(t, ok)
	assert.Equal(t, domain.Scalar(4), value)
}

func TestEnv_Variable_ReportsAnUnknownName(t *testing.T) {
	_, ok := domain.NewEnv(1, samples).Variable("x")

	assert.False(t, ok)
}

func TestEnv_Reseed_ReplaysTheSameStreamAsAFreshEnv(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("1 to 2", env)
	require.NoError(t, err)

	env.Reseed(9)
	reseeded, err := domain.Evaluate("1 to 2", env)
	require.NoError(t, err)
	fresh, err := domain.Evaluate("1 to 2", domain.NewEnv(9, samples))
	require.NoError(t, err)

	assert.Equal(t, fresh, reseeded)
}

func TestEnv_WithSamples_DrawsTheNewCountAndForgetsVariables(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("x = 1", env)
	require.NoError(t, err)

	smaller := env.WithSamples(3)

	_, known := smaller.Variable("x")
	assert.False(t, known)
	assert.Len(t, evaluateIn(t, smaller, "1 to 2"), 3)
	assert.Equal(t, domain.Scalar(3.141592653589793), evaluateIn(t, smaller, "pi"))
}

func evaluateIn(t *testing.T, env domain.Env, program string) domain.Value {
	t.Helper()
	results, err := domain.Evaluate(program, env)
	require.NoError(t, err)
	return results[len(results)-1].Value
}

func TestEnv_WithSamples_KeepsUserFunctions(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("double(x) = x * 2", env)
	require.NoError(t, err)

	assert.Equal(t, domain.Scalar(8), evaluateIn(t, env.WithSamples(3), "double(4)"))
}

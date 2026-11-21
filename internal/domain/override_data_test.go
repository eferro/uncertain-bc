package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestOverrideData_BindsTheValuesAsAResampledSampleSet(t *testing.T) {
	env := domain.NewEnv(1, samples)
	require.NoError(t, domain.OverrideData("lead", []float64{3, 5, 8, 13}, env))

	value, _ := env.Variable("lead")

	assert.Equal(t, evaluateSamples(t, "[3, 5, 8, 13]"), value)
}

func TestOverrideData_SkipsLaterAssignmentsToTheName(t *testing.T) {
	env := domain.NewEnv(1, samples)
	require.NoError(t, domain.OverrideData("lead", []float64{7}, env))

	results, err := domain.Evaluate("lead = 1 to 2\nmax(lead)", env)

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(7), results[1].Value)
}

func TestOverrideData_Fails_WhenThereAreNoValues(t *testing.T) {
	env := domain.NewEnv(1, samples)

	err := domain.OverrideData("lead", nil, env)

	require.EqualError(t, err, "no values to resample")
	_, bound := env.Variable("lead")
	assert.False(t, bound)
}

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

const samples = 5000

func evaluate(t *testing.T, program string) domain.Value {
	t.Helper()
	value, err := domain.Evaluate(program, domain.NewEnv(1, samples))
	require.NoError(t, err)
	return value
}

func evaluateSamples(t *testing.T, program string) domain.SampleSet {
	t.Helper()
	value := evaluate(t, program)
	require.IsType(t, domain.SampleSet{}, value)
	return value.(domain.SampleSet)
}

func assertNear(t *testing.T, want, got, intervalWidth float64) {
	t.Helper()
	assert.InDelta(t, want, got, 0.03*intervalWidth)
}

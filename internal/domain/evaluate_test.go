package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ReturnsTheNumber_WhenTheProgramIsANumberLiteral(t *testing.T) {
	cases := map[string]float64{
		"42":     42,
		"3.14":   3.14,
		".5":     0.5,
		"1e6":    1e6,
		"2.5e-3": 0.0025,
		"1E3":    1000,
		"1e+3":   1000,
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			assert.Equal(t, domain.Scalar(want), evaluate(t, program))
		})
	}
}

func TestEvaluate_DrawsTheSameSamples_WhenTheSeedIsTheSame(t *testing.T) {
	first, err := domain.Evaluate("10 to 20", domain.NewEnv(7, samples))
	require.NoError(t, err)
	second, err := domain.Evaluate("10 to 20", domain.NewEnv(7, samples))
	require.NoError(t, err)

	assert.Equal(t, first, second)
}

func TestEvaluate_DrawsDifferentSamples_WhenTheSeedChanges(t *testing.T) {
	first, err := domain.Evaluate("10 to 20", domain.NewEnv(1, samples))
	require.NoError(t, err)
	second, err := domain.Evaluate("10 to 20", domain.NewEnv(2, samples))
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}

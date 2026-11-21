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
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			value, err := domain.Evaluate(program)

			require.NoError(t, err)
			assert.Equal(t, want, value)
		})
	}
}

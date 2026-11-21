package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ReturnsTheNumber_WhenTheProgramIsANumberLiteral(t *testing.T) {
	value, err := domain.Evaluate("42")

	require.NoError(t, err)
	assert.Equal(t, 42.0, value)
}

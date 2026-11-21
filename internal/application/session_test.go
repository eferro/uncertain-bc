package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSession_EvaluatesAStatement(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	value, err := session.Evaluate("1 + 2")

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(3), value)
}

func TestSession_DrawsTheConfiguredNumberOfSamples(t *testing.T) {
	session := application.NewSession(application.Settings{Seed: 1, Samples: 10})

	value, err := session.Evaluate("1 to 2")

	require.NoError(t, err)
	assert.Len(t, value, 10)
}

func TestSession_KeepsOneRandomStreamAcrossStatements(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	first, err := session.Evaluate("1 to 2")
	require.NoError(t, err)
	second, err := session.Evaluate("1 to 2")
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}

func TestDefaultSettings_UseSeedOneAndFiveThousandSamples(t *testing.T) {
	assert.Equal(t, application.Settings{Seed: 1, Samples: 5000}, application.DefaultSettings())
}

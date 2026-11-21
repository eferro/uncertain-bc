package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSession_EvaluatesAProgram(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	results, err := session.Evaluate("1 + 2")

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{{Line: 1, Value: domain.Scalar(3)}}, results)
}

func TestSession_DrawsTheConfiguredNumberOfSamples(t *testing.T) {
	session := application.NewSession(application.Settings{Seed: 1, Samples: 10})

	results, err := session.Evaluate("1 to 2")

	require.NoError(t, err)
	assert.Len(t, results[0].Value, 10)
}

func TestSession_KeepsOneRandomStreamAcrossPrograms(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	first, err := session.Evaluate("1 to 2")
	require.NoError(t, err)
	second, err := session.Evaluate("1 to 2")
	require.NoError(t, err)

	assert.NotEqual(t, first[0].Value, second[0].Value)
}

func TestSession_KeepsVariablesAcrossPrograms(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("x = 20")
	require.NoError(t, err)

	results, err := session.Evaluate("x + 1")

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(21), results[0].Value)
}

func TestDefaultSettings_UseSeedOneAndFiveThousandSamples(t *testing.T) {
	assert.Equal(t, application.Settings{Seed: 1, Samples: 5000}, application.DefaultSettings())
}

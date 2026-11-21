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

func TestSession_ListsTheUserVariablesInTheOrderTheyWereFirstAssigned(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("y = 2; x = 1; y = 3; 4")
	require.NoError(t, err)

	assert.Equal(t, []application.Variable{
		{Name: "y", Value: domain.Scalar(3)},
		{Name: "x", Value: domain.Scalar(1)},
	}, session.Variables())
}

func TestSession_Reseed_RestartsTheRandomStreamAndKeepsVariables(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("x = 5; 1 to 2")
	require.NoError(t, err)

	session.Reseed(2)
	results, err := session.Evaluate("x; 1 to 2")
	require.NoError(t, err)

	fresh, err := application.NewSession(application.Settings{Seed: 2, Samples: 5000}).Evaluate("1 to 2")
	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(5), results[0].Value)
	assert.Equal(t, fresh[0].Value, results[1].Value)
}

func TestSession_SetSamples_ChangesTheSampleCountAndClearsVariables(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("x = 1")
	require.NoError(t, err)

	session.SetSamples(10)

	results, err := session.Evaluate("1 to 2")
	require.NoError(t, err)
	assert.Len(t, results[0].Value, 10)
	assert.Empty(t, session.Variables())
	_, err = session.Evaluate("x")
	assert.Error(t, err)
}

func TestSession_RecordsOnlyTheProgramsThatSucceeded(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	for _, program := range []string{"x = 2", "y", "  ", "x * 3"} {
		_, _ = session.Evaluate(program)
	}

	assert.Equal(t, []string{"x = 2", "x * 3"}, session.Transcript())
}

func TestSession_SetSamples_StartsANewTranscript(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("x = 2")
	require.NoError(t, err)

	session.SetSamples(10)

	assert.Empty(t, session.Transcript())
}

func TestSession_Override_WinsOverTheProgramsAssignment(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	require.NoError(t, session.Override("team", "5"))

	results, err := session.Evaluate("team = 3; team * 2")

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(10), results[1].Value)
}

func TestSession_Override_ListsTheVariableOnce_WhenOverriddenTwice(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	require.NoError(t, session.Override("team", "5"))
	require.NoError(t, session.Override("team", "6"))

	assert.Equal(t, []application.Variable{{Name: "team", Value: domain.Scalar(6)}}, session.Variables())
}

func TestSession_Override_ListsNothing_WhenItFails(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	err := session.Override("team", "nobody")

	assert.Error(t, err)
	assert.Empty(t, session.Variables())
}

func TestSession_OverrideData_WinsOverTheProgramsAssignment(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	require.NoError(t, session.OverrideData("lead", []float64{4}))

	results, err := session.Evaluate("lead = 3; max(lead)")

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(4), results[1].Value)
}

func TestSession_OverrideData_ListsTheVariable(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	require.NoError(t, session.OverrideData("lead", []float64{4, 6}))

	variables := session.Variables()
	require.Len(t, variables, 1)
	assert.Equal(t, "lead", variables[0].Name)
	assert.Equal(t, 6.0, variables[0].Value.(domain.SampleSet).Quantile(1))
}

func TestSession_OverrideData_ListsNothing_WhenItFails(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	err := session.OverrideData("lead", nil)

	assert.Error(t, err)
	assert.Empty(t, session.Variables())
}

func TestSession_AssertsAgainstTheSessionsVariables(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate("team = 5")
	require.NoError(t, err)

	check, err := session.Assert("team < 3")

	require.NoError(t, err)
	assert.Equal(t, domain.Check{Value: domain.Probability(0), Left: domain.Scalar(5), Operator: "<", Right: domain.Scalar(3)}, check)
}

func TestSession_KeepsAssertionsOutOfTheTranscript(t *testing.T) {
	session := application.NewSession(application.DefaultSettings())

	_, err := session.Assert("1 < 2")

	require.NoError(t, err)
	assert.Empty(t, session.Transcript())
}

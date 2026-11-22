package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/application"
)

func sessionWith(t *testing.T, program string) *application.Session {
	t.Helper()
	session := application.NewSession(application.DefaultSettings())
	_, err := session.Evaluate(program)
	require.NoError(t, err)
	return session
}

func driverNames(drivers []application.Driver) []string {
	names := make([]string, len(drivers))
	for i, d := range drivers {
		names[i] = d.Name
	}
	return names
}

func TestSensitivity_RanksTheVariablesByHowStronglyTheyMoveTheTarget(t *testing.T) {
	session := sessionWith(t, "noise = 1 to 2; a = 1 to 10; opposite = 0 - a; total = 10 * a + noise")

	drivers, err := session.Sensitivity("total")

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "opposite", "noise"}, driverNames(drivers))
	assert.Greater(t, drivers[0].Correlation, 0.9)
	assert.Equal(t, -drivers[0].Correlation, drivers[1].Correlation)
	assert.Less(t, drivers[2].Correlation, 0.5)
}

func TestSensitivity_LeavesOutPlainNumbersConstantsAndTheTargetItself(t *testing.T) {
	session := sessionWith(t, "rate = 3; chance = P(1 to 2 > 1.5); flat = (1 to 2) * 0; x = 1 to 2; y = x * rate")

	drivers, err := session.Sensitivity(" y ")

	require.NoError(t, err)
	assert.Equal(t, []string{"x"}, driverNames(drivers))
}

func TestSensitivity_ComparesAnExpressionWithEveryVariable(t *testing.T) {
	session := sessionWith(t, "x = 1 to 2; y = 3 to 4")

	drivers, err := session.Sensitivity("x")

	require.NoError(t, err)
	assert.Equal(t, []string{"y"}, driverNames(drivers))
}

func TestSensitivity_FailsForAPlainNumber(t *testing.T) {
	session := sessionWith(t, "x = 1 to 2")

	_, err := session.Sensitivity("P(x > 1)")

	assert.ErrorIs(t, err, application.ErrCertainTarget)
}

func TestSensitivity_Fails_WhenNoOtherVariableIsUncertain(t *testing.T) {
	session := sessionWith(t, "x = 1 to 2; y = 3")

	_, err := session.Sensitivity("x")

	assert.ErrorIs(t, err, application.ErrNoDrivers)
}

func TestSensitivity_ReportsAnExpressionThatCannotBeEvaluated(t *testing.T) {
	session := sessionWith(t, "x = 1 to 2")

	_, err := session.Sensitivity("nobody")

	assert.ErrorContains(t, err, "unknown variable 'nobody'")
}

func TestSensitivity_LeavesTheSessionAsItWas(t *testing.T) {
	session := sessionWith(t, "x = 1 to 2; y = x * 2")

	_, err := session.Sensitivity("y + 1")
	require.NoError(t, err)

	assert.Len(t, session.VariablesAndLast(), 2)
	assert.Equal(t, []string{"x = 1 to 2; y = x * 2"}, session.Transcript())
}

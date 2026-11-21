package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func override(t *testing.T, env domain.Env, name, expr string) {
	t.Helper()
	require.NoError(t, domain.Override(name, expr, env))
}

func TestOverride_BindsTheVariable(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "2 + 3")

	value, ok := env.Variable("team")

	assert.True(t, ok)
	assert.Equal(t, domain.Scalar(5), value)
}

func TestOverride_SkipsLaterAssignmentsToTheName_AndReportsTheOverride(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "5")

	results, err := domain.Evaluate("team = 3\nteam * 2", env)

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{
		{Line: 1, Name: "team", Value: domain.Scalar(5)},
		{Line: 2, Value: domain.Scalar(10)},
	}, results)
}

func TestOverride_DoesNotEvaluateTheSkippedAssignment(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "5")

	_, err := domain.Evaluate("team = 10 to 5", env)

	assert.NoError(t, err)
}

func TestOverride_LeavesOtherAssignmentsAlone(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "5")

	results, err := domain.Evaluate("rate = 3", env)

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{{Line: 1, Name: "rate", Value: domain.Scalar(3)}}, results)
}

func TestOverride_StillEvaluatesExpressions_WhenLastIsOverridden(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "last", "5")

	results, err := domain.Evaluate("1 + 1", env)

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(2), results[0].Value)
}

func TestOverride_ReplacesAnEarlierOverride(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "5")
	override(t, env, "team", "team + 1")

	value, _ := env.Variable("team")

	assert.Equal(t, domain.Scalar(6), value)
}

func TestOverride_ReportsAComparisonAsAProbability(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "chance", "1 < 2")

	value, _ := env.Variable("chance")

	assert.Equal(t, domain.Probability(1), value)
}

func TestOverride_IsForgotten_WhenTheSampleCountChanges(t *testing.T) {
	env := domain.NewEnv(1, samples)
	override(t, env, "team", "5")

	results, err := domain.Evaluate("team = 3", env.WithSamples(10))

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(3), results[0].Value)
}

func TestOverride_Fails_WhenTheExpressionIsNotOneExpression(t *testing.T) {
	cases := map[string]domain.SyntaxError{
		"":       {Line: 1, Col: 1, Msg: "unexpected end of input"},
		"1 +":    {Line: 1, Col: 4, Msg: "unexpected end of input"},
		"1; 2":   {Line: 1, Col: 2, Msg: "unexpected ';'"},
		"x = 2":  {Line: 1, Col: 3, Msg: "unexpected '='"},
		"1 $":    {Line: 1, Col: 3, Msg: "unexpected character '$'"},
		"f(x)=x": {Line: 1, Col: 5, Msg: "unexpected '='"},
	}
	for expr, want := range cases {
		t.Run(expr, func(t *testing.T) {
			env := domain.NewEnv(1, samples)

			err := domain.Override("team", expr, env)

			var syntaxErr *domain.SyntaxError
			require.ErrorAs(t, err, &syntaxErr)
			assert.Equal(t, want, *syntaxErr)
			_, bound := env.Variable("team")
			assert.False(t, bound)
		})
	}
}

func TestOverride_Fails_WhenTheExpressionCannotBeEvaluated(t *testing.T) {
	env := domain.NewEnv(1, samples)

	err := domain.Override("team", "nobody", env)

	var evalErr *domain.EvalError
	require.ErrorAs(t, err, &evalErr)
	assert.Equal(t, domain.EvalError{Line: 1, Col: 1, Msg: "unknown variable 'nobody'"}, *evalErr)
	_, bound := env.Variable("team")
	assert.False(t, bound)
}

func TestIsName_AcceptsIdentifiersOnly(t *testing.T) {
	cases := map[string]bool{
		"team": true, "_x9": true, "Rate": true,
		"": false, " team": false, "team ": false, "2x": false, "to": false, "a b": false, "a-b": false, "x;y": false, "é": false,
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, domain.IsName(name))
		})
	}
}

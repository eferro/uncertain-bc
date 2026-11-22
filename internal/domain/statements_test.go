package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ReturnsOneResultPerStatement(t *testing.T) {
	results, err := domain.Evaluate("1 + 2; 3 * 4\n5", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{
		{Line: 1, Source: "1 + 2", Value: domain.Scalar(3)},
		{Line: 1, Source: "3 * 4", Value: domain.Scalar(12)},
		{Line: 2, Source: "5", Value: domain.Scalar(5)},
	}, results)
}

func TestEvaluate_SkipsBlankLinesAndComments(t *testing.T) {
	results, err := domain.Evaluate("# header\n\n  ;\n1 # trailing\r\n", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{{Line: 4, Source: "1", Value: domain.Scalar(1)}}, results)
}

func TestEvaluate_BindsAVariable_WhenAStatementIsAnAssignment(t *testing.T) {
	results, err := domain.Evaluate("x = 2\nx * 3", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{
		{Line: 1, Name: "x", Source: "2", Value: domain.Scalar(2)},
		{Line: 2, Source: "x * 3", Value: domain.Scalar(6)},
	}, results)
}

func TestEvaluate_KeepsTheSourceOfEachExpressionAsWritten(t *testing.T) {
	results, err := domain.Evaluate("x  =  2 * 3   # six\n  P(x  > 5) ;max(x, 1)\n", domain.NewEnv(1, samples))

	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, []string{"2 * 3", "P(x  > 5)", "max(x, 1)"}, sources(results))
}

func sources(results []domain.Result) []string {
	texts := make([]string, 0, len(results))
	for _, result := range results {
		texts = append(texts, result.Source)
	}
	return texts
}

func TestEvaluate_LetsAVariableBeReassigned(t *testing.T) {
	assert.Equal(t, domain.Scalar(2), evaluate(t, "x = 1; x = x + 1; x"))
}

func TestEvaluate_KeepsTheSamplesOfAVariable_SoXMinusXIsZero(t *testing.T) {
	s := evaluateSamples(t, "x = 1 to 5\nx - x")

	assert.Equal(t, 0.0, s.Quantile(0))
	assert.Equal(t, 0.0, s.Quantile(1))
}

func TestEvaluate_DrawsIndependentSamples_ForDifferentVariables(t *testing.T) {
	s := evaluateSamples(t, "x = 1 to 5\ny = 1 to 5\nx - y")

	assert.Greater(t, s.Quantile(0.95)-s.Quantile(0.05), 1.0)
}

func TestEvaluate_KeepsVariablesAcrossCallsWithTheSameEnv(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("x = 2", env)
	require.NoError(t, err)

	results, err := domain.Evaluate("x + 1", env)

	require.NoError(t, err)
	assert.Equal(t, domain.Scalar(3), results[0].Value)
}

func TestEvaluate_Fails_WhenAVariableIsUnknown(t *testing.T) {
	_, err := domain.Evaluate("1\n2 + y", domain.NewEnv(1, samples))

	var evalErr *domain.EvalError
	require.ErrorAs(t, err, &evalErr)
	assert.Equal(t, domain.EvalError{Line: 2, Col: 5, Msg: "unknown variable 'y'"}, *evalErr)
}

func TestEvaluate_ReturnsTheResultsBeforeAnEvaluationError(t *testing.T) {
	results, err := domain.Evaluate("1\n10 to 5\n3", domain.NewEnv(1, samples))

	var evalErr *domain.EvalError
	require.ErrorAs(t, err, &evalErr)
	assert.Equal(t, domain.EvalError{Line: 2, Col: 4, Msg: "range start 10 must be below its end 5"}, *evalErr)
	assert.Equal(t, []domain.Result{{Line: 1, Source: "1", Value: domain.Scalar(1)}}, results)
}

func TestEvaluate_EvaluatesNothing_WhenAnyStatementHasASyntaxError(t *testing.T) {
	results, err := domain.Evaluate("x = 1\n2\n(3", domain.NewEnv(1, samples))

	var syntaxErr *domain.SyntaxError
	require.ErrorAs(t, err, &syntaxErr)
	assert.Equal(t, domain.SyntaxError{Line: 3, Col: 3, Msg: "unexpected end of input"}, *syntaxErr)
	assert.Empty(t, results)
}

func TestEvaluate_TreatsEveryIdentifierAsAVariable(t *testing.T) {
	for _, name := range []string{"forty", "a", "z", "A", "Z", "_x9"} {
		t.Run(name, func(t *testing.T) {
			_, err := domain.Evaluate(name, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, "unknown variable '"+name+"'", evalErr.Msg)
		})
	}
}

func TestEvaluate_ReturnsNoResults_WhenTheProgramIsEmpty(t *testing.T) {
	results, err := domain.Evaluate("", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestEvaluate_AcceptsWindowsLineEndings(t *testing.T) {
	results, err := domain.Evaluate("1\r\n2\r\n", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{
		{Line: 1, Source: "1", Value: domain.Scalar(1)},
		{Line: 2, Source: "2", Value: domain.Scalar(2)},
	}, results)
}

package domain_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestUserFunction_IsDefinedSilentlyAndCalledLikeABuiltIn(t *testing.T) {
	results, err := domain.Evaluate("cost(n) = n * 2\ncost(3)", domain.NewEnv(1, samples))

	require.NoError(t, err)
	assert.Equal(t, []domain.Result{{Line: 2, Source: "cost(3)", Value: domain.Scalar(6)}}, results)
}

func TestUserFunction_TakesSeveralParameters(t *testing.T) {
	assert.Equal(t, domain.Scalar(6), evaluate(t, "area(w, h) = w * h; area(2, 3)"))
}

func TestUserFunction_TakesNoParameters(t *testing.T) {
	assert.Equal(t, domain.Scalar(42), evaluate(t, "answer() = 42; answer()"))
}

func TestUserFunction_SeesGlobalVariablesAsTheyAreWhenCalled(t *testing.T) {
	assert.Equal(t, domain.Scalar(40), evaluate(t, "rate = 10; pay(h) = h * rate; rate = 20; pay(2)"))
}

func TestUserFunction_ParametersShadowGlobalsWithoutChangingThem(t *testing.T) {
	assert.Equal(t, domain.Scalar(2), evaluate(t, "x = 100; f(x) = x + 1; f(1)"))
	assert.Equal(t, domain.Scalar(100), evaluate(t, "x = 100; f(x) = x + 1; f(1); x"))
}

func TestUserFunction_DrawsFreshSamplesOnEveryCall(t *testing.T) {
	s := evaluateSamples(t, "noisy(n) = n * (1 to 2); noisy(1) - noisy(1)")

	assert.Greater(t, s.Quantile(0.95)-s.Quantile(0.05), 0.2)
}

func TestUserFunction_AcceptsUncertainArguments(t *testing.T) {
	s := evaluateSamples(t, "double(x) = x * 2; double(1 to 3)")

	assertNear(t, 2, s.Quantile(0.05), 4)
	assertNear(t, 6, s.Quantile(0.95), 4)
}

func TestUserFunction_CanBeRedefined(t *testing.T) {
	assert.Equal(t, domain.Scalar(6), evaluate(t, "f(x) = x; f(x) = x * 2; f(3)"))
}

func TestUserFunction_CanCallOtherFunctions(t *testing.T) {
	assert.Equal(t, domain.Scalar(12), evaluate(t, "double(x) = x * 2; quad(x) = double(double(x)); quad(3)"))
}

func TestUserFunction_ReportsMistakes(t *testing.T) {
	cases := map[string]domain.EvalError{
		"f(x) = x; f(1, 2)": {Line: 1, Col: 11, Msg: "f expects 1 argument, got 2"},
		"mean(x) = x":       {Line: 1, Col: 1, Msg: "cannot redefine built-in function 'mean'"},
		"f(x, x) = x":       {Line: 1, Col: 6, Msg: "parameter 'x' appears twice"},
		"f(x) = y; f(1)":    {Line: 1, Col: 8, Msg: "unknown variable 'y'"},
		"f(x) = f(x); f(1)": {Line: 1, Col: 8, Msg: "function calls nest too deeply (does f call itself?)"},
		"f(x) = g(x); g(1)": {Line: 1, Col: 14, Msg: "unknown function 'g'"},
	}
	for program, want := range cases {
		t.Run(program, func(t *testing.T) {
			_, err := domain.Evaluate(program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, want, *evalErr)
		})
	}
}

func TestUserFunction_KeepsDefinitionsAcrossCallsWithTheSameEnv(t *testing.T) {
	env := domain.NewEnv(1, samples)
	_, err := domain.Evaluate("triple(x) = x * 3", env)
	require.NoError(t, err)

	assert.Equal(t, domain.Scalar(9), evaluateIn(t, env, "triple(3)"))
}

func chainOfCalls(depth int) string {
	var program strings.Builder
	for i := 1; i < depth; i++ {
		fmt.Fprintf(&program, "f%d(x) = f%d(x)\n", i, i+1)
	}
	fmt.Fprintf(&program, "f%d(x) = x\nf1(7)", depth)
	return program.String()
}

func TestUserFunction_AllowsCallsNestedSixtyFourDeep(t *testing.T) {
	assert.Equal(t, domain.Scalar(7), evaluate(t, chainOfCalls(64)))
}

func TestUserFunction_RejectsCallsNestedSixtyFiveDeep(t *testing.T) {
	_, err := domain.Evaluate(chainOfCalls(65), domain.NewEnv(1, samples))

	var evalErr *domain.EvalError
	require.ErrorAs(t, err, &evalErr)
	assert.Contains(t, evalErr.Msg, "function calls nest too deeply")
}

package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestRange_DrawsOneSamplePerMonteCarloIteration(t *testing.T) {
	assert.Len(t, evaluateSamples(t, "10 to 20"), samples)
}

func TestRange_IsLognormal_WhenBothEndsArePositive(t *testing.T) {
	s := evaluateSamples(t, "10 to 20")

	assert.InDelta(t, 10, s.Quantile(0.05), 0.3)
	assert.InDelta(t, math.Sqrt(10*20), s.Quantile(0.5), 0.3)
	assert.InDelta(t, 20, s.Quantile(0.95), 0.6)
}

func TestRange_IsNormal_WhenAnEndIsNotPositive(t *testing.T) {
	s := evaluateSamples(t, "0 to 10")

	assert.InDelta(t, 0, s.Quantile(0.05), 0.3)
	assert.InDelta(t, 5, s.Quantile(0.5), 0.3)
	assert.InDelta(t, 10, s.Quantile(0.95), 0.3)
}

func TestRange_BindsLooserThanAddition(t *testing.T) {
	s := evaluateSamples(t, "5 + 5 to 10 + 10")

	assert.InDelta(t, math.Sqrt(10*20), s.Quantile(0.5), 0.3)
}

func TestRange_Fails_WhenItsStartIsNotBelowItsEnd(t *testing.T) {
	cases := []struct {
		program string
		want    domain.EvalError
	}{
		{"10 to 5", domain.EvalError{Line: 1, Col: 4, Msg: "range start 10 must be below its end 5"}},
		{"5 to 5", domain.EvalError{Line: 1, Col: 3, Msg: "range start 5 must be below its end 5"}},
	}
	for _, c := range cases {
		t.Run(c.program, func(t *testing.T) {
			_, err := domain.Evaluate(c.program, domain.NewEnv(1, samples))

			var evalErr *domain.EvalError
			require.ErrorAs(t, err, &evalErr)
			assert.Equal(t, c.want, *evalErr)
		})
	}
}

func TestEvalError_ReadsAsLineColumnAndMessage(t *testing.T) {
	err := &domain.EvalError{Line: 1, Col: 4, Msg: "range start 10 must be below its end 5"}

	assert.EqualError(t, err, "line 1, col 4: range start 10 must be below its end 5")
}

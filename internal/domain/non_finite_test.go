package domain_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSampleSet_Quantile_IgnoresNonFiniteSamples(t *testing.T) {
	s := domain.SampleSet{1, math.NaN(), 3, math.Inf(1), 2, math.Inf(-1)}

	assert.Equal(t, 1.0, s.Quantile(0))
	assert.Equal(t, 2.0, s.Quantile(0.5))
	assert.Equal(t, 3.0, s.Quantile(1))
}

func TestSampleSet_Quantile_IsNaN_WhenNoSampleIsFinite(t *testing.T) {
	assert.True(t, math.IsNaN(domain.SampleSet{math.NaN(), math.Inf(1)}.Quantile(0.5)))
}

func TestSampleSet_Dropped_CountsNonFiniteSamples(t *testing.T) {
	assert.Equal(t, 3, domain.SampleSet{1, math.NaN(), 3, math.Inf(1), 2, math.Inf(-1)}.Dropped())
	assert.Equal(t, 0, domain.SampleSet{1, 2}.Dropped())
}

func TestSummaries_IgnoreNonFiniteSamples(t *testing.T) {
	program := "x = 0 to 1; y = x / (x > 0.5); "
	sd := 1 / (2 * 1.6448536269514722)
	tailMean := 0.5 + sd*math.Sqrt(2/math.Pi)

	assertNear(t, tailMean, float64(evaluate(t, program+"mean(y)").(domain.Scalar)), 1)
	assertNear(t, 0.5+sd*0.6745, float64(evaluate(t, program+"median(y)").(domain.Scalar)), 1)
	assert.False(t, math.IsNaN(float64(evaluate(t, program+"stdev(y)").(domain.Scalar))))
}

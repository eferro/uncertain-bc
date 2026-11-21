package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSampleSet_Quantile_InterpolatesBetweenSortedSamples(t *testing.T) {
	cases := []struct {
		samples domain.SampleSet
		q       float64
		want    float64
	}{
		{domain.SampleSet{3, 1, 2}, 0, 1},
		{domain.SampleSet{3, 1, 2}, 0.5, 2},
		{domain.SampleSet{3, 1, 2}, 1, 3},
		{domain.SampleSet{1, 2}, 0.5, 1.5},
		{domain.SampleSet{10, 0, 20, 30}, 0.25, 7.5},
	}
	for _, c := range cases {
		assert.InDelta(t, c.want, c.samples.Quantile(c.q), 1e-9, "q=%v of %v", c.q, c.samples)
	}
}

func TestSampleSet_Quantile_LeavesTheSamplesInTheirOriginalOrder(t *testing.T) {
	s := domain.SampleSet{3, 1, 2}

	s.Quantile(0.5)

	assert.Equal(t, domain.SampleSet{3, 1, 2}, s)
}

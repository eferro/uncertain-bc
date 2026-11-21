package cli

import (
	"math"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func evenlySpread(n int) domain.SampleSet {
	s := make(domain.SampleSet, n)
	for i := range s {
		s[i] = float64(i)
	}
	return s
}

func twoClusters(n int) domain.SampleSet {
	s := make(domain.SampleSet, n)
	for i := n / 2; i < n; i++ {
		s[i] = 10
	}
	return s
}

func TestSparkline_IsSixteenCellsWide(t *testing.T) {
	assert.Equal(t, 16, utf8.RuneCountInString(sparkline(evenlySpread(1000), unicodeBlocks)))
}

func TestSparkline_IsFlat_WhenSamplesAreEvenlySpread(t *testing.T) {
	assert.Equal(t, "████████████████", sparkline(evenlySpread(1000), unicodeBlocks))
}

func TestSparkline_ShowsBothEnds_WhenSamplesFormTwoClusters(t *testing.T) {
	assert.Equal(t, "█▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", sparkline(twoClusters(1000), unicodeBlocks))
}

func TestSparkline_ScalesHeightsToTheFullestCell(t *testing.T) {
	s := evenlySpread(1000)
	for range 1000 {
		s = append(s, 500)
	}

	line := []rune(sparkline(s, unicodeBlocks))

	assert.Equal(t, '█', line[8])
	assert.Equal(t, '▁', line[0])
}

func TestSparkline_ClipsToTheFirstAndNinetyNinthPercentiles(t *testing.T) {
	s := evenlySpread(1000)
	s[0], s[999] = -1e9, 1e9

	assert.Equal(t, "████████████████", sparkline(s, unicodeBlocks))
}

func TestSparkline_UsesTheGivenGlyphs(t *testing.T) {
	assert.Equal(t, "#              #", sparkline(twoClusters(1000), asciiBlocks))
}

func TestSparkline_IgnoresNonFiniteSamples(t *testing.T) {
	s := evenlySpread(1000)
	s[500] = math.NaN()

	assert.Equal(t, "████████████████", sparkline(s, unicodeBlocks))
}

package cli

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistogramChart_HasTwentyLabelledRowsAndAFooter(t *testing.T) {
	lines := strings.Split(strings.TrimSuffix(histogramChart(evenlySpread(1000), unicodeStyle), "\n"), "\n")

	require.Len(t, lines, 21)
	assert.Equal(t, "   34.5 ┤██████████████████████████████████████████████████", lines[0])
	assert.Equal(t, "    965 ┤██████████████████████████████████████████████████", lines[19])
	assert.Equal(t, "(P1 to P99 shown, 20 samples outside)", lines[20])
}

func TestHistogramChart_ScalesBarsToTheFullestRow(t *testing.T) {
	lines := strings.Split(histogramChart(twoClusters(1000), unicodeStyle), "\n")

	assert.True(t, strings.HasSuffix(lines[0], "┤"+strings.Repeat("█", 50)))
	assert.True(t, strings.HasSuffix(lines[10], "┤"))
	assert.True(t, strings.HasSuffix(lines[19], "┤"+strings.Repeat("█", 50)))
}

func TestHistogramChart_UsesTheAsciiStyle(t *testing.T) {
	first := strings.Split(histogramChart(evenlySpread(1000), asciiStyle), "\n")[0]

	assert.Equal(t, "   34.5 |"+strings.Repeat("#", 50), first)
}

func TestHistogramChart_DoesNotCountNonFiniteSamplesAsOutside(t *testing.T) {
	s := evenlySpread(1000)
	s[500], s[501] = math.NaN(), math.Inf(1)

	assert.Contains(t, histogramChart(s, unicodeStyle), "(P1 to P99 shown, 20 samples outside)")
}

func TestHistogramChart_AlignsTheAxis_WhenLabelsAreLong(t *testing.T) {
	s := evenlySpread(1000)
	for i := range s {
		s[i] *= 1e-7
	}

	lines := strings.Split(strings.TrimSuffix(histogramChart(s, unicodeStyle), "\n"), "\n")
	for _, line := range lines[:20] {
		assert.Equal(t, strings.Index(lines[0], "┤"), strings.Index(line, "┤"), line)
	}
}

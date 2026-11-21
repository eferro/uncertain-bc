package cli

import (
	"fmt"
	"math"
	"strings"

	"github.com/eferro/uncertain-bc/internal/domain"
)

const (
	histogramRows     = 20
	histogramBarWidth = 50
	minLabelWidth     = 7
)

// histogramChart draws one row per bin over [P1, P99] and says how many
// samples fall outside, so heavy tails are visible but cannot squash the
// chart (SPEC §7, VISION §8).
func histogramChart(s domain.SampleSet, st style) string {
	lo, hi := s.Quantile(0.01), s.Quantile(0.99)
	counts := histogram(s, lo, hi, histogramRows)
	labels := rowLabels(lo, hi)
	width := labelWidth(labels)
	fullest := max(1, maxOf(counts))
	var chart strings.Builder
	for i, count := range counts {
		bar := strings.Repeat(st.bar, int(math.Round(float64(count)/float64(fullest)*histogramBarWidth)))
		fmt.Fprintf(&chart, "%*s %s%s\n", width, labels[i], st.axis, bar)
	}
	fmt.Fprintf(&chart, "(P1 to P99 shown, %d samples outside)\n", len(s)-s.Dropped()-sum(counts))
	return chart.String()
}

func rowLabels(lo, hi float64) []string {
	step := (hi - lo) / histogramRows
	labels := make([]string, histogramRows)
	for i := range labels {
		labels[i] = significant(lo + step*(float64(i)+0.5))
	}
	return labels
}

func labelWidth(labels []string) int {
	width := minLabelWidth
	for _, label := range labels {
		width = max(width, len(label))
	}
	return width
}

func sum(counts []int) int {
	total := 0
	for _, c := range counts {
		total += c
	}
	return total
}

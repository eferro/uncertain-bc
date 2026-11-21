package cli

import (
	"math"

	"github.com/eferro/uncertain-bc/internal/domain"
)

const sparklineCells = 16

var (
	unicodeBlocks = []rune("▁▂▃▄▅▆▇█")
	asciiBlocks   = []rune(" .:-=+*#")
)

// sparkline draws a histogram over [P1, P99], so a few extreme samples
// cannot squash the shape of the rest (SPEC §6).
func sparkline(s domain.SampleSet, glyphs []rune) string {
	counts := histogram(s, s.Quantile(0.01), s.Quantile(0.99), sparklineCells)
	fullest := max(1, maxOf(counts))
	line := make([]rune, sparklineCells)
	for i, count := range counts {
		level := math.Round(float64(count) / float64(fullest) * float64(len(glyphs)-1))
		line[i] = glyphs[int(level)]
	}
	return string(line)
}

func histogram(s domain.SampleSet, lo, hi float64, cells int) []int {
	counts := make([]int, cells)
	width := hi - lo
	for _, x := range s {
		if !(x >= lo && x <= hi) {
			continue
		}
		cell := cells - 1
		if width > 0 {
			cell = min(int((x-lo)/width*float64(cells)), cells-1)
		}
		counts[cell]++
	}
	return counts
}

func maxOf(counts []int) int {
	fullest := 0
	for _, c := range counts {
		fullest = max(fullest, c)
	}
	return fullest
}

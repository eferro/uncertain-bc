package cli

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/eferro/uncertain-bc/internal/domain"
)

const (
	plotWidth   = 640.0
	plotHeight  = 360.0
	plotBins    = 40
	plotLeft    = 40.0
	plotRight   = 20.0
	plotTop     = 50.0
	plotBottom  = 40.0
	plotTicks   = 4
	plotAreaW   = plotWidth - plotLeft - plotRight
	plotAreaH   = plotHeight - plotTop - plotBottom
	plotBaseY   = plotTop + plotAreaH
	plotPalette = `.bar{fill:#4c78a8}.axis{stroke:#333}text{font:12px sans-serif;fill:#333}.title{font-size:16px}`
)

type chartText struct {
	x, y          float64
	class, anchor string
	text          string
}

type chartBar struct {
	x, y, width, height float64
}

// chart is the layout shared by every export format: a histogram over
// [P1, P99] titled with the expression and its summary (SPEC §7).
type chart struct {
	texts []chartText
	bars  []chartBar
}

func layoutChart(s domain.SampleSet, title string, st style) chart {
	lo, hi := s.Quantile(0.01), s.Quantile(0.99)
	counts := histogram(s, lo, hi, plotBins)
	outside := fmt.Sprintf("P1 to P99 shown, %d samples outside", len(s)-s.Dropped()-sum(counts))
	texts := []chartText{
		{plotLeft, 22, "title", "start", title},
		{plotLeft, 40, "summary", "start", summaryNumbers(s, st)},
		{plotWidth - plotRight, plotHeight - 8, "note", "end", outside},
	}
	return chart{texts: append(texts, tickTexts(lo, hi)...), bars: layoutBars(counts)}
}

func layoutBars(counts []int) []chartBar {
	fullest := float64(max(1, maxOf(counts)))
	width := plotAreaW / plotBins
	bars := make([]chartBar, len(counts))
	for i, count := range counts {
		height := float64(count) / fullest * plotAreaH
		bars[i] = chartBar{plotLeft + float64(i)*width, plotBaseY - height, width - 1, height}
	}
	return bars
}

func tickTexts(lo, hi float64) []chartText {
	ticks := make([]chartText, 0, plotTicks+1)
	for k := 0; k <= plotTicks; k++ {
		x := plotLeft + float64(k)*plotAreaW/plotTicks
		ticks = append(ticks, chartText{x, plotBaseY + 16, "tick", "middle", significant(lo + float64(k)*(hi-lo)/plotTicks)})
	}
	return ticks
}

func svgHistogram(s domain.SampleSet, title string) string {
	c := layoutChart(s, title, unicodeStyle)
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`,
		plotWidth, plotHeight, plotWidth, plotHeight)
	fmt.Fprintf(&svg, "<style>%s</style>\n", plotPalette)
	svg.WriteString(`<rect width="100%" height="100%" fill="white"/>` + "\n")
	for _, b := range c.bars {
		fmt.Fprintf(&svg, `<rect class="bar" x="%.1f" y="%.1f" width="%.1f" height="%.1f"/>`+"\n", b.x, b.y, b.width, b.height)
	}
	fmt.Fprintf(&svg, `<line class="axis" x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`+"\n",
		plotLeft, plotBaseY, plotLeft+plotAreaW, plotBaseY)
	for _, t := range c.texts {
		fmt.Fprintf(&svg, `<text class="%s" x="%.1f" y="%.1f" text-anchor="%s">`, t.class, t.x, t.y, t.anchor)
		_ = xml.EscapeText(&svg, []byte(t.text))
		svg.WriteString("</text>\n")
	}
	svg.WriteString("</svg>\n")
	return svg.String()
}

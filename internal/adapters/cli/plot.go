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

// svgHistogram draws a standalone SVG histogram over [P1, P99], titled with
// the expression and its one-line summary (SPEC §7).
func svgHistogram(s domain.SampleSet, title string) string {
	lo, hi := s.Quantile(0.01), s.Quantile(0.99)
	counts := histogram(s, lo, hi, plotBins)
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`,
		plotWidth, plotHeight, plotWidth, plotHeight)
	fmt.Fprintf(&svg, "<style>%s</style>\n", plotPalette)
	svg.WriteString(`<rect width="100%" height="100%" fill="white"/>` + "\n")
	writeText(&svg, plotLeft, 22, "title", "start", title)
	writeText(&svg, plotLeft, 40, "summary", "start", summaryNumbers(s, unicodeStyle))
	writeBars(&svg, counts)
	writeAxis(&svg, lo, hi)
	outside := fmt.Sprintf("P1 to P99 shown, %d samples outside", len(s)-s.Dropped()-sum(counts))
	writeText(&svg, plotWidth-plotRight, plotHeight-8, "note", "end", outside)
	svg.WriteString("</svg>\n")
	return svg.String()
}

func writeBars(svg *strings.Builder, counts []int) {
	fullest := float64(max(1, maxOf(counts)))
	width := plotAreaW / plotBins
	for i, count := range counts {
		height := float64(count) / fullest * plotAreaH
		fmt.Fprintf(svg, `<rect class="bar" x="%.1f" y="%.1f" width="%.1f" height="%.1f"/>`+"\n",
			plotLeft+float64(i)*width, plotBaseY-height, width-1, height)
	}
}

func writeAxis(svg *strings.Builder, lo, hi float64) {
	fmt.Fprintf(svg, `<line class="axis" x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`+"\n",
		plotLeft, plotBaseY, plotLeft+plotAreaW, plotBaseY)
	for k := 0; k <= plotTicks; k++ {
		x := plotLeft + float64(k)*plotAreaW/plotTicks
		writeText(svg, x, plotBaseY+16, "tick", "middle", significant(lo+float64(k)*(hi-lo)/plotTicks))
	}
}

func writeText(svg *strings.Builder, x, y float64, class, anchor, text string) {
	fmt.Fprintf(svg, `<text class="%s" x="%.1f" y="%.1f" text-anchor="%s">`, class, x, y, anchor)
	_ = xml.EscapeText(svg, []byte(text))
	svg.WriteString("</text>\n")
}

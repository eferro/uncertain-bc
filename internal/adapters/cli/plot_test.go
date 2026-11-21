package cli

import (
	"encoding/xml"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type svgElements struct {
	bars  int
	texts []string
}

func parseSVG(t *testing.T, svg string) svgElements {
	t.Helper()
	var found svgElements
	decoder := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			return found
		}
		require.NoError(t, err)
		found = collect(t, decoder, tok, found)
	}
}

func collect(t *testing.T, decoder *xml.Decoder, tok xml.Token, found svgElements) svgElements {
	t.Helper()
	start, ok := tok.(xml.StartElement)
	if !ok {
		return found
	}
	if start.Name.Local == "rect" && hasClass(start, "bar") {
		found.bars++
	}
	if start.Name.Local == "text" {
		var text string
		require.NoError(t, decoder.DecodeElement(&text, &start))
		found.texts = append(found.texts, text)
	}
	return found
}

func hasClass(start xml.StartElement, class string) bool {
	for _, attr := range start.Attr {
		if attr.Name.Local == "class" && attr.Value == class {
			return true
		}
	}
	return false
}

func TestSVGHistogram_IsWellFormedWithFortyBars(t *testing.T) {
	assert.Equal(t, 40, parseSVG(t, svgHistogram(evenlySpread(1000), "x")).bars)
}

func TestSVGHistogram_ShowsTitleSummaryTicksAndClipping(t *testing.T) {
	texts := parseSVG(t, svgHistogram(evenlySpread(1000), "cost <& more>")).texts

	assert.Contains(t, texts, "cost <& more>")
	assert.Contains(t, texts, "≈ 500   90% [50 … 949]")
	assert.Contains(t, texts, "9.99")
	assert.Contains(t, texts, "989")
	assert.Contains(t, texts, "P1 to P99 shown, 20 samples outside")
}

func TestSVGHistogram_ScalesBarsToTheTallestOne(t *testing.T) {
	svg := svgHistogram(twoClusters(1000), "x")

	assert.Equal(t, 2, strings.Count(svg, `height="270.0"`))
	assert.Equal(t, 38, strings.Count(svg, `height="0.0"`))
}

func TestSVGHistogram_DoesNotCountNonFiniteSamplesAsOutside(t *testing.T) {
	s := evenlySpread(1000)
	s[500] = math.NaN()

	assert.Contains(t, parseSVG(t, svgHistogram(s, "x")).texts, "P1 to P99 shown, 20 samples outside")
}

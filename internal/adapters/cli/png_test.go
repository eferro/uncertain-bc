package cli

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isWhite(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0xffff && g == 0xffff && b == 0xffff
}

func darkPixelsIn(img image.Image, area image.Rectangle) int {
	dark := 0
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r < 0x8000 {
				dark++
			}
		}
	}
	return dark
}

func TestPNGHistogram_IsSixHundredFortyByThreeHundredSixtyOnWhite(t *testing.T) {
	img := pngHistogram(evenlySpread(1000), "x")

	assert.Equal(t, image.Rect(0, 0, 640, 360), img.Bounds())
	assert.True(t, isWhite(img.At(1, 1)))
}

func TestPNGHistogram_DrawsBarsWhereSamplesAre(t *testing.T) {
	img := pngHistogram(twoClusters(1000), "x")

	assert.False(t, isWhite(img.At(int(plotLeft)+5, int(plotBaseY)-10)))
	assert.True(t, isWhite(img.At(int(plotLeft+plotAreaW/2), int(plotBaseY)-10)))
	assert.False(t, isWhite(img.At(int(plotLeft+plotAreaW)-8, int(plotBaseY)-10)))
}

func TestPNGHistogram_WritesTheTitleAndTicks(t *testing.T) {
	img := pngHistogram(evenlySpread(1000), "a long title")

	assert.Greater(t, darkPixelsIn(img, image.Rect(int(plotLeft), 8, int(plotLeft)+90, 26)), 50)
	assert.Greater(t, darkPixelsIn(img, image.Rect(int(plotLeft)-20, int(plotBaseY)+4, int(plotLeft)+20, int(plotBaseY)+20)), 10)
}

func TestPNGHistogram_GrowsBarsUpFromTheAxis(t *testing.T) {
	s := twoClusters(1000)
	for i := 500; i < 750; i++ {
		s[i] = 0
	}
	right := int(plotLeft+plotAreaW) - 8

	img := pngHistogram(s, "x")

	assert.False(t, isWhite(img.At(right, int(plotBaseY)-10)))
	assert.True(t, isWhite(img.At(right, int(plotTop)+10)))
}

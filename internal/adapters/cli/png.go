package cli

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/eferro/uncertain-bc/internal/domain"
)

var (
	barColour  = color.RGBA{R: 0x4c, G: 0x78, B: 0xa8, A: 0xff}
	inkColour  = color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}
	paperColor = color.White
)

// pngHistogram renders the chart layout as an image. Its bitmap font only
// covers ASCII, so the summary uses the ASCII style (ADR 0024).
func pngHistogram(s domain.SampleSet, title string) image.Image {
	c := layoutChart(s, title, asciiStyle)
	img := image.NewRGBA(image.Rect(0, 0, int(plotWidth), int(plotHeight)))
	draw.Draw(img, img.Bounds(), image.NewUniform(paperColor), image.Point{}, draw.Src)
	for _, b := range c.bars {
		fill(img, image.Rect(int(b.x), int(b.y), int(b.x+b.width), int(plotBaseY)), barColour)
	}
	fill(img, image.Rect(int(plotLeft), int(plotBaseY), int(plotLeft+plotAreaW), int(plotBaseY)+1), inkColour)
	for _, t := range c.texts {
		write(img, t)
	}
	return img
}

func fill(img draw.Image, area image.Rectangle, c color.Color) {
	draw.Draw(img, area, image.NewUniform(c), image.Point{}, draw.Src)
}

func write(img draw.Image, t chartText) {
	pen := &font.Drawer{Dst: img, Src: image.NewUniform(inkColour), Face: basicfont.Face7x13}
	x := fixed.I(int(t.x))
	switch t.anchor {
	case "middle":
		x -= pen.MeasureString(t.text) / 2
	case "end":
		x -= pen.MeasureString(t.text)
	}
	pen.Dot = fixed.Point26_6{X: x, Y: fixed.I(int(t.y))}
	pen.DrawString(t.text)
}

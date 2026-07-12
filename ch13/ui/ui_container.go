package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UIContainer draws content clipped to its bounds. Children are guaranteed to stay inside.
type UIContainer struct {
	X, Y          float64
	Width, Height float64
	BgColor       color.Color
	BorderColor   color.Color
	BorderWidth   float64
}

// NewUIContainer creates a container with the given bounds.
func NewUIContainer(x, y, w, h float64) *UIContainer {
	return &UIContainer{
		X:           x,
		Y:           y,
		Width:       w,
		Height:      h,
		BgColor:     color.RGBA{35, 35, 55, 255},
		BorderColor: color.RGBA{100, 120, 180, 255},
		BorderWidth: 2.0,
	}
}

// Draw draws the container: creates a canvas, fills background, calls drawFn for content,
// then blits to target. Everything drawn in drawFn is contained within the container bounds.
func (c *UIContainer) Draw(target *ebiten.Image, drawFn func(canvas *ebiten.Image)) {
	canvas := ebiten.NewImage(int(c.Width), int(c.Height))

	// Sfondo
	if c.BgColor != nil {
		r, g, bl, a := c.BgColor.RGBA()
		vector.DrawFilledRect(canvas, 0, 0, float32(c.Width), float32(c.Height),
			color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}, true)
	}

	// Content (local coordinates: 0,0 = top-left corner)
	drawFn(canvas)

	// Copy the canvas onto the target (Round forces integer coordinates = same result each frame)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Round(c.X), math.Round(c.Y))
	target.DrawImage(canvas, op)

	// Bordo
	if c.BorderWidth > 0 && c.BorderColor != nil {
		cx, cy := math.Round(c.X), math.Round(c.Y)
		vector.StrokeRect(target, float32(cx), float32(cy), float32(c.Width), float32(c.Height),
			float32(c.BorderWidth), c.BorderColor, true)
	}
}

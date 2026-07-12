package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UIBar draws a progress bar (health, XP) as background + foreground rectangles.
type UIBar struct {
	X, Y          float64
	Width, Height float64
	Progress      float64 // 0.0 to 1.0
	BgColor       color.Color
	FgColor       color.Color
}

// NewUIBar creates a bar with default colors.
func NewUIBar(x, y, width, height float64) *UIBar {
	return &UIBar{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Progress: 1.0,
		BgColor:  color.RGBA{30, 30, 50, 255},
		FgColor:  color.RGBA{100, 180, 255, 255},
	}
}

// SetProgress clamps progress to [0, 1].
func (b *UIBar) SetProgress(p float64) {
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	b.Progress = p
}

// colorToRGBA8 converts a color.Color to 8-bit-per-channel RGBA for vector.DrawFilledRect.
// Color.RGBA() returns 16-bit channels; we keep the high byte (same as before).
func colorToRGBA8(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// Draw renders the bar on the target.
func (b *UIBar) Draw(target *ebiten.Image) {
	vector.DrawFilledRect(
		target,
		float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height),
		colorToRGBA8(b.BgColor),
		true,
	)
	fillW := b.Width * b.Progress
	if fillW > 0 {
		vector.DrawFilledRect(
			target,
			float32(b.X), float32(b.Y), float32(fillW), float32(b.Height),
			colorToRGBA8(b.FgColor),
			true,
		)
	}
}

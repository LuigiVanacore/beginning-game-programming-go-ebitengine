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

// Draw renders the bar on the target.
func (b *UIBar) Draw(target *ebiten.Image) {
	r, g, bl, a := b.BgColor.RGBA()
	vector.DrawFilledRect(target, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height),
		color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}, true)
	fillWidth := b.Width * b.Progress
	if fillWidth > 0 {
		r, g, bl, a := b.FgColor.RGBA()
		vector.DrawFilledRect(target, float32(b.X), float32(b.Y), float32(fillWidth), float32(b.Height),
			color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}, true)
	}
}

package ui

import (
	. "book/code/ch08/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// XPBar draws progress toward the next level in screen space.
type XPBar struct {
	bar *UIBar
}

// NewXPBar creates the XP bar with fixed layout and colors.
func NewXPBar() *XPBar {
	b := NewUIBar(10, 30, 200, 10)
	b.FgColor = color.RGBA{R: 255, G: 200, B: 64, A: 255}
	return &XPBar{bar: b}
}

// Draw sets fill from current XP toward XPBaseLevel * level and renders the bar.
func (x *XPBar) Draw(screen *ebiten.Image, xp int, level int) {
	if x == nil || x.bar == nil {
		return
	}
	need := GameSettings.XPBaseLevel * level
	xpProgress := 0.0
	if need > 0 {
		xpProgress = float64(xp) / float64(need)
		if xpProgress > 1 {
			xpProgress = 1
		}
	}
	x.bar.SetProgress(xpProgress)
	x.bar.Draw(screen)
}

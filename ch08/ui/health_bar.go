package ui

import "github.com/hajimehoshi/ebiten/v2"

// HealthBar draws the player HP bar in screen space.
type HealthBar struct {
	bar *UIBar
}

// NewHealthBar creates the health bar with fixed layout (position and size).
func NewHealthBar() *HealthBar {
	return &HealthBar{
		bar: NewUIBar(10, 10, 200, 14),
	}
}

// Draw updates the fill from current HP ratio (0..1) and renders the bar.
func (h *HealthBar) Draw(screen *ebiten.Image, hpRatio float64) {
	if h == nil || h.bar == nil {
		return
	}
	h.bar.SetProgress(hpRatio)
	h.bar.Draw(screen)
}

package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	PopupDuration  = 72        // frames
	PopupRisePxPerFrame = 0.6666666666666666 // vertical motion per Update step
)

// UIPopup shows text that rises and fades over time.
type UIPopup struct {
	Text      string
	X, Y      float64
	StartY    float64
	Elapsed   float64
	Duration  float64
	RiseSpeed float64
	Visible   bool
}

// NewUIPopup creates a popup.
func NewUIPopup(text string, duration, riseSpeed float64) *UIPopup {
	return &UIPopup{
		Text:      text,
		Duration:  duration,
		RiseSpeed: riseSpeed,
	}
}

// Show displays the popup at screen position (x, y).
func (p *UIPopup) Show(x, y float64) {
	p.X = x
	p.Y = y
	p.StartY = y
	p.Elapsed = 0
	p.Visible = true
}

// Update advances the popup. Returns false when expired.
func (p *UIPopup) Update() bool {
	if !p.Visible {
		return false
	}
	p.Elapsed += 1
	p.Y = p.StartY - p.RiseSpeed*p.Elapsed
	if p.Elapsed >= p.Duration {
		p.Visible = false
		return false
	}
	return true
}

// Draw renders the popup text.
func (p *UIPopup) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}
	face := text.NewGoXFace(basicfont.Face7x13)
	w, _ := text.Measure(p.Text, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(p.X-float64(w)/2, p.Y)
	alpha := 1.0
	if p.Elapsed > p.Duration*0.7 {
		alpha = (p.Duration - p.Elapsed) / (p.Duration * 0.3)
	}
	opts.ColorScale.Reset()
	opts.ColorScale.Scale(1, 1, 1, float32(alpha))
	opts.ColorScale.ScaleWithColor(color.RGBA{255, 220, 100, 255})
	text.Draw(screen, p.Text, face, opts)
}

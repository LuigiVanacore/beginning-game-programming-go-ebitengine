package ui

import (
	. "book/code/ch12/internal/core"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// GameResultOverlay draws centered game-result text (GAME OVER).
// The image is created lazily on first draw.
type GameResultOverlay struct {
	img     *ebiten.Image
	once    sync.Once
	message string
	clr     color.RGBA
}

// NewGameOverOverlay builds a red "GAME OVER" overlay.
func NewGameOverOverlay() *GameResultOverlay {
	return &GameResultOverlay{
		message: "GAME OVER",
		clr:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
	}
}

// Draw renders the overlay centered on screen; safe to call every frame.
func (o *GameResultOverlay) Draw(screen *ebiten.Image) {
	if o == nil {
		return
	}
	o.once.Do(func() {
		face := text.NewGoXFace(basicfont.Face7x13)
		w, h := text.Measure(o.message, face, 0)
		o.img = ebiten.NewImage(int(w)+2, int(h)+2)
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(1, 1)
		opts.ColorScale.ScaleWithColor(o.clr)
		text.Draw(o.img, o.message, face, opts)
	})
	if o.img == nil {
		return
	}
	b := o.img.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -10)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(float64(GameSettings.ScreenWidth)/2, float64(GameSettings.ScreenHeight)/2)
	screen.DrawImage(o.img, op)
}

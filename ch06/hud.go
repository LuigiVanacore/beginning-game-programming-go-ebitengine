package game

import (
	. "book/code/ch06/internal/core"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// HUD groups the interface into an overlay above the game world.
// In this chapter it handles only the game over; you can add bars, texts, menus, etc.
type HUD struct {
	gameOver gameOverOverlay
}

// NewHUD creates the HUD (UI state for this game session).
func NewHUD() *HUD {
	return &HUD{}
}

// DrawGameOver draws the centered "GAME OVER" text (rasterized only once).
func (h *HUD) DrawGameOver(screen *ebiten.Image) {
	if h == nil {
		return
	}
	h.gameOver.draw(screen)
}

// gameOverOverlay: lazy-init the image holding the text.
type gameOverOverlay struct {
	img  *ebiten.Image
	once sync.Once
}

func (o *gameOverOverlay) draw(screen *ebiten.Image) {
	o.once.Do(func() {
		face := text.NewGoXFace(basicfont.Face7x13)
		w, h := text.Measure("GAME OVER", face, 0)
		o.img = ebiten.NewImage(int(w)+2, int(h)+2)
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(1, 1)
		opts.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 255})
		text.Draw(o.img, "GAME OVER", face, opts)
	})
	if o.img == nil {
		return
	}
	b := o.img.Bounds()
	iw := float64(b.Dx())
	ih := float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-iw/2, -ih/2)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(float64(GameSettings.ScreenWidth)/2, float64(GameSettings.ScreenHeight)/2)
	screen.DrawImage(o.img, op)
}

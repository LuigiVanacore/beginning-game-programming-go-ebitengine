package game

import (
	. "book/code/ch07/internal/core"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// HUD overlay: for Ch7, only the game-over screen.
type HUD struct {
	gameOver gameOverOverlay
}

func NewHUD() *HUD {
	return &HUD{}
}

func (h *HUD) DrawGameOver(screen *ebiten.Image) {
	if h == nil {
		return
	}
	h.gameOver.draw(screen)
}

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
		opts.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		text.Draw(o.img, "GAME OVER", face, opts)
	})
	b := o.img.Bounds()
	iw := float64(b.Dx())
	ih := float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-iw/2, -ih/2)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(float64(GameSettings.ScreenWidth)/2, float64(GameSettings.ScreenHeight)/2)
	screen.DrawImage(o.img, op)
}

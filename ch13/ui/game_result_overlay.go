package ui

import (
	. "book/code/ch13/internal/core"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

// Button geometry for the game-over "New Game" button, in screen space.
const (
	newGameButtonWidth   = 160.0
	newGameButtonHeight  = 36.0
	newGameButtonOffsetY = 44.0 // distance below screen center, under the GAME OVER text
)

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

	// "New Game" button under the message. Hovering brightens it, matching the
	// upgrade-panel and menu button style.
	bx, by, bw, bh := o.NewGameButtonBounds()
	fill := color.RGBA{60, 100, 160, 255}
	border := color.RGBA{100, 150, 220, 255}
	mx, my := ebiten.CursorPosition()
	if o.NewGameButtonContains(float64(mx), float64(my)) {
		fill = color.RGBA{85, 135, 210, 255}
		border = color.RGBA{140, 190, 255, 255}
	}
	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(bw), float32(bh), fill, true)
	vector.StrokeRect(screen, float32(bx), float32(by), float32(bw), float32(bh), 1, border, true)
	DrawLabel(screen, UILabel{
		Text:               "NEW GAME",
		X:                  bx + bw/2,
		Y:                  by,
		Color:              color.RGBA{255, 255, 255, 255},
		CenterHorizontally: true,
	}, bw, bh)
}

// NewGameButtonBounds returns the screen-space rect (x, y, w, h) of the New Game button.
func (o *GameResultOverlay) NewGameButtonBounds() (x, y, w, h float64) {
	x = float64(GameSettings.ScreenWidth)/2 - newGameButtonWidth/2
	y = float64(GameSettings.ScreenHeight)/2 + newGameButtonOffsetY
	return x, y, newGameButtonWidth, newGameButtonHeight
}

// NewGameButtonContains reports whether the screen point (sx, sy) is inside the New Game button.
func (o *GameResultOverlay) NewGameButtonContains(sx, sy float64) bool {
	x, y, w, h := o.NewGameButtonBounds()
	return sx >= x && sx <= x+w && sy >= y && sy <= y+h
}

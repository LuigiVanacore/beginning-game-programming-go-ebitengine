package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds the game state and implements the ebiten.Game interface.
type Game struct {
	gopherImage *ebiten.Image
}

// NewGame creates a new Game with the given gopher image.
func NewGame(gopherImage *ebiten.Image) *Game {
	return &Game{gopherImage: gopherImage}
}

// Update is called every tick (default 60 times per second).
func (g *Game) Update() error {
	return nil
}

// Draw is called every frame to render the game to the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	gopherWidth := float64(g.gopherImage.Bounds().Dx())
	gopherHeight := float64(g.gopherImage.Bounds().Dy())
	x := (float64(Settings.ScreenWidth) - gopherWidth) / 2
	y := (float64(Settings.ScreenHeight) - gopherHeight) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(g.gopherImage, op)
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Settings.ScreenWidth, Settings.ScreenHeight
}

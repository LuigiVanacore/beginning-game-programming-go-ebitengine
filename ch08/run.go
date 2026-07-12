package game


import (
	. "book/code/ch08/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run creates the game and starts the Ebitengine loop.
func Run() error {
	g := NewGame()
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 8: UI, Health, XP, Level Up")
	return ebiten.RunGame(g)
}

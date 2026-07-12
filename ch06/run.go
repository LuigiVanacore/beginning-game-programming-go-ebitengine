package game


import (
	. "book/code/ch06/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run creates the game and starts the Ebitengine loop.
func Run() error {
	g := NewGame()
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 6: Enemy and Collisions")
	return ebiten.RunGame(g)
}

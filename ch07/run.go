package game


import (
	. "book/code/ch07/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run creates the game and starts the Ebitengine loop.
func Run() error {
	g := NewGame()
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 7: Weapons and Projectiles")
	return ebiten.RunGame(g)
}

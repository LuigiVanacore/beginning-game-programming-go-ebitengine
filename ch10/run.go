package game


import (
	. "book/code/ch10/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run creates the game and starts the Ebitengine loop.
func Run() error {
	game := NewGame()
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 10: Weapon Upgrade UI")
	return ebiten.RunGame(game)
}

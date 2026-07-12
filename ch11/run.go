package game


import (
	. "book/code/ch11/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run creates the game and starts the Ebitengine loop.
func Run() error {
	game := NewGame()
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 11: Gopher Survivor")
	return ebiten.RunGame(game)
}

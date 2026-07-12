package game

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	. "book/code/ch02/internal/core"
	"book/code/ch02/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// Run loads assets, builds the scene, and starts the game loop.
func Run() error {
	decoded, _, err := image.Decode(bytes.NewReader(assets.GopherPNG))
	if err != nil {
		return fmt.Errorf("failed to decode embedded gopher image: %w", err)
	}
	img := ebiten.NewImageFromImage(decoded)

	// Create engine and world
	engine := NewEngine()
	world := engine.World()

	// Create a Node2D container and a sprite as its child (centerPivot so position is the sprite center)
	logo := NewNode2D("logo")
	logo.SetPosition(320, 240)
	sprite := NewSprite("gopher", img, 0, true)
	sprite.SetPosition(0, 0)
	logo.AddChildren(sprite)
	world.AddNodeToDefaultLayer(logo)

	// Run the game
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 2: Scene Graph Framework")

	game := &Game{
		engine: engine,
		logo:   logo,
	}
	return ebiten.RunGame(game)
}

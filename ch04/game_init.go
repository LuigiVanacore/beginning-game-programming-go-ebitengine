package game

import (
	ast "book/code/ch04/assets"
	. "book/code/ch04/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	floorTextureKey = "floor_tile"
	playerTextureKey = "player"
	floorName       = "floor"
	floorLayerIndex = 0
)

func loadTextures(rm *ResourceManager) {
	if rm == nil {
		return
	}
	rm.UseEmbeddedFS(ast.FS)
	rm.LoadTexture(ast.FloorTile, floorTextureKey)
	rm.LoadTexture(ast.Player, playerTextureKey)
}

func registerPlayerInput(inp *InputManager) {
	if inp == nil {
		return
	}
	// Arrow keys and WASD both drive movement: an action can bind multiple buttons.
	inp.AddAction("move_up", NewKeyRawInputButton(ebiten.KeyArrowUp))
	inp.AddAction("move_up", NewKeyRawInputButton(ebiten.KeyW))
	inp.AddAction("move_down", NewKeyRawInputButton(ebiten.KeyArrowDown))
	inp.AddAction("move_down", NewKeyRawInputButton(ebiten.KeyS))
	inp.AddAction("move_left", NewKeyRawInputButton(ebiten.KeyArrowLeft))
	inp.AddAction("move_left", NewKeyRawInputButton(ebiten.KeyA))
	inp.AddAction("move_right", NewKeyRawInputButton(ebiten.KeyArrowRight))
	inp.AddAction("move_right", NewKeyRawInputButton(ebiten.KeyD))
}

func setupFloor(world *World, rm *ResourceManager) {
	if world == nil || rm == nil {
		return
	}
	floorTex, _ := rm.GetTexture(floorTextureKey)
	floor := NewSprite(floorName, floorTex, floorLayerIndex, false)
	floor.SetPosition(0, 0)
	floor.SetPivot(0, 0)
	floor.SetScale(float64(GameSettings.ScreenWidth)/float64(GameSettings.TileSize), float64(GameSettings.ScreenHeight)/float64(GameSettings.TileSize))
	world.AddNodeToLayer(floor, floorLayerIndex)
}

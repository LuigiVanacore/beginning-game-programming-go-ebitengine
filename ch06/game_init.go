package game

import (
	ast "book/code/ch06/assets"
	. "book/code/ch06/internal/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// sceneConfig groups the static scene parameters (tilemap).
type sceneConfig struct {
	tilemapNodeName string
}

var defaultSceneConfig = sceneConfig{
	tilemapNodeName: "tilemap",
}

func loadTextures(rm *ResourceManager) {
	if rm == nil {
		return
	}
	rm.UseEmbeddedFS(ast.FS)
	rm.LoadTexture(ast.Spritesheet, "spritesheet")
	rm.LoadTexture(ast.Player, "player")
	rm.LoadTexture(ast.Enemy, "enemy")
}

// registerPlayerInput binds the movement actions to the arrow keys (RawInputButton).
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

func setupTilemap(world *World, rm *ResourceManager, cfg sceneConfig) {
	if world == nil || rm == nil {
		return
	}
	spritesheetTex, _ := rm.GetTexture("spritesheet")
	// The three floor tiles are stacked with a 1px transparent gap, so the vertical
	// cell spacing is 1 (offsetY). offsetX stays 0 (a single column).
	tileset := NewTileset(spritesheetTex, GameSettings.TileSize, GameSettings.TileSize, 0, 1)
	cam := world.Camera()
	rawMap, err := ast.FS.ReadFile(ast.FloorMap)
	if err != nil {
		log.Fatalf("floor.map: %v", err)
	}
	pattern, err := ParseTilePattern(string(rawMap), FloorTileIndex)
	if err != nil {
		log.Fatalf("floor.map: %v", err)
	}
	tilemap := NewInfiniteTilemapNode(
		cfg.tilemapNodeName, tileset, pattern, DrawLayerBackground, cam, 1,
	)
	world.AddNodeToLayer(tilemap, DrawLayerBackground)
}

// wirePlayerGameOver sets the callback: contact with an enemy -> game over.
func wirePlayerGameOver(player *Player, game *Game) {
	if player == nil || game == nil {
		return
	}
	player.Collider.SetOnCollide(func(other *Collider) {
		game.gameOver = true
	})
}

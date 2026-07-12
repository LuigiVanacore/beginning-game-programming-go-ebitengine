package game

import (
	ast "book/code/ch07/assets"
	. "book/code/ch07/internal/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// sceneConfig holds the scene parameters (tilemap).
type sceneConfig struct {
	tilemapNodeName string
}

var defaultSceneConfig = sceneConfig{
	tilemapNodeName: NameTilemap,
}

func loadTextures(rm *ResourceManager) {
	if rm == nil {
		return
	}
	rm.UseEmbeddedFS(ast.FS)
	rm.LoadTexture(ast.Spritesheet, SpritesheetTexture)
	rm.LoadTexture(ast.Player, PlayerTexture)
	rm.LoadTexture(ast.Enemy, EnemyTexture)
	rm.LoadTexture(ast.Cursor, CursorTexture)
	rm.LoadTexture(ast.BloodyKnife, BloodyKnifeTexture)
}

func registerPlayerInput(inp *InputManager) {
	if inp == nil {
		return
	}
	// Arrow keys and WASD both drive movement: an action can bind multiple buttons.
	inp.AddAction(ActionMoveUp, NewKeyRawInputButton(ebiten.KeyArrowUp))
	inp.AddAction(ActionMoveUp, NewKeyRawInputButton(ebiten.KeyW))
	inp.AddAction(ActionMoveDown, NewKeyRawInputButton(ebiten.KeyArrowDown))
	inp.AddAction(ActionMoveDown, NewKeyRawInputButton(ebiten.KeyS))
	inp.AddAction(ActionMoveLeft, NewKeyRawInputButton(ebiten.KeyArrowLeft))
	inp.AddAction(ActionMoveLeft, NewKeyRawInputButton(ebiten.KeyA))
	inp.AddAction(ActionMoveRight, NewKeyRawInputButton(ebiten.KeyArrowRight))
	inp.AddAction(ActionMoveRight, NewKeyRawInputButton(ebiten.KeyD))
}

func setupTilemap(world *World, rm *ResourceManager, cfg sceneConfig) {
	if world == nil || rm == nil {
		return
	}
	spritesheetTex, _ := rm.GetTexture(SpritesheetTexture)
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

// createSession builds the player, cursor, and weapon after the tilemap.
func createSession(engine *Engine, rm *ResourceManager, world *World) (
	player *Player,
	cursor *Cursor,
	weapon *BloodyKnifeWeapon,
) {
	cam := world.Camera()
	player = NewPlayer(engine)
	cursorTex, _ := rm.GetTexture(CursorTexture)
	cursor = NewCursor(NameCursor, cursorTex, 2, cam)
	world.AddNodeToLayer(cursor, DrawLayerCursor)
	world.Camera().SetFollow(player)

	knifeTex, _ := rm.GetTexture(BloodyKnifeTexture)
	pool := NewProjectilePool(engine, knifeTex, KnifeRadius,
		KnifeProjectileSpeedPxPerFrame, DrawLayerPlayer, NameBloodyKnife, ProjectilePoolSize)
	weapon = NewBloodyKnifeWeapon(pool, KnifeCooldown)
	return player, cursor, weapon
}

func wirePlayerGameOver(player *Player, game *Game) {
	if player == nil || game == nil {
		return
	}
	player.Collider.SetOnCollide(func(other *Collider) {
		game.gameOver = true
	})
}

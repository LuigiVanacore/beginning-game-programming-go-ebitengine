package game

import (
	ast "book/code/ch09/assets"
	. "book/code/ch09/internal/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	rm.LoadTexture(ast.Potion, PotionTexture)
	rm.LoadTexture(ast.SacredBook, SacredBookTexture)
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

func createSession(engine *Engine, rm *ResourceManager, world *World) (
	player *Player,
	cursor *Cursor,
) {
	cam := world.Camera()
	player = NewPlayer(engine)
	cursorTex, _ := rm.GetTexture(CursorTexture)
	cursor = NewCursor(NameCursor, cursorTex, 2, cam)
	player.SetCursor(cursor)
	world.AddNodeToLayer(cursor, DrawLayerCursor)
	world.Camera().SetFollow(player)

	attachPlayerWeapons(engine, rm, player)
	return player, cursor
}

// attachPlayerWeapons adds weapon Node2D nodes as children of player.WeaponsRoot (Ch9).
func attachPlayerWeapons(engine *Engine, rm *ResourceManager, player *Player) {
	if player == nil || player.WeaponsRoot == nil {
		return
	}
	knifeTex, _ := rm.GetTexture(BloodyKnifeTexture)
	pool := NewProjectilePool(engine, knifeTex, KnifeRadius,
		KnifeProjectileSpeedPxPerFrame, DrawLayerPlayer, NameBloodyKnife, ProjectilePoolSize)
	knife := NewBloodyKnifeWeapon(pool, KnifeCooldown)
	player.MountWeapon(knife)

	if tex, ok := rm.GetTexture(SacredBookTexture); ok {
		player.MountWeapon(NewSacredBook(engine, player, OrbitWeaponDistance, tex))
	}
	player.MountWeapon(NewHolyShieldWeapon(engine, player, HolyShieldRadius))
}

// wirePlayerCallbacks registers enemy damage, weapon combat queues, and orb collection.
func wirePlayerCallbacks(player *Player, game *Game) {
	if player == nil || game == nil {
		return
	}
	player.SetWeaponHit(func(a, b *Collider) {
		if a != nil {
			game.removalQueue = append(game.removalQueue, a)
		}
		if b != nil {
			game.removalQueue = append(game.removalQueue, b)
		}
	})
	player.Collider.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() == LayerEnemy {
			game.player.DamageFromEnemyContact()
		}
	})
	if player.PickupCollider != nil {
		player.PickupCollider.SetOnCollide(func(other *Collider) {
			if other.GetCollisionMask().GetIdentity() != LayerPickup {
				return
			}
			for _, orb := range game.pickups.Orbs() {
				if orb.Col == other {
					game.pickups.QueueOrbCollect(orb)
					return
				}
			}
			for _, pot := range game.pickups.Potions() {
				if pot.Col == other {
					game.pickups.QueuePotionCollect(pot)
					return
				}
			}
		})
	}
}

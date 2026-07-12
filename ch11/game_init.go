package game

import (
	ast "book/code/ch11/assets"
	. "book/code/ch11/internal/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// registerPlayerInput binds the arrow keys and WASD to the movement actions.
// Both key sets drive the same action, so an action can fire from multiple buttons.
// Carried over unchanged from Chapter 10.
func registerPlayerInput(inp *InputManager) {
	if inp == nil {
		return
	}
	inp.AddAction(ActionMoveUp, NewKeyRawInputButton(ebiten.KeyArrowUp))
	inp.AddAction(ActionMoveUp, NewKeyRawInputButton(ebiten.KeyW))
	inp.AddAction(ActionMoveDown, NewKeyRawInputButton(ebiten.KeyArrowDown))
	inp.AddAction(ActionMoveDown, NewKeyRawInputButton(ebiten.KeyS))
	inp.AddAction(ActionMoveLeft, NewKeyRawInputButton(ebiten.KeyArrowLeft))
	inp.AddAction(ActionMoveLeft, NewKeyRawInputButton(ebiten.KeyA))
	inp.AddAction(ActionMoveRight, NewKeyRawInputButton(ebiten.KeyArrowRight))
	inp.AddAction(ActionMoveRight, NewKeyRawInputButton(ebiten.KeyD))
}

// sceneConfig groups the scene-construction parameters (carried over from Chapter 10).
type sceneConfig struct {
	tilemapNodeName string
}

var defaultSceneConfig = sceneConfig{
	tilemapNodeName: NameTilemap,
}

// setupTilemap builds the infinite floor from spritesheet.png and floor.map and adds
// it to the background layer. It is the same helper as Chapter 10; the survival layer
// this chapter adds does not touch the floor.
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

// loadTextures registers all textures used by the chapter (called once from NewGame).
func loadTextures(rm *ResourceManager) {
	if rm == nil {
		return
	}
	rm.UseEmbeddedFS(ast.FS)
	rm.LoadTexture(ast.Spritesheet, SpritesheetTexture)
	rm.LoadTexture(ast.Player, PlayerTexture)
	rm.LoadTexture(ast.Enemy, EnemyTexture)
	rm.LoadTexture(ast.Ghost, GhostTexture)
	rm.LoadTexture(ast.Spider, SpiderTexture)
	rm.LoadTexture(ast.Bat, BatTexture)
	rm.LoadTexture(ast.DarkWizard, DarkWizardTexture)
	rm.LoadTexture(ast.Cyclops, CyclopsTexture)
	rm.LoadTexture(ast.Cursor, CursorTexture)
	rm.LoadTexture(ast.BloodyKnife, BloodyKnifeTexture)
	rm.LoadTexture(ast.Potion, PotionTexture)
	rm.LoadTexture(ast.SacredBook, SacredBookTexture)
	rm.LoadTexture(ast.FlyingAxe, FlyingAxeTexture)
	rm.LoadTexture(ast.Armor, "armor")
	rm.LoadTexture(ast.Boots, "boots")
	rm.LoadTexture(ast.Gem, "gem")
	rm.LoadTexture(ast.Skull, "skull")
	rm.LoadTexture(ast.Ring, "ring")
	rm.LoadTexture(ast.Shield, HolyShieldTexture)
}

// wirePlayerCallbacks connects damage, weapon hit, and pickup callbacks to the Game.
// Called once from NewGame after the Game struct is fully initialised.
func wirePlayerCallbacks(player *Player, game *Game) {
	if player == nil || game == nil {
		return
	}
	player.SetWeaponHit(func(proj, enemyCol *Collider, dmg float64) {
		// The projectile is always spent on impact; queue it for release.
		if proj != nil {
			game.removalQueue = append(game.removalQueue, proj)
		}
		if enemyCol == nil {
			return
		}
		// Enemies now have hit points: subtract the weapon's damage and only
		// queue the enemy for removal once its HP reaches zero.
		enemy := game.enemyManager.FindByCollider(enemyCol)
		if enemy == nil {
			return
		}
		enemy.TakeDamage(dmg)
		if enemy.IsDead() {
			game.removalQueue = append(game.removalQueue, enemyCol)
		}
	})
	player.Collider.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() != LayerEnemy {
			return
		}
		// Each kind hits for a different amount, scaled by the difficulty tier.
		if enemy := game.enemyManager.FindByCollider(other); enemy != nil {
			game.player.DamageFromEnemyContact(enemy.Damage)
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

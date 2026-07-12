package game

import (
	. "book/code/ch07/internal/core"
	en "book/code/ch07/enemy"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds the session state: player, weapons, enemies, HUD.
type Game struct {
	engine        *Engine
	world         *World
	player        *Player
	cursor        *Cursor
	weaponManager *WeaponManager
	enemyManager  *en.EnemyManager
	removalQueue  []*Collider // colliders to remove after CheckCollision (this frame)
	gameOver      bool
	hud           *HUD
}

func NewGame() *Game {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)

	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)

	player, cursor, weapon := createSession(engine, rm, world)

	game := &Game{
		engine:        engine,
		world:         world,
		player:        player,
		cursor:        cursor,
		weaponManager: NewWeaponManager(weapon),
		enemyManager:  en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false), nil)),
		removalQueue:  make([]*Collider, 0),
		gameOver:      false,
		hud:           NewHUD(),
	}
	wirePlayerGameOver(player, game)
	return game
}

func (g *Game) Update() error {
	if g.gameOver {
		return g.engine.Update()
	}

	playerX, playerY := g.updatePlay()
	g.enemyManager.Update(playerX, playerY, g.engine)
	cursorX, cursorY := g.cursor.GetWorldPosition().X(), g.cursor.GetWorldPosition().Y()
	g.weaponManager.Update(playerX, playerY, cursorX, cursorY, func(projectileCol, hitCollider *Collider) {
		g.removalQueue = append(g.removalQueue, projectileCol, hitCollider)
	})
	g.engine.CollisionManager().CheckCollision()
	g.processRemovals()

	return g.engine.Update()
}

func (g *Game) updatePlay() (playerX, playerY float64) {
	g.player.Update()
	g.world.Camera().Update()
	g.cursor.Update()
	p := g.player.GetWorldPosition()
	return p.X(), p.Y()
}

func (g *Game) processRemovals() {
	seen := make(map[*Collider]struct{}, len(g.removalQueue))
	for _, collider := range g.removalQueue {
		g.processQueuedCollider(collider, seen)
	}
	g.clearRemovalQueue()
}

// shouldProcessRemoval returns false for nil or duplicate collider pointers this frame.
func shouldProcessRemoval(collider *Collider, seen map[*Collider]struct{}) bool {
	if collider == nil {
		return false
	}
	if _, dup := seen[collider]; dup {
		return false
	}
	seen[collider] = struct{}{}
	return true
}

// processQueuedCollider dispatches one enqueued collider to projectile, enemy, or fallback removal.
func (g *Game) processQueuedCollider(collider *Collider, seen map[*Collider]struct{}) {
	if !shouldProcessRemoval(collider, seen) {
		return
	}
	if g.tryReleaseQueuedProjectile(collider) {
		return
	}
	if g.removeQueuedEnemy(collider) {
		return
	}
	g.removeQueuedOrphanCollider(collider)
}

// tryReleaseQueuedProjectile returns the collider to the knife pool when it belongs to an active projectile.
func (g *Game) tryReleaseQueuedProjectile(collider *Collider) bool {
	return g.weaponManager.TryReleaseProjectileByCollider(collider)
}

// removeQueuedEnemy removes an enemy whose body collider was enqueued. Returns false if not an enemy.
func (g *Game) removeQueuedEnemy(collider *Collider) bool {
	enemy := g.enemyManager.FindByCollider(collider)
	if enemy == nil {
		return false
	}
	g.engine.CollisionManager().RemoveCollider(collider)
	g.world.RemoveNode(enemy)
	g.enemyManager.Remove(enemy)
	return true
}

// removeQueuedOrphanCollider removes a collider that is neither a live projectile nor an enemy.
// Rare in Chapter 7; dangerous for pooled projectile colliders if deduplication fails.
func (g *Game) removeQueuedOrphanCollider(collider *Collider) {
	g.world.RemoveNode(collider)
	g.engine.CollisionManager().RemoveCollider(collider)
}

func (g *Game) clearRemovalQueue() {
	g.removalQueue = g.removalQueue[:0]
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
	if g.gameOver {
		g.hud.DrawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

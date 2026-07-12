package game


import (
	. "book/code/ch09/internal/core"
	en "book/code/ch09/enemy"
	pkup "book/code/ch09/pickups"
	. "book/code/ch09/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game: the level session (enemies, pickups, HUD). The weapons live on the Player (a Node2D under WeaponsRoot).
type Game struct {
	engine            *Engine
	world             *World
	player            *Player
	cursor            *Cursor
	enemyManager      *en.EnemyManager
	removalQueue      []*Collider
	gameOver bool
	pickups  *pkup.PickupManager
	hud      *HUD
}

func NewGame() *Game {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)

	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)

	player, cursor := createSession(engine, rm, world)

	game := &Game{
		engine:            engine,
		world:             world,
		player:            player,
		cursor:            cursor,
		enemyManager:      en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false), nil)),
		removalQueue: make([]*Collider, 0),
		gameOver:     false,
		pickups:      pkup.NewPickupManager(),
		hud:          NewHUD(),
	}
	wirePlayerCallbacks(player, game)
	return game
}

func (g *Game) Update() error {
	if g.gameOver {
		return g.engine.Update()
	}
	if g.player.IsDead() {
		g.gameOver = true
		return g.engine.Update()
	}

	playerX, playerY := g.updatePlay()
	g.enemyManager.Update(playerX, playerY, g.engine)
	g.engine.CollisionManager().CheckCollision()
	g.processRemovals()
	g.pickups.CollectOrbs(g.world, g.engine.CollisionManager(), func(v int) {
		g.player.XP += v
	})
	g.pickups.CollectPotions(g.world, g.engine.CollisionManager(), func() {
		g.player.HP += PlayerMaxHP / 2
		if g.player.HP > PlayerMaxHP {
			g.player.HP = PlayerMaxHP
		}
	})
	g.updateLevelProgress()

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
	// Same-frame multi-hit: one projectile can overlap several enemies; each pair fires
	// onCollide, so the same *Collider may appear twice. Skip duplicates so we never
	// RemoveNode(projectileCollider) after pool Release (would detach collider from pooled Projectile).
	seen := make(map[*Collider]struct{}, len(g.removalQueue))
	for _, collider := range g.removalQueue {
		if collider == nil {
			continue
		}
		if _, dup := seen[collider]; dup {
			continue
		}
		seen[collider] = struct{}{}
		if g.player.TryReleaseProjectileByCollider(collider) {
			continue
		}
		enemy := g.enemyManager.FindByCollider(collider)
		if enemy != nil {
			g.spawnLootAtEnemyDeath(collider)
			g.engine.CollisionManager().RemoveCollider(collider)
			g.world.RemoveNode(enemy)
			g.enemyManager.Remove(enemy)
			continue
		}
		g.world.RemoveNode(collider)
		g.engine.CollisionManager().RemoveCollider(collider)
	}
	g.removalQueue = g.removalQueue[:0]
}

func (g *Game) spawnLootAtEnemyDeath(collider *Collider) {
	pos := collider.GetWorldPosition()
	g.pickups.AddOrb(pkup.CreateOrb(g.engine, pos.X(), pos.Y(), pkup.OrbXPValue))
	g.pickups.MaybeDropPotion(g.engine, pos.X(), pos.Y())
}

func (g *Game) updateLevelProgress() {
	xpNeeded := GameSettings.XPBaseLevel * g.player.Level
	if g.player.XP >= xpNeeded {
		g.player.XP -= xpNeeded
		g.player.Level++
		g.hud.ShowLevelUp(float64(GameSettings.ScreenWidth)/2, float64(GameSettings.ScreenHeight)/2-40)
	}
	g.hud.UpdateLevelPopup()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
	g.hud.DrawGameplay(screen, g.player.HP/PlayerMaxHP, g.player.XP, g.player.Level)
	if g.gameOver {
		g.hud.DrawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

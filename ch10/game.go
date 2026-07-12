package game

import (
	en "book/code/ch10/enemy"
	. "book/code/ch10/internal/core"
	pkup "book/code/ch10/pickups"
	. "book/code/ch10/ui"
)

// xpNeededForLevel returns XP required to reach the next level (exponential curve).
func xpNeededForLevel(level int) int {
	if level < 1 {
		return GameSettings.XPBaseLevel
	}
	f := float64(GameSettings.XPBaseLevel)
	for i := 1; i < level; i++ {
		f *= GameSettings.XPGrowthFactor
	}
	return int(f)
}

func upgradeChoicesFromOptions(opts []UpgradeOption, g *Game) ([]UpgradeChoice, []func()) {
	choices := make([]UpgradeChoice, len(opts))
	applies := make([]func(), len(opts))
	for i, o := range opts {
		o := o
		choices[i] = UpgradeChoice{WeaponName: o.WeaponName, UpgradeDesc: o.UpgradeDesc, IconKey: o.IconKey}
		applies[i] = func() { o.Apply(g) }
	}
	return choices, applies
}

// Game holds all session state. It is built on the ch09 base (Player, Cursor, WeaponManager,
// EnemyManager, PickupManager) and extends it with ch10's upgrade system and survival timer.
type Game struct {
	// --- Infrastructure ---
	engine  *Engine
	world   *World
	player  *Player
	cursor  *Cursor
	rm      *ResourceManager
	weapons *WeaponManager

	// --- Enemies & pickups ---
	enemyManager *en.EnemyManager
	removalQueue []*Collider
	pickups      *pkup.PickupManager

	// --- Game state ---
	gameOver      bool
	elapsedFrames int // survival time, counted up once per Update

	// --- UI ---
	hud *HUD
}

// NewGame creates and wires the full game session.
func NewGame() *Game {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)

	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)

	player, cursor := createSession(engine, rm, world)

	game := &Game{
		engine:  engine,
		world:   world,
		player:  player,
		cursor:  cursor,
		rm:      rm,
		weapons: NewWeaponManager(),

		enemyManager: en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false), nil)),
		removalQueue: make([]*Collider, 0),
		pickups:      pkup.NewPickupManager(),

		gameOver:      false,
		elapsedFrames: 0,

		hud: NewHUD(rm),
	}

	attachPlayerWeapons(engine, rm, player, game)

	wirePlayerCallbacks(player, game)
	return game
}

// Update runs one game frame.
func (g *Game) Update() error {
	if g.gameOver {
		return g.engine.Update()
	}
	if g.player.IsDead() {
		g.gameOver = true
		return g.engine.Update()
	}
	if g.hud.IsChoosingUpgrade() {
		g.hud.HandleInput()
		return g.engine.Update()
	}

	g.elapsedFrames++

	playerX, playerY := g.updatePlay()
	g.weapons.Update(g.engine, g.player, g.cursor)
	g.enemyManager.Update(playerX, playerY, g.engine)
	g.engine.CollisionManager().CheckCollision()
	g.processRemovals()
	g.pickups.CollectOrbs(g.world, g.engine.CollisionManager(), func(v int) {
		g.player.XP += int(float64(v) * g.player.XPBonusMult)
	})
	g.pickups.CollectPotions(g.world, g.engine.CollisionManager(), func() {
		g.player.HP += g.player.MaxHP / 2
		if g.player.HP > g.player.MaxHP {
			g.player.HP = g.player.MaxHP
		}
	})
	g.checkLevelUp()

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
	seen := make(map[*Collider]bool, len(g.removalQueue))
	for _, collider := range g.removalQueue {
		if seen[collider] {
			continue
		}
		seen[collider] = true
		if g.weapons.TryReleaseProjectileByCollider(collider) {
			continue
		}
		enemy := g.enemyManager.FindByCollider(collider)
		if enemy != nil {
			pos := collider.GetWorldPosition()
			g.pickups.AddOrb(pkup.CreateOrb(g.engine, pos.X(), pos.Y(), pkup.OrbXPValue))
			g.pickups.MaybeDropPotion(g.engine, g.rm, pos.X(), pos.Y())
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

func (g *Game) checkLevelUp() {
	need := xpNeededForLevel(g.player.Level)
	if g.player.XP < need {
		if !g.hud.IsChoosingUpgrade() {
			g.hud.UpdatePopup()
		}
		return
	}
	g.player.XP -= need
	g.player.Level++
	opts := PickRandomUpgrades(3)
	c, a := upgradeChoicesFromOptions(opts, g)
	g.hud.TriggerLevelUp(c, a)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

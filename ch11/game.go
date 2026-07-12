package game


import (
	. "book/code/ch11/internal/core"
	en "book/code/ch11/enemy"
	pkup "book/code/ch11/pickups"
	. "book/code/ch11/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds all session state: infrastructure, entities, weapons, upgrades, and HUD.
// It extends the ch10 base with the weapon unlock system (new in ch11).
type Game struct {
	// --- Infrastructure ---
	engine *Engine
	player *Player
	cursor  *Cursor
	weapons *WeaponManager

	// --- Enemies & pickups ---
	enemyManager *en.EnemyManager
	removalQueue []*Collider
	pickups      *pkup.PickupManager

	// --- Non-weapon upgrade multipliers (boots / gem; weapon stats live in WeaponManager) ---
	playerSpeedMult float64
	xpBonusMult     float64

	// --- Game state ---
	gameOver     bool
	elapsedFrames int // survival time, counted up once per Update

	// --- UI ---
	hud *HUD

	// upgradeCount caps all upgrade picks (weapons + bonus items)
	upgradeCount int
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
	g.enemyManager.Update(playerX, playerY, g.elapsedSeconds(), g.engine)
	if g.enemyManager.EnemiesGrewStronger() {
		g.hud.TriggerEnemyGrew(g.elapsedFrames)
	}
	g.engine.CollisionManager().CheckCollision()
	g.processRemovals()
	g.pickups.CollectOrbs(g.engine.World(), g.engine.CollisionManager(), func(v int) {
		g.player.XP += int(float64(v) * g.xpBonusMult)
	})
	g.pickups.CollectPotions(g.engine.World(), g.engine.CollisionManager(), func() {
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
	g.engine.World().Camera().Update()
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
			g.pickups.AddOrb(pkup.CreateOrb(g.engine, pos.X(), pos.Y(), enemy.XPValue))
			g.pickups.MaybeDropPotion(g.engine, g.engine.ResourceManager(), pos.X(), pos.Y())
			g.engine.CollisionManager().RemoveCollider(collider)
			g.engine.World().RemoveNode(enemy)
			g.enemyManager.Remove(enemy)
			continue
		}
		g.engine.World().RemoveNode(collider)
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
	opts := PickUpgradesForLevelUp(g)
	if len(opts) == 0 {
		return
	}
	c, a := upgradeChoicesFromOptions(opts, g)
	g.hud.TriggerLevelUp(c, a)
}

// elapsedSeconds converts the survival frame counter to seconds using the tick rate.
func (g *Game) elapsedSeconds() float64 {
	tps := ebiten.TPS()
	if tps < 1 {
		tps = 1
	}
	return float64(g.elapsedFrames) / float64(tps)
}

// Layout satisfies ebiten.Game.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

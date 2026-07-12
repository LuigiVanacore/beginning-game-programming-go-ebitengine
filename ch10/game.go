package game

import (
	en "book/code/ch10/enemy"
	. "book/code/ch10/internal/core"
	pkup "book/code/ch10/pickups"
	. "book/code/ch10/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	hud             *HUD
	gameOverOverlay *GameResultOverlay
}

// NewGame creates the Game and builds its first session.
func NewGame() *Game {
	g := &Game{}
	g.start()
	return g
}

// start builds a fresh session on g: a new engine, world, player, weapons, enemies,
// pickups, and HUD. NewGame calls it once; restart calls it again after a game over,
// so a New Game click begins from a clean state. Because it assigns to the existing g
// rather than allocating a new *Game, the callbacks wired below keep pointing at the
// same Game that Ebitengine drives.
func (g *Game) start() {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)

	world := engine.World()
	registerPlayerInput(engine.Input())
	setupTilemap(world, rm, defaultSceneConfig)

	player, cursor := createSession(engine, rm, world)

	g.engine = engine
	g.world = world
	g.player = player
	g.cursor = cursor
	g.rm = rm
	g.weapons = NewWeaponManager()

	g.enemyManager = en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false), nil))
	g.removalQueue = make([]*Collider, 0)
	g.pickups = pkup.NewPickupManager()

	g.gameOver = false
	g.elapsedFrames = 0

	g.hud = NewHUD(rm)
	g.gameOverOverlay = NewGameOverOverlay()

	attachPlayerWeapons(engine, rm, player, g)
	wirePlayerCallbacks(player, g)
}

// restart begins a fresh run in place. Ebitengine keeps calling Update and Draw on the
// same *Game pointer, so rebuilding g's fields is all it takes to start over.
func (g *Game) restart() {
	g.start()
}

// Update runs one game frame.
func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if g.gameOverOverlay.NewGameButtonContains(float64(mx), float64(my)) {
				g.restart()
			}
		}
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

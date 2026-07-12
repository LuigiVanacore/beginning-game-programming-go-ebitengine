package game

import (
	en "book/code/ch11/enemy"
	. "book/code/ch11/internal/core"
	pkup "book/code/ch11/pickups"
	. "book/code/ch11/ui"
	"math/rand"
)

// NewGame creates the Game and builds its first session.
func NewGame() *Game {
	g := &Game{}
	g.start()
	return g
}

// start builds a fresh session on g: engine, textures, infinite tilemap, and one random
// weapon. NewGame calls it once; restart calls it again after a game over, so a New Game
// click begins from a clean state. Assigning to the existing g (rather than allocating a
// new *Game) keeps the speed-multiplier pointer and the player callbacks valid across
// restarts.
func (g *Game) start() {
	engine := NewEngine()
	rm := engine.ResourceManager()
	loadTextures(rm)
	world := engine.World()
	registerPlayerInput(engine.Input())

	setupTilemap(world, rm, defaultSceneConfig)

	player := NewPlayer(engine)

	cam := world.Camera()
	cursorTex, _ := rm.GetTexture(CursorTexture)
	cursor := NewCursor(NameCursor, cursorTex, 2, cam)
	player.SetCursor(cursor)
	world.AddNodeToLayer(cursor, DrawLayerCursor)
	cam.SetFollow(player)

	wm := NewWeaponLoadout(engine, player)

	g.engine = engine
	g.player = player
	g.cursor = cursor
	g.weapons = wm

	g.enemyManager = en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false)))
	g.removalQueue = make([]*Collider, 0)
	g.pickups = pkup.NewPickupManager()

	g.playerSpeedMult = 1.0
	g.xpBonusMult = 1.0

	g.gameOver = false
	g.elapsedFrames = 0
	g.upgradeCount = 0

	g.hud = NewHUD(engine.ResourceManager())
	g.gameOverOverlay = NewGameOverOverlay()

	// Wire speed multiplier pointer after game struct is stable.
	player.SetSpeedMult(&g.playerSpeedMult)

	// Mount one random starting weapon (the others are unlocked via upgrades).
	switch rand.Intn(4) {
	case WeaponKnife:
		wm.KnifeUnlocked = true
		wm.Mount(wm.Knife, player.WeaponsRoot)
	case WeaponFlyingAxe:
		wm.FlyingAxeUnlocked = true
		wm.Mount(wm.Axe, player.WeaponsRoot)
	case WeaponSacredBook:
		wm.SacredBookUnlocked = true
		if wm.SacredBook != nil {
			wm.Mount(wm.SacredBook, player.WeaponsRoot)
		}
	default: // WeaponHolyShield
		wm.HolyShieldUnlocked = true
		wm.Mount(wm.HolyShield, player.WeaponsRoot)
	}

	wirePlayerCallbacks(player, g)
}

// restart begins a fresh run in place. Ebitengine keeps calling Update and Draw on the
// same *Game pointer, so rebuilding g's fields is all it takes to start over.
func (g *Game) restart() {
	g.start()
}

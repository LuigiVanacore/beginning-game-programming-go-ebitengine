package game

import (
	en "book/code/ch12/enemy"
	. "book/code/ch12/internal/core"
	pkup "book/code/ch12/pickups"
	. "book/code/ch12/ui"
	"math/rand"
)

// NewGame creates the Game and builds its first session. The state machine calls it
// for every new run — from the main menu and from the game-over screen — and hands the
// result to App.SetGame.
func NewGame() *Game {
	g := &Game{}
	g.start()
	return g
}

// start builds a fresh session on g: engine, textures, infinite tilemap, and one random
// weapon. It carries over from Chapter 11, where it also backed an in-place restart;
// Chapter 12 drops that restart because the state machine now starts fresh runs through
// NewGame instead. Music is owned by the game state's Enter/Exit, so start does not touch
// the track here.
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

package game

import (
	en "book/code/ch13/enemy"
	. "book/code/ch13/internal/core"
	pkup "book/code/ch13/pickups"
	. "book/code/ch13/ui"
	"math/rand"
)

// NewGame creates a fresh game session: engine, textures, infinite tilemap, and one random weapon.
func NewGame() *Game {
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

	game := &Game{
		engine:  engine,
		player:  player,
		cursor:  cursor,
		weapons: wm,

		enemyManager: en.NewEnemyManager(en.NewEnemySpawner(NewTimer(0, false))),
		removalQueue: make([]*Collider, 0),
		pickups:      pkup.NewPickupManager(),

		playerSpeedMult: 1.0,
		xpBonusMult:     1.0,

		gameOver:      false,
		elapsedFrames: 0,

		hud: NewHUD(engine.ResourceManager()),

		particles:    NewParticleSystem(),
		floatingText: NewFloatingTextSystem(),
	}

	// Wire speed multiplier pointer after game struct is stable.
	player.SetSpeedMult(&game.playerSpeedMult)

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

	wirePlayerCallbacks(player, game)

	// Music is owned by the game state's Enter/Exit, not by the session, so a new
	// session does not start the track here.
	return game
}

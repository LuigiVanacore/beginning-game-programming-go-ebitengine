package core

// Settings holds all tunable parameters for the game.
// Changing a value here is the only edit needed to adjust the configuration.
type Settings struct {
	ScreenWidth    int
	ScreenHeight   int
	TileSize       int
	XPBaseLevel    int     // XP needed for first level up (1→2)
	XPGrowthFactor float64 // each level requires this much more XP than the previous
}

// GameSettings is the single source of truth for all game parameters.
var GameSettings = Settings{
	ScreenWidth:    640,
	ScreenHeight:   480,
	TileSize:       16,
	XPBaseLevel:    30,
	XPGrowthFactor: 1.5,
}

// Weapon unlock system (new in ch11)
const MaxUpgrades = 10
const MaxStatLevelPerWeapon = 3

// Weapon IDs for random start (new in ch11)
const (
	WeaponKnife      = 0
	WeaponFlyingAxe  = 1
	WeaponSacredBook = 2
	WeaponHolyShield = 3
)

// Texture file paths live in package assets (embed); use assets.Tilemap etc. when loading.

const (
	DrawLayerBackground = 0
	DrawLayerPlayer     = 1
	DrawLayerCursor     = 2
)

// Texture resource keys (ResourceManager keys)
const (
	SpritesheetTexture = "spritesheet"
	PlayerTexture      = "player"
	EnemyTexture       = "enemy"
	CursorTexture      = "cursor"
	BloodyKnifeTexture = "bloody_knife"
	FlyingAxeTexture   = "flying_axe"
	PotionTexture      = "potion"
	SacredBookTexture  = "sacred_book"
	HolyShieldTexture  = "holy_shield" // new in ch11

	// Monster texture keys, ordered from the weakest to the strongest kind (new in ch11).
	GhostTexture      = "ghost"
	SpiderTexture     = "spider"
	BatTexture        = "bat"
	DarkWizardTexture = "dark_wizard"
	CyclopsTexture    = "cyclops"
)

// Input action names
const (
	ActionMoveUp    = "move_up"
	ActionMoveDown  = "move_down"
	ActionMoveLeft  = "move_left"
	ActionMoveRight = "move_right"
	ActionPause     = "pause" // new in ch12: leaves the gameplay state for the pause menu
)

// Options holds player-adjustable settings exposed by the Options menu (new in ch12).
// The state machine reads and mutates GameOptions; the audio subsystem consults the
// volume fields when it plays back, and Fullscreen is applied through ebiten.SetFullscreen.
type Options struct {
	MasterVolume float64 // 0.0 (muted) .. 1.0 (full)
	SFXEnabled   bool
	Fullscreen   bool
}

// GameOptions is the single source of truth for player-adjustable settings.
var GameOptions = Options{
	MasterVolume: 0.8,
	SFXEnabled:   true,
	Fullscreen:   false,
}

// Node names
const (
	NameTilemap       = "tilemap"
	NamePlayer        = "player"
	NamePlayerPickup  = "player_pickup"
	NamePlayerSprite  = "player_sprite"
	NameCursor        = "cursor"
	NameEnemy         = "enemy"
	NameEnemySprite   = "enemy_sprite"
	NameBloodyKnife   = "bloody_knife"
	NameFlyingAxe     = "flying_axe"
	NameCamera        = "camera"
	NameRoot          = "root"
	NameLayerRoot     = "layer_root"
	NameProjectileDef = "projectile"
	NameWeaponsRoot   = "weapons"
)

// FloorTileIndex maps the integer indices used in floor.map to
// tileset cell coordinates (tilesetCol, tilesetRow). The floor tileset
// (spritesheet.png) stacks three 16x16 variants vertically, so index 0 is the
// plain floor, index 1 adds crates, and index 2 adds scattered stones.
var FloorTileIndex = [][2]int{
	{0, 0}, // 0: plain floor
	{0, 1}, // 1: floor with crates
	{0, 2}, // 2: floor with stones
}

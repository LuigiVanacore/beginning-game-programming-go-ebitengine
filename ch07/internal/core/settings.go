package core

// Settings holds all tunable parameters for the game.
// Changing a value here is the only edit needed to adjust the configuration.
type Settings struct {
	ScreenWidth  int
	ScreenHeight int
	TileSize     int // size of one tile in pixels
}

// GameSettings is the single source of truth for all game parameters.
var GameSettings = Settings{
	ScreenWidth:  640,
	ScreenHeight: 480,
	TileSize:     16,
}

// Texture file paths live in package assets (embed); use assets.Spritesheet etc. when loading.

// Render layers (fixed indices, not tunable).
const (
	DrawLayerBackground = 0
	DrawLayerPlayer     = 1
	DrawLayerCursor     = 2
)

// Texture resource keys.
const (
	SpritesheetTexture = "spritesheet"
	PlayerTexture      = "player"
	EnemyTexture       = "enemy"
	CursorTexture      = "cursor"
	BloodyKnifeTexture = "bloody_knife"
)

// Input action names.
const (
	ActionMoveUp    = "move_up"
	ActionMoveDown  = "move_down"
	ActionMoveLeft  = "move_left"
	ActionMoveRight = "move_right"
)

// Node names.
const (
	NameTilemap       = "tilemap"
	NamePlayer        = "player"
	NamePlayerSprite  = "player_sprite"
	NameCursor        = "cursor"
	NameEnemy         = "enemy"
	NameEnemySprite   = "enemy_sprite"
	NameBloodyKnife   = "bloody_knife"
	NameCamera        = "camera"
	NameRoot          = "root"
	NameLayerRoot     = "layer_root"
	NameProjectileDef = "projectile"
)

// FloorTileIndex maps the integer indices used in floor.map to
// tileset cell coordinates (tilesetCol, tilesetRow).
var FloorTileIndex = [][2]int{{0, 0}, {0, 1}, {0, 2}}

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

// Texture file paths live in package assets (embed); use assets.FloorTile etc. when loading.

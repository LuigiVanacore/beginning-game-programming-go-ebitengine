package core

// Settings holds all tunable parameters for the game.
// Changing a value here is the only edit needed to adjust the configuration.
type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

// GameSettings is the single source of truth for all game parameters (exported for the game package).
var GameSettings = Settings{
	ScreenWidth:  640,
	ScreenHeight: 480,
}

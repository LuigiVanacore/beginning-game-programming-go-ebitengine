package game

// GameSettings holds all tunable parameters for the game.
// Changing a value here is the only edit needed to adjust the configuration.
type GameSettings struct {
	ScreenWidth  int
	ScreenHeight int
}

// Settings is the single source of truth for all game parameters.
var Settings = GameSettings{
	ScreenWidth:  640,
	ScreenHeight: 480,
}

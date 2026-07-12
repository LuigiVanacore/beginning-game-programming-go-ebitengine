package game

import (
	. "book/code/ch12/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// App implements ebiten.Game and owns the StateMachine. Gameplay sessions are
// created on demand by NewGame() when the player starts a run.
type App struct {
	sm   *StateMachine
	game *Game
}

// NewApp creates the app and state machine, registers the screens, and starts at
// the main menu. No gameplay session exists until the player picks "NEW GAME".
func NewApp() *App {
	loadAudio() // decode and register every sound once, before the loop starts
	app := &App{}
	app.sm = NewStateMachine(app)
	registerGameStates(app.sm)
	app.sm.SwitchTo(StateIDMainMenu)
	return app
}

// registerGameStates wires the menu, gameplay, pause, and options screens.
func registerGameStates(sm *StateMachine) {
	sm.Register(StateIDMainMenu, &StateMainMenuImpl{})
	sm.Register(StateIDGame, &StateGameImpl{})
	sm.Register(StateIDPause, &StatePauseImpl{})
	sm.Register(StateIDOptions, &StateOptionsImpl{})
}

// Run configures the window, applies saved options, and runs the Ebitengine loop.
// It replaces the package-level Run() bootstrap used up to Chapter 11.
func (a *App) Run() error {
	ebiten.SetWindowSize(GameSettings.ScreenWidth, GameSettings.ScreenHeight)
	ebiten.SetWindowTitle("Chapter 12: Gopher Survivor — State machine")
	ebiten.SetFullscreen(GameOptions.Fullscreen)

	err := ebiten.RunGame(a)
	if err == ErrQuit {
		return nil
	}
	return err
}

// SetGame implements core.AppContext. Passing nil discards the current session.
func (a *App) SetGame(g any) {
	if g == nil {
		a.game = nil
		return
	}
	a.game = g.(*Game)
}

// Gameplay implements core.AppContext. It returns the live session, or nil when
// none is running (for example while the player sits in the main menu).
func (a *App) Gameplay() GameSession {
	if a.game == nil {
		return nil
	}
	return a.game
}

// Update delegates to the state machine.
func (a *App) Update() error {
	return a.sm.Update()
}

// Draw paints a solid backdrop and then the current state.
func (a *App) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(GameSettings.ScreenWidth), float32(GameSettings.ScreenHeight),
		color.RGBA{15, 15, 30, 255}, true)
	a.sm.Draw(screen)
}

// Layout returns the logical screen size.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return GameSettings.ScreenWidth, GameSettings.ScreenHeight
}

package game

import (
	. "book/code/ch12/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// StateGameImpl runs the live gameplay session. The session lives on the App and is
// reached through the state machine's AppContext; the only state kept here is a click
// edge tracker for the game-over "New Game" button.
type StateGameImpl struct {
	mouseWasPressed bool
}

// Enter starts the looping run music and seeds the click tracker so a click that began
// on a previous screen does not leak in. Because the machine calls Exit when it leaves
// this state, the music plays only while the run is on screen.
func (s *StateGameImpl) Enter(sm *StateMachine) {
	s.mouseWasPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	Audio().PlayMusic(SoundMusic)
}

// Exit stops the music, so pausing, opening options, or returning to the menu all
// silence the track. Re-entering the state (for example, CONTINUE from pause)
// starts it again.
func (s *StateGameImpl) Exit(sm *StateMachine) {
	Audio().StopMusic()
}

func (s *StateGameImpl) Update(sm *StateMachine) error {
	g := sm.App().Gameplay()
	if g == nil {
		// No session (should not happen from the menu path); fall back to the menu.
		sm.SwitchTo(StateIDMainMenu)
		return nil
	}
	// When the run is over, the game-over overlay's New Game button starts a fresh run
	// through the machine — the same path the main menu uses — rather than the in-place
	// restart of Chapters 10 and 11. The overlay still draws the button; the click is
	// handled here so the state machine stays in charge of application flow.
	if gg, ok := g.(*Game); ok && gg.gameOver {
		if risingClick(&s.mouseWasPressed) {
			mx, my := cursorF()
			if gg.gameOverOverlay.NewGameButtonContains(mx, my) {
				sm.App().SetGame(NewGame())
				sm.SwitchTo(StateIDGame)
				return nil
			}
		}
		return g.Update()
	}
	if g.IsPausePressed() {
		sm.SwitchTo(StateIDPause)
		return nil
	}
	return g.Update()
}

func (s *StateGameImpl) Draw(sm *StateMachine, screen *ebiten.Image) {
	if g := sm.App().Gameplay(); g != nil {
		g.Draw(screen)
	}
}

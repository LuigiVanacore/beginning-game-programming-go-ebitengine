package game

import (
	. "book/code/ch13/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// StateGameImpl runs the live gameplay session. It owns no state of its own; the
// session lives on the App and is reached through the state machine's AppContext.
type StateGameImpl struct{}

// Enter starts the looping run music. Because the machine calls Exit when it leaves
// this state, the music plays only while the run is on screen.
func (s *StateGameImpl) Enter(sm *StateMachine) {
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

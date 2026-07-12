package game

import (
	. "book/code/ch12/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// StatePauseImpl freezes the run and overlays a menu: resume, options, restart, or exit.
// The gameplay session is kept alive on the App, so "CONTINUE" resumes exactly where
// the player left off.
type StatePauseImpl struct {
	mouseWasPressed bool
}

func (s *StatePauseImpl) Enter(sm *StateMachine) {
	s.mouseWasPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// Exit has nothing to tear down; the frozen session lives on the App.
func (s *StatePauseImpl) Exit(sm *StateMachine) {}

func (s *StatePauseImpl) Update(sm *StateMachine) error {
	cx := float64(GameSettings.ScreenWidth) / 2
	continueY, optionsY, newY, exitY := 150.0, 200.0, 250.0, 300.0

	if risingClick(&s.mouseWasPressed) {
		mx, my := cursorF()
		switch {
		case inButton(mx, my, cx, continueY):
			sm.SwitchTo(StateIDGame)
		case inButton(mx, my, cx, optionsY):
			openOptions(sm, StateIDPause)
		case inButton(mx, my, cx, newY):
			sm.App().SetGame(NewGame())
			sm.SwitchTo(StateIDGame)
		case inButton(mx, my, cx, exitY):
			sm.App().SetGame(nil)
			sm.SwitchTo(StateIDMainMenu)
		}
	}
	return nil
}

func (s *StatePauseImpl) Draw(sm *StateMachine, screen *ebiten.Image) {
	// Paint the frozen world beneath the overlay so the player keeps their bearings.
	if g := sm.App().Gameplay(); g != nil {
		g.Draw(screen)
	}
	dimScreen(screen, 140)

	cx := float64(GameSettings.ScreenWidth) / 2
	drawCenteredTitle(screen, "PAUSE", cx, 100, color.White)

	mx, my := cursorF()
	continueY, optionsY, newY, exitY := 150.0, 200.0, 250.0, 300.0
	left := cx - menuButtonWidth/2
	drawMenuButton(screen, left, continueY, menuButtonWidth, menuButtonHeight, "CONTINUE", inButton(mx, my, cx, continueY))
	drawMenuButton(screen, left, optionsY, menuButtonWidth, menuButtonHeight, "OPTIONS", inButton(mx, my, cx, optionsY))
	drawMenuButton(screen, left, newY, menuButtonWidth, menuButtonHeight, "NEW GAME", inButton(mx, my, cx, newY))
	drawMenuButton(screen, left, exitY, menuButtonWidth, menuButtonHeight, "EXIT", inButton(mx, my, cx, exitY))
}

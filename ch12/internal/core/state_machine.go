package core

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

// ErrQuit is returned from Update when the user requests to exit.
var ErrQuit = errors.New("quit")

// StateID identifies a state.
type StateID string

const (
	StateIDMainMenu StateID = "mainmenu"
	StateIDGame     StateID = "game"
	StateIDPause    StateID = "pause"
	StateIDOptions  StateID = "options" // new in ch12
)

// GameSession is the gameplay loop exposed to states (package game implements this).
type GameSession interface {
	Update() error
	Draw(*ebiten.Image)
	IsPausePressed() bool
}

// AppContext is what the state machine holds (implemented by game.App).
type AppContext interface {
	SetGame(g any)
	Gameplay() GameSession
}

// State is implemented by each screen (menu, game, pause, options).
// Enter runs once when the state becomes current; Exit runs once as the machine
// leaves it, so a state can own a resource for its whole lifetime — the game state
// starts the music in Enter and stops it in Exit.
type State interface {
	Enter(sm *StateMachine)
	Exit(sm *StateMachine)
	Update(sm *StateMachine) error
	Draw(sm *StateMachine, screen *ebiten.Image)
}

// StateMachine manages states.
type StateMachine struct {
	app     AppContext
	current State
	states  map[StateID]State
}

// NewStateMachine creates a state machine bound to app.
func NewStateMachine(app AppContext) *StateMachine {
	return &StateMachine{
		app:    app,
		states: make(map[StateID]State),
	}
}

// Register adds a state.
func (sm *StateMachine) Register(id StateID, state State) {
	sm.states[id] = state
}

// State returns the registered state for id, or nil if none is registered.
// Callers use it to configure a state (such as the options screen's return target)
// before switching to it.
func (sm *StateMachine) State(id StateID) State {
	return sm.states[id]
}

// SwitchTo changes the current state, calling Exit on the state being left and
// Enter on the one being entered, so each state can set up and tear down cleanly.
func (sm *StateMachine) SwitchTo(id StateID) {
	if s, ok := sm.states[id]; ok {
		if sm.current != nil {
			sm.current.Exit(sm)
		}
		sm.current = s
		sm.current.Enter(sm)
	}
}

// Update delegates to the current state.
func (sm *StateMachine) Update() error {
	if sm.current == nil {
		return nil
	}
	return sm.current.Update(sm)
}

// Draw delegates to the current state.
func (sm *StateMachine) Draw(screen *ebiten.Image) {
	if sm.current == nil {
		return
	}
	sm.current.Draw(sm, screen)
}

// App returns the app host for states.
func (sm *StateMachine) App() AppContext {
	return sm.app
}

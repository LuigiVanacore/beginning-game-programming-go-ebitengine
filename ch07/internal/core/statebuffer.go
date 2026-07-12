package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputState holds pressed state for one frame and the previous frame.
type InputState struct {
	IsPressed  bool
	WasPressed bool
}

// StateBuffer tracks keyboard and mouse button states for Action-based input.
type StateBuffer struct {
	keyStates   map[ebiten.Key]InputState
	mouseStates map[ebiten.MouseButton]InputState
}

// NewStateBuffer creates a new StateBuffer.
func NewStateBuffer() *StateBuffer {
	return &StateBuffer{
		keyStates:   make(map[ebiten.Key]InputState),
		mouseStates: make(map[ebiten.MouseButton]InputState),
	}
}

// Update polls ebiten and updates all button states.
func (s *StateBuffer) Update() {
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		old := s.keyStates[key]
		s.keyStates[key] = InputState{
			IsPressed:  ebiten.IsKeyPressed(key),
			WasPressed: old.IsPressed,
		}
	}
	for btn := ebiten.MouseButtonLeft; btn <= ebiten.MouseButtonRight; btn++ {
		old := s.mouseStates[btn]
		s.mouseStates[btn] = InputState{
			IsPressed:  ebiten.IsMouseButtonPressed(btn),
			WasPressed: old.IsPressed,
		}
	}
}

// IsActionActive returns whether the action is currently active according to its mode.
func (s *StateBuffer) IsActionActive(action Action) bool {
	switch action.GetActionType() {
	case KeyAction:
		state := s.keyStates[action.GetKey()]
		return s.checkState(state, action.GetMode())
	case MouseButtonAction:
		state := s.mouseStates[action.GetMouseButton()]
		return s.checkState(state, action.GetMode())
	}
	return false
}

func (s *StateBuffer) checkState(state InputState, mode ActionMode) bool {
	if mode.Has(ActionHold) && state.IsPressed {
		return true
	}
	if mode.Has(ActionPressOnce) && state.IsPressed && !state.WasPressed {
		return true
	}
	if mode.Has(ActionReleaseOnce) && !state.IsPressed && state.WasPressed {
		return true
	}
	return false
}

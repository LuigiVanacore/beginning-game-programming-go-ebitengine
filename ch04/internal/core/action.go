package core

import "github.com/hajimehoshi/ebiten/v2"

// ActionMode is a bit mask for when an action triggers (hold, press once, release once).
type ActionMode uint

const (
	ActionHold      ActionMode = 1 << iota
	ActionPressOnce
	ActionReleaseOnce
)

// Has returns true if the given flag is set in the mode.
func (m ActionMode) Has(flag ActionMode) bool {
	return m&flag != 0
}

// ActionType identifies the kind of input (key, mouse, gamepad).
type ActionType int

const (
	KeyAction ActionType = iota
	MouseButtonAction
	GamepadButtonAction
)

// Action represents a single input trigger bound to a key, mouse button, or gamepad button.
// The mode controls when IsActionActive returns true (hold, press once, or release once).
type Action struct {
	actionType    ActionType
	mode          ActionMode
	key           ebiten.Key
	mouseButton   ebiten.MouseButton
	gamepadID     ebiten.GamepadID
	gamepadButton ebiten.GamepadButton
}

// NewKeyAction creates an Action bound to a keyboard key.
func NewKeyAction(key ebiten.Key, mode ActionMode) Action {
	return Action{
		actionType: KeyAction,
		mode:       mode,
		key:        key,
	}
}

// NewMouseAction creates an Action bound to a mouse button.
func NewMouseAction(btn ebiten.MouseButton, mode ActionMode) Action {
	return Action{
		actionType:  MouseButtonAction,
		mode:        mode,
		mouseButton: btn,
	}
}

// NewGamepadAction creates an Action bound to a gamepad button.
func NewGamepadAction(id ebiten.GamepadID, btn ebiten.GamepadButton, mode ActionMode) Action {
	return Action{
		actionType:    GamepadButtonAction,
		mode:          mode,
		gamepadID:     id,
		gamepadButton: btn,
	}
}

// GetActionType returns the action's input type.
func (a Action) GetActionType() ActionType {
	return a.actionType
}

// GetMode returns the action's trigger mode.
func (a Action) GetMode() ActionMode {
	return a.mode
}

// GetKey returns the key (valid only for KeyAction).
func (a Action) GetKey() ebiten.Key {
	return a.key
}

// GetMouseButton returns the mouse button (valid only for MouseButtonAction).
func (a Action) GetMouseButton() ebiten.MouseButton {
	return a.mouseButton
}

// GetGamepadButton returns (id, button) for GamepadButtonAction.
func (a Action) GetGamepadButton() (ebiten.GamepadID, ebiten.GamepadButton) {
	return a.gamepadID, a.gamepadButton
}

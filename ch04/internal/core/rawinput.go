package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputButtonType identifies keyboard, mouse, or gamepad input.
type InputButtonType int

const (
	ButtonKeyboard InputButtonType = iota
	ButtonMouse
	ButtonGamepad
)

// RawInputButton represents a single physical input (key, mouse button, or gamepad button).
// It abstracts over the input hardware so the game can bind actions to different devices.
type RawInputButton struct {
	key          ebiten.Key
	mouseButton  ebiten.MouseButton
	gamepadID    ebiten.GamepadID
	gamepadBtn   ebiten.GamepadButton
	buttonType   InputButtonType
}

// NewKeyRawInputButton creates a RawInputButton for a keyboard key.
func NewKeyRawInputButton(key ebiten.Key) RawInputButton {
	return RawInputButton{
		key:        key,
		buttonType: ButtonKeyboard,
	}
}

// NewMouseRawInputButton creates a RawInputButton for a mouse button.
func NewMouseRawInputButton(btn ebiten.MouseButton) RawInputButton {
	return RawInputButton{
		mouseButton: btn,
		buttonType:  ButtonMouse,
	}
}

// NewGamepadRawInputButton creates a RawInputButton for a gamepad button.
func NewGamepadRawInputButton(id ebiten.GamepadID, btn ebiten.GamepadButton) RawInputButton {
	return RawInputButton{
		gamepadID:  id,
		gamepadBtn: btn,
		buttonType: ButtonGamepad,
	}
}

// IsPressed returns true if the button is currently held down.
func (r RawInputButton) IsPressed() bool {
	switch r.buttonType {
	case ButtonKeyboard:
		return ebiten.IsKeyPressed(r.key)
	case ButtonMouse:
		return ebiten.IsMouseButtonPressed(r.mouseButton)
	case ButtonGamepad:
		return ebiten.IsGamepadButtonPressed(r.gamepadID, r.gamepadBtn)
	}
	return false
}

// IsJustPressed returns true only on the first frame the button is pressed.
func (r RawInputButton) IsJustPressed() bool {
	switch r.buttonType {
	case ButtonKeyboard:
		return inpututil.IsKeyJustPressed(r.key)
	case ButtonMouse:
		return inpututil.IsMouseButtonJustPressed(r.mouseButton)
	case ButtonGamepad:
		return inpututil.IsGamepadButtonJustPressed(r.gamepadID, r.gamepadBtn)
	}
	return false
}

// IsJustReleased returns true only on the frame the button is released.
func (r RawInputButton) IsJustReleased() bool {
	switch r.buttonType {
	case ButtonKeyboard:
		return inpututil.IsKeyJustReleased(r.key)
	case ButtonMouse:
		return inpututil.IsMouseButtonJustReleased(r.mouseButton)
	case ButtonGamepad:
		return inpututil.IsGamepadButtonJustReleased(r.gamepadID, r.gamepadBtn)
	}
	return false
}

// PressDuration returns how many frames the button has been held (0 if not pressed).
func (r RawInputButton) PressDuration() int {
	switch r.buttonType {
	case ButtonKeyboard:
		return inpututil.KeyPressDuration(r.key)
	case ButtonMouse:
		return inpututil.MouseButtonPressDuration(r.mouseButton)
	case ButtonGamepad:
		return inpututil.GamepadButtonPressDuration(r.gamepadID, r.gamepadBtn)
	}
	return 0
}

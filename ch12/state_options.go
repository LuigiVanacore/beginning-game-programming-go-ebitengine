package game

import (
	. "book/code/ch12/internal/core"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const volumeStep = 0.1

// StateOptionsImpl edits GameOptions: master volume, an SFX toggle, and fullscreen.
// It is reachable from both the main menu and the pause menu; returnTo records which
// screen to go back to when the player presses BACK.
type StateOptionsImpl struct {
	mouseWasPressed bool
	returnTo        StateID
}

func (s *StateOptionsImpl) Enter(sm *StateMachine) {
	if s.returnTo == "" {
		s.returnTo = StateIDMainMenu
	}
	s.mouseWasPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// Exit has nothing to tear down; the settings live on GameOptions.
func (s *StateOptionsImpl) Exit(sm *StateMachine) {}

// optionsLayout returns the clickable rectangles, so Update and Draw stay in sync.
// Each rect is [x, y, w, h].
func optionsLayout() (volMinus, volPlus, sfx, full, back [4]float64) {
	cx := float64(GameSettings.ScreenWidth) / 2
	left := cx - menuButtonWidth/2
	volMinus = [4]float64{cx - 130, 148, 36, 32}
	volPlus = [4]float64{cx + 94, 148, 36, 32}
	sfx = [4]float64{left, 208, menuButtonWidth, menuButtonHeight}
	full = [4]float64{left, 258, menuButtonWidth, menuButtonHeight}
	back = [4]float64{left, 320, menuButtonWidth, menuButtonHeight}
	return
}

func (s *StateOptionsImpl) Update(sm *StateMachine) error {
	if !risingClick(&s.mouseWasPressed) {
		return nil
	}
	mx, my := cursorF()
	volMinus, volPlus, sfx, full, back := optionsLayout()
	switch {
	case inRect(mx, my, volMinus):
		GameOptions.MasterVolume = clamp01(GameOptions.MasterVolume - volumeStep)
		Audio().SetMasterVolume(GameOptions.MasterVolume) // push the new level to the audio system
	case inRect(mx, my, volPlus):
		GameOptions.MasterVolume = clamp01(GameOptions.MasterVolume + volumeStep)
		Audio().SetMasterVolume(GameOptions.MasterVolume)
	case inRect(mx, my, sfx):
		GameOptions.SFXEnabled = !GameOptions.SFXEnabled
	case inRect(mx, my, full):
		GameOptions.Fullscreen = !GameOptions.Fullscreen
		ebiten.SetFullscreen(GameOptions.Fullscreen)
	case inRect(mx, my, back):
		sm.SwitchTo(s.returnTo)
	}
	return nil
}

func (s *StateOptionsImpl) Draw(sm *StateMachine, screen *ebiten.Image) {
	dimScreen(screen, 180)
	cx := float64(GameSettings.ScreenWidth) / 2
	drawCenteredTitle(screen, "OPTIONS", cx, 90, color.RGBA{255, 220, 100, 255})

	mx, my := cursorF()
	volMinus, volPlus, sfx, full, back := optionsLayout()

	// Volume: a -/+ pair with a percentage readout centered between them.
	drawSmallButton(screen, volMinus, "-", inRect(mx, my, volMinus))
	drawSmallButton(screen, volPlus, "+", inRect(mx, my, volPlus))
	drawCenteredTitle(screen, fmt.Sprintf("VOLUME  %d%%", int(GameOptions.MasterVolume*100+0.5)), cx, 170, color.White)

	drawRectButton(screen, sfx, "SFX: "+onOff(GameOptions.SFXEnabled), inRect(mx, my, sfx))
	drawRectButton(screen, full, "FULLSCREEN: "+onOff(GameOptions.Fullscreen), inRect(mx, my, full))
	drawRectButton(screen, back, "BACK", inRect(mx, my, back))
}

// --- helpers local to the options screen ---

// inRect adapts pointInRect to a [x,y,w,h] rectangle.
func inRect(px, py float64, r [4]float64) bool {
	return pointInRect(px, py, r[0], r[1], r[2], r[3])
}

// drawRectButton draws a menu button described by a [x,y,w,h] rectangle.
func drawRectButton(screen *ebiten.Image, r [4]float64, label string, hovered bool) {
	drawMenuButton(screen, r[0], r[1], r[2], r[3], label, hovered)
}

// drawSmallButton draws a compact square control (the volume steppers).
func drawSmallButton(screen *ebiten.Image, r [4]float64, label string, hovered bool) {
	fill := color.RGBA{60, 100, 160, 255}
	border := color.RGBA{100, 150, 220, 255}
	if hovered {
		fill = color.RGBA{85, 135, 210, 255}
		border = color.RGBA{140, 190, 255, 255}
	}
	vector.DrawFilledRect(screen, float32(r[0]), float32(r[1]), float32(r[2]), float32(r[3]), fill, true)
	vector.StrokeRect(screen, float32(r[0]), float32(r[1]), float32(r[2]), float32(r[3]), 2, border, true)

	face := text.NewGoXFace(basicfont.Face7x13)
	lw, _ := text.Measure(label, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(r[0]+r[2]/2-lw/2, r[1]+(r[3]-13)/2)
	opts.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, label, face, opts)
}

// clamp01 keeps a value inside [0, 1].
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// onOff renders a boolean toggle as ON/OFF.
func onOff(b bool) string {
	if b {
		return "ON"
	}
	return "OFF"
}

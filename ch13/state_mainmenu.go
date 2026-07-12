package game

import (
	. "book/code/ch13/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Menu button geometry shared by the menu-like states.
const (
	menuButtonWidth  = 160.0
	menuButtonHeight = 32.0
)

// StateMainMenuImpl is the first screen. It offers a new run, the options menu, and exit.
type StateMainMenuImpl struct {
	mouseWasPressed bool
}

func (s *StateMainMenuImpl) Enter(sm *StateMachine) {
	// Seed the edge tracker so a click that started on a previous screen does not
	// leak into this one on the first frame.
	s.mouseWasPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// Exit has nothing to tear down; the menu owns no resource.
func (s *StateMainMenuImpl) Exit(sm *StateMachine) {}

func (s *StateMainMenuImpl) Update(sm *StateMachine) error {
	cx := float64(GameSettings.ScreenWidth) / 2
	newY, optionsY, exitY := 180.0, 230.0, 280.0

	clicked := risingClick(&s.mouseWasPressed)
	mx, my := cursorF()
	if clicked && inButton(mx, my, cx, newY) {
		Audio().Play(SoundClick)
		sm.App().SetGame(NewGame())
		sm.SwitchTo(StateIDGame)
		return nil
	}
	if clicked && inButton(mx, my, cx, optionsY) {
		Audio().Play(SoundClick)
		openOptions(sm, StateIDMainMenu)
		return nil
	}
	if clicked && inButton(mx, my, cx, exitY) {
		Audio().Play(SoundClick)
		return ErrQuit
	}
	return nil
}

func (s *StateMainMenuImpl) Draw(sm *StateMachine, screen *ebiten.Image) {
	cx := float64(GameSettings.ScreenWidth) / 2

	drawCenteredTitle(screen, "Gopher Survivor", cx, 100, color.RGBA{255, 220, 100, 255})

	mx, my := cursorF()
	newY, optionsY, exitY := 180.0, 230.0, 280.0
	drawMenuButton(screen, cx-menuButtonWidth/2, newY, menuButtonWidth, menuButtonHeight, "NEW GAME", inButton(mx, my, cx, newY))
	drawMenuButton(screen, cx-menuButtonWidth/2, optionsY, menuButtonWidth, menuButtonHeight, "OPTIONS", inButton(mx, my, cx, optionsY))
	drawMenuButton(screen, cx-menuButtonWidth/2, exitY, menuButtonWidth, menuButtonHeight, "EXIT", inButton(mx, my, cx, exitY))
}

// --- shared menu helpers (used by every menu-like state) ---

// risingClick reports the frame the left mouse button transitions to pressed and
// updates the caller's edge tracker. It gives menu states a rising-edge click
// without a dedicated input action, mirroring the HUD's upgrade-panel handling.
func risingClick(was *bool) bool {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	just := pressed && !*was
	*was = pressed
	return just
}

// cursorF returns the mouse position as float64.
func cursorF() (float64, float64) {
	mx, my := ebiten.CursorPosition()
	return float64(mx), float64(my)
}

// inButton reports whether (mx,my) is inside a centered button at column cx, top y.
func inButton(mx, my, cx, y float64) bool {
	left := cx - menuButtonWidth/2
	return mx >= left && mx <= left+menuButtonWidth && my >= y && my <= y+menuButtonHeight
}

// pointInRect reports whether (px,py) lies inside the rectangle (x,y,w,h). Used for
// the small +/- and toggle controls on the options screen.
func pointInRect(px, py, x, y, w, h float64) bool {
	return px >= x && px <= x+w && py >= y && py <= y+h
}

// dimScreen paints a translucent black veil over the whole screen (0..255 alpha),
// used to darken the frozen world behind the pause and options overlays.
func dimScreen(screen *ebiten.Image, alpha uint8) {
	vector.DrawFilledRect(screen, 0, 0, float32(GameSettings.ScreenWidth), float32(GameSettings.ScreenHeight),
		color.RGBA{0, 0, 0, alpha}, true)
}

// openOptions configures the options screen to return to `from` and switches to it.
func openOptions(sm *StateMachine, from StateID) {
	if opt, ok := sm.State(StateIDOptions).(*StateOptionsImpl); ok {
		opt.returnTo = from
	}
	sm.SwitchTo(StateIDOptions)
}

// drawCenteredTitle draws a horizontally centered heading.
func drawCenteredTitle(screen *ebiten.Image, s string, cx, y float64, clr color.Color) {
	face := text.NewGoXFace(basicfont.Face7x13)
	w, _ := text.Measure(s, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(cx-w/2, y)
	opts.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, s, face, opts)
}

// drawMenuButton draws a filled, outlined button with a centered label.
func drawMenuButton(screen *ebiten.Image, x, y, w, h float64, label string, hovered bool) {
	fill := color.RGBA{60, 100, 160, 255}
	border := color.RGBA{100, 150, 220, 255}
	if hovered {
		fill = color.RGBA{85, 135, 210, 255}
		border = color.RGBA{140, 190, 255, 255}
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), fill, true)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, border, true)

	face := text.NewGoXFace(basicfont.Face7x13)
	lw, _ := text.Measure(label, face, 0)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x+w/2-lw/2, y+(h-13)/2)
	opts.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, label, face, opts)
}

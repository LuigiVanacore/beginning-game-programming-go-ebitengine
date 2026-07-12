package ui

import (
	. "book/code/ch12/internal/core"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// enemyGrewWarnFrames is how long the "Enemies Grow Stronger" warning stays on
// screen after a difficulty tier increase; warnBlinkHalfPeriod sets its blink rate.
const (
	enemyGrewWarnFrames = 180 // ~3 seconds at 60 TPS
	warnBlinkHalfPeriod = 12  // frames shown, then frames hidden
)

// HUD owns all in-game UI: health bar, XP bar, level-up popup, and the upgrade-choice overlay.
type HUD struct {
	rm           *ResourceManager
	healthBar    *UIBar
	xpBar        *UIBar
	levelUpPopup *UIPopup

	choosingUpgrade bool
	upgradePanels   []*UpgradePanel
	upgradeApplies  []func()
	mouseWasPressed bool

	warnUntilFrame int // elapsed frame at which the "Enemies Grow Stronger" warning disappears
}

// NewHUD creates the HUD with bars and popup initialised from GameSettings.
func NewHUD(rm *ResourceManager) *HUD {
	healthBar := NewUIBar(10, 10, 200, 14)
	xpBar := NewUIBar(10, 30, 200, 10)
	xpBar.FgColor = color.RGBA{R: 255, G: 200, B: 64, A: 255}
	return &HUD{
		rm:           rm,
		healthBar:    healthBar,
		xpBar:        xpBar,
		levelUpPopup: NewUIPopup("LEVEL UP", PopupDuration, PopupRiseSpeed),
	}
}

// IsChoosingUpgrade returns true while the player is selecting an upgrade.
func (h *HUD) IsChoosingUpgrade() bool { return h.choosingUpgrade }

// TriggerLevelUp starts the upgrade-choice screen with parallel choices and apply callbacks.
func (h *HUD) TriggerLevelUp(choices []UpgradeChoice, applies []func()) {
	h.choosingUpgrade = true
	h.upgradeApplies = applies
	h.upgradePanels = h.buildPanels(choices)
}

// HandleInput processes upgrade-panel mouse clicks and runs the matching apply callback.
func (h *HUD) HandleInput() {
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mousePressed && !h.mouseWasPressed {
		mx, my := ebiten.CursorPosition()
		for i, pan := range h.upgradePanels {
			if pan == nil || !pan.ButtonContainsPoint(float64(mx), float64(my)) {
				continue
			}
			if i < len(h.upgradeApplies) {
				h.upgradeApplies[i]()
			}
			h.choosingUpgrade = false
			h.upgradePanels = nil
			h.upgradeApplies = nil
			break
		}
	}
	h.mouseWasPressed = mousePressed
}

// UpdatePopup advances the level-up popup animation (call when not choosing upgrade).
func (h *HUD) UpdatePopup() { h.levelUpPopup.Update() }

// TriggerEnemyGrew flashes the "Enemies Grow Stronger" warning starting from the given
// elapsed frame. Call it when the difficulty tier advances.
func (h *HUD) TriggerEnemyGrew(elapsedFrames int) {
	h.warnUntilFrame = elapsedFrames + enemyGrewWarnFrames
}

// Draw renders health bar, XP bar, survival timer, warning, and level-up popup.
func (h *HUD) Draw(screen *ebiten.Image, hpRatio, xpPct float64, elapsedFrames int) {
	h.healthBar.SetProgress(hpRatio)
	h.healthBar.Draw(screen)

	h.xpBar.SetProgress(clamp01(xpPct))
	h.xpBar.Draw(screen)

	h.drawTimer(screen, elapsedFrames)
	h.drawGrowWarning(screen, elapsedFrames)
	h.levelUpPopup.Draw(screen)
}

// DrawUpgradeOverlay renders the full-screen upgrade-choice overlay.
func (h *HUD) DrawUpgradeOverlay(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(GameSettings.ScreenWidth), float32(GameSettings.ScreenHeight),
		color.RGBA{R: 15, G: 15, B: 30, A: 255}, true)
	mx, my := ebiten.CursorPosition()
	for _, p := range h.upgradePanels {
		if p == nil {
			continue
		}
		p.ButtonHovered = p.ButtonContainsPoint(float64(mx), float64(my))
		p.Draw(screen)
	}
}

// drawTimer renders the survival clock. It counts up from 0:00, so a longer run
// shows a larger value: the goal is to beat your own best time.
func (h *HUD) drawTimer(screen *ebiten.Image, elapsedFrames int) {
	if elapsedFrames < 0 {
		elapsedFrames = 0
	}
	tps := ebiten.TPS()
	if tps < 1 {
		tps = 1
	}
	secsTotal := elapsedFrames / tps
	const secondsPerMinute = 60
	mins := secsTotal / secondsPerMinute
	secs := secsTotal % secondsPerMinute
	DrawLabel(screen, UILabel{
		Text:               fmt.Sprintf("%d:%02d", mins, secs),
		X:                  float64(GameSettings.ScreenWidth) / 2,
		Y:                  8,
		Color:              color.RGBA{R: 255, G: 255, B: 255, A: 255},
		CenterHorizontally: true,
	}, 80, 16)
}

// drawGrowWarning shows a red, blinking "Enemies Grow Stronger" label under the timer
// while the warning window is open. The blink is driven by the elapsed frame count.
func (h *HUD) drawGrowWarning(screen *ebiten.Image, elapsedFrames int) {
	if elapsedFrames >= h.warnUntilFrame {
		return
	}
	// Blink: show for warnBlinkHalfPeriod frames, then hide for the same span.
	if (elapsedFrames/warnBlinkHalfPeriod)%2 == 1 {
		return
	}
	DrawLabel(screen, UILabel{
		Text:               "Enemies Grow Stronger",
		X:                  float64(GameSettings.ScreenWidth) / 2,
		Y:                  26,
		Color:              color.RGBA{R: 255, G: 40, B: 40, A: 255},
		CenterHorizontally: true,
	}, 200, 16)
}

func (h *HUD) buildPanels(options []UpgradeChoice) []*UpgradePanel {
	rm := h.rm
	panelW := upgradePanelWidth
	panelH := upgradePanelHeight
	gap := upgradePanelGap
	totalW := panelW*float64(len(options)) + gap*float64(len(options)-1)
	startX := (float64(GameSettings.ScreenWidth) - totalW) / 2
	y := (float64(GameSettings.ScreenHeight) - panelH) / 2

	panels := make([]*UpgradePanel, 0, len(options))
	for i, opt := range options {
		x := startX + float64(i)*(panelW+gap)
		pan := NewUpgradePanel(x, y, panelW, panelH, opt, rm, nil)
		panels = append(panels, pan)
	}
	return panels
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

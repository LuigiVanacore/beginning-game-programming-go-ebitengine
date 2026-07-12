package ui

import "github.com/hajimehoshi/ebiten/v2"

// HUD composes chapter 8 UI widgets; each widget type lives in its own file.
type HUD struct {
	healthBar    *HealthBar
	xpBar        *XPBar
	levelUpPopup *UIPopup
	gameOver     GameOverLabel
}

// NewHUD creates the HUD with bars and popup initialised from GameSettings.
func NewHUD() *HUD {
	return &HUD{
		healthBar:    NewHealthBar(),
		xpBar:        NewXPBar(),
		levelUpPopup: NewUIPopup("LEVEL UP", PopupDuration, PopupRisePxPerFrame),
	}
}

// DrawGameplay draws health, XP, and level-up popup (no game-over overlay).
func (h *HUD) DrawGameplay(screen *ebiten.Image, hpRatio float64, xp int, level int) {
	if h == nil {
		return
	}
	h.healthBar.Draw(screen, hpRatio)
	h.xpBar.Draw(screen, xp, level)
	h.levelUpPopup.Draw(screen)
}

// DrawGameOver draws the centered game-over label on top of the frame.
func (h *HUD) DrawGameOver(screen *ebiten.Image) {
	if h == nil {
		return
	}
	h.gameOver.Draw(screen)
}

// ShowLevelUp delegates to the level-up popup.
func (h *HUD) ShowLevelUp(x, y float64) {
	if h == nil || h.levelUpPopup == nil {
		return
	}
	h.levelUpPopup.Show(x, y)
}

// UpdateLevelPopup advances the level-up popup animation.
func (h *HUD) UpdateLevelPopup() {
	if h == nil || h.levelUpPopup == nil {
		return
	}
	h.levelUpPopup.Update()
}

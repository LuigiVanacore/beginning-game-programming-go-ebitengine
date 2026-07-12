package game

import (
	en "book/code/ch13/enemy"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw renders the game. When choosing an upgrade the world is hidden and the overlay is shown.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.hud.IsChoosingUpgrade() {
		g.hud.DrawUpgradeOverlay(screen)
		return
	}
	g.engine.Draw(screen)

	// VFX pass: drawn in world space over the scene but under the HUD (ch13). The
	// camera position converts world coordinates to screen coordinates.
	cam := g.engine.World().Camera()
	camX, camY := cam.GetPosition().X(), cam.GetPosition().Y()
	g.drawHitFlashes(screen, camX, camY)
	if g.particles != nil {
		g.particles.Draw(screen, camX, camY)
	}
	if g.floatingText != nil {
		g.floatingText.Draw(screen, camX, camY)
	}

	need := xpNeededForLevel(g.player.Level)
	xpPct := 0.0
	if need > 0 {
		xpPct = float64(g.player.XP) / float64(need)
	}
	g.hud.Draw(screen, g.player.HP/g.player.MaxHP, xpPct, g.elapsedFrames)
	if g.gameOver {
		// Survival mode ends only when the player dies; the frozen HUD timer above
		// keeps showing how long the player lasted. The overlay draws the New Game
		// button; the game state handles the click (see state_game.go).
		g.gameOverOverlay.Draw(screen)
	}
}

// drawHitFlashes paints a fading white disc over each enemy currently flashing from
// a recent hit, reading its HitFlashFrames countdown (ch13).
func (g *Game) drawHitFlashes(screen *ebiten.Image, camX, camY float64) {
	for _, e := range g.enemyManager.Enemies() {
		if e == nil || e.HitFlashFrames <= 0 {
			continue
		}
		p := e.GetWorldPosition()
		alpha := uint8(160 * float64(e.HitFlashFrames) / float64(en.HitFlashDuration))
		vector.DrawFilledCircle(screen, float32(p.X()-camX), float32(p.Y()-camY), 14,
			color.RGBA{255, 255, 255, alpha}, true)
	}
}

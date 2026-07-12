package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Draw renders the game. When choosing an upgrade the world is hidden and the overlay is shown.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.hud.IsChoosingUpgrade() {
		g.hud.DrawUpgradeOverlay(screen)
		return
	}
	g.engine.Draw(screen)
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

package enemy

import "math"

const (
	tierDurationSeconds = 60.0 // a new difficulty tier every minute
	statGrowthPerTier   = 1.5   // HP and contact damage multiply by this each tier
	xpGrowthPerTier     = 1.15  // experience grows, but slower than combat strength
)

// Difficulty tracks how strong the monsters are. It derives a tier from the
// elapsed survival time, then exposes two multipliers: one for combat stats
// (HP and damage) that grows quickly, and one for experience that grows slowly.
type Difficulty struct {
	tier int
}

// Update recomputes the tier from the elapsed survival time. It returns true on
// the single frame the tier increases, so the HUD can flash "Enemies Grow Stronger".
func (d *Difficulty) Update(elapsedSeconds float64) bool {
	newTier := int(elapsedSeconds / tierDurationSeconds)
	if newTier > d.tier {
		d.tier = newTier
		return true
	}
	return false
}

// Tier is the current difficulty step (0 for the first three minutes).
func (d *Difficulty) Tier() int { return d.tier }

// StatScale multiplies enemy HP and contact damage. It grows exponentially, so
// monsters become dangerous quickly.
func (d *Difficulty) StatScale() float64 { return math.Pow(statGrowthPerTier, float64(d.tier)) }

// XPScale multiplies the experience reward. It uses a smaller base than StatScale,
// so rewards rise more gently than the threat does.
func (d *Difficulty) XPScale() float64 { return math.Pow(xpGrowthPerTier, float64(d.tier)) }

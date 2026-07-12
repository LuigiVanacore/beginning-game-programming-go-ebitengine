package enemy

import . "book/code/ch11/internal/core"

// EnemyType identifies a monster kind. The values are ordered from the weakest
// (Ghost) to the strongest (Cyclops), so a larger value always means a tougher
// enemy. Cyclops is the solitary mini-boss.
type EnemyType int

const (
	Ghost EnemyType = iota
	Spider
	Bat
	DarkWizard
	Cyclops
)

// enemyStats holds the tier-0 attributes for one monster kind. The spawner scales
// baseHP and baseDamage by the difficulty stat multiplier and baseXP by the
// (slower) XP multiplier before building an enemy.
type enemyStats struct {
	textureKey string
	baseHP     float64 // hit points before difficulty scaling
	baseDamage float64 // HP drained from the player per contact frame, before scaling
	baseXP     int     // experience granted on death, before scaling
	scale      float64 // sprite scale factor
	radius     float64 // collision-circle radius
}

// enemyStatTable is the single source of truth for every monster kind.
// Reading down the table, each kind is tougher and grants more experience than
// the one above it; the Cyclops is the outsized mini-boss.
var enemyStatTable = map[EnemyType]enemyStats{
	Ghost:      {GhostTexture, 3, 0.10, 5, 2, enemyRadius},
	Spider:     {SpiderTexture, 6, 0.15, 8, 2, enemyRadius},
	Bat:        {BatTexture, 10, 0.20, 12, 2, enemyRadius},
	DarkWizard: {DarkWizardTexture, 18, 0.30, 20, 2, enemyRadius},
	Cyclops:    {CyclopsTexture, 120, 0.60, 150, 3, cyclopsRadius},
}

// statsFor returns the base stats for a kind, defaulting to Ghost for safety.
func statsFor(kind EnemyType) enemyStats {
	if st, ok := enemyStatTable[kind]; ok {
		return st
	}
	return enemyStatTable[Ghost]
}

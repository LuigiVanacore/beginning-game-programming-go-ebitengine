package game

import (
	. "book/code/ch11/internal/core"
	. "book/code/ch11/ui"

	"math/rand"
)

// xpNeededForLevel returns XP required to reach the next level (exponential curve).
func xpNeededForLevel(level int) int {
	if level < 1 {
		return GameSettings.XPBaseLevel
	}
	f := float64(GameSettings.XPBaseLevel)
	for i := 1; i < level; i++ {
		f *= GameSettings.XPGrowthFactor
	}
	return int(f)
}

func upgradeChoicesFromOptions(opts []UpgradeOption, g *Game) ([]UpgradeChoice, []func()) {
	choices := make([]UpgradeChoice, len(opts))
	applies := make([]func(), len(opts))
	for i, o := range opts {
		o := o
		choices[i] = UpgradeChoice{WeaponName: o.WeaponName, UpgradeDesc: o.UpgradeDesc, IconKey: o.IconKey}
		applies[i] = func() { o.Apply(g) }
	}
	return choices, applies
}

// randFloat64 returns a random float in [0,1).
func randFloat64() float64 { return rand.Float64() }

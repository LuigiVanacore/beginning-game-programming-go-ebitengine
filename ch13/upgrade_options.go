package game

import (
	. "book/code/ch13/internal/core"

	"math/rand"
)

// UpgradeOption represents a single upgrade the player can choose.
// IconKey empty = no icon (stat-only: Health, Damage, XP).
type UpgradeOption struct {
	WeaponName  string
	UpgradeDesc string
	IconKey     string
	Apply       func(g *Game)
}

// upgradeStrategy is the Strategy pattern in Go: each function encapsulates
// eligibility + the Apply closure for one upgrade card.
type upgradeStrategy func(*Game) (UpgradeOption, bool)

// bonusIconKeys are IconKeys for one-time bonus upgrades (Armor, Boots, Gem, Skull, Ring).
var bonusIconKeys = map[string]bool{"armor": true, "boots": true, "gem": true, "skull": true, "ring": true}

// upgradeStrategies is the catalog (Registry) of all standard upgrade strategies, in evaluation order.
var upgradeStrategies = []upgradeStrategy{
	strategyKnifeUnlock,
	strategyKnifeSpeed,
	strategyKnifeCooldown,
	strategyFlyingAxeUnlock,
	strategyFlyingAxeSpeed,
	strategyFlyingAxeCooldown,
	strategyFlyingAxeRotation,
	strategySacredBookUnlock,
	strategySacredBookSpeed,
	strategyHolyShieldUnlock,
	strategyHolyShieldRadius,
	strategyBonusArmor,
	strategyBonusBoots,
	strategyBonusGem,
	strategyBonusSkull,
	strategyBonusRing,
}

// PickUpgradesForLevelUp returns upgrade options for the current level-up.
// It ensures at least one bonus item appears when any is available.
func PickUpgradesForLevelUp(g *Game) []UpgradeOption {
	pool := buildUpgradePool(g)
	if len(pool) == 0 {
		return nil
	}
	if len(pool) <= 2 {
		idx := rand.Intn(len(pool))
		return []UpgradeOption{pool[idx]}
	}
	n := 3
	if len(pool) < n {
		n = len(pool)
	}
	var bonuses, weapons []UpgradeOption
	for _, opt := range pool {
		if bonusIconKeys[opt.IconKey] {
			bonuses = append(bonuses, opt)
		} else {
			weapons = append(weapons, opt)
		}
	}
	var result []UpgradeOption
	if len(bonuses) > 0 {
		idx := rand.Intn(len(bonuses))
		result = append(result, bonuses[idx])
		bonuses = append(bonuses[:idx], bonuses[idx+1:]...)
		pool = append(weapons, bonuses...)
	} else {
		pool = weapons
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	need := n - len(result)
	for i := 0; i < need && i < len(pool); i++ {
		result = append(result, pool[i])
	}
	return result
}

func buildUpgradePool(g *Game) []UpgradeOption {
	if g.upgradeCount >= MaxUpgrades {
		return fallbackUpgrades(g)
	}
	pool := make([]UpgradeOption, 0, 32)
	for _, strat := range upgradeStrategies {
		if opt, ok := strat(g); ok {
			pool = append(pool, opt)
		}
	}
	if len(pool) == 0 {
		return fallbackUpgrades(g)
	}
	return pool
}

// The concrete strategies live in upgrade_weapons.go (weapon cards) and
// upgrade_bonus.go (one-time bonus items); both register into upgradeStrategies.

func fallbackUpgrades(g *Game) []UpgradeOption {
	return []UpgradeOption{
		{WeaponName: "Health", UpgradeDesc: "+20% Max HP", IconKey: "",
			Apply: func(g *Game) {
				g.player.MaxHP *= 1.2
				if g.player.HP > g.player.MaxHP {
					g.player.HP = g.player.MaxHP
				}
			}},
		{WeaponName: "Damage", UpgradeDesc: "+15% Weapon Damage", IconKey: "",
			Apply: func(g *Game) { g.weapons.WeaponDamageMult *= 1.15 }},
		{WeaponName: "Experience", UpgradeDesc: "+25% XP Bonus", IconKey: "",
			Apply: func(g *Game) { g.xpBonusMult *= 1.25 }},
	}
}

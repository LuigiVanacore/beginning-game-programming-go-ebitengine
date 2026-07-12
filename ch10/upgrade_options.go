package game

import (
	"math/rand"
)

// UpgradeOption represents a single weapon upgrade the player can choose.
type UpgradeOption struct {
	WeaponName  string
	UpgradeDesc string // e.g. "+15% Speed", "-15% Cooldown"
	IconKey     string // texture name for the weapon icon
	Apply       func(g *Game)
}

// upgradeFactory builds one card from the static pool (Registry + Strategy-style builders:
// each factory is a small “strategy” with no eligibility — ch11+ adds gating per Game state).
type upgradeFactory func() UpgradeOption

// upgradeFactories lists all possible upgrades in a fixed order; PickRandomUpgrades samples without replacement.
var upgradeFactories = []upgradeFactory{
	newKnifeSpeedUpgrade,
	newKnifeCooldownUpgrade,
	newSacredBookSpeedUpgrade,
	newFlyingAxeSpeedUpgrade,
	newFlyingAxeCooldownUpgrade,
	newFlyingAxeDamageUpgrade,
	newArmorUpgrade,
	newBootsUpgrade,
	newGemUpgrade,
	newSkullUpgrade,
	newRingUpgrade,
}

func newKnifeSpeedUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Bloody Knife", UpgradeDesc: "+15% Speed", IconKey: "bloody_knife",
		Apply: func(g *Game) { g.weapons.Knife.SpeedMult *= 1.15 },
	}
}

func newKnifeCooldownUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Bloody Knife", UpgradeDesc: "-15% Cooldown", IconKey: "bloody_knife",
		Apply: func(g *Game) { g.weapons.Knife.CooldownMult *= 0.85 },
	}
}

func newSacredBookSpeedUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Sacred Book", UpgradeDesc: "+20% Orbit Speed", IconKey: "sacred_book",
		Apply: func(g *Game) {
			if g.weapons.SacredBook != nil {
				g.weapons.SacredBook.SpeedMult *= 1.2
			}
		},
	}
}

func newFlyingAxeSpeedUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "+15% Speed", IconKey: "flying_axe",
		Apply: func(g *Game) { g.weapons.Axe.SpeedMult *= 1.15 },
	}
}

func newFlyingAxeCooldownUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "-15% Cooldown", IconKey: "flying_axe",
		Apply: func(g *Game) { g.weapons.Axe.CooldownMult *= 0.85 },
	}
}

func newFlyingAxeDamageUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "+15% Damage", IconKey: "flying_axe",
		Apply: func(g *Game) { g.weapons.Axe.DamageMult *= 1.15 },
	}
}

func newArmorUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Defense Armor", UpgradeDesc: "+20% Max HP", IconKey: "armor",
		Apply: func(g *Game) {
			g.player.MaxHP *= 1.2
			if g.player.HP > g.player.MaxHP {
				g.player.HP = g.player.MaxHP
			}
		},
	}
}

func newBootsUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Speed Boots", UpgradeDesc: "+20% Move Speed", IconKey: "boots",
		Apply: func(g *Game) { g.player.SpeedMult *= 1.2 },
	}
}

func newGemUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Experience Gem", UpgradeDesc: "+25% XP Bonus", IconKey: "gem",
		Apply: func(g *Game) { g.player.XPBonusMult *= 1.25 },
	}
}

func newSkullUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Deadly Skull", UpgradeDesc: "+15% Weapon Damage", IconKey: "skull",
		Apply: func(g *Game) { g.weapons.DamageMult *= 1.15 },
	}
}

func newRingUpgrade() UpgradeOption {
	return UpgradeOption{
		WeaponName: "Ring of Power", UpgradeDesc: "-10% All Cooldowns", IconKey: "ring",
		Apply: func(g *Game) { g.weapons.GlobalCooldownMult *= 0.9 },
	}
}

// AllUpgrades is the full pool of possible upgrades (same entries as upgradeFactories).
// Populated in init for compatibility with listings that iterate the slice directly.
var AllUpgrades []UpgradeOption

func init() {
	AllUpgrades = make([]UpgradeOption, len(upgradeFactories))
	for i, f := range upgradeFactories {
		AllUpgrades[i] = f()
	}
}

// PickRandomUpgrades returns n unique upgrades chosen at random from the pool (for level-up choice).
func PickRandomUpgrades(n int) []UpgradeOption {
	m := len(upgradeFactories)
	if n > m {
		n = m
	}
	order := rand.Perm(m)
	out := make([]UpgradeOption, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, upgradeFactories[order[i]]())
	}
	return out
}

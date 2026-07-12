package game

import (
	. "book/code/ch12/internal/core"
)

// Bonus-item upgrade strategies: one-time pickups gated by a per-item Has* flag
// and the global MaxUpgrades cap. Registered in upgradeStrategies (upgrade_options.go).

func strategyBonusArmor(g *Game) (UpgradeOption, bool) {
	if g.upgradeCount >= MaxUpgrades || g.player.HasArmor {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Defense Armor", UpgradeDesc: "+20% Max HP", IconKey: "armor",
		Apply: func(g *Game) {
			g.player.HasArmor = true
			g.player.MaxHP *= 1.2
			if g.player.HP > g.player.MaxHP {
				g.player.HP = g.player.MaxHP
			}
			g.upgradeCount++
		},
	}, true
}

func strategyBonusBoots(g *Game) (UpgradeOption, bool) {
	if g.upgradeCount >= MaxUpgrades || g.player.HasBoots {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Speed Boots", UpgradeDesc: "+20% Move Speed", IconKey: "boots",
		Apply: func(g *Game) { g.player.HasBoots = true; g.playerSpeedMult *= 1.2; g.upgradeCount++ },
	}, true
}

func strategyBonusGem(g *Game) (UpgradeOption, bool) {
	if g.upgradeCount >= MaxUpgrades || g.player.HasGem {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Experience Gem", UpgradeDesc: "+25% XP Bonus", IconKey: "gem",
		Apply: func(g *Game) { g.player.HasGem = true; g.xpBonusMult *= 1.25; g.upgradeCount++ },
	}, true
}

func strategyBonusSkull(g *Game) (UpgradeOption, bool) {
	if g.upgradeCount >= MaxUpgrades || g.player.HasSkull {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Deadly Skull", UpgradeDesc: "+15% Weapon Damage", IconKey: "skull",
		Apply: func(g *Game) { g.player.HasSkull = true; g.weapons.WeaponDamageMult *= 1.15; g.upgradeCount++ },
	}, true
}

func strategyBonusRing(g *Game) (UpgradeOption, bool) {
	if g.upgradeCount >= MaxUpgrades || g.player.HasRing {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Ring of Power", UpgradeDesc: "-10% All Cooldowns", IconKey: "ring",
		Apply: func(g *Game) { g.player.HasRing = true; g.weapons.GlobalCooldownMult *= 0.9; g.upgradeCount++ },
	}, true
}

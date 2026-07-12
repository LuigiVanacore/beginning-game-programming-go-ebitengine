package game

import (
	. "book/code/ch11/internal/core"
)

// Weapon upgrade strategies: each returns an UpgradeOption plus whether the card
// is eligible for the current game state. They are registered in upgradeStrategies
// (upgrade_options.go) and evaluated in order when the pool is built.

// --- Knife ---

func strategyKnifeUnlock(g *Game) (UpgradeOption, bool) {
	if g.weapons.KnifeUnlocked {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Bloody Knife", UpgradeDesc: "Unlock +15% Speed", IconKey: "bloody_knife",
		Apply: func(g *Game) {
			w := g.weapons
			w.KnifeUnlocked = true
			w.Knife.SpeedMult *= 1.15
			w.KnifeSpeedLvl++
			g.upgradeCount++
			g.weapons.Mount(w.Knife, g.player.WeaponsRoot)
		},
	}, true
}

func strategyKnifeSpeed(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.KnifeUnlocked || w.KnifeSpeedLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Bloody Knife", UpgradeDesc: "+15% Speed", IconKey: "bloody_knife",
		Apply: func(g *Game) { g.weapons.Knife.SpeedMult *= 1.15; g.weapons.KnifeSpeedLvl++; g.upgradeCount++ },
	}, true
}

func strategyKnifeCooldown(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.KnifeUnlocked || w.KnifeCooldownLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Bloody Knife", UpgradeDesc: "-15% Cooldown", IconKey: "bloody_knife",
		Apply: func(g *Game) { g.weapons.Knife.CooldownMult *= 0.85; g.weapons.KnifeCooldownLvl++; g.upgradeCount++ },
	}, true
}

// --- Flying Axe ---

func strategyFlyingAxeUnlock(g *Game) (UpgradeOption, bool) {
	if g.weapons.FlyingAxeUnlocked {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "Unlock +15% Speed", IconKey: "flying_axe",
		Apply: func(g *Game) {
			w := g.weapons
			w.FlyingAxeUnlocked = true
			w.Axe.SpeedMult *= 1.15
			w.FlyingAxeSpeedLvl++
			g.upgradeCount++
			g.weapons.Mount(w.Axe, g.player.WeaponsRoot)
		},
	}, true
}

func strategyFlyingAxeSpeed(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.FlyingAxeUnlocked || w.FlyingAxeSpeedLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "+15% Speed", IconKey: "flying_axe",
		Apply: func(g *Game) { g.weapons.Axe.SpeedMult *= 1.15; g.weapons.FlyingAxeSpeedLvl++; g.upgradeCount++ },
	}, true
}

func strategyFlyingAxeCooldown(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.FlyingAxeUnlocked || w.FlyingAxeCooldownLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "-15% Cooldown", IconKey: "flying_axe",
		Apply: func(g *Game) {
			g.weapons.Axe.CooldownMult *= 0.85
			g.weapons.FlyingAxeCooldownLvl++
			g.upgradeCount++
		},
	}, true
}

func strategyFlyingAxeRotation(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.FlyingAxeUnlocked || w.FlyingAxeRotationLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Flying Axe", UpgradeDesc: "+25% Rotation", IconKey: "flying_axe",
		Apply: func(g *Game) {
			g.weapons.Axe.RotationMult *= 1.25
			g.weapons.FlyingAxeRotationLvl++
			g.upgradeCount++
		},
	}, true
}

// --- Sacred Book ---

func strategySacredBookUnlock(g *Game) (UpgradeOption, bool) {
	if g.weapons.SacredBookUnlocked {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Sacred Book", UpgradeDesc: "Unlock +20% Orbit Speed", IconKey: "sacred_book",
		Apply: func(g *Game) {
			w := g.weapons
			w.SacredBookUnlocked = true
			w.SacredBookSpeedLvl++
			g.upgradeCount++
			if w.SacredBook != nil {
				w.SacredBook.SpeedMult *= 1.2
				g.weapons.Mount(w.SacredBook, g.player.WeaponsRoot)
			}
		},
	}, true
}

func strategySacredBookSpeed(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.SacredBookUnlocked || w.SacredBookSpeedLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Sacred Book", UpgradeDesc: "+20% Orbit Speed", IconKey: "sacred_book",
		Apply: func(g *Game) {
			if g.weapons.SacredBook != nil {
				g.weapons.SacredBook.SpeedMult *= 1.2
			}
			g.weapons.SacredBookSpeedLvl++
			g.upgradeCount++
		},
	}, true
}

// --- Holy Shield ---

func strategyHolyShieldUnlock(g *Game) (UpgradeOption, bool) {
	if g.weapons.HolyShieldUnlocked {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Holy Shield", UpgradeDesc: "Unlock +20% Radius", IconKey: "holy_shield",
		Apply: func(g *Game) {
			w := g.weapons
			w.HolyShieldUnlocked = true
			w.HolyShieldRadiusMult = 1.2
			w.HolyShieldRadiusLvl++
			g.upgradeCount++
			w.HolyShield = NewHolyShieldWeapon(g.engine, g.player, HolyShieldRadius*w.HolyShieldRadiusMult)
			g.weapons.Mount(w.HolyShield, g.player.WeaponsRoot)
		},
	}, true
}

func strategyHolyShieldRadius(g *Game) (UpgradeOption, bool) {
	w := g.weapons
	if !w.HolyShieldUnlocked || w.HolyShieldRadiusLvl >= MaxStatLevelPerWeapon {
		return UpgradeOption{}, false
	}
	return UpgradeOption{
		WeaponName: "Holy Shield", UpgradeDesc: "+20% Radius", IconKey: "holy_shield",
		Apply: func(g *Game) {
			w := g.weapons
			w.HolyShieldRadiusMult *= 1.2
			w.HolyShieldRadiusLvl++
			g.upgradeCount++
			if w.HolyShield != nil {
				// Unregister drops it from the update list and removes its ring collider;
				// DetachChild then stops it drawing before the larger shield is built.
				g.weapons.Unregister(w.HolyShield)
				g.player.WeaponsRoot.DetachChild(w.HolyShield)
			}
			w.HolyShield = NewHolyShieldWeapon(g.engine, g.player, HolyShieldRadius*w.HolyShieldRadiusMult)
			g.weapons.Mount(w.HolyShield, g.player.WeaponsRoot)
		},
	}, true
}

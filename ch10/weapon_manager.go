package game

import . "book/code/ch10/internal/core"

// WeaponManager owns the active weapon list, a typed reference to each weapon, and the
// shared multipliers that affect every weapon. Per-weapon multipliers (knife speed,
// axe rotation, etc.) now live on the weapon structs themselves; upgrades reach them
// through the typed references below.
type WeaponManager struct {
	weapons []PlayerWeapon

	// Typed references so an upgrade can target a specific weapon's own stats.
	Knife      *BloodyKnifeWeapon
	SacredBook *SacredBook
	HolyShield *HolyShieldWeapon
	Axe        *FlyingAxeWeapon

	// Shared multipliers — affect every weapon.
	GlobalCooldownMult float64
	DamageMult         float64
}

// NewWeaponManager creates an empty manager with the shared multipliers at their neutral value (1.0).
func NewWeaponManager() *WeaponManager {
	return &WeaponManager{
		weapons:            make([]PlayerWeapon, 0),
		GlobalCooldownMult: 1.0,
		DamageMult:         1.0,
	}
}

// Mount adds a weapon to the scene graph under parent and registers it for updates.
func (m *WeaponManager) Mount(w PlayerWeapon, parent SceneNode) {
	if w == nil || parent == nil {
		return
	}
	parent.AddChildren(w)
	m.weapons = append(m.weapons, w)
}

// Update calls UpdateWeapon on every mounted weapon, passing the manager so weapons
// can read their current multipliers.
func (m *WeaponManager) Update(engine *Engine, player *Player, cursor *Cursor) {
	for _, w := range m.weapons {
		w.UpdateWeapon(engine, player, cursor, m)
	}
}

// TryReleaseProjectileByCollider returns a pooled projectile to its weapon when hit.
func (m *WeaponManager) TryReleaseProjectileByCollider(c *Collider) bool {
	for _, w := range m.weapons {
		if pc, ok := w.(ProjectileCarrier); ok && pc.TryReleaseProjectile(c) {
			return true
		}
	}
	return false
}

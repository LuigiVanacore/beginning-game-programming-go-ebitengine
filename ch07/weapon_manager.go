package game

import . "book/code/ch07/internal/core"

// WeaponManager owns the BloodyKnifeWeapon and drives its update each frame.
type WeaponManager struct {
	bloodyKnife *BloodyKnifeWeapon
}

// NewWeaponManager wraps an existing BloodyKnifeWeapon.
func NewWeaponManager(knife *BloodyKnifeWeapon) *WeaponManager {
	return &WeaponManager{bloodyKnife: knife}
}

// Update fires new projectiles toward the cursor and advances all active ones.
func (m *WeaponManager) Update(playerX, playerY, cursorX, cursorY float64, onHit func(*Collider, *Collider)) {
	m.bloodyKnife.Update(playerX, playerY, cursorX, cursorY, onHit)
	m.bloodyKnife.UpdateProjectiles(playerX, playerY)
}

// TryReleaseProjectileByCollider returns true and releases the projectile if c belongs to this weapon.
func (m *WeaponManager) TryReleaseProjectileByCollider(c *Collider) bool {
	return m.bloodyKnife.TryReleaseProjectileByCollider(c)
}

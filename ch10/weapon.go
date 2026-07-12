package game

import . "book/code/ch10/internal/core"

// PlayerWeapon is a Node2D mounted under player.WeaponsRoot; WeaponManager drives UpdateWeapon each frame.
// wm gives weapons read-only access to the shared multipliers owned by WeaponManager.
type PlayerWeapon interface {
	SceneNode
	UpdateWeapon(engine *Engine, player *Player, cursor *Cursor, wm *WeaponManager)
}

// ProjectileCarrier lets removal processing return pooled projectiles to the owning weapon.
type ProjectileCarrier interface {
	TryReleaseProjectile(c *Collider) bool
}

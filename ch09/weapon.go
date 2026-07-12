package game

import . "book/code/ch09/internal/core"

// PlayerWeapon is a Node2D under player.WeaponsRoot; Player drives UpdateWeapon each frame.
type PlayerWeapon interface {
	SceneNode
	UpdateWeapon(engine *Engine, player *Player, cursor *Cursor)
}

// ProjectileCarrier lets removal processing return pooled projectiles to the owning weapon.
type ProjectileCarrier interface {
	TryReleaseProjectile(c *Collider) bool
}

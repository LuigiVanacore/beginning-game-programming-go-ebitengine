package game

import (
	. "book/code/ch11/internal/core"
)

// playerWeapon is a Node2D mounted under player.WeaponsRoot; Player drives Update via updateWeapon.
type playerWeapon interface {
	SceneNode
	updateWeapon(engine *Engine, player *Player, cursor *Cursor)
}

// projectileCarrier lets removal processing return pooled projectiles to the owning weapon.
type projectileCarrier interface {
	tryReleaseProjectile(c *Collider) bool
}

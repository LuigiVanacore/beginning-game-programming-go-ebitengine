package game

import (
	. "book/code/ch13/internal/core"
	"math"
	"time"
)

const (
	flyingAxeProjectileSpeedPxPerFrame = 5.333333333333333 // displacement per Update step
	flyingAxeRadius                    = 8.0               // collision radius
	flyingAxeCooldown                  = 2.0               // seconds between throws
	flyingAxeRotationRadPerFrame       = 0.13333333333333333 // spin per Update step
	flyingAxeMaxDist                   = 800.0        // max distance from player before despawn
	flyingAxeDamage                    = 5.0          // HP removed from an enemy per hit (before the shared damage multiplier)
)

// FlyingAxeWeapon is a Node2D under player.WeaponsRoot.
// On each cooldown tick it throws an axe that spins while flying horizontally
// (direction determined by cursor position relative to player).
// SpeedMult, CooldownMult, and RotationMult are this weapon's own stats; the shared
// multipliers stay on WeaponManager and are read through pointers.
type FlyingAxeWeapon struct {
	Node2D
	pool           *ProjectilePool
	timer          *Timer
	projectiles    []*Projectile
	SpeedMult      float64  // this weapon's own speed multiplier
	CooldownMult   float64  // this weapon's own cooldown multiplier
	RotationMult   float64  // this weapon's own spin multiplier
	globalCoolMult *float64 // pointer to WeaponManager.GlobalCooldownMult (shared)
	damageMult     *float64 // pointer to WeaponManager.WeaponDamageMult (shared)
}

// NewFlyingAxeWeapon creates the weapon with its own multipliers at 1.0 and pointers
// to the shared multipliers owned by WeaponManager.
func NewFlyingAxeWeapon(
	pool *ProjectilePool,
	globalCoolMult, damageMult *float64,
) *FlyingAxeWeapon {
	return &FlyingAxeWeapon{
		Node2D:         *NewNode2D(NameFlyingAxe + "_weapon"),
		pool:           pool,
		timer:          NewTimer(time.Duration(flyingAxeCooldown*float64(time.Second)), true).Start(),
		projectiles:    make([]*Projectile, 0),
		SpeedMult:      1.0,
		CooldownMult:   1.0,
		RotationMult:   1.0,
		globalCoolMult: globalCoolMult,
		damageMult:     damageMult,
	}
}

func (w *FlyingAxeWeapon) updateWeapon(_ *Engine, player *Player, cursor *Cursor) {
	if w == nil || player == nil || cursor == nil {
		return
	}
	// Update cooldown from current multipliers.
	cd := flyingAxeCooldown * w.CooldownMult * (*w.globalCoolMult)
	w.timer.SetDuration(time.Duration(cd * float64(time.Second)))

	if w.timer.Update() {
		px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
		cx := cursor.GetWorldPosition().X()
		dir := 1.0
		if cx < px {
			dir = -1.0
		}
		speed := flyingAxeProjectileSpeedPxPerFrame * w.SpeedMult * (*w.damageMult)
		rotSpeed := flyingAxeRotationRadPerFrame * w.RotationMult
		proj := w.pool.AcquireAt(px, py, dir*speed, 0, rotSpeed, func(projectileCol, hitCollider *Collider) {
			player.QueueWeaponHit(projectileCol, hitCollider, flyingAxeDamage*(*w.damageMult))
		})
		if proj != nil {
			w.projectiles = append(w.projectiles, proj)
		}
	}
	px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
	w.moveProjectiles(px, py)
}

func (w *FlyingAxeWeapon) moveProjectiles(playerX, playerY float64) {
	if w == nil {
		return
	}
	keep := make([]*Projectile, 0, len(w.projectiles))
	for _, proj := range w.projectiles {
		x, y := proj.GetPosition().X(), proj.GetPosition().Y()
		proj.SetPosition(x+proj.Vx, y+proj.Vy)
		// Apply spin rotation.
		if proj.RotationSpeed != 0 {
			proj.SetRotation(proj.GetRotation() + proj.RotationSpeed)
		}
		dist := math.Hypot(x-playerX, y-playerY)
		if dist < flyingAxeMaxDist {
			keep = append(keep, proj)
		} else {
			w.pool.Release(proj)
		}
	}
	w.projectiles = keep
}

// ForEachProjectile calls fn with each active projectile's world position and
// velocity, so the caller can emit a particle trail behind it (ch13).
func (w *FlyingAxeWeapon) ForEachProjectile(fn func(x, y, vx, vy float64)) {
	if w == nil {
		return
	}
	for _, proj := range w.projectiles {
		p := proj.GetPosition()
		fn(p.X(), p.Y(), proj.Vx, proj.Vy)
	}
}

func (w *FlyingAxeWeapon) tryReleaseProjectile(c *Collider) bool {
	if w == nil || c == nil {
		return false
	}
	for i, proj := range w.projectiles {
		if proj.GetCollider() == c {
			w.pool.Release(proj)
			w.projectiles = append(w.projectiles[:i], w.projectiles[i+1:]...)
			return true
		}
	}
	return false
}

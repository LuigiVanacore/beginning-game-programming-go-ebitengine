package game

import (
	. "book/code/ch10/internal/core"
	"math"
	"time"
)

const (
	// FlyingAxeSpeed is the base horizontal throw displacement per frame (Update step).
	FlyingAxeSpeed         = 5.333333333333333
	FlyingAxeRadius        = 8.0   // collision radius
	FlyingAxeCooldown      = 2.0   // seconds between throws
	flyingAxeRotationSpeed = 8.0   // rad/s
	flyingAxeMaxDist       = 800.0 // max distance from player before despawn
)

// FlyingAxeWeapon is a Node2D under player.WeaponsRoot.
// On each cooldown tick it throws an axe that spins while flying horizontally.
// SpeedMult, CooldownMult, DamageMult, and RotationMult are this weapon's own tunable
// stats, read each frame so an upgrade takes effect immediately.
type FlyingAxeWeapon struct {
	Node2D
	pool         *ProjectilePool
	timer        *Timer
	projectiles  []*Projectile
	SpeedMult    float64
	CooldownMult float64
	DamageMult   float64
	RotationMult float64
}

// NewFlyingAxeWeapon creates the weapon with its multipliers at the neutral value (1.0).
func NewFlyingAxeWeapon(pool *ProjectilePool) *FlyingAxeWeapon {
	return &FlyingAxeWeapon{
		Node2D:       *NewNode2D(NameFlyingAxe + "_weapon"),
		pool:         pool,
		timer:        NewTimer(time.Duration(FlyingAxeCooldown*float64(time.Second)), true).Start(),
		projectiles:  make([]*Projectile, 0),
		SpeedMult:    1.0,
		CooldownMult: 1.0,
		DamageMult:   1.0,
		RotationMult: 1.0,
	}
}

// UpdateWeapon runs cooldown, throw, spin, and despawn logic.
func (w *FlyingAxeWeapon) UpdateWeapon(_ *Engine, player *Player, cursor *Cursor, wm *WeaponManager) {
	if w == nil || player == nil || cursor == nil {
		return
	}
	cd := FlyingAxeCooldown * w.CooldownMult * wm.GlobalCooldownMult
	w.timer.SetDuration(time.Duration(cd * float64(time.Second)))

	if w.timer.Update() {
		px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
		cx := cursor.GetWorldPosition().X()
		dir := 1.0
		if cx < px {
			dir = -1.0
		}
		speed := FlyingAxeSpeed * w.SpeedMult * w.DamageMult * wm.DamageMult
		rotSpeed := flyingAxeRotationSpeed * w.RotationMult
		proj := w.pool.AcquireAt(px, py, dir*speed, 0, rotSpeed, func(projectileCol, hitCollider *Collider) {
			player.QueueWeaponHit(projectileCol, hitCollider)
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

// TryReleaseProjectile returns a pooled projectile to the pool if c matches an active axe.
func (w *FlyingAxeWeapon) TryReleaseProjectile(c *Collider) bool {
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

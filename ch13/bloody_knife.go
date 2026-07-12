package game

import (
	. "book/code/ch13/internal/core"
	"math"
	"time"
)

const (
	knifeProjectileSpeedPxPerFrame = 6.666666666666667 // projectile displacement per Update step
	knifeRadius                    = 6.0          // collision radius
	knifeCooldown                  = 1.2          // seconds between shots
	knifeMaxDist                   = 800.0        // max distance from player before despawn
	knifeDamage                    = 3.0          // HP removed from an enemy per hit (before the shared damage multiplier)
)

// BloodyKnifeWeapon is a Node2D under player.WeaponsRoot.
// SpeedMult and CooldownMult are this weapon's own stats (scaled by upgrades).
// The shared multipliers stay on WeaponManager and are read through pointers.
type BloodyKnifeWeapon struct {
	Node2D
	pool           *ProjectilePool
	timer          *Timer
	projectiles    []*Projectile
	SpeedMult      float64  // this weapon's own speed multiplier
	CooldownMult   float64  // this weapon's own cooldown multiplier
	globalCoolMult *float64 // pointer to WeaponManager.GlobalCooldownMult (shared)
	damageMult     *float64 // pointer to WeaponManager.WeaponDamageMult (shared)
}

// NewBloodyKnifeWeapon creates the weapon with its own multipliers at 1.0 and pointers
// to the shared multipliers owned by WeaponManager.
func NewBloodyKnifeWeapon(
	pool *ProjectilePool,
	globalCoolMult, damageMult *float64,
) *BloodyKnifeWeapon {
	return &BloodyKnifeWeapon{
		Node2D:         *NewNode2D(NameBloodyKnife + "_weapon"),
		pool:           pool,
		timer:          NewTimer(time.Duration(knifeCooldown*float64(time.Second)), true).Start(),
		projectiles:    make([]*Projectile, 0),
		SpeedMult:      1.0,
		CooldownMult:   1.0,
		globalCoolMult: globalCoolMult,
		damageMult:     damageMult,
	}
}

func (w *BloodyKnifeWeapon) updateWeapon(_ *Engine, player *Player, cursor *Cursor) {
	if w == nil || player == nil || cursor == nil {
		return
	}
	// Update timer duration to reflect current cooldown multipliers.
	cd := knifeCooldown * w.CooldownMult * (*w.globalCoolMult)
	w.timer.SetDuration(time.Duration(cd * float64(time.Second)))

	if w.timer.Update() {
		px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
		cx, cy := cursor.GetWorldPosition().X(), cursor.GetWorldPosition().Y()
		speed := knifeProjectileSpeedPxPerFrame * w.SpeedMult * (*w.damageMult)
		ux, uy := unitDirection2D(px, py, cx, cy)
		vx, vy := ux*speed, uy*speed
		// Rotate sprite to face direction: already handled by Projectile.reset via atan2Pi2.
		angle := math.Atan2(vy, vx) + math.Pi/2
		_ = angle // projectile pool sets rotation in reset()
		proj := w.pool.AcquireAt(px, py, vx, vy, 0, func(projectileCol, hitCollider *Collider) {
			player.QueueWeaponHit(projectileCol, hitCollider, knifeDamage*(*w.damageMult))
		})
		if proj != nil {
			w.projectiles = append(w.projectiles, proj)
		}
	}
	px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
	w.moveProjectiles(px, py)
}

func (w *BloodyKnifeWeapon) moveProjectiles(playerX, playerY float64) {
	if w == nil {
		return
	}
	keep := make([]*Projectile, 0, len(w.projectiles))
	for _, proj := range w.projectiles {
		x, y := proj.GetPosition().X(), proj.GetPosition().Y()
		proj.SetPosition(x+proj.Vx, y+proj.Vy)
		dist := math.Hypot(x-playerX, y-playerY)
		if dist < knifeMaxDist {
			keep = append(keep, proj)
		} else {
			w.pool.Release(proj)
		}
	}
	w.projectiles = keep
}

// ForEachProjectile calls fn with each active projectile's world position and
// velocity, so the caller can emit a particle trail behind it (ch13).
func (w *BloodyKnifeWeapon) ForEachProjectile(fn func(x, y, vx, vy float64)) {
	if w == nil {
		return
	}
	for _, proj := range w.projectiles {
		p := proj.GetPosition()
		fn(p.X(), p.Y(), proj.Vx, proj.Vy)
	}
}

func (w *BloodyKnifeWeapon) tryReleaseProjectile(c *Collider) bool {
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

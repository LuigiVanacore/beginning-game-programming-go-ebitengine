package game

import (
	. "book/code/ch10/internal/core"
	"math"
	"time"
)

const (
	// KnifeSpeed is the base projectile displacement per frame (Update step).
	KnifeSpeed    = 6.666666666666667
	KnifeRadius   = 6.0   // collision radius
	KnifeCooldown = 1.2   // seconds between shots
	knifeMaxDist  = 800.0 // max distance from player before despawn
)

// BloodyKnifeWeapon is a Node2D under player.WeaponsRoot.
// SpeedMult and CooldownMult are this weapon's own tunable stats; upgrades scale them
// directly and UpdateWeapon reads them each frame, so changes take effect immediately.
type BloodyKnifeWeapon struct {
	Node2D
	pool         *ProjectilePool
	timer        *Timer
	projectiles  []*Projectile
	SpeedMult    float64
	CooldownMult float64
}

// NewBloodyKnifeWeapon creates the weapon with its multipliers at the neutral value (1.0).
func NewBloodyKnifeWeapon(pool *ProjectilePool) *BloodyKnifeWeapon {
	return &BloodyKnifeWeapon{
		Node2D:       *NewNode2D(NameBloodyKnife + "_weapon"),
		pool:         pool,
		timer:        NewTimer(time.Duration(KnifeCooldown*float64(time.Second)), true).Start(),
		projectiles:  make([]*Projectile, 0),
		SpeedMult:    1.0,
		CooldownMult: 1.0,
	}
}

// UpdateWeapon runs cooldown, firing, and projectile motion.
func (w *BloodyKnifeWeapon) UpdateWeapon(_ *Engine, player *Player, cursor *Cursor, wm *WeaponManager) {
	if w == nil || player == nil || cursor == nil {
		return
	}
	cd := KnifeCooldown * w.CooldownMult * wm.GlobalCooldownMult
	w.timer.SetDuration(time.Duration(cd * float64(time.Second)))

	if w.timer.Update() {
		px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
		cx, cy := cursor.GetWorldPosition().X(), cursor.GetWorldPosition().Y()
		speed := KnifeSpeed * w.SpeedMult * wm.DamageMult
		ux, uy := unitDirection2D(px, py, cx, cy)
		proj := w.pool.AcquireAt(px, py, ux*speed, uy*speed, 0, func(projectileCol, hitCollider *Collider) {
			player.QueueWeaponHit(projectileCol, hitCollider)
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

// TryReleaseProjectile returns a pooled projectile to the pool if c matches an active knife.
func (w *BloodyKnifeWeapon) TryReleaseProjectile(c *Collider) bool {
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

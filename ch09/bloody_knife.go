package game

import (
	. "book/code/ch09/internal/core"
	"math"
	"time"
)

const (
	// KnifeProjectileSpeedPxPerFrame is displacement per Update step for the bloody knife projectile.
	KnifeProjectileSpeedPxPerFrame = 6.666666666666667
	KnifeRadius                    = 6.0   // collision radius
	KnifeCooldown                  = 1.2   // seconds between shots
	knifeMaxDist                   = 800.0 // max distance from player before returning to pool
)

// BloodyKnifeWeapon is a Node2D under player.WeaponsRoot: pool, cooldown, active projectiles.
type BloodyKnifeWeapon struct {
	Node2D
	pool        *ProjectilePool
	timer       *Timer
	projectiles []*Projectile
}

// NewBloodyKnifeWeapon creates the spawner with pool and timer.
func NewBloodyKnifeWeapon(pool *ProjectilePool, cooldownSeconds float64) *BloodyKnifeWeapon {
	return &BloodyKnifeWeapon{
		Node2D:      *NewNode2D(NameBloodyKnife + "_weapon"),
		pool:        pool,
		timer:       NewTimer(time.Duration(cooldownSeconds*float64(time.Second)), true).Start(),
		projectiles: make([]*Projectile, 0),
	}
}

// UpdateWeapon runs cooldown, firing, and projectile motion.
func (w *BloodyKnifeWeapon) UpdateWeapon(_ *Engine, player *Player, cursor *Cursor) {
	if w == nil || player == nil || cursor == nil {
		return
	}
	px, py := player.GetWorldPosition().X(), player.GetWorldPosition().Y()
	cx, cy := cursor.GetWorldPosition().X(), cursor.GetWorldPosition().Y()
	if w.timer.Update() {
		proj := w.pool.Acquire(px, py, cx, cy, func(projectileCol, hitCollider *Collider) {
			player.QueueWeaponHit(projectileCol, hitCollider)
		})
		if proj != nil {
			w.projectiles = append(w.projectiles, proj)
		}
	}
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
		dist := math.Sqrt((x-playerX)*(x-playerX) + (y-playerY)*(y-playerY))
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

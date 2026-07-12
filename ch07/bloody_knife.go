package game

import (
	. "book/code/ch07/internal/core"
	"math"
	"time"
)

const (
	// KnifeProjectileSpeedPxPerFrame is displacement per Update step for the bloody knife projectile.
	KnifeProjectileSpeedPxPerFrame = 6
	KnifeRadius                    = 6.0   // collision radius
	KnifeCooldown                  = 1.2   // seconds between shots
	knifeMaxDist                   = 800.0 // max distance from player before returning to pool
)

// BloodyKnifeWeapon manages the pool, cooldown, and active projectiles; it embeds Node2D as the weapon's logical root (Ch7).
type BloodyKnifeWeapon struct {
	Node2D
	pool        *ProjectilePool
	timer       *Timer
	projectiles []*Projectile
}

// NewBloodyKnifeWeapon creates the spawner with pool and timer.
func NewBloodyKnifeWeapon(
	pool *ProjectilePool,
	cooldownSeconds float64,
) *BloodyKnifeWeapon {
	return &BloodyKnifeWeapon{
		Node2D:      *NewNode2D(NameBloodyKnife + "_weapon"),
		pool:        pool,
		timer:       NewTimer(time.Duration(cooldownSeconds*float64(time.Second)), true).Start(),
		projectiles: make([]*Projectile, 0),
	}
}

// Update checks the cooldown and, when it expires, fires a projectile toward the cursor.
func (w *BloodyKnifeWeapon) Update(
	playerX, playerY, cursorX, cursorY float64,
	onHit func(proj *Collider, other *Collider),
) {
	if w == nil || !w.timer.Update() {
		return
	}
	proj := w.pool.Acquire(playerX, playerY, cursorX, cursorY, onHit)
	if proj != nil {
		w.projectiles = append(w.projectiles, proj)
	}
}

// UpdateProjectiles moves the active knives and returns to the pool those too far from the player.
func (w *BloodyKnifeWeapon) UpdateProjectiles(playerX, playerY float64) {
	if w == nil {
		return
	}
	newProj := make([]*Projectile, 0, len(w.projectiles))
	for _, proj := range w.projectiles {
		x, y := proj.GetPosition().X(), proj.GetPosition().Y()
		proj.SetPosition(x+proj.Vx, y+proj.Vy)
		dist := math.Sqrt((x-playerX)*(x-playerX) + (y-playerY)*(y-playerY))
		if dist < knifeMaxDist {
			newProj = append(newProj, proj)
		} else {
			w.pool.Release(proj)
		}
	}
	w.projectiles = newProj
}

// TryReleaseProjectileByCollider removes the projectile with collider c from the active list and returns it to the pool.
func (w *BloodyKnifeWeapon) TryReleaseProjectileByCollider(c *Collider) bool {
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

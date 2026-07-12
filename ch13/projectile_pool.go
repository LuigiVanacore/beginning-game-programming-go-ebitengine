package game

import (
	. "book/code/ch13/internal/core"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	defaultPoolCapacity = 32
	projectilePoolSize  = 32 // pre-allocated projectile slots
)

// ProjectilePool preallocates N Projectiles and reuses them to avoid allocations during gameplay.
type ProjectilePool struct {
	projectiles []*Projectile
	available   []*Projectile
	engine      *Engine
	texture     *ebiten.Image
	radius      float64
	speed       float64
	layer       int
	namePrefix  string
}

// NewProjectilePool creates a pool with capacity preallocated projectiles.
func NewProjectilePool(
	engine *Engine,
	texture *ebiten.Image,
	radius, speed float64,
	layer int,
	namePrefix string,
	capacity int,
) *ProjectilePool {
	capacity = normalizedPoolCapacity(capacity)
	namePrefix = normalizedProjectileNamePrefix(namePrefix)

	pool := &ProjectilePool{
		engine:     engine,
		texture:    texture,
		radius:     radius,
		speed:      speed,
		layer:      layer,
		namePrefix: namePrefix,
	}
	pool.projectiles, pool.available = preallocateProjectiles(texture, radius, namePrefix, capacity)
	return pool
}

func normalizedPoolCapacity(capacity int) int {
	if capacity <= 0 {
		return defaultPoolCapacity
	}
	return capacity
}

func normalizedProjectileNamePrefix(namePrefix string) string {
	if namePrefix == "" {
		return NameProjectileDef
	}
	return namePrefix
}

func preallocateProjectiles(texture *ebiten.Image, radius float64, namePrefix string, capacity int) (all, free []*Projectile) {
	all = make([]*Projectile, 0, capacity)
	free = make([]*Projectile, 0, capacity)
	for i := 0; i < capacity; i++ {
		name := projectileInstanceName(namePrefix, i, capacity)
		proj := newProjectile(name, texture, radius)
		if proj != nil {
			all = append(all, proj)
			free = append(free, proj)
		}
	}
	return all, free
}

func projectileInstanceName(namePrefix string, index, capacity int) string {
	if capacity <= 1 {
		return namePrefix
	}
	return fmt.Sprintf("%s_%d", namePrefix, index)
}

// Acquire takes a projectile from the pool, aims toward (targetX, targetY), and adds it to the world.
// Returns nil if the pool is exhausted or texture is nil.
func (p *ProjectilePool) Acquire(
	spawnX, spawnY, targetX, targetY float64,
	onHit func(proj *Collider, other *Collider),
) *Projectile {
	if !p.canAcquire() {
		return nil
	}
	proj := p.popAvailable()
	vx, vy := p.velocityToward(spawnX, spawnY, targetX, targetY)
	proj.reset(spawnX, spawnY, vx, vy, 0)
	p.wireOnHit(proj, onHit)
	p.addToWorld(proj)
	return proj
}

// AcquireAt takes a projectile from the pool with explicit velocity and rotation speed (new in ch10).
// Use this when weapons compute their own velocity (e.g. to apply speed multipliers).
func (p *ProjectilePool) AcquireAt(
	spawnX, spawnY, vx, vy, rotSpeed float64,
	onHit func(proj *Collider, other *Collider),
) *Projectile {
	if !p.canAcquire() {
		return nil
	}
	proj := p.popAvailable()
	proj.reset(spawnX, spawnY, vx, vy, rotSpeed)
	p.wireOnHit(proj, onHit)
	p.addToWorld(proj)
	return proj
}

func (p *ProjectilePool) canAcquire() bool {
	return len(p.available) > 0 && p.texture != nil
}

func (p *ProjectilePool) popAvailable() *Projectile {
	i := len(p.available) - 1
	proj := p.available[i]
	p.available = p.available[:i]
	return proj
}

func (p *ProjectilePool) velocityToward(spawnX, spawnY, targetX, targetY float64) (vx, vy float64) {
	ux, uy := unitDirection2D(spawnX, spawnY, targetX, targetY)
	return ux * p.speed, uy * p.speed
}

// unitDirection2D returns a normalized direction from (fromX,fromY) toward (toX,toY), or (1,0) if too short.
func unitDirection2D(fromX, fromY, toX, toY float64) (ux, uy float64) {
	dx := toX - fromX
	dy := toY - fromY
	dist := math.Hypot(dx, dy)
	if dist < 1 {
		return 1, 0
	}
	return dx / dist, dy / dist
}

func (p *ProjectilePool) wireOnHit(proj *Projectile, onHit func(*Collider, *Collider)) {
	if onHit == nil {
		return
	}
	col := proj.GetCollider()
	if col == nil {
		return
	}
	col.SetOnCollide(func(other *Collider) {
		onHit(col, other)
	})
}

func (p *ProjectilePool) addToWorld(proj *Projectile) {
	p.engine.World().AddNodeToLayer(proj, p.layer)
	p.engine.CollisionManager().AddCollider(proj.GetCollider())
}

// Release removes the projectile from the world and CollisionManager and returns it to the pool.
func (p *ProjectilePool) Release(proj *Projectile) {
	if proj == nil {
		return
	}
	p.engine.World().RemoveNode(proj)
	p.removeColliderFromEngine(proj)
	p.available = append(p.available, proj)
}

func (p *ProjectilePool) removeColliderFromEngine(proj *Projectile) {
	col := proj.GetCollider()
	if col == nil {
		return
	}
	p.engine.CollisionManager().RemoveCollider(col)
	col.SetOnCollide(nil)
}

package game

import (
	. "book/code/ch07/internal/core"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Projectile is a generic projectile: it embeds Node2D and has Collider and Sprite as children.
// Used by the object pool and configured by individual weapons (e.g. BloodyKnifeWeapon).
type Projectile struct {
	Node2D
	Vx float64
	Vy float64
}

// GetCollider returns the Collider child (for CollisionManager and callbacks).
func (p *Projectile) GetCollider() *Collider {
	for _, c := range p.GetChildren() {
		if col, ok := c.(*Collider); ok {
			return col
		}
	}
	return nil
}

// newProjectile creates a Projectile not yet added to the world. Used by the pool.
func newProjectile(name string, texture *ebiten.Image, radius float64) *Projectile {
	if texture == nil {
		return nil
	}
	proj := &Projectile{
		Node2D: *NewNode2D(name),
	}
	proj.SetPosition(0, 0)

	shape := NewCollisionCircle(radius)
	mask := NewCollisionMask(LayerProjectile, LayerEnemy)
	// Not registered here: ProjectilePool adds the collider in Acquire when the instance is active.
	collider := NewColliderNode(name+"_collider", shape, mask)
	collider.SetPosition(0, 0)
	proj.AddChildren(collider)

	sprite := NewSprite(name+"_sprite", texture, 0, true)
	sprite.SetScale(2, 2)
	proj.AddChildren(sprite)

	return proj
}

// reset sets position, velocity, and rotation for reuse from the pool.
func (p *Projectile) reset(spawnX, spawnY, vx, vy float64) {
	angle := 0.0
	if vx != 0 || vy != 0 {
		angle = atan2Pi2(vy, vx)
	}
	p.SetPosition(spawnX, spawnY)
	p.SetRotation(angle)
	p.Vx = vx
	p.Vy = vy
}

// atan2Pi2 returns math.Atan2(dy, dx) + math.Pi/2 (sprite orientation).
func atan2Pi2(dy, dx float64) float64 {
	return math.Atan2(dy, dx) + math.Pi/2
}

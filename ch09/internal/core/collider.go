package core

// Collider is a Node2D with a collision shape and mask. Create via CollisionManager.NewCollider.
type Collider struct {
	Node2D
	shape     *CollisionCircle
	mask      CollisionMask
	onCollide func(other *Collider) // called when collision with another Collider
}

// SetOnCollide sets the callback for when this collider hits another.
func (c *Collider) SetOnCollide(fn func(other *Collider)) {
	c.onCollide = fn
}

// GetShape returns the collision shape.
func (c *Collider) GetShape() *CollisionCircle { return c.shape }

// GetCollisionMask returns the collision mask.
func (c *Collider) GetCollisionMask() CollisionMask { return c.mask }

// CanCollideWith returns true if we should collide with the other participant.
func (c *Collider) CanCollideWith(other *Collider) bool {
	return c.mask.CanCollideWith(other.mask)
}

// NewColliderNode builds a collider without registering it (e.g. projectiles).
func NewColliderNode(name string, shape *CollisionCircle, mask CollisionMask) *Collider {
	return &Collider{Node2D: *NewNode2D(name), shape: shape, mask: mask}
}

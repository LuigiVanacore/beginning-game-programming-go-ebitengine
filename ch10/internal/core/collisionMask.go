package core

// Layer identifiers for collision masks (bitset-style).
const (
	LayerPlayer       = 1 << iota
	LayerEnemy
	LayerProjectile
	LayerPickup
	LayerPlayerPickup // dedicated wide pickup-radius collider on the Player
)

// CollisionMask defines which layers this collider is on and which it collides with.
type CollisionMask struct {
	identity    uint
	collidesWith uint
}

// NewCollisionMask creates a mask: identity is this collider's layer, collidesWith are layers to collide with.
func NewCollisionMask(identity uint, collidesWith uint) CollisionMask {
	return CollisionMask{identity: identity, collidesWith: collidesWith}
}

// CanCollideWith returns true if we should collide with the other's identity.
func (m CollisionMask) CanCollideWith(other CollisionMask) bool {
	return (m.collidesWith & other.identity) != 0
}

// GetIdentity returns this collider's layer.
func (m CollisionMask) GetIdentity() uint { return m.identity }

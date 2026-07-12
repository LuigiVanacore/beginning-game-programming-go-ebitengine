package core

// CollisionManager holds colliders and checks collisions each frame.
type CollisionManager struct {
	colliders []*Collider
}

// NewCollisionManager creates an empty manager.
func NewCollisionManager() *CollisionManager {
	return &CollisionManager{
		colliders: make([]*Collider, 0),
	}
}

// AddCollider registers a collider.
func (m *CollisionManager) AddCollider(c *Collider) {
	if c == nil {
		return
	}
	m.colliders = append(m.colliders, c)
}

// NewCollider builds a Collider with the given name, circle shape, and mask, and registers it.
func (m *CollisionManager) NewCollider(name string, shape *CollisionCircle, mask CollisionMask) *Collider {
	if m == nil {
		return nil
	}
	c := &Collider{
		Node2D: *NewNode2D(name),
		shape:  shape,
		mask:   mask,
	}
	m.AddCollider(c)
	return c
}

// RemoveCollider removes a collider from the manager.
func (m *CollisionManager) RemoveCollider(c *Collider) {
	for i, col := range m.colliders {
		if col == c {
			m.colliders[i] = m.colliders[len(m.colliders)-1]
			m.colliders = m.colliders[:len(m.colliders)-1]
			return
		}
	}
}

// circleSample is a collider's world-space circle: centre (x, y) and radius r.
type circleSample struct {
	x, y, r float64
}

// sampleCircle reads world centre and radius from c.
func sampleCircle(c *Collider) circleSample {
	pos := c.GetWorldPosition()
	return circleSample{
		x: pos.X(),
		y: pos.Y(),
		r: c.shape.GetRadius(),
	}
}

// pairWantsCollision reports whether both layer masks allow a and b to interact.
func pairWantsCollision(a, b *Collider) bool {
	return a.CanCollideWith(b) && b.CanCollideWith(a)
}

// samplesOverlap returns true when two circle samples penetrate (strict < sum radii).
func samplesOverlap(a, b circleSample) bool {
	return circlesOverlap(a.x, a.y, a.r, b.x, b.y, b.r)
}

// notifyPair calls onCollide on each collider that registered a callback.
func notifyPair(a, b *Collider) {
	if a.onCollide != nil {
		a.onCollide(b)
	}
	if b.onCollide != nil {
		b.onCollide(a)
	}
}

// checkColliderPair runs mask filter, circle test, and callbacks for one unordered pair.
// sampleA is precomputed for a so the outer loop does not re-read a's transform per j.
func checkColliderPair(a *Collider, sampleA circleSample, b *Collider) {
	if b == nil {
		return
	}
	if !pairWantsCollision(a, b) {
		return
	}
	if !samplesOverlap(sampleA, sampleCircle(b)) {
		return
	}
	notifyPair(a, b)
}

// CheckCollision tests all pairs and calls OnCollide when overlapping.
func (m *CollisionManager) CheckCollision() {
	for i := 0; i < len(m.colliders); i++ {
		a := m.colliders[i]
		if a == nil {
			continue
		}
		sampleA := sampleCircle(a)
		for j := i + 1; j < len(m.colliders); j++ {
			checkColliderPair(a, sampleA, m.colliders[j])
		}
	}
}

// CombineIDs produces a deterministic key for a pair (order-independent).
func combineCollisionIDs(id1, id2 uint64) uint64 {
	if id1 < id2 {
		return (id1 << 32) | id2
	}
	return (id2 << 32) | id1
}

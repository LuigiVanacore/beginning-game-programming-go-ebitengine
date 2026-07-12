package core

// CollisionCircle is a circle shape for collision (center at origin, radius).
type CollisionCircle struct {
	radius float64
}

// NewCollisionCircle creates a circle with the given radius (center is at node position).
func NewCollisionCircle(radius float64) *CollisionCircle {
	return &CollisionCircle{radius: radius}
}

// GetRadius returns the circle radius.
func (c *CollisionCircle) GetRadius() float64 { return c.radius }

// circlesOverlap returns true if two circles overlap (distance < r1 + r2).
func circlesOverlap(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	distSq := dx*dx + dy*dy
	sum := r1 + r2
	return distSq < sum*sum
}

package core

// Transformable is implemented by nodes with spatial data (e.g. Node2D).
// The engine uses it to compute world position for drawing.
type Transformable interface {
	GetTransform() Transform
	SetTransform(t Transform)
	GetWorldTransform() Transform
}

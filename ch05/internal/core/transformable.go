package core

type Transformable interface {
	GetTransform() Transform
	SetTransform(t Transform)
	GetWorldTransform() Transform
}

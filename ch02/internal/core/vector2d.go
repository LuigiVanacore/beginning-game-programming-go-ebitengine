package core

import "math"

// Vector2D represents a 2D vector with X and Y components.
type Vector2D struct {
	x, y float64
}

func NewVector2D(x, y float64) Vector2D       { return Vector2D{x: x, y: y} }
func ZeroVector2D() Vector2D                  { return Vector2D{0, 0} }
func (v Vector2D) X() float64                  { return v.x }
func (v Vector2D) Y() float64                  { return v.y }
func (v *Vector2D) SetX(x float64)            { v.x = x }
func (v *Vector2D) SetY(y float64)             { v.y = y }
func (v *Vector2D) SetPosition(other Vector2D) { v.x, v.y = other.x, other.y }

func (v Vector2D) RotateVector(radians float64) Vector2D {
	s, c := math.Sin(radians), math.Cos(radians)
	return Vector2D{v.x*c - v.y*s, v.x*s + v.y*c}
}

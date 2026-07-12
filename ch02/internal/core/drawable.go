package core

import "github.com/hajimehoshi/ebiten/v2"

// Drawable is implemented by nodes that can be drawn.
// They must have a transform (Transformable) and provide GetLayer and Draw.
type Drawable interface {
	Transformable
	GetLayer() int
	Draw(target *ebiten.Image, op *ebiten.DrawImageOptions)
}

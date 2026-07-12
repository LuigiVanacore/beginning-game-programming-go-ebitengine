package core

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Transformable
	GetLayer() int
	Draw(target *ebiten.Image, op *ebiten.DrawImageOptions)
}

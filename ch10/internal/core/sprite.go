package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Node2D
	texture *ebiten.Image
	layer   int
	visible bool
}

func NewSprite(name string, texture *ebiten.Image, layer int) *Sprite {
	s := &Sprite{
		Node2D:  *NewNode2D(name),
		texture: texture,
		layer:   layer,
		visible: true,
	}
	if texture != nil {
		s.SetPivotToCenter()
	}
	return s
}

func (s *Sprite) GetTexture() *ebiten.Image    { return s.texture }
func (s *Sprite) SetTexture(tex *ebiten.Image) { s.texture = tex }
func (s *Sprite) GetLayer() int                { return s.layer }
func (s *Sprite) SetLayer(l int)               { s.layer = l }
func (s *Sprite) GetVisible() bool             { return s.visible }
func (s *Sprite) SetVisible(v bool)             { s.visible = v }

func (s *Sprite) SetPivotToCenter() {
	if s.texture == nil {
		return
	}
	w := float64(s.texture.Bounds().Dx())
	h := float64(s.texture.Bounds().Dy())
	s.SetPivot(w/2, h/2)
}

func (s *Sprite) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	if s.texture == nil || !s.visible {
		return
	}
	target.DrawImage(s.texture, op)
}

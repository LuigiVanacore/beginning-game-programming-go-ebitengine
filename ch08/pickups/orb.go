package pickups

import (
	. "book/code/ch08/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	OrbRadius  = 6.0
	OrbXPValue = 10
)

var maskPickup = NewCollisionMask(LayerPickup, LayerPlayer|LayerPlayerPickup)

// Orb is an XP orb: root collider in the world with an attached sprite.
type Orb struct {
	Col   *Collider
	spr   *Sprite
	Value int
}

func newOrbImage() *ebiten.Image {
	img := ebiten.NewImage(12, 12)
	orbColor := color.RGBA{64, 140, 255, 255}
	ApplyDrawShape(img, ShapeFilledCircle(6, 6, 5), orbColor)
	return img
}

// CreateOrb spawns an XP orb at the given world position.
func CreateOrb(engine *Engine, x, y float64, value int) *Orb {
	world := engine.World()
	colMgr := engine.CollisionManager()
	shape := NewCollisionCircle(OrbRadius)
	col := colMgr.NewCollider("orb", shape, maskPickup)
	col.SetPosition(x, y)
	spr := NewSprite("orb_sprite", newOrbImage(), 0, true)
	spr.SetScale(1.5, 1.5)
	col.AddChildren(spr)
	world.AddNodeToLayer(col, DrawLayerPlayer)
	return &Orb{Col: col, spr: spr, Value: value}
}

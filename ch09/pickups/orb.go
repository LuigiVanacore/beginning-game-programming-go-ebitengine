package pickups

import (
	. "book/code/ch09/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	OrbRadius  = 6.0
	OrbXPValue = 10
)

var maskPickup = NewCollisionMask(LayerPickup, LayerPlayer|LayerPlayerPickup)

// OrbEnt is an XP orb node: collider + sprite under a root Node2D.
type OrbEnt struct {
	Node2D
	Col   *Collider
	Value int
}

func newOrbImage() *ebiten.Image {
	img := ebiten.NewImage(12, 12)
	orbColor := color.RGBA{64, 140, 255, 255}
	ApplyDrawShape(img, ShapeFilledCircle(6, 6, 5), orbColor)
	return img
}

// CreateOrb spawns an XP orb at the given world position.
func CreateOrb(engine *Engine, x, y float64, value int) *OrbEnt {
	world := engine.World()
	colMgr := engine.CollisionManager()
	o := &OrbEnt{Node2D: *NewNode2D("orb"), Value: value}
	o.SetPosition(x, y)
	shape := NewCollisionCircle(OrbRadius)
	col := colMgr.NewCollider("orb_col", shape, maskPickup)
	spr := NewSprite("orb_sprite", newOrbImage(), 0, true)
	spr.SetScale(1.5, 1.5)
	col.AddChildren(spr)
	o.AddChildren(col)
	o.Col = col
	world.AddNodeToLayer(o, DrawLayerPlayer)
	return o
}

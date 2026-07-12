package pickups

import (
	. "book/code/ch10/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	OrbRadius  = 6.0
	OrbXPValue = 10
)

// Orbs and potions collide with the player: the wide pickup collider
// (LayerPlayerPickup) grabs them early, the body collider (LayerPlayer) is the
// fallback. Collision matching is mutual, so both layers must appear in
// collidesWith for pickups to register.
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
	vector.DrawFilledCircle(img, 6, 6, 5, orbColor, true)
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
	spr := NewSprite("orb_sprite", newOrbImage(), 0)
	spr.SetPivotToCenter()
	spr.SetScale(1.0, 1.0)
	col.AddChildren(spr)
	o.AddChildren(col)
	o.Col = col
	world.AddNodeToLayer(o, DrawLayerPlayer)
	return o
}

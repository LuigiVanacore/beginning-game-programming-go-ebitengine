package pickups

import (
	. "book/code/ch12/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	potionRadius     = 10.0
	potionDropChance = 0.05
)

// PotionEnt is a health potion pickup node.
type PotionEnt struct {
	Node2D
	Col *Collider
}

// CreatePotion spawns a health potion at the given world position.
func CreatePotion(engine *Engine, x, y float64, tex *ebiten.Image) *PotionEnt {
	world := engine.World()
	colMgr := engine.CollisionManager()
	p := &PotionEnt{Node2D: *NewNode2D("potion")}
	p.SetPosition(x, y)
	shape := NewCollisionCircle(potionRadius)
	col := colMgr.NewCollider("potion_col", shape, maskPickup)
	spr := NewSprite("potion_sprite", tex, 0)
	spr.SetPivotToCenter()
	spr.SetScale(2, 2)
	col.AddChildren(spr)
	p.AddChildren(col)
	p.Col = col
	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

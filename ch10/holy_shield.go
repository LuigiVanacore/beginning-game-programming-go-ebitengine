package game

import (
	. "book/code/ch10/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const holyShieldBorder = 3.0 // white border stroke width in pixels

// HolyShieldRadius is the Holy Shield ring radius, kept smaller than the Sacred Book
// orbit (OrbitWeaponDistance) so the shield hugs the player.
const HolyShieldRadius = 30.0

func newHolyShieldCircleImage(radius float64) *ebiten.Image {
	size := int(radius*2 + holyShieldBorder*2)
	img := ebiten.NewImage(size, size)
	cx := float32(size) / 2
	cy := float32(size) / 2
	vector.StrokeCircle(img, cx, cy, float32(radius), holyShieldBorder, color.White, true)
	return img
}

// HolyShieldWeapon is a Node2D under player.WeaponsRoot; ring collider damages enemies on contact.
type HolyShieldWeapon struct {
	Node2D
	col *Collider
}

// NewHolyShieldWeapon builds the Holy Shield; caller adds it to the scene via WeaponManager.Mount.
func NewHolyShieldWeapon(engine *Engine, player *Player, radius float64) *HolyShieldWeapon {
	w := &HolyShieldWeapon{
		Node2D: *NewNode2D("holy_shield_weapon"),
	}
	colMgr := engine.CollisionManager()
	shape := NewCollisionCircle(radius)
	mask := NewCollisionMask(LayerProjectile, LayerEnemy)
	col := colMgr.NewCollider("holy_shield", shape, mask)
	col.SetPosition(0, 0)
	sprite := NewSprite("holy_shield_sprite", newHolyShieldCircleImage(radius), 0)
	col.AddChildren(sprite)
	w.AddChildren(col)
	w.col = col

	col.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() == LayerEnemy {
			player.QueueWeaponHit(nil, other)
		}
	})
	return w
}

// UpdateWeapon is a no-op; the shield is static relative to the player.
func (_ *HolyShieldWeapon) UpdateWeapon(_ *Engine, _ *Player, _ *Cursor, _ *WeaponManager) {
}

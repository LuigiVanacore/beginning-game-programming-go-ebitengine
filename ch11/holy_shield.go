package game

import (
	. "book/code/ch11/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const holyShieldBorder = 3.0 // white border stroke width in pixels

const holyShieldDamage = 0.5 // HP removed from an enemy per contact frame

// HolyShieldRadius is the Holy Shield base ring radius, kept smaller than the Sacred Book
// orbit (OrbitWeaponDistance) so the shield hugs the player. Upgrades scale it.
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
// Col is exported so weaponColliders() can hand it to the WeaponManager on mount/unmount (new in ch11).
type HolyShieldWeapon struct {
	Node2D
	Col *Collider // exposed to WeaponManager via weaponColliders()
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
	// NewCollider registers the collider immediately. Deactivate it until the weapon
	// is mounted (WeaponManager.Mount re-adds it); otherwise an unlocked-but-unmounted
	// shield would damage enemies from the world origin.
	colMgr.RemoveCollider(col)
	col.SetPosition(0, 0)
	sprite := NewSprite("holy_shield_sprite", newHolyShieldCircleImage(radius), 0)
	col.AddChildren(sprite)
	w.AddChildren(col)
	w.Col = col

	col.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() == LayerEnemy {
			player.QueueWeaponHit(nil, other, holyShieldDamage)
		}
	})
	return w
}

func (_ *HolyShieldWeapon) updateWeapon(_ *Engine, _ *Player, _ *Cursor) {
	// Static relative to player; collision-only weapon.
}

// weaponColliders exposes the ring collider so WeaponManager can register it on
// mount and remove it on unmount.
func (w *HolyShieldWeapon) weaponColliders() []*Collider {
	if w == nil || w.Col == nil {
		return nil
	}
	return []*Collider{w.Col}
}

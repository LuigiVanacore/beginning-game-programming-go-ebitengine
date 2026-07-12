package game

import (
	. "book/code/ch11/internal/core"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	sacredBookAnglePerFrame = 1.0 / 30.0 // orbit angle (radians) advanced per frame
	sacredBookColliderRadius   = 3.0
	sacredBookDamage           = 0.6 // HP removed from an enemy per contact frame
	// OrbitWeaponDistance is the Sacred Book orbit radius (the Holy Shield uses its own, smaller HolyShieldRadius).
	OrbitWeaponDistance = 48.0
)

// SacredBook is a Node2D under player.WeaponsRoot; child collider orbits the player.
// SpeedMult is this weapon's own orbit-speed multiplier, scaled by upgrades.
type SacredBook struct {
	Node2D
	col         *Collider
	angle       float64
	localRadius float64
	SpeedMult   float64
}

// NewSacredBook builds the Sacred Book with its multiplier at 1.0; caller adds it via WeaponManager.Mount.
func NewSacredBook(engine *Engine, player *Player, orbitRadius float64, tex *ebiten.Image) *SacredBook {
	w := &SacredBook{
		Node2D:      *NewNode2D("sacred_book"),
		angle:       0,
		localRadius: orbitRadius,
		SpeedMult:   1.0,
	}
	colMgr := engine.CollisionManager()
	shape := NewCollisionCircle(sacredBookColliderRadius)
	mask := NewCollisionMask(LayerProjectile, LayerEnemy)
	col := colMgr.NewCollider("sacred_book", shape, mask)
	// Deactivate until the weapon is mounted (WeaponManager.Mount re-adds it), so an
	// unlocked-but-unmounted book cannot damage enemies from its resting position.
	colMgr.RemoveCollider(col)
	col.SetPosition(orbitRadius, 0)
	spr := NewSprite("sacred_book_sprite", tex, 0)
	spr.SetScale(0.5, 0.5)
	col.AddChildren(spr)
	w.AddChildren(col)
	w.col = col

	col.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() == LayerEnemy {
			player.QueueWeaponHit(nil, other, sacredBookDamage)
		}
	})
	return w
}

// weaponColliders exposes the orbiting collider so WeaponManager can register it on
// mount and remove it on unmount.
func (w *SacredBook) weaponColliders() []*Collider {
	if w == nil || w.col == nil {
		return nil
	}
	return []*Collider{w.col}
}

func (w *SacredBook) updateWeapon(_ *Engine, _ *Player, _ *Cursor) {
	if w == nil {
		return
	}
	w.angle += sacredBookAnglePerFrame * w.SpeedMult
	x := w.localRadius * math.Cos(w.angle)
	y := w.localRadius * math.Sin(w.angle)
	w.col.SetPosition(x, y)
}

package game

import (
	. "book/code/ch10/internal/core"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	sacredBookAnglePerFrame  = 1.0 / 30.0 // orbit angle (radians) advanced per frame
	sacredBookColliderRadius = 3.0
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

// NewSacredBook builds the Sacred Book with its multiplier at the neutral value (1.0).
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
	col.SetPosition(orbitRadius, 0)
	spr := NewSprite("sacred_book_sprite", tex, 0)
	spr.SetScale(0.5, 0.5)
	col.AddChildren(spr)
	w.AddChildren(col)
	w.col = col

	col.SetOnCollide(func(other *Collider) {
		if other.GetCollisionMask().GetIdentity() == LayerEnemy {
			player.QueueWeaponHit(nil, other)
		}
	})
	return w
}

// UpdateWeapon advances the orbit using the book's own SpeedMult.
func (w *SacredBook) UpdateWeapon(_ *Engine, _ *Player, _ *Cursor, _ *WeaponManager) {
	if w == nil {
		return
	}
	w.angle += sacredBookAnglePerFrame * w.SpeedMult
	x := w.localRadius * math.Cos(w.angle)
	y := w.localRadius * math.Sin(w.angle)
	w.col.SetPosition(x, y)
}

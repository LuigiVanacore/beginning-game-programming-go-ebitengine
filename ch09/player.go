package game

import (
	. "book/code/ch09/internal/core"

	"math"
)
const (
	playerSpeed        = 4.0       // pixels per Update step
	playerRadius       = 14.0      // body vs enemies (matches scaled sprite roughly)
	playerPickupRadius = 20.0      // wider circle so orbs are collected before overlapping the sprite
	PlayerMaxHP        = 100.0
	enemyDamagePerFrame = 0.2 // HP lost per Update step while touching an enemy
)

// Player: embeds Node2D; sprite, colliders, and weapons root are direct children (Ch9).
type Player struct {
	Node2D
	engine         *Engine
	Collider       *Collider
	PickupCollider *Collider
	Sprite         *Sprite
	WeaponsRoot    *Node2D
	cursor         *Cursor
	weaponHit      func(a, b *Collider)
	HP             float64
	XP             int
	Level          int
}

// NewPlayer creates the player; the world and resources come from the engine.
func NewPlayer(engine *Engine) *Player {
	rm := engine.ResourceManager()
	world := engine.World()
	colMgr := engine.CollisionManager()

	p := &Player{
		Node2D: *NewNode2D(NamePlayer),
		engine: engine,
	}
	p.SetPosition(0, 0)

	playerTex, _ := rm.GetTexture(PlayerTexture)
	sprite := NewSprite(NamePlayerSprite, playerTex, 0, true)
	sprite.SetPivotToCenter()
	sprite.SetScale(1, 1)
	p.AddChildren(sprite)

	bodyShape := NewCollisionCircle(playerRadius)
	bodyMask := NewCollisionMask(LayerPlayer, LayerEnemy)
	c := colMgr.NewCollider(NamePlayer, bodyShape, bodyMask)
	c.SetPosition(0, 0)
	p.AddChildren(c)

	pickupShape := NewCollisionCircle(playerPickupRadius)
	pickupMask := NewCollisionMask(LayerPlayerPickup, LayerPickup)
	pickup := colMgr.NewCollider(NamePlayerPickup, pickupShape, pickupMask)
	pickup.SetPosition(0, 0)
	p.AddChildren(pickup)

	weaponsRoot := NewNode2D(NameWeaponsRoot)
	p.AddChildren(weaponsRoot)

	p.Sprite = sprite
	p.Collider = c
	p.PickupCollider = pickup
	p.WeaponsRoot = weaponsRoot
	p.HP = PlayerMaxHP
	p.XP = 0
	p.Level = 1

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// SetCursor links the aiming cursor (called from createSession).
func (p *Player) SetCursor(c *Cursor) { p.cursor = c }

// SetWeaponHit sets how weapon hits append to the combat removal queue (wired from Game).
func (p *Player) SetWeaponHit(fn func(a, b *Collider)) { p.weaponHit = fn }

// MountWeapon adds a weapon Node2D under WeaponsRoot. The node must implement PlayerWeapon.
func (p *Player) MountWeapon(w SceneNode) {
	if p == nil || p.WeaponsRoot == nil || w == nil {
		return
	}
	if _, ok := w.(PlayerWeapon); !ok {
		return
	}
	p.WeaponsRoot.AddChildren(w)
}

// QueueWeaponHit is called by weapons to enqueue projectile/enemy colliders for removal processing.
func (p *Player) QueueWeaponHit(a, b *Collider) {
	if p.weaponHit == nil {
		return
	}
	p.weaponHit(a, b)
}

// Update applies input, movement, and the mounted weapons.
func (p *Player) Update() {
	if p == nil || p.engine == nil {
		return
	}
	p.updateMovement(p.engine.Input())
	p.updateWeapons()
}

func (p *Player) updateWeapons() {
	if p.WeaponsRoot == nil || p.cursor == nil {
		return
	}
	eng := p.engine
	for _, ch := range p.WeaponsRoot.GetChildren() {
		if w, ok := ch.(PlayerWeapon); ok {
			w.UpdateWeapon(eng, p, p.cursor)
		}
	}
}

// TryReleaseProjectileByCollider returns a pooled projectile to the weapon that owns it, if any.
func (p *Player) TryReleaseProjectileByCollider(c *Collider) bool {
	if p == nil || p.WeaponsRoot == nil {
		return false
	}
	for _, ch := range p.WeaponsRoot.GetChildren() {
		if pc, ok := ch.(ProjectileCarrier); ok && pc.TryReleaseProjectile(c) {
			return true
		}
	}
	return false
}

func (p *Player) updateMovement(inp *InputManager) {
	dx, dy := 0.0, 0.0
	if inp.IsActionPressed(ActionMoveUp) {
		dy -= 1
	}
	if inp.IsActionPressed(ActionMoveDown) {
		dy += 1
	}
	if inp.IsActionPressed(ActionMoveLeft) {
		dx -= 1
	}
	if inp.IsActionPressed(ActionMoveRight) {
		dx += 1
	}
	if dx == 0 && dy == 0 {
		return
	}
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx /= length
		dy /= length
	}
	x, y := p.GetPosition().X(), p.GetPosition().Y()
	p.SetPosition(x+dx*playerSpeed, y+dy*playerSpeed)
}

// DamageFromEnemyContact applies per-frame damage while overlapping an enemy (Ch8).
func (p *Player) DamageFromEnemyContact() {
	p.HP -= enemyDamagePerFrame
	if p.HP < 0 {
		p.HP = 0
	}
}

// IsDead returns true when HP has reached zero.
func (p *Player) IsDead() bool {
	return p.HP <= 0
}

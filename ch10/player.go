package game


import (
	. "book/code/ch10/internal/core"

	"math"
)

const (
	playerSpeedPxPerFrame = 3.3333333333333335 // horizontal move per Update step
	playerRadius          = 14.0               // collision circle radius for the player body
	playerPickupRadius    = 20.0               // wider circle so orbs are collected before overlapping the sprite
	playerMaxHP           = 100.0
	enemyDamagePerFrame   = 0.2 // HP lost per Update step while touching an enemy
)

// Player embeds Node2D; sprite, colliders, and weapons root are direct children.
// MaxHP is now a mutable field so upgrades can scale it (new in ch10).
type Player struct {
	Node2D
	engine         *Engine
	Collider       *Collider
	PickupCollider *Collider
	Sprite         *Sprite
	WeaponsRoot    *Node2D
	cursor         *Cursor
	weaponHit      func(a, b *Collider)
	HP          float64
	MaxHP       float64 // mutable max HP; upgrades may increase this (ch10)
	XP          int
	Level       int
	SpeedMult   float64 // movement speed multiplier; increased by Speed Boots upgrade
	XPBonusMult float64 // XP gain multiplier; increased by Experience Gem upgrade
}

// NewPlayer creates the player; adds it to the world.
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
	sprite := NewSprite(NamePlayerSprite, playerTex, 0)
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
	p.MaxHP = playerMaxHP
	p.HP = playerMaxHP
	p.XP = 0
	p.Level = 1
	p.SpeedMult = 1.0
	p.XPBonusMult = 1.0

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// SetCursor links the aiming cursor (called from createSession).
func (p *Player) SetCursor(c *Cursor) { p.cursor = c }

// SetWeaponHit sets how weapon hits append to the combat removal queue (wired from Game).
func (p *Player) SetWeaponHit(fn func(a, b *Collider)) { p.weaponHit = fn }

// QueueWeaponHit is called by weapons to enqueue projectile/enemy colliders for removal.
func (p *Player) QueueWeaponHit(a, b *Collider) {
	if p.weaponHit == nil {
		return
	}
	p.weaponHit(a, b)
}

// Update applies input and movement. Weapon updates are driven by WeaponManager (ch10).
func (p *Player) Update() {
	if p == nil || p.engine == nil {
		return
	}
	p.updateMovement(p.engine.Input())
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
	p.SetPosition(x+dx*playerSpeedPxPerFrame*p.SpeedMult, y+dy*playerSpeedPxPerFrame*p.SpeedMult)
}

// DamageFromEnemyContact applies damage while overlapping an enemy.
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

package game


import (
	. "book/code/ch12/internal/core"

	"math"
)

const (
	playerSpeedPxPerFrame = 3.3333333333333335 // horizontal move per Update step
	playerRadius          = 14.0               // collision circle radius for the player body
	playerPickupRadius    = 20.0               // wider pickup-radius collider
	playerMaxHP           = 100.0
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
	weaponHit      func(a, b *Collider, dmg float64)
	HP             float64
	MaxHP          float64  // mutable max HP; upgrades may increase this (ch10)
	XP             int
	Level          int
	speedMult      *float64 // pointer to Game.playerSpeedMult; set via SetSpeedMult

	// One-time bonus items from the level-up panel (ch11).
	HasArmor, HasBoots, HasGem, HasSkull, HasRing bool
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

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// SetCursor links the aiming cursor (called from createSession).
func (p *Player) SetCursor(c *Cursor) { p.cursor = c }

// SetWeaponHit sets how weapon hits are resolved: the projectile collider a, the
// enemy collider b, and the damage the weapon deals (wired from Game).
func (p *Player) SetWeaponHit(fn func(a, b *Collider, dmg float64)) { p.weaponHit = fn }

// SetSpeedMult wires the player movement multiplier to a Game-owned float (new in ch10).
func (p *Player) SetSpeedMult(mult *float64) { p.speedMult = mult }

// QueueWeaponHit is called by weapons to resolve a hit: release the projectile a
// and apply dmg to the enemy behind collider b.
func (p *Player) QueueWeaponHit(a, b *Collider, dmg float64) {
	if p.weaponHit == nil {
		return
	}
	p.weaponHit(a, b, dmg)
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
	mult := 1.0
	if p.speedMult != nil {
		mult = *p.speedMult
	}
	x, y := p.GetPosition().X(), p.GetPosition().Y()
	p.SetPosition(x+dx*playerSpeedPxPerFrame*mult, y+dy*playerSpeedPxPerFrame*mult)
}

// DamageFromEnemyContact applies an enemy's contact damage while overlapping it.
// Different enemy kinds hit for different amounts, and that amount grows with the
// difficulty tier, so the caller passes the touching enemy's current damage.
func (p *Player) DamageFromEnemyContact(amount float64) {
	p.HP -= amount
	if p.HP < 0 {
		p.HP = 0
	}
}

// IsDead returns true when HP has reached zero.
func (p *Player) IsDead() bool {
	return p.HP <= 0
}

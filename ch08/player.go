package game


import (
	. "book/code/ch08/internal/core"

	"math"
)

const (
	playerSpeed        = 4.0        // pixels per Update step
	playerRadius       = 14.0       // body vs enemies (matches scaled sprite roughly)
	playerPickupRadius = 20.0       // wider circle so orbs are collected before overlapping the sprite
	// PlayerMaxHP is exported for HUD (package ui) and game logic.
	PlayerMaxHP        = 100.0
	enemyDamagePerFrame = 0.2 // HP lost per Update step while touching an enemy
)

// Player: embeds Node2D; sprite, body collider, and pickup collider are direct children.
type Player struct {
	Node2D
	engine         *Engine
	Collider       *Collider
	PickupCollider *Collider
	Sprite         *Sprite
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

	p.Sprite = sprite
	p.Collider = c
	p.PickupCollider = pickup
	p.HP = PlayerMaxHP
	p.XP = 0
	p.Level = 1

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// Update applies input and movement in world space.
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

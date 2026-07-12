package game

import (
	. "book/code/ch06/internal/core"
	"math"
)

const (
	playerSpeed  = 4.0  // pixels per Update step
	playerRadius = 14.0 // collision circle radius for the player body
)

// Player: embeds Node2D; sprite and collider are direct children (same idea as Enemy).
type Player struct {
	Node2D
	engine   *Engine
	Collider *Collider
	Sprite   *Sprite
}

// NewPlayer creates the player and registers its hierarchy with the world and CollisionManager.
func NewPlayer(engine *Engine) *Player {
	rm := engine.ResourceManager()
	world := engine.World()
	colMgr := engine.CollisionManager()

	p := &Player{
		Node2D: *NewNode2D("player"),
		engine: engine,
	}
	p.SetPosition(0, 0)

	playerTex, _ := rm.GetTexture("player")
	sprite := NewSprite("player_sprite", playerTex, 0, true)
	sprite.SetPivotToCenter()
	sprite.SetScale(1, 1)
	p.AddChildren(sprite)

	shape := NewCollisionCircle(playerRadius)
	mask := NewCollisionMask(LayerPlayer, LayerEnemy)
	c := colMgr.NewCollider("player", shape, mask)
	c.SetPosition(0, 0)
	p.AddChildren(c)

	p.Sprite = sprite
	p.Collider = c

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// Update applies arrow-key input and movement to the root Node2D.
func (p *Player) Update() {
	if p == nil || p.engine == nil {
		return
	}
	p.updateMovement(p.engine.Input())
}

func (p *Player) updateMovement(inp *InputManager) {
	dx, dy := 0.0, 0.0
	if inp.IsActionPressed("move_up") {
		dy -= 1
	}
	if inp.IsActionPressed("move_down") {
		dy += 1
	}
	if inp.IsActionPressed("move_left") {
		dx -= 1
	}
	if inp.IsActionPressed("move_right") {
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

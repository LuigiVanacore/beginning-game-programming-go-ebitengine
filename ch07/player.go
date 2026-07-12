package game


import (
	. "book/code/ch07/internal/core"

	"math"
)

const (
	playerSpeed  = 4.0  // pixels per Update step
	playerRadius = 14.0 // collision circle radius for the player body
)

// Player: embeds Node2D; sprite and body collider are direct children.
type Player struct {
	Node2D
	engine   *Engine
	Collider *Collider
	Sprite   *Sprite
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

	shape := NewCollisionCircle(playerRadius)
	mask := NewCollisionMask(LayerPlayer, LayerEnemy)
	c := colMgr.NewCollider(NamePlayer, shape, mask)
	c.SetPosition(0, 0)
	p.AddChildren(c)

	p.Sprite = sprite
	p.Collider = c

	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// Update applies the configured actions (settings) and movement in world space.
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

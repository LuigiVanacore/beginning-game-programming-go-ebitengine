package game

import (
	. "book/code/ch05/internal/core"
	"math"
)

const playerSpeed = 4.0 // pixels per Update step

// Player is the playable character: embeds Node2D; sprite is a direct child (world space; camera follows).
type Player struct {
	Node2D
}

// NewPlayer builds the player and registers it on LayerPlayer.
func NewPlayer(world *World, rm *ResourceManager) *Player {
	p := &Player{
		Node2D: *NewNode2D("player"),
	}
	p.SetPosition(0, 0)
	playerTex, _ := rm.GetTexture("player")
	sprite := NewSprite("player_sprite", playerTex, 0, true)
	sprite.SetPivotToCenter()
	sprite.SetScale(1, 1)
	p.AddChildren(sprite)
	world.AddNodeToLayer(p, DrawLayerPlayer)
	return p
}

// UpdateMovement reads arrow-key actions and moves the player in world space.
func (p *Player) UpdateMovement(inp *InputManager) {
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

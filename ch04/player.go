package game

import (
	"math"

	. "book/code/ch04/internal/core"
)

const playerSpeed = 4.0 // pixels per Update step

const (
	playerNodeName        = "player"
	playerSpriteName      = "player_sprite"
	playerSpriteLayer     = 0
	playerWorldLayerIndex = 1
)

// Player is the playable character: embeds Node2D; sprite is a direct child.
type Player struct {
	Node2D
}

// NewPlayer builds the player and registers it on layer 1.
func NewPlayer(world *World, rm *ResourceManager, centerX, centerY float64) *Player {
	p := &Player{
		Node2D: *NewNode2D(playerNodeName),
	}
	p.SetPosition(centerX, centerY)
	playerTex, _ := rm.GetTexture(playerTextureKey)
	sprite := NewSprite(playerSpriteName, playerTex, playerSpriteLayer, true)
	sprite.SetPosition(0, 0)
	sprite.SetScale(1, 1)
	p.AddChildren(sprite)
	world.AddNodeToLayer(p, playerWorldLayerIndex)
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

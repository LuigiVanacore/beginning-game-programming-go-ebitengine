package game

import . "book/code/ch03/internal/core"

const (
	playerNodeName        = "player"
	playerSpriteName      = "player_sprite"
	playerSpriteLayer     = 0
	playerWorldLayerIndex = 1
)

// Player is the playable character: embeds Node2D; sprite is a direct child (Ch3 — no collision yet).
type Player struct {
	Node2D
}

// NewPlayer builds the player node, attaches the sprite, and adds it to the world on layer 1.
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

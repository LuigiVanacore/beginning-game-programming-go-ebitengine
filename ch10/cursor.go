package game

import (
	. "book/code/ch10/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

// Cursor is a Node2D that follows the mouse in world space for aiming.
// It holds a reference to the camera for coordinate conversion.
type Cursor struct {
	Node2D
	camera *Camera
}

// NewCursor creates a Cursor with a Sprite child for aiming.
// The sprite is scaled. Add the returned cursor to the world with AddNodeToLayer.
func NewCursor(name string, texture *ebiten.Image, scale float64, camera *Camera) *Cursor {
	c := &Cursor{
		Node2D: *NewNode2D(name),
		camera: camera,
	}
	if texture != nil {
		sprite := NewSprite(name+"_sprite", texture, 0)
		sprite.SetScale(scale, scale)
		c.AddChildren(sprite)
	}
	return c
}

// Update sets the cursor position in world space to match the mouse
// position on screen. Call this each frame in Game.Update.
func (c *Cursor) Update() {
	mx, my := ebiten.CursorPosition()
	wx, wy := c.camera.GetWorldCoords(mx, my)
	c.SetPosition(wx, wy)
}

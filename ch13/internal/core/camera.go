package core

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera provides a viewport into the world and can follow a Node2D.
type Camera struct {
	Node2D
	width        uint
	height       uint
	surface      *ebiten.Image
	nodeToFollow Transformable
	shakeMag     float64 // current screen-shake magnitude in pixels; decays each frame (ch13)
}

// NewCamera creates a camera with the given dimensions.
func NewCamera(w, h uint) *Camera {
	c := &Camera{
		Node2D: *NewNode2D(NameCamera),
		width:  w,
		height: h,
	}
	c.surface = ebiten.NewImage(int(w), int(h))
	return c
}

// GetSurface returns the offscreen image the scene is drawn to.
func (c *Camera) GetSurface() *ebiten.Image {
	return c.surface
}

// GetWidth returns the camera viewport width in pixels.
func (c *Camera) GetWidth() uint { return c.width }

// GetHeight returns the camera viewport height in pixels.
func (c *Camera) GetHeight() uint { return c.height }

// SetFollow sets the node to follow. Pass nil to disable.
func (c *Camera) SetFollow(node Transformable) {
	c.nodeToFollow = node
}

// Update updates camera position to follow the target node (center on screen).
func (c *Camera) Update() {
	if c.nodeToFollow == nil {
		return
	}
	wt := c.nodeToFollow.GetWorldTransform()
	px := wt.GetPosition().X()
	py := wt.GetPosition().Y()
	// Center the target in the view: camera top-left = target center - half screen
	c.SetPosition(px-float64(c.width)/2, py-float64(c.height)/2)

	// Apply a decaying screen shake on top of the follow position (ch13). Because
	// the world and particles are both drawn relative to the camera position, they
	// shake together while the screen-space HUD stays still.
	if c.shakeMag > 0.4 {
		ox := (rand.Float64()*2 - 1) * c.shakeMag
		oy := (rand.Float64()*2 - 1) * c.shakeMag
		c.SetPosition(c.GetPosition().X()+ox, c.GetPosition().Y()+oy)
		c.shakeMag *= 0.85
	} else {
		c.shakeMag = 0
	}
}

// Shake starts (or strengthens) a brief screen shake. magnitude is in pixels and
// decays to zero over a few frames. Call it on impactful events such as an enemy
// death or the player taking damage.
func (c *Camera) Shake(magnitude float64) {
	if magnitude > c.shakeMag {
		c.shakeMag = magnitude
	}
}

// ApplyOffset modifies op so world coords are drawn relative to camera position.
func (c *Camera) ApplyOffset(op *ebiten.DrawImageOptions) {
	px := c.GetPosition().X()
	py := c.GetPosition().Y()
	op.GeoM.Translate(-px, -py)
}

// DrawToScreen draws the camera surface to the target (screen).
func (c *Camera) DrawToScreen(target *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	target.DrawImage(c.surface, op)
}

// GetWorldCoords converts screen coordinates (e.g. mouse) to world coordinates.
func (c *Camera) GetWorldCoords(screenX, screenY int) (worldX, worldY float64) {
	px := c.GetPosition().X()
	py := c.GetPosition().Y()
	return float64(screenX) + px, float64(screenY) + py
}

package game

import (
	"math"

	. "book/code/ch02/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

const rotationSpeed = 0.02 // radians advanced per Update step

// Game drives the rotating logo demo.
type Game struct {
	engine *Engine
	logo   *Node2D
}

func (g *Game) Update() error {
	// Rotate the logo continuously (radians; 2*pi = full rotation)
	g.logo.SetRotation(g.logo.GetRotation() + rotationSpeed)
	// Wrap angle to avoid precision loss over long runs
	if g.logo.GetRotation() >= 2*math.Pi {
		g.logo.SetRotation(0)
	}
	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

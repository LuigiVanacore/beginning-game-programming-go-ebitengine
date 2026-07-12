package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Engine is the central hub of the framework.
// It owns the World and delegates Update/Draw to it.
type Engine struct {
	world *World
}

func NewEngine() *Engine {
	return &Engine{
		world: NewWorld(),
	}
}

func (e *Engine) World() *World {
	return e.world
}

func (e *Engine) Update() error {
	e.world.Update()
	return nil
}

func (e *Engine) Draw(target *ebiten.Image) {
	e.world.Draw(target)
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

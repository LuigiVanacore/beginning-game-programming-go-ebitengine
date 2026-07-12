package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Engine struct {
	world   *World
	resource *ResourceManager
}

func NewEngine() *Engine {
	return &Engine{
		world:    NewWorld(640, 480),
		resource: NewResourceManager(),
	}
}

func (e *Engine) World() *World {
	return e.world
}

func (e *Engine) ResourceManager() *ResourceManager {
	return e.resource
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

package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Engine struct {
	world      *World
	input      *InputManager
	resource   *ResourceManager
	collisions *CollisionManager
}

func NewEngine() *Engine {
	return &Engine{
		world:      NewWorld(640, 480),
		input:      NewInputManager(),
		resource:   NewResourceManager(),
		collisions: NewCollisionManager(),
	}
}

// CollisionManager returns the collision subsystem (registration and broad-phase checks).
func (e *Engine) CollisionManager() *CollisionManager {
	if e == nil {
		return nil
	}
	return e.collisions
}

func (e *Engine) World() *World {
	return e.world
}

func (e *Engine) Input() *InputManager {
	return e.input
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

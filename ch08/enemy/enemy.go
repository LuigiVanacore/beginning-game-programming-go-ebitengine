package enemy

import (
	. "book/code/ch08/internal/core"
	"math"
)

const (
	enemySpeed  = 1.4  // pixels per frame
	enemyRadius = 12.0 // collision circle radius for each enemy
)

// Enemy: a root Node2D; the Collider and Sprite are children (as in Ch6).
type Enemy struct {
	Node2D
	Collider *Collider
	Sprite   *Sprite
}

// NewEnemy registers the hierarchy through the engine.
func NewEnemy(engine *Engine, x, y float64) *Enemy {
	rm := engine.ResourceManager()
	world := engine.World()
	colMgr := engine.CollisionManager()

	e := &Enemy{
		Node2D: *NewNode2D(NameEnemy),
	}
	e.SetPosition(x, y)

	enemyShape := NewCollisionCircle(enemyRadius)
	enemyMask := NewCollisionMask(LayerEnemy, LayerPlayer|LayerProjectile)
	c := colMgr.NewCollider(NameEnemy+"_body", enemyShape, enemyMask)
	c.SetPosition(0, 0)

	enemyTex, _ := rm.GetTexture(EnemyTexture)
	enemySprite := NewSprite(NameEnemySprite, enemyTex, 0, true)
	enemySprite.SetScale(2, 2)

	e.AddChildren(c)
	e.AddChildren(enemySprite)

	e.Collider = c
	e.Sprite = enemySprite

	world.AddNodeToLayer(e, DrawLayerPlayer)
	return e
}

// Update chases toward (targetX, targetY) by moving the root.
func (e *Enemy) Update(targetX, targetY float64) {
	if e == nil {
		return
	}
	ex, ey := e.GetPosition().X(), e.GetPosition().Y()
	toX := targetX - ex
	toY := targetY - ey
	dist := math.Sqrt(toX*toX + toY*toY)
	if dist > 1 {
		toX /= dist
		toY /= dist
		e.SetPosition(ex+toX*enemySpeed, ey+toY*enemySpeed)
	}
}

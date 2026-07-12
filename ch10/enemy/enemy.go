package enemy

import (
	. "book/code/ch10/internal/core"
	"math"
	"math/rand"
)

const (
	enemySpeed  = 1.4 // pixels per frame
	enemyRadius = 12.0 // collision circle radius for each enemy
)

// SpawnPos is a map position used when placing the first enemy away from the player.
type SpawnPos struct {
	X, Y float64
}

// SpawnEnemyFarFrom returns a position near a random map edge, at least minDist from (px, py).
func SpawnEnemyFarFrom(px, py, mapW, mapH float64) SpawnPos {
	const minDist = 200.0
	for try := 0; try < 50; try++ {
		x, y := randomMapEdgePosition(mapW, mapH)
		dx, dy := x-px, y-py
		if math.Sqrt(dx*dx+dy*dy) >= minDist {
			return SpawnPos{X: x, Y: y}
		}
	}
	return SpawnPos{X: 0, Y: 0}
}

func randomMapEdgePosition(mapW, mapH float64) (x, y float64) {
	const edgeDepth = 80.0
	switch rand.Intn(4) {
	case 0: // top
		return rand.Float64() * mapW, rand.Float64() * edgeDepth
	case 1: // bottom
		return rand.Float64() * mapW, mapH - edgeDepth + rand.Float64()*edgeDepth
	case 2: // left
		return rand.Float64() * edgeDepth, rand.Float64() * mapH
	default: // right
		return mapW - edgeDepth + rand.Float64()*edgeDepth, rand.Float64() * mapH
	}
}

// Enemy: a root Node2D; the collider and sprite are children (aligned with Ch6-8).
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
		Node2D: *NewNode2D("enemy"),
	}
	e.SetPosition(x, y)

	enemyShape := NewCollisionCircle(enemyRadius)
	enemyMask := NewCollisionMask(LayerEnemy, LayerPlayer|LayerProjectile)
	c := colMgr.NewCollider("enemy_body", enemyShape, enemyMask)
	c.SetPosition(0, 0)

	enemyTex, _ := rm.GetTexture("enemy")
	enemySprite := NewSprite("enemy_sprite", enemyTex, 0)
	enemySprite.SetPivotToCenter()
	enemySprite.SetScale(2, 2)

	e.AddChildren(c)
	e.AddChildren(enemySprite)

	e.Collider = c
	e.Sprite = enemySprite

	world.AddNodeToLayer(e, DrawLayerPlayer)
	return e
}

// Update moves the enemy toward (targetX, targetY) by enemySpeed pixels each frame.
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

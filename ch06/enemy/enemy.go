package enemy

import (
	. "book/code/ch06/internal/core"
	"math"
	"math/rand"
)

const (
	enemySpeed   = 1.4   // pixels per frame
	enemyRadius  = 12.0  // collision circle radius for each enemy
	spawnMinDist = 200.0 // minimum distance from player when spawning an enemy
)

// Enemy is the enemy entity: embeds Node2D as scene root; Collider and Sprite are child components.
type Enemy struct {
	Node2D
	Collider *Collider
	Sprite   *Sprite
}

// NewEnemy builds an enemy at (x, y): root Node2D, then collider + sprite as children, registered in world/colMgr.
func NewEnemy(engine *Engine, x, y float64) *Enemy {
	world := engine.World()
	rm := engine.ResourceManager()
	colMgr := engine.CollisionManager()

	e := &Enemy{
		Node2D: *NewNode2D("enemy"),
	}
	e.SetPosition(x, y)

	shape := NewCollisionCircle(enemyRadius)
	mask := NewCollisionMask(LayerEnemy, LayerPlayer)
	c := colMgr.NewCollider("enemy_collider", shape, mask)
	c.SetPosition(0, 0)

	enemyTex, _ := rm.GetTexture("enemy")
	sprite := NewSprite("enemy_sprite", enemyTex, 0, true)
	sprite.SetPivotToCenter()
	sprite.SetScale(2, 2)

	e.AddChildren(c)
	e.AddChildren(sprite)

	e.Collider = c
	e.Sprite = sprite

	world.AddNodeToLayer(e, DrawLayerPlayer)
	return e
}

// NewEnemyFarFromPlayer spawns one enemy at random direction spawnMinDist from (px, py).
func NewEnemyFarFromPlayer(engine *Engine, px, py float64) *Enemy {
	pos := spawnEnemyFarFrom(px, py)
	return NewEnemy(engine, pos.x, pos.y)
}

func spawnEnemyFarFrom(px, py float64) spawnPos {
	angle := rand.Float64() * 2 * math.Pi
	x := px + math.Cos(angle)*spawnMinDist
	y := py + math.Sin(angle)*spawnMinDist
	return spawnPos{x, y}
}

type spawnPos struct {
	x, y float64
}

// Update moves the enemy's root toward the target (e.g. the player's world position).
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

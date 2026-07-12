package enemy

import (
	. "book/code/ch12/internal/core"
	"math"
)

const (
	enemySpeed    = 1.4  // pixels per frame
	enemyRadius   = 12.0 // collision circle radius for a regular enemy
	cyclopsRadius = 20.0 // wider collision circle for the mini-boss
)

// Enemy: a root Node2D; the collider and sprite are children (aligned with Ch6-8).
// Each enemy now carries combat stats: kinds differ in HP, contact damage, and
// the experience they grant, and the spawner scales those values over time.
type Enemy struct {
	Node2D
	Collider *Collider
	Sprite   *Sprite

	Type     EnemyType
	HP       float64 // remaining hit points; the enemy dies at zero
	MaxHP    float64 // hit points at spawn (after difficulty scaling)
	Damage   float64 // HP drained from the player per contact frame
	XPValue  int     // experience granted to the player on death
}

// NewEnemy builds an enemy of the given kind and registers its hierarchy.
// statScale multiplies HP and contact damage; xpScale (grown more slowly)
// multiplies the experience reward. Both come from the current difficulty tier.
func NewEnemy(engine *Engine, kind EnemyType, x, y, statScale, xpScale float64) *Enemy {
	rm := engine.ResourceManager()
	world := engine.World()
	colMgr := engine.CollisionManager()

	st := statsFor(kind)

	e := &Enemy{
		Node2D:  *NewNode2D("enemy"),
		Type:    kind,
		MaxHP:   st.baseHP * statScale,
		Damage:  st.baseDamage * statScale,
		XPValue: int(float64(st.baseXP) * xpScale),
	}
	e.HP = e.MaxHP
	e.SetPosition(x, y)

	enemyShape := NewCollisionCircle(st.radius)
	enemyMask := NewCollisionMask(LayerEnemy, LayerPlayer|LayerProjectile)
	c := colMgr.NewCollider("enemy_body", enemyShape, enemyMask)
	c.SetPosition(0, 0)

	enemyTex, _ := rm.GetTexture(st.textureKey)
	enemySprite := NewSprite("enemy_sprite", enemyTex, 0)
	enemySprite.SetPivotToCenter()
	enemySprite.SetScale(st.scale, st.scale)

	e.AddChildren(c)
	e.AddChildren(enemySprite)

	e.Collider = c
	e.Sprite = enemySprite

	world.AddNodeToLayer(e, DrawLayerPlayer)
	return e
}

// TakeDamage subtracts weapon damage from the enemy's hit points.
func (e *Enemy) TakeDamage(amount float64) {
	if e == nil {
		return
	}
	e.HP -= amount
	if e.HP < 0 {
		e.HP = 0
	}
}

// IsDead reports whether the enemy has lost all its hit points.
func (e *Enemy) IsDead() bool {
	return e != nil && e.HP <= 0
}

// Update chases toward (targetX, targetY) on the root.
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

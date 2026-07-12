package enemy

import . "book/code/ch12/internal/core"

// EnemyManager owns the active enemy list, the spawner, and the difficulty curve.
// It drives spawn waves and per-frame movement in one place.
type EnemyManager struct {
	enemies    []*Enemy
	spawner    *EnemySpawner
	difficulty Difficulty
	grew       bool // true on the frame the difficulty tier last advanced
}

// NewEnemyManager creates an EnemyManager with the given spawner and an empty enemy list.
func NewEnemyManager(spawner *EnemySpawner) *EnemyManager {
	return &EnemyManager{
		enemies: make([]*Enemy, 0),
		spawner: spawner,
	}
}

// Add appends a pre-built enemy to the active list.
func (m *EnemyManager) Add(e *Enemy) {
	m.enemies = append(m.enemies, e)
}

// Update advances the difficulty, runs the spawn wave, and moves every active enemy
// toward (playerX, playerY). elapsedSeconds is the survival time so far.
func (m *EnemyManager) Update(playerX, playerY, elapsedSeconds float64, engine *Engine) {
	m.grew = m.difficulty.Update(elapsedSeconds)
	m.spawner.Update(playerX, playerY, elapsedSeconds, &m.difficulty, engine, m.Add)
	for _, e := range m.enemies {
		if e != nil {
			e.Update(playerX, playerY)
		}
	}
}

// EnemiesGrewStronger reports whether the difficulty tier advanced this frame.
func (m *EnemyManager) EnemiesGrewStronger() bool { return m.grew }

// FindByCollider returns the enemy whose Collider matches c, or nil.
func (m *EnemyManager) FindByCollider(c *Collider) *Enemy {
	for _, e := range m.enemies {
		if e != nil && e.Collider == c {
			return e
		}
	}
	return nil
}

// Remove deletes enemy from the active list.
func (m *EnemyManager) Remove(enemy *Enemy) {
	for i, e := range m.enemies {
		if e == enemy {
			m.enemies = append(m.enemies[:i], m.enemies[i+1:]...)
			return
		}
	}
}

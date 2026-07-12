package enemy

import (
	. "book/code/ch09/internal/core"
	"math/rand"
	"time"
)

const (
	spawnIntervalBase      = 2.0   // initial seconds between waves
	spawnIntervalMin       = 0.5   // minimum seconds between waves (floor)
	spawnMargin            = 80.0  // pixels beyond half-screen for off-screen spawning
	spawnEscalationSeconds = 45.0  // wave interval used to compute escalation
	spawnIntervalDecay     = 0.002 // seconds subtracted from interval per wave-second
)

func randomEdgePosition(playerPosX, playerPosY float64) (x, y float64) {
	edge := rand.Intn(4)
	halfW := float64(GameSettings.ScreenWidth)/2 + spawnMargin
	halfH := float64(GameSettings.ScreenHeight)/2 + spawnMargin
	switch edge {
	case 0: // top
		x, y = playerPosX+(rand.Float64()*2-1)*halfW, playerPosY-halfH
	case 1: // bottom
		x, y = playerPosX+(rand.Float64()*2-1)*halfW, playerPosY+halfH
	case 2: // left
		x, y = playerPosX-halfW, playerPosY+(rand.Float64()*2-1)*halfH
	default: // right
		x, y = playerPosX+halfW, playerPosY+(rand.Float64()*2-1)*halfH
	}
	return x, y
}

// avgSecondsPerWaveApprox maps wave index to the old time-based escalation (interval + extra enemies).
const avgSecondsPerWaveApprox = 1.25

// EnemySpawner schedules waves using a Timer for the interval between spawns.
// The enemy field describes which kind to instantiate (nil = default CreateEnemy for this chapter).
type EnemySpawner struct {
	timer     *Timer
	enemy     *Enemy
	waveIndex int // completed waves; incremented only when the wave timer elapses
}

// NewEnemySpawner creates a spawner with the given wave timer and enemy kind (nil = default CreateEnemy).
func NewEnemySpawner(timer *Timer, enemy *Enemy) *EnemySpawner {
	return &EnemySpawner{
		timer: timer,
		enemy: enemy,
	}
}

// Update uses player position for spawn placement. Wave timing and escalation use only the internal
// Timer and waveIndex (incremented each time the timer elapses and a wave spawns).
func (s *EnemySpawner) Update(playerPosX, playerPosY float64,
	engine *Engine, addEnemy func(*Enemy)) {
	if s == nil || addEnemy == nil || s.timer == nil {
		return
	}
	s.refreshWaveTimerDuration()
	if !s.timer.IsEnded() {
		return
	}
	s.spawnWaveOnTimerElapsed(engine, playerPosX, playerPosY, addEnemy)
}

func (s *EnemySpawner) refreshWaveTimerDuration() {
	s.timer.EnsureStarted()
	interval := s.currentSpawnIntervalSeconds()
	s.timer.SetDuration(time.Duration(interval * float64(time.Second)))
}

func (s *EnemySpawner) currentSpawnIntervalSeconds() float64 {
	effectiveSeconds := float64(s.waveIndex) * avgSecondsPerWaveApprox
	interval := spawnIntervalBase - effectiveSeconds*spawnIntervalDecay
	if interval < spawnIntervalMin {
		return spawnIntervalMin
	}
	return interval
}

func (s *EnemySpawner) enemiesToSpawnThisWave() int {
	effectiveSeconds := float64(s.waveIndex) * avgSecondsPerWaveApprox
	extra := int(effectiveSeconds / spawnEscalationSeconds)
	if extra > 6 {
		extra = 6
	}
	return 2 + rand.Intn(2) + extra
}

// spawnWaveOnTimerElapsed runs when the wave timer has elapsed: restart timer, advance wave, spawn enemies.
func (s *EnemySpawner) spawnWaveOnTimerElapsed(engine *Engine, playerPosX, playerPosY float64, addEnemy func(*Enemy)) {
	n := s.enemiesToSpawnThisWave()
	s.timer.Restart()
	s.waveIndex++
	for i := 0; i < n; i++ {
		addEnemy(s.spawnAtEdge(engine, playerPosX, playerPosY))
	}
}

func (s *EnemySpawner) spawnAtEdge(engine *Engine, playerPosX, playerPosY float64) *Enemy {
	x, y := randomEdgePosition(playerPosX, playerPosY)
	return s.instantiateEnemy(engine, x, y)
}

func (s *EnemySpawner) instantiateEnemy(engine *Engine, x, y float64) *Enemy {
	// s.enemy will select the kind when this chapter supports multiple enemy types.
	return CreateEnemy(engine, x, y)
}

// CreateEnemy is the default factory used by the spawner (one visual kind per chapter).
func CreateEnemy(engine *Engine, x, y float64) *Enemy {
	return NewEnemy(engine, x, y)
}

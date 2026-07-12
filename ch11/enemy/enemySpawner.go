package enemy

import (
	. "book/code/ch11/internal/core"
	"math/rand"
	"time"
)

const (
	spawnIntervalBase      = 2.0   // initial seconds between waves
	spawnIntervalMin       = 0.5   // minimum seconds between waves (floor)
	spawnMargin            = 80.0  // pixels beyond half-screen for off-screen spawning
	spawnEscalationSeconds = 45.0  // wave interval used to compute escalation
	spawnIntervalDecay     = 0.002 // seconds subtracted from interval per wave-second
	cyclopsIntervalSeconds = 300.0 // the Cyclops mini-boss appears every 5 minutes
	typeUnlockBaseSeconds  = 30.0  // Spider at 30s, Bat at 90s, DarkWizard at 210s: base*(2^k-1), gaps doubling
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
// It also releases the solitary Cyclops mini-boss on its own five-minute cadence.
type EnemySpawner struct {
	timer         *Timer
	waveIndex     int // completed waves; incremented only when the wave timer elapses
	bossesSpawned int // Cyclops mini-bosses released so far
}

// NewEnemySpawner creates a spawner driven by the given wave timer.
func NewEnemySpawner(timer *Timer) *EnemySpawner {
	return &EnemySpawner{
		timer: timer,
	}
}

// Update uses player position for spawn placement. Wave timing and escalation use the internal
// Timer and waveIndex; the difficulty selects which kinds spawn and how strong they are.
// The Cyclops mini-boss is checked first and, when due, spawns alone for that frame.
func (s *EnemySpawner) Update(playerPosX, playerPosY, elapsedSeconds float64,
	diff *Difficulty, engine *Engine, addEnemy func(*Enemy)) {
	if s == nil || addEnemy == nil || s.timer == nil {
		return
	}
	if s.trySpawnBoss(elapsedSeconds, diff, engine, playerPosX, playerPosY, addEnemy) {
		return
	}
	s.refreshWaveTimerDuration()
	if !s.timer.IsEnded() {
		return
	}
	s.spawnWaveOnTimerElapsed(elapsedSeconds, diff, engine, playerPosX, playerPosY, addEnemy)
}

// trySpawnBoss releases one Cyclops each time a five-minute mark is passed. The
// boss arrives alone: when this returns true, no regular wave spawns that frame.
func (s *EnemySpawner) trySpawnBoss(elapsedSeconds float64, diff *Difficulty,
	engine *Engine, playerPosX, playerPosY float64, addEnemy func(*Enemy)) bool {
	due := int(elapsedSeconds / cyclopsIntervalSeconds)
	if due <= s.bossesSpawned {
		return false
	}
	s.bossesSpawned = due
	x, y := randomEdgePosition(playerPosX, playerPosY)
	addEnemy(createEnemyOfType(engine, Cyclops, x, y, diff))
	return true
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
// Each enemy's kind is drawn from the kinds unlocked by the survival time, biased toward the stronger ones.
func (s *EnemySpawner) spawnWaveOnTimerElapsed(elapsedSeconds float64, diff *Difficulty, engine *Engine, playerPosX, playerPosY float64, addEnemy func(*Enemy)) {
	n := s.enemiesToSpawnThisWave()
	s.timer.Restart()
	s.waveIndex++
	kinds := unlockedTypes(elapsedSeconds)
	for i := 0; i < n; i++ {
		x, y := randomEdgePosition(playerPosX, playerPosY)
		addEnemy(createEnemyOfType(engine, pickWeighted(kinds), x, y, diff))
	}
}

// unlockedTypes returns the regular kinds available at the given survival time, from
// weakest upward. Ghost is always available; each further kind unlocks after a
// doubling gap — Spider at 30s, Bat at 90s, DarkWizard at 210s (base*(2^k - 1)) — so
// new kinds arrive quickly at first, then ever more rarely. This is independent of the
// stat difficulty tier. The Cyclops is never part of this pool: it spawns on its own
// boss timer.
func unlockedTypes(elapsedSeconds float64) []EnemyType {
	regular := []EnemyType{Ghost, Spider, Bat, DarkWizard}
	n := 1 // Ghost is always unlocked
	for k := 1; k < len(regular); k++ {
		steps := (1 << k) - 1 // 1, 3, 7 -> thresholds 30s, 90s, 210s
		if elapsedSeconds < typeUnlockBaseSeconds*float64(steps) {
			break
		}
		n++
	}
	return regular[:n]
}

// pickWeighted chooses a kind at random, giving the stronger (later) kinds a higher
// weight so that tougher monsters gradually dominate the waves.
func pickWeighted(kinds []EnemyType) EnemyType {
	total := 0
	for i := range kinds {
		total += i + 1
	}
	r := rand.Intn(total)
	for i, k := range kinds {
		r -= i + 1
		if r < 0 {
			return k
		}
	}
	return kinds[len(kinds)-1]
}

// createEnemyOfType builds an enemy of the given kind, scaling its stats and XP by the current tier.
func createEnemyOfType(engine *Engine, kind EnemyType, x, y float64, diff *Difficulty) *Enemy {
	return NewEnemy(engine, kind, x, y, diff.StatScale(), diff.XPScale())
}

package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle represents a single particle with position, velocity, and lifetime.
type Particle struct {
	x, y     float64
	vx, vy   float64
	life     float64
	maxLife  float64
	color    color.RGBA
	endColor color.RGBA // fades to this (typically more transparent)
	size     float64
}

// ParticleSystem manages active particles and provides emit methods. It lives on
// Game, updates every frame, and draws in world space behind the HUD (new in ch13).
type ParticleSystem struct {
	particles []Particle
}

// NewParticleSystem creates a new particle system with preallocated capacity.
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0, 256),
	}
}

// --- Blood ---

// EmitBlood spawns a red burst (8-16 particles) at the given world position.
// Used when an enemy dies.
func (ps *ParticleSystem) EmitBlood(worldX, worldY float64) {
	n := 8 + rand.Intn(9)
	for i := 0; i < n; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 80 + rand.Float64()*120
		ps.emitParticle(worldX, worldY, angle, speed, 0.3, 0.4, 2.0, 3.0,
			color.RGBA{180, 0, 0, 255}, color.RGBA{180, 0, 0, 0})
	}
}

// EmitBloodSmall spawns a smaller burst (3-6 particles). Used on each weapon hit
// and when the player takes contact damage.
func (ps *ParticleSystem) EmitBloodSmall(worldX, worldY float64) {
	n := 3 + rand.Intn(4)
	for i := 0; i < n; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 40 + rand.Float64()*50
		ps.emitParticle(worldX, worldY, angle, speed, 0.2, 0.2, 1.5, 1.5,
			color.RGBA{180, 0, 0, 255}, color.RGBA{180, 0, 0, 0})
	}
}

// --- Weapon trails ---

// EmitProjectileTrail spawns a light trail behind a projectile (knife, flying axe).
// worldX, worldY is the projectile position; vx, vy its velocity (trail drifts opposite).
func (ps *ParticleSystem) EmitProjectileTrail(worldX, worldY, vx, vy float64) {
	n := 1 + rand.Intn(2)
	speed := math.Sqrt(vx*vx + vy*vy)
	if speed < 1 {
		speed = 1
	}
	backX, backY := -vx/speed*25, -vy/speed*25
	trailColor := color.RGBA{220, 230, 255, 200}
	trailEnd := color.RGBA{200, 210, 255, 0}
	for i := 0; i < n; i++ {
		ox := (rand.Float64() - 0.5) * 6
		oy := (rand.Float64() - 0.5) * 6
		ps.particles = append(ps.particles, Particle{
			x:        worldX + ox,
			y:        worldY + oy,
			vx:       backX + (rand.Float64()-0.5)*25,
			vy:       backY + (rand.Float64()-0.5)*25,
			life:     0.15 + rand.Float64()*0.12,
			maxLife:  0.15 + rand.Float64()*0.12,
			color:    trailColor,
			endColor: trailEnd,
			size:     1.2 + rand.Float64()*1.0,
		})
	}
}

// EmitOrbitTrail spawns trail particles for the Sacred Book: they start at the book
// position and drift backward along the orbit tangent (no radial spread).
func (ps *ParticleSystem) EmitOrbitTrail(worldX, worldY, tangentVx, tangentVy float64) {
	n := 1 + rand.Intn(2)
	speed := math.Sqrt(tangentVx*tangentVx + tangentVy*tangentVy)
	if speed < 1 {
		speed = 1
	}
	trailSpeed := 25 + rand.Float64()*20
	backX := -tangentVx / speed * trailSpeed
	backY := -tangentVy / speed * trailSpeed
	trailColor := color.RGBA{220, 230, 255, 200}
	trailEnd := color.RGBA{200, 210, 255, 0}
	for i := 0; i < n; i++ {
		ps.particles = append(ps.particles, Particle{
			x:        worldX,
			y:        worldY,
			vx:       backX,
			vy:       backY,
			life:     0.15 + rand.Float64()*0.12,
			maxLife:  0.15 + rand.Float64()*0.12,
			color:    trailColor,
			endColor: trailEnd,
			size:     1.2 + rand.Float64()*1.0,
		})
	}
}

// --- Level up ---

// EmitLevelUp spawns an upward golden fountain (12-20 particles) at the position.
func (ps *ParticleSystem) EmitLevelUp(worldX, worldY float64) {
	n := 12 + rand.Intn(9)
	for i := 0; i < n; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*math.Pi
		speed := 100 + rand.Float64()*80
		vx := math.Cos(angle) * speed * 0.3
		vy := math.Sin(angle) * speed
		life := 0.5 + rand.Float64()*0.5
		size := 2.0 + rand.Float64()*4
		ps.particles = append(ps.particles, Particle{
			x:        worldX,
			y:        worldY,
			vx:       vx,
			vy:       vy,
			life:     life,
			maxLife:  life,
			color:    color.RGBA{255, 220, 0, 255},
			endColor: color.RGBA{255, 220, 0, 0},
			size:     size,
		})
	}
}

// emitParticle is the helper for burst-style emissions (blood).
func (ps *ParticleSystem) emitParticle(x, y float64, angle, speed, lifeMin, lifeMax, sizeMin, sizeMax float64, col, endCol color.RGBA) {
	vx := math.Cos(angle) * speed
	vy := math.Sin(angle) * speed
	life := lifeMin + rand.Float64()*lifeMax
	size := sizeMin + rand.Float64()*sizeMax
	ps.particles = append(ps.particles, Particle{
		x: x, y: y, vx: vx, vy: vy,
		life: life, maxLife: life,
		color: col, endColor: endCol, size: size,
	})
}

// Update advances every particle and removes the expired ones. It reuses the
// backing array (particles[:0]) so a steady stream of effects makes no garbage.
func (ps *ParticleSystem) Update() {
	tps := ebiten.TPS()
	if tps < 1 {
		tps = 1
	}
	inv := 1.0 / float64(tps)
	newParticles := ps.particles[:0]
	for _, p := range ps.particles {
		p.x += p.vx * inv
		p.y += p.vy * inv
		p.life -= inv
		if p.life > 0 {
			newParticles = append(newParticles, p)
		}
	}
	ps.particles = newParticles
}

// Draw renders every particle. camX, camY are the camera's world position (its
// top-left), so world coordinates map to screen coordinates.
func (ps *ParticleSystem) Draw(screen *ebiten.Image, camX, camY float64) {
	for _, p := range ps.particles {
		sx := p.x - camX
		sy := p.y - camY
		t := 1.0 - p.life/p.maxLife
		col := color.RGBA{
			R: lerpU8(p.color.R, p.endColor.R, t),
			G: lerpU8(p.color.G, p.endColor.G, t),
			B: lerpU8(p.color.B, p.endColor.B, t),
			A: lerpU8(p.color.A, p.endColor.A, t),
		}
		half := float32(p.size / 2)
		vector.DrawFilledRect(screen, float32(sx)-half, float32(sy)-half, float32(p.size), float32(p.size), col, true)
	}
}

// lerpU8 linearly interpolates two 8-bit channels by t in [0,1].
func lerpU8(a, b uint8, t float64) uint8 {
	return uint8(float64(a)*(1-t) + float64(b)*t)
}

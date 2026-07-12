package game

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// floatingText is a short-lived damage number that rises and fades.
type floatingText struct {
	x, y    float64
	vy      float64
	life    float64
	maxLife float64
	label   string
	color   color.RGBA
}

// FloatingTextSystem manages floating damage numbers in world space (new in ch13).
type FloatingTextSystem struct {
	items []floatingText
	face  text.Face
}

// NewFloatingTextSystem creates the system with a shared bitmap font face.
func NewFloatingTextSystem() *FloatingTextSystem {
	return &FloatingTextSystem{
		items: make([]floatingText, 0, 64),
		face:  text.NewGoXFace(basicfont.Face7x13),
	}
}

// AddDamage spawns a rising damage number just above the given world position.
func (fts *FloatingTextSystem) AddDamage(worldX, worldY, amount float64) {
	const life = 0.6
	fts.items = append(fts.items, floatingText{
		x:       worldX + (rand.Float64()-0.5)*8,
		y:       worldY - 8,
		vy:      -40, // pixels per second, upward
		life:    life,
		maxLife: life,
		label:   strconv.Itoa(int(amount + 0.5)),
		color:   color.RGBA{255, 240, 120, 255},
	})
}

// Update advances every number and drops the expired ones.
func (fts *FloatingTextSystem) Update() {
	tps := ebiten.TPS()
	if tps < 1 {
		tps = 1
	}
	inv := 1.0 / float64(tps)
	kept := fts.items[:0]
	for _, it := range fts.items {
		it.y += it.vy * inv
		it.life -= inv
		if it.life > 0 {
			kept = append(kept, it)
		}
	}
	fts.items = kept
}

// Draw renders every number, fading it out over its lifetime. camX, camY are the
// camera's world position so world coordinates map to the screen.
func (fts *FloatingTextSystem) Draw(screen *ebiten.Image, camX, camY float64) {
	for _, it := range fts.items {
		c := it.color
		c.A = uint8(float64(c.A) * (it.life / it.maxLife))
		op := &text.DrawOptions{}
		op.GeoM.Translate(it.x-camX, it.y-camY)
		op.ColorScale.ScaleWithColor(c)
		text.Draw(screen, it.label, fts.face, op)
	}
}

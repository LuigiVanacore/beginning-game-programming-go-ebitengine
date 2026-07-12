package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawShape draws a colored figure into dst (e.g. a small texture for a sprite).
type DrawShape func(dst *ebiten.Image, c color.Color)

// ApplyDrawShape runs the given shape painter with color c on dst.
func ApplyDrawShape(dst *ebiten.Image, shape DrawShape, c color.Color) {
	if dst == nil || shape == nil {
		return
	}
	shape(dst, c)
}

// ShapeFilledCircle returns a DrawShape that fills a circle at (cx, cy) with radius r.
func ShapeFilledCircle(cx, cy, r float32) DrawShape {
	return func(dst *ebiten.Image, c color.Color) {
		vector.DrawFilledCircle(dst, cx, cy, r, colorToRGBA8(c), true)
	}
}

// ShapeFilledRect returns a DrawShape that fills an axis-aligned rectangle at (x, y) with size (w, h).
func ShapeFilledRect(x, y, w, h float32) DrawShape {
	return func(dst *ebiten.Image, c color.Color) {
		vector.DrawFilledRect(dst, x, y, w, h, colorToRGBA8(c), true)
	}
}

func colorToRGBA8(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

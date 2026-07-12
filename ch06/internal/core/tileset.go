package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tileset holds a texture and the parameters to extract tiles from it.
type Tileset struct {
	texture     *ebiten.Image
	tileWidth   int
	tileHeight  int
	offsetX     int // spacing between tiles (e.g. 1)
	offsetY     int
}

// NewTileset creates a Tileset from a texture with the given cell size and offset.
func NewTileset(texture *ebiten.Image, tileWidth, tileHeight, offsetX, offsetY int) *Tileset {
	return &Tileset{
		texture:    texture,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		offsetX:    offsetX,
		offsetY:    offsetY,
	}
}

// GetTileRect returns the source rectangle for the tile at cell (col, row).
func (t *Tileset) GetTileRect(col, row int) image.Rectangle {
	sx := col * (t.tileWidth + t.offsetX)
	sy := row * (t.tileHeight + t.offsetY)
	return image.Rect(sx, sy, sx+t.tileWidth, sy+t.tileHeight)
}

// GetTexture returns the tileset texture.
func (t *Tileset) GetTexture() *ebiten.Image {
	return t.texture
}

// TileWidth returns the tile width.
func (t *Tileset) TileWidth() int { return t.tileWidth }

// TileHeight returns the tile height.
func (t *Tileset) TileHeight() int { return t.tileHeight }

package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// InfiniteTilemapNode tiles a fixed 2D pattern across the world. Only cells
// visible in the camera viewport are drawn, so memory is O(pattern), not O(map).
type InfiniteTilemapNode struct {
	Node2D
	tileset     *Tileset
	layer       int
	pattern     [][][2]int
	camera      *Camera
	marginTiles int
}

// NewInfiniteTilemapNode creates an infinite tilemap that renders only visible tiles.
// pattern[row][col] stores a tileset cell (tilesetCol, tilesetRow) and repeats forever.
// marginTiles adds extra tiles beyond viewport edges (default 1) to avoid seams.
func NewInfiniteTilemapNode(name string, tileset *Tileset, pattern [][][2]int, layer int, camera *Camera, marginTiles int) *InfiniteTilemapNode {
	if marginTiles < 1 {
		marginTiles = 1
	}
	return &InfiniteTilemapNode{
		Node2D:      *NewNode2D(name),
		tileset:     tileset,
		layer:       layer,
		pattern:     pattern,
		camera:      camera,
		marginTiles: marginTiles,
	}
}

// ParseTilePattern converts a floor.map text file into the 2D pattern slice
// used by NewInfiniteTilemapNode. Each non-comment, non-blank line is one
// tile row; space-separated integers are indices into tileIndex.
func ParseTilePattern(data string, tileIndex [][2]int) ([][][2]int, error) {
	var rows [][][2]int
	width := -1
	for lineNumber, line := range strings.Split(strings.TrimSpace(data), "\n") {
		cols, ok, err := parseTilePatternLine(line, lineNumber+1, tileIndex)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if width == -1 {
			width = len(cols)
		} else if len(cols) != width {
			return nil, fmt.Errorf("line %d: got %d columns, want %d", lineNumber+1, len(cols), width)
		}
		rows = append(rows, cols)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("tile pattern is empty")
	}
	return rows, nil
}

func parseTilePatternLine(line string, lineNumber int, tileIndex [][2]int) ([][2]int, bool, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, false, nil
	}
	var cols [][2]int
	for _, token := range strings.Fields(line) {
		cell, err := parseTilePatternCell(token, lineNumber, tileIndex)
		if err != nil {
			return nil, false, err
		}
		cols = append(cols, cell)
	}
	return cols, true, nil
}

func parseTilePatternCell(token string, lineNumber int, tileIndex [][2]int) ([2]int, error) {
	index, err := strconv.Atoi(token)
	if err != nil {
		return [2]int{}, fmt.Errorf("line %d: tile index %q: %w", lineNumber, token, err)
	}
	if index < 0 || index >= len(tileIndex) {
		return [2]int{}, fmt.Errorf("line %d: tile index %d out of range", lineNumber, index)
	}
	return tileIndex[index], nil
}

// patternAt returns the tile cell for world coordinates (col, row), wrapping by
// pattern dimensions. The modulo expression handles negative coordinates.
func (t *InfiniteTilemapNode) patternAt(col, row int) ([2]int, bool) {
	height := len(t.pattern)
	if height == 0 || len(t.pattern[0]) == 0 {
		return [2]int{}, false
	}
	width := len(t.pattern[0])
	patternRow := ((row % height) + height) % height
	patternCol := ((col % width) + width) % width
	return t.pattern[patternRow][patternCol], true
}

func (t *InfiniteTilemapNode) GetLayer() int  { return t.layer }
func (t *InfiniteTilemapNode) SetLayer(l int) { t.layer = l }

func (t *InfiniteTilemapNode) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	if t.tileset == nil || t.tileset.GetTexture() == nil {
		return
	}
	tex := t.tileset.GetTexture()
	tw := float64(t.tileset.TileWidth())
	th := float64(t.tileset.TileHeight())

	minCol, minRow, maxCol, maxRow := t.visibleCellRange(tw, th)

	for row := minRow; row <= maxRow; row++ {
		for col := minCol; col <= maxCol; col++ {
			t.drawFloorTileAt(target, op, tex, tw, th, col, row)
		}
	}
}

// viewportWorldRect returns the camera view in world space. If camera is nil, all zeros (same as before).
func (t *InfiniteTilemapNode) viewportWorldRect() (camX, camY, viewW, viewH float64) {
	if t.camera == nil {
		return 0, 0, 0, 0
	}
	camX = t.camera.GetPosition().X()
	camY = t.camera.GetPosition().Y()
	viewW = float64(t.camera.GetWidth())
	viewH = float64(t.camera.GetHeight())
	return camX, camY, viewW, viewH
}

// visibleCellRange is the inclusive [minCol..maxCol] x [minRow..maxRow] of tile cells to draw.
func (t *InfiniteTilemapNode) visibleCellRange(tw, th float64) (minCol, minRow, maxCol, maxRow int) {
	camX, camY, viewW, viewH := t.viewportWorldRect()
	m := t.marginTiles
	minCol = int(camX/tw) - m
	minRow = int(camY/th) - m
	maxCol = int((camX+viewW)/tw) + m
	maxRow = int((camY+viewH)/th) + m
	return minCol, minRow, maxCol, maxRow
}

func (t *InfiniteTilemapNode) drawFloorTileAt(target *ebiten.Image, op *ebiten.DrawImageOptions,
	tex *ebiten.Image, tw, th float64, col, row int) {
	cell, ok := t.patternAt(col, row)
	if !ok {
		return
	}
	srcRect := t.tileset.GetTileRect(cell[0], cell[1])
	subImg := tex.SubImage(srcRect).(*ebiten.Image)

	tileOp := &ebiten.DrawImageOptions{}
	tileOp.GeoM.Translate(float64(col)*tw, float64(row)*th)
	tileOp.GeoM.Concat(op.GeoM)
	target.DrawImage(subImg, tileOp)
}

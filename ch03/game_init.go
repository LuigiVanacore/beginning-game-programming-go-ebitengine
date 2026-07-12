package game

import (
	ast "book/code/ch03/assets"
	. "book/code/ch03/internal/core"
)

const (
	floorTextureKey = "floor_tile"
	playerTextureKey = "player"
	floorName       = "floor"
	floorLayerIndex = 0
)

func loadTextures(rm *ResourceManager) {
	if rm == nil {
		return
	}
	rm.UseEmbeddedFS(ast.FS)
	rm.LoadTexture(ast.FloorTile, floorTextureKey)
	rm.LoadTexture(ast.Player, playerTextureKey)
}

func setupFloor(world *World, rm *ResourceManager) {
	if world == nil || rm == nil {
		return
	}
	tex, _ := rm.GetTexture(floorTextureKey)
	if tex == nil {
		return
	}
	floor := NewSprite(floorName, tex, floorLayerIndex, false)
	floor.SetPosition(0, 0)
	floor.SetPivot(0, 0)
	b := tex.Bounds()
	tw, th := float64(b.Dx()), float64(b.Dy())
	if tw > 0 && th > 0 {
		floor.SetScale(float64(GameSettings.ScreenWidth)/tw, float64(GameSettings.ScreenHeight)/th)
	}
	world.AddNodeToLayer(floor, floorLayerIndex)
}

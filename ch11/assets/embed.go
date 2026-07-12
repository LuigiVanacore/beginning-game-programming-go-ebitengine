package assets

import "embed"

// FS holds every runtime asset: the floor pattern and all sprites (floor tileset,
// player, monsters, weapons, pickups, and upgrade icons) under sprites/.
//
//go:embed sprites floor.map
var FS embed.FS

package assets

import "embed"

// FS holds every runtime asset: the floor pattern, all sprites (floor tileset,
// player, monsters, weapons, pickups, and upgrade icons) under sprites/, and the
// background music under audio/ (new in ch12).
//
//go:embed sprites floor.map audio
var FS embed.FS

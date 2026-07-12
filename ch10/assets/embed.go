package assets

import "embed"

// FS embeds the sprites/ directory (all PNGs, including the upgrade icons)
// and the floor.map level data.
//
//go:embed sprites floor.map
var FS embed.FS

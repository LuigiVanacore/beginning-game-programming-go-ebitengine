package assets

import "embed"

// FS is the embedded sprites/ tree. Paths in paths.go use the "sprites/..." prefix.
//
//go:embed sprites
var FS embed.FS

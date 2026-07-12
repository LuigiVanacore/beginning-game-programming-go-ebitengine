package core

import (
	"os"
	"path/filepath"
)

// ResolveAssetPath returns a path to rel that exists, by walking up from the current
// working directory (e.g. starting in cmd/ finds assets/ in the chapter folder).
// If rel is absolute, it is returned unchanged. If nothing is found, rel is returned as-is.
func ResolveAssetPath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	dir, err := os.Getwd()
	if err != nil {
		return rel
	}
	rel = filepath.FromSlash(rel)
	for {
		candidate := filepath.Join(dir, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return rel
}

package core

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ResourceManager loads and caches textures by name.
type ResourceManager struct {
	textures map[string]*ebiten.Image
	fsys     fs.FS // optional embedded tree (embed.FS); nil = load from disk via ResolveAssetPath
}

// NewResourceManager creates an empty resource manager.
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		textures: make(map[string]*ebiten.Image),
	}
}

// UseEmbeddedFS sets the filesystem used for LoadTexture (e.g. embed.FS from the assets package).
// Call before any LoadTexture when using go:embed.
func (r *ResourceManager) UseEmbeddedFS(fsys fs.FS) {
	r.fsys = fsys
}

// LoadTexture loads an image from path and stores it under the given name.
// Path is relative to the embedded FS root when UseEmbeddedFS was called, or a disk path otherwise.
func (r *ResourceManager) LoadTexture(path, name string) {
	var img *ebiten.Image
	var err error
	if r.fsys != nil {
		var data []byte
		data, err = fs.ReadFile(r.fsys, path)
		if err != nil {
			log.Fatalf("failed to read embedded texture %q at %q: %v", name, path, err)
		}
		var decoded image.Image
		decoded, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatalf("failed to decode texture %q: %v", name, err)
		}
		img = ebiten.NewImageFromImage(decoded)
	} else {
		path = ResolveAssetPath(path)
		img, _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatalf("failed to load texture %q from %q: %v", name, path, err)
		}
	}
	r.textures[name] = img
}

// GetTexture returns the texture stored under name.
// The boolean reports whether a texture with that name was loaded.
func (r *ResourceManager) GetTexture(name string) (*ebiten.Image, bool) {
	tex, ok := r.textures[name]
	return tex, ok
}

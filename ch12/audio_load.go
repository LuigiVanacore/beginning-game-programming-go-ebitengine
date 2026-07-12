package game

import (
	ast "book/code/ch12/assets"
	. "book/code/ch12/internal/core"
	"log"
)

// soundFiles maps each SoundID to its embedded WAV path. Adding a new sound is a
// two-line change: a SoundID constant in internal/core/audio.go and an entry here.
var soundFiles = map[SoundID]string{
	SoundMusic: ast.MusicTrack,
}

// loadAudio decodes and registers every sound once at startup. It is called from
// NewApp before the game loop begins. A missing or undecodable file is logged and
// skipped, so the game still runs (silently for that sound) rather than crashing.
func loadAudio() {
	am := Audio()
	for id, path := range soundFiles {
		data, err := ast.FS.ReadFile(path)
		if err != nil {
			log.Printf("audio: read %s: %v", path, err)
			continue
		}
		if err := am.RegisterSound(id, data); err != nil {
			log.Printf("audio: decode %s: %v", path, err)
		}
	}
}

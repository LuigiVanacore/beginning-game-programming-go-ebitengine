package core

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// audioSampleRate is the process-wide mixing rate. Every decoded stream is
// resampled to it, so all sounds share the single audio.Context Ebitengine allows.
const audioSampleRate = 44100

// SoundID names a registered sound. Gameplay code selects a sound by its constant
// and plays it without knowing the underlying file (new in ch12).
type SoundID string

const (
	SoundHit       SoundID = "hit"
	SoundPickup    SoundID = "pickup"
	SoundLevelUp   SoundID = "levelup"
	SoundClick     SoundID = "click"
	SoundMusic     SoundID = "music"
	SoundEnemyGrow SoundID = "enemy_grow"
)

// AudioManager owns the single audio.Context, the decoded PCM for every sound,
// the currently playing music track, and its own playback settings. It is a
// self-contained subsystem: it never reads global game state, so callers push the
// master volume and SFX toggle in through SetMasterVolume / SetSFXEnabled.
type AudioManager struct {
	ctx          *audio.Context
	sfx          map[SoundID][]byte
	music        *audio.Player
	masterVolume float64 // 0..1, applied to every sound and the music track
	sfxEnabled   bool    // when false, Play is a no-op (music is unaffected)
}

// globalAudio is the process-wide manager. Ebitengine permits only one
// audio.Context per process, so the manager is a singleton.
var globalAudio *AudioManager

// Audio returns the process-wide audio manager, creating it (and its audio.Context)
// on first use. Its defaults match the game's initial settings; the options screen
// keeps them in sync afterwards.
func Audio() *AudioManager {
	if globalAudio == nil {
		globalAudio = &AudioManager{
			ctx:          audio.NewContext(audioSampleRate),
			sfx:          make(map[SoundID][]byte),
			masterVolume: 0.8,
			sfxEnabled:   true,
		}
	}
	return globalAudio
}

// RegisterSound decodes a WAV and stores its PCM under id for cheap replay.
// Registering the same id again replaces the previous sound.
func (a *AudioManager) RegisterSound(id SoundID, wavBytes []byte) error {
	stream, err := wav.DecodeWithSampleRate(audioSampleRate, bytes.NewReader(wavBytes))
	if err != nil {
		return err
	}
	pcm, err := io.ReadAll(stream)
	if err != nil {
		return err
	}
	a.sfx[id] = pcm
	return nil
}

// Play starts a one-shot sound effect selected by its constant. It honours the
// manager's own SFX toggle and master volume, and is a no-op when the sound was
// never registered, so call sites never need to guard the call.
func (a *AudioManager) Play(id SoundID) {
	if a == nil || !a.sfxEnabled {
		return
	}
	pcm, ok := a.sfx[id]
	if !ok {
		return
	}
	// NewPlayerFromBytes returns an independent player each call, so repeated or
	// overlapping effects (several hits in one frame) mix correctly.
	p := a.ctx.NewPlayerFromBytes(pcm)
	p.SetVolume(a.masterVolume)
	p.Play()
}

// PlayMusic loops the sound registered under id as background music, replacing any
// track already playing. It honours the master volume but ignores the SFX toggle.
func (a *AudioManager) PlayMusic(id SoundID) {
	if a == nil {
		return
	}
	a.StopMusic()
	pcm, ok := a.sfx[id]
	if !ok {
		return
	}
	loop := audio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	p, err := a.ctx.NewPlayer(loop)
	if err != nil {
		return
	}
	p.SetVolume(a.masterVolume)
	p.Play()
	a.music = p
}

// StopMusic stops and releases the current music track, if any.
func (a *AudioManager) StopMusic() {
	if a == nil || a.music == nil {
		return
	}
	a.music.Pause()
	a.music.Close()
	a.music = nil
}

// SetMasterVolume sets the volume applied to every sound and updates the playing
// music track at once, so an options change takes effect immediately.
func (a *AudioManager) SetMasterVolume(v float64) {
	if a == nil {
		return
	}
	a.masterVolume = v
	if a.music != nil {
		a.music.SetVolume(v)
	}
}

// SetSFXEnabled turns one-shot sound effects on or off. Music is unaffected.
func (a *AudioManager) SetSFXEnabled(enabled bool) {
	if a != nil {
		a.sfxEnabled = enabled
	}
}

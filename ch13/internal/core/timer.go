package core

import "time"

// Timer measures elapsed time and supports looping (like ebiten_extended).
type Timer struct {
	startTime time.Time
	duration  time.Duration
	looped    bool
}

// NewTimer creates a timer with the given duration. isLooping: when ended, auto-restart.
func NewTimer(duration time.Duration, isLooping bool) *Timer {
	return &Timer{duration: duration, looped: isLooping}
}

// Start begins the countdown.
func (t *Timer) Start() *Timer {
	t.startTime = time.Now()
	return t
}

// IsEnded returns true if the duration has elapsed.
func (t *Timer) IsEnded() bool {
	return time.Since(t.startTime) >= t.duration
}

// SetDuration updates the target duration (e.g. variable spawn interval each frame).
func (t *Timer) SetDuration(d time.Duration) {
	if t == nil {
		return
	}
	t.duration = d
}

// EnsureStarted calls Start if the timer was never started.
func (t *Timer) EnsureStarted() {
	if t != nil && t.startTime.IsZero() {
		t.Start()
	}
}

// Restart resets the timer and returns elapsed time.
func (t *Timer) Restart() time.Duration {
	elapsed := time.Since(t.startTime)
	t.startTime = time.Now()
	return elapsed
}

// Update checks if ended; if looping, restarts. Returns true when the timer just ended this frame.
func (t *Timer) Update() bool {
	if t.startTime.IsZero() {
		return false
	}
	if !t.IsEnded() {
		return false
	}
	if t.looped {
		t.Restart()
	}
	return true
}

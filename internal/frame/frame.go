package frame

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	start = time.Now()
)

const (
	// Give an arbitrary head start before glitching everything to account
	// for assets initialization / procedural generation
	headStart = 3 * time.Second
)

// Chilling reports whether we're chilling or ready for glitch bombing
func Chilling() bool {
	return time.Since(start) < headStart
}

func TPS() uint64 {
	// Return fixed tps for glitch logic
	return 60
	tps := ebiten.TPS()
	if tps == ebiten.SyncWithFPS {
		return 60
	}
	return uint64(tps)
}

func Current() uint64 {
	elapsed := time.Since(start)
	duration := time.Second / time.Duration(TPS())
	return uint64(elapsed / duration)
}

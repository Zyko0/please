package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/frame"
)

func NewScreenEventNoop(rng *rand.Rand) *ScreenEvent {
	return &ScreenEvent{
		shader: assets.ShaderNoop,

		Name:     "Noop",
		Duration: uint64(config.ScreenEventDuration * float64(frame.TPS())),
	}
}

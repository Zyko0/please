package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/frame"
)

func NewScreenEventRelief(rng *rand.Rand) *ScreenEvent {
	return &ScreenEvent{
		shader: assets.ShaderRelief,

		Name:     "Relief",
		Duration: uint64(config.ScreenEventDuration * float64(frame.TPS())),
	}
}

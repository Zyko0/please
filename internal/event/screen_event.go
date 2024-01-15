package event

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenEvent struct {
	shader *ebiten.Shader

	Name     string
	Duration uint64
}

func (se *ScreenEvent) Shader() *ebiten.Shader {
	return se.shader
}

func (se *ScreenEvent) Update() {
	if se == nil || se.Duration == 0 {
		return
	}

	se.Duration--
}

func (se *ScreenEvent) Expired() bool {
	return se == nil || se.Duration == 0
}

var (
	ScreenEventPool = []func(*rand.Rand) *ScreenEvent{
		NewScreenEventNoop,
		NewScreenEventRelief,
		NewScreenEventProjection,
		//NewScreenEventChromaticAb,
	}
)

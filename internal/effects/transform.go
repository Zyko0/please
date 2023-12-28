package effects

import (
	"math"
	"math/rand"

	"github.com/Zyko0/please/internal/frame"
	"github.com/hajimehoshi/ebiten/v2"
)

type Transform interface {
	Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64)
}

type TransformTriggerMode byte

const (
	TransformTriggerModeEach TransformTriggerMode = iota
	TransformTriggerModeRng
	TransformTriggerModeCount
)

type TransformInterpolation byte

const (
	TransformInterpolationNone TransformInterpolation = iota
	TransformInterpolationLinear
	TransformInterpolationSine
	TransformInterpolationCosine
	TransformInterpolationTriSine
	TransformInterpolationTriCosine
	TransformInterpolationCount
)

type TransformTrigger struct {
	Mode  TransformTriggerMode
	Value float64
}

type TransformFunc func([]ebiten.Vertex, *rand.Rand)

type BaseTransform struct {
	Trigger       TransformTrigger
	Interpolation TransformInterpolation
	Duration      uint64
}

func (bt *BaseTransform) ensureTransform(rng *rand.Rand, tick uint64, active *uint64) (float32, bool) {
	// Try to activate the transformation
	if *active < tick {
		switch bt.Trigger.Mode {
		case TransformTriggerModeEach:
			if tick%uint64(bt.Trigger.Value) == 0 {
				*active = tick + bt.Duration
			}
		case TransformTriggerModeRng:
			if rng.Float64() < bt.Trigger.Value {
				*active = tick + bt.Duration
			}
		}
	}
	// If not active
	if *active < tick {
		return 0, false
	}
	// Interpolate
	c := 1 - float32(*active-tick)/float32(bt.Duration)
	switch bt.Interpolation {
	case TransformInterpolationNone:
		c = 1
	case TransformInterpolationLinear:
	case TransformInterpolationSine:
		c = float32(math.Sin(
			float64(c) * math.Pi,
		))
	case TransformInterpolationCosine:
		c = float32(math.Cos(
			float64(c) * math.Pi,
		))
	case TransformInterpolationTriSine:
		c = -abs(float32(math.Mod(float64(c)+0.5, 2)-1))*2 + 1
	case TransformInterpolationTriCosine:
		c = -abs(float32(math.Mod(float64(c)+1, 2)-1))*2 + 1
	}

	return c, true
}

func NewTransformFunc(t Transform) TransformFunc {
	start := frame.Current()
	active := uint64(0)
	return func(vertices []ebiten.Vertex, rng *rand.Rand) {
		tick := frame.Current() - start
		t.Apply(vertices, rng, tick, &active)
	}
}

package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/effects"
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/heuristics"
	"github.com/hajimehoshi/ebiten/v2"
)

func randomSign(rng *rand.Rand) float32 {
	if rng.Intn(2) == 0 {
		return -1
	}
	return 1
}

func randomTransformInterpolation(rng *rand.Rand) effects.TransformInterpolation {
	interp := rng.Intn(int(effects.TransformInterpolationCount))
	return effects.TransformInterpolation(interp)
}

func randomXY(rng *rand.Rand, factor, base float32) (float32, float32) {
	var dx, dy float32
	if rng.Intn(2) == 0 {
		dx = rng.Float32() * factor * base
	}
	if rng.Intn(2) == 0 {
		dy = rng.Float32() * factor * base
	}

	return dx, dy
}

func randomRGB(rng *rand.Rand) (float32, float32, float32) {
	var r, g, b float32

	//if config.EffectCumulativeTransform() {
	r = rng.Float32() * config.EffectColorFactor
	g = rng.Float32() * config.EffectColorFactor
	b = rng.Float32() * config.EffectColorFactor
	/*} else {
		switch rng.Intn(3) {
		case 0:
			r = rng.Float32() * config.EffectColorFactor()
		case 1:
			g = rng.Float32() * config.EffectColorFactor()
		case 2:
			b = rng.Float32() * config.EffectColorFactor()
		}
	}*/

	return r, g, b
}

func newRandomEffect(rng *rand.Rand) NewEffectFunc {
	return func(_ bool) *effects.Effect {
		duration := config.EffectTransformDuration * float64(frame.TPS())
		transforms := []effects.TransformFunc{}
		for i := 0; i < config.EffectTransformsCount; i++ {
			base := effects.BaseTransform{
				Interpolation: randomTransformInterpolation(rng),
				Duration:      uint64(duration),
			}
			var t effects.Transform
			switch rng.Intn(5) {
			case 0:
				dx, dy := randomXY(rng, config.EffectTranslationFactor, 50)
				t = &effects.Translate{
					BaseTransform: base,
					Dx:            dx,
					Dy:            dy,
				}
			case 1:
				dx, dy := randomXY(rng, 1, config.EffectScaleFactor)
				t = &effects.Scale{
					BaseTransform: base,
					Dx:            dx,
					Dy:            dy,
				}
			case 2:
				t = &effects.Rotate{
					BaseTransform: base,
					Sx:            randomSign(rng),
					Sy:            randomSign(rng),
					CMin:          0,
					CMax:          1,
				}
			case 3:
				r, g, b := randomRGB(rng)
				t = &effects.RotateColor{
					BaseTransform: base,
					R:             r,
					G:             g,
					B:             b,
				}
			case 4:
				r, g, b := randomRGB(rng)
				t = &effects.ScaleColor{
					BaseTransform: base,
					R:             r,
					G:             g,
					B:             b,
				}
			}
			transforms = append(transforms, effects.NewTransformFunc(t))
		}
		var interval = uint64(1)
		if config.EffectFrequency > 0 && config.EffectFrequency <= 1 {
			interval = uint64(1 / config.EffectFrequency)
		}

		return effects.New(
			rng,
			interval,
			transforms,
			assets.ShaderNoop,
			nil, nil,
		)
	}
}

func NewNoopEffect(_ bool) *effects.Effect {
	return effects.NewNoopEffect()
}

func NewEventDefault(rng *rand.Rand) *Event {
	return &Event{
		effectInstances: [heuristics.Count]int{},
		Name:            "Default",
		EffectInstancers: [heuristics.Count]NewEffectFunc{
			heuristics.Player:     newRandomEffect(rng),
			heuristics.Enemy:      newRandomEffect(rng),
			heuristics.Resource:   newRandomEffect(rng),
			heuristics.Projectile: newRandomEffect(rng),
			heuristics.Block:      newRandomEffect(rng),
			heuristics.UI:         NewNoopEffect,
			heuristics.Text:       newRandomEffect(rng),
			heuristics.Unknown:    NewNoopEffect,
		},
		Start:    frame.Current(),
		Duration: uint64(config.EventDuration * float64(ebiten.TPS())),
	}
}

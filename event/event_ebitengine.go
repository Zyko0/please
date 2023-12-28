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

var one float32 = 1

func newEbitenginePlayerEffect(rng *rand.Rand) NewEffectFunc {
	return func(variation bool) *effects.Effect {
		transforms := []effects.TransformFunc{
			effects.NewTransformFunc(&effects.ScaleLeast{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationLinear,
					Duration:      60,
				},
				Dx: 32,
				Dy: 32,
			}),
			effects.NewTransformFunc(&effects.Rotate{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationTriSine,
					Duration:      60,
				},
				Sx:   1,
				Sy:   1,
				CMin: -0.125,
				CMax: 0.125,
			}),
			effects.NewTransformFunc(&effects.TranslateColor{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationNone,
					Duration:      20,
				},
				R: 1,
				G: 1,
				B: 1,
				A: &one,
			}),
		}

		return effects.New(
			rng,
			1,
			transforms,
			assets.ShaderNoop,
			assets.GopherImage, nil,
		)
	}
}

func newEbitengineEnemyEffect(rng *rand.Rand) NewEffectFunc {
	return func(variation bool) *effects.Effect {
		var g, b float32
		if variation {
			g = rng.Float32() * 0.25
			b = rng.Float32() * 0.5
		}
		transforms := []effects.TransformFunc{
			effects.NewTransformFunc(&effects.ScaleLeast{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationLinear,
					Duration:      60,
				},
				Dx: 32,
				Dy: 32,
			}),
			effects.NewTransformFunc(&effects.Rotate{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationTriSine,
					Duration:      20,
				},
				Sx:   1,
				Sy:   1,
				CMin: -0.1,
				CMax: 0.1,
			}),
			effects.NewTransformFunc(&effects.Translate{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationSine,
					Duration:      20,
				},
				Dx: 0,
				Dy: -20,
			}),
			effects.NewTransformFunc(&effects.TranslateColor{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationNone,
					Duration:      20,
				},
				R: 1,
				G: g,
				B: b,
				A: &one,
			}),
		}

		return effects.New(
			rng,
			1,
			transforms,
			assets.ShaderNoop,
			assets.GopherAnim0Image, nil,
		)
	}
}

func newEbitengineResourceEffect(rng *rand.Rand) NewEffectFunc {
	return func(variation bool) *effects.Effect {
		img := assets.EbitenRes0Image
		if variation {
			switch rng.Intn(3) {
			case 0:
				img = assets.EbitenRes1Image
			case 1:
				img = assets.EbitenRes2Image
			case 2:
				img = assets.EbitenRes3Image
			}
		}
		transforms := []effects.TransformFunc{
			effects.NewTransformFunc(&effects.ScaleLeast{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationLinear,
					Duration:      60,
				},
				Dx: 16,
				Dy: 16,
			}),
			effects.NewTransformFunc(&effects.TranslateColor{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationNone,
					Duration:      20,
				},
				R: 1,
				G: 1,
				B: 1,
				A: &one,
			}),
		}

		return effects.New(
			rng,
			1,
			transforms,
			assets.ShaderNoop,
			img, nil,
		)
	}
}

func newEbitengineProjectileEffect(rng *rand.Rand) NewEffectFunc {
	return func(_ bool) *effects.Effect {
		transforms := []effects.TransformFunc{
			effects.NewTransformFunc(&effects.ScaleLeast{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationLinear,
					Duration:      120,
				},
				Dx: 16,
				Dy: 16,
			}),
			effects.NewTransformFunc(&effects.Rotate{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationLinear,
					Duration:      60,
				},
				Sx:   randomSign(rng),
				Sy:   randomSign(rng),
				CMin: 0,
				CMax: 1,
			}),
			effects.NewTransformFunc(&effects.TranslateColor{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationNone,
					Duration:      20,
				},
				R: 1,
				G: 1,
				B: 1,
				A: &one,
			}),
		}

		return effects.New(
			rng,
			1,
			transforms,
			assets.ShaderNoop,
			assets.EbitenImage, nil,
		)
	}
}

func newEbitengineBlockEffect(rng *rand.Rand) NewEffectFunc {
	return func(variation bool) *effects.Effect {
		img := assets.EbitenBlock0Image
		if variation {
			switch rng.Intn(4) {
			case 0:
				img = assets.EbitenBlock1Image
			case 1:
				img = assets.EbitenBlock2Image
			case 2:
				img = assets.EbitenBlock3Image
			case 3:
				img = assets.EbitenBlock4Image
			}
		}
		transforms := []effects.TransformFunc{
			effects.NewTransformFunc(&effects.TranslateColor{
				BaseTransform: effects.BaseTransform{
					Trigger: effects.TransformTrigger{
						Mode:  effects.TransformTriggerModeEach,
						Value: 1,
					},
					Interpolation: effects.TransformInterpolationNone,
					Duration:      20,
				},
				R: 1,
				G: 1,
				B: 1,
				A: &one,
			}),
		}

		return effects.New(
			rng,
			1,
			transforms,
			assets.ShaderNoop,
			img, nil,
		)
	}
}

func NewEbitengineEvent(rng *rand.Rand) *Event {
	return &Event{
		effectInstances: [heuristics.Count]int{},
		Name:            "Ebitengine",
		EffectInstancers: [heuristics.Count]NewEffectFunc{
			heuristics.Player:     newEbitenginePlayerEffect(rng),
			heuristics.Enemy:      newEbitengineEnemyEffect(rng),
			heuristics.Resource:   newEbitengineResourceEffect(rng),
			heuristics.Projectile: newEbitengineProjectileEffect(rng),
			heuristics.Block:      newEbitengineBlockEffect(rng),
			heuristics.UI:         NewNoopEffect,
			heuristics.Text:       newRandomEffect(rng),
			heuristics.Unknown:    NewNoopEffect,
		},
		Start:    frame.Current(),
		Duration: uint64(config.EventDuration() * float64(ebiten.TPS())),
	}
}

package please

import (
	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/locker"
)

type Mode string

const (
	None    Mode = "NONE"
	Default Mode = "DEFAULT"
	Medium  Mode = "MEDIUM"
	Unsafe  Mode = "UNSAFE"
)

func SetMode(mode Mode) {
	locker.Lock()
	defer locker.Unlock()

	config.Noop = false
	switch mode {
	case None:
		config.Noop = true
	case Default:
		config.EffectTransformsCount = 1
		config.EffectTransformDuration = 2
		config.EffectTranslationFactor = 1
		config.EffectScaleFactor = 2
		config.EffectColorFactor = 1
		config.EffectFrequency = 0.25
		config.EventDuration = 5
		config.ScreenEventDuration = 7.5
		config.UniformFactor = 1
	case Medium:
		config.EffectTransformsCount = 2
		config.EffectTransformDuration = 1.75
		config.EffectTranslationFactor = 1.5
		config.EffectScaleFactor = 4
		config.EffectColorFactor = 1.25
		config.EffectFrequency = 0.4
		config.EventDuration = 4
		config.ScreenEventDuration = 6
		config.UniformFactor = 2
	case Unsafe:
		config.EffectTransformsCount = 5
		config.EffectTransformDuration = 1.25
		config.EffectTranslationFactor = 1.5
		config.EffectScaleFactor = 4
		config.EffectColorFactor = 1.5
		config.EffectFrequency = 0.8
		config.EventDuration = 2.5
		config.ScreenEventDuration = 4
		config.UniformFactor = 8
	}
}

package runtime

import (
	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/effects"
	"github.com/Zyko0/please/internal/graphics"
)

var (
	warnings []*effects.WarningBubble

	ebitengineWarnings = map[string]bool{}
	etxtWarning        = false
	solarLuneWarning   = false
)

func pickWarningXY() (float64, float64) {
	var x, y float64
	if img := graphics.Screen(); img != nil {
		w, h := img.Bounds().Dx(), img.Bounds().Dy()
		x, y = float64(w)*rng.Float64(), float64(h)*rng.Float64()
	}

	return x, y
}

func RecordEbitengine(str string) {
	if _, ok := ebitengineWarnings[str]; !ok {
		x, y := pickWarningXY()
		ebitengineWarnings[str] = true
		warnings = append(warnings, effects.NewWarningBubble(
			x, y, str, assets.HajimeHoshiImage,
		))
	}
}

func RecordEtxt() {
	if !etxtWarning {
		x, y := pickWarningXY()
		etxtWarning = true
		warnings = append(warnings, effects.NewWarningBubble(
			x, y, "Hi! Welcome to ebitengine: github.com/tinne26/etxt", assets.DogcowImage,
		))
	}
}

func RecordSolarLune() {
	if !solarLuneWarning {
		x, y := pickWarningXY()
		solarLuneWarning = true
		warnings = append(warnings, effects.NewWarningBubble(
			x, y, "Cool library 8) github.com/SolarLune", assets.SolarLuneImage,
		))
	}
}

// TODO: more

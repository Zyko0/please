package runtime

import (
	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/graphics"
)

var (
	warnings []*graphics.WarningBubble

	ebitengineWarnings = map[string]bool{}
	etxtWarning        = false
	solarLuneWarning   = false
)

func RecordEbitengine(str string) {
	if _, ok := ebitengineWarnings[str]; !ok {
		ebitengineWarnings[str] = true
		warnings = append(warnings, graphics.NewWarningBubble(
			str, assets.HajimeHoshiImage,
		))
	}
}

func RecordEtxt() {
	if !etxtWarning {
		etxtWarning = true
		warnings = append(warnings, graphics.NewWarningBubble(
			"Hi! Welcome to ebitengine: github.com/tinne26/etxt", assets.DogcowImage,
		))
	}
}

func RecordSolarLune() {
	if !solarLuneWarning {
		solarLuneWarning = true
		warnings = append(warnings, graphics.NewWarningBubble(
			"Cool library 8) github.com/SolarLune", assets.SolarLuneImage,
		))
	}
}

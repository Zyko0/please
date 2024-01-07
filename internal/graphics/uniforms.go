package graphics

import "github.com/Zyko0/please/internal/frame"

func EffectUniforms() map[string]any {
	tick := frame.Current()
	tps := frame.TPS()
	return map[string]any{
		"ITime": float32(tick / tps),
		"FTime": float32(tick%tps) / float32(tps),
	}
}
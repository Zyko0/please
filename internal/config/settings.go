package config

var (
	effectTransformsCount     = 1
	effectTransformDuration   = float64(2)
	effectTranslationFactor   = float32(1)
	effectColorFactor         = float32(1)
	effectFrequency           = float64(0.2)
	effectCumulativeTransform = false
	eventDuration             = float64(5)
)

func EffectTransformsCount() int {
	return effectTransformsCount
}

func SetEffectTransformsCount(count int) {
	effectTransformsCount = count
}

func EffectTransformDuration() float64 {
	return effectTransformDuration
}

func SetEffectTransformDuration(seconds float64) {
	effectTransformDuration = seconds
}

func EffectTranslationFactor() float32 {
	return effectTranslationFactor
}

func SetEffectTranslationFactor(factor float32) {
	effectTranslationFactor = factor
}

func EffectColorFactor() float32 {
	return effectColorFactor
}

func SetEffectColorFactor(factor float32) {
	effectColorFactor = factor
}

func EffectFrequency() float64 {
	return effectFrequency
}

func SetEffectFrequency(frequency float64) {
	effectFrequency = min(frequency, 1)
}

func EffectCumulativeTransform() bool {
	return effectCumulativeTransform
}

func SetEffectCumulativeTransform(cumulative bool) {
	effectCumulativeTransform = cumulative
}

func EventDuration() float64 {
	return eventDuration
}

func SetEventDuration(seconds float64) {
	eventDuration = seconds
}

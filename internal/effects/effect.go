package effects

import (
	"math/rand"

	"github.com/Zyko0/please/internal/assets"
	"github.com/Zyko0/please/internal/frame"
	"github.com/hajimehoshi/ebiten/v2"
)

type Effect struct {
	counter  uint64
	interval uint64
	rng      *rand.Rand

	transforms []TransformFunc
	shader     *ebiten.Shader
	img        *ebiten.Image
	textFunc   func(string) string
}

func New(
	rng *rand.Rand,
	interval uint64,
	transforms []TransformFunc,
	shader *ebiten.Shader,
	img *ebiten.Image,
	textFunc func(string) string,
) *Effect {
	return &Effect{
		counter:  0,
		interval: interval,
		rng:      rng,

		transforms: transforms,
		img:        img,
		shader:     shader,
		textFunc:   textFunc,
	}
}

func NewNoopEffect() *Effect {
	return &Effect{
		counter:  0,
		interval: 1,
		shader:   assets.ShaderNoop,
	}
}

func (e *Effect) IsNoop() bool {
	return e.shader == assets.ShaderNoop && e.rng == nil
}

func (e *Effect) ResetCounter() {
	e.counter = 0
}

func (e *Effect) UpdateCounter() {
	e.counter++
}

func (e *Effect) Active() bool {
	return e.counter%e.interval == 0
}

func (e *Effect) Image() *ebiten.Image {
	return e.img
}

func (e *Effect) Shader() *ebiten.Shader {
	return e.shader
}

func (e *Effect) ApplyTransformations(vertices []ebiten.Vertex) {
	for _, t := range e.transforms {
		t(vertices, e.rng)
	}
}

func (e Effect) ApplyText(s string) string {
	if e.textFunc == nil {
		return s
	}
	return e.textFunc(s)
}

func (e *Effect) Uniforms() map[string]any {
	tick := frame.Current()
	tps := frame.TPS()
	return map[string]any{
		"ITime": float32(tick / tps),
		"FTime": float32(tick%tps) / float32(tps),
	}
}

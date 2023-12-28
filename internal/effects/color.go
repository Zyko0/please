package effects

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScaleColor

type ScaleColor struct {
	BaseTransform
	R, G, B float32
	A       *float32
}

func (s *ScaleColor) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if s == nil {
		return
	}
	// Try to activate the transformation
	c, ok := s.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	for i := range vertices {
		vertices[i].ColorR *= (c * s.R)
		vertices[i].ColorG *= (c * s.G)
		vertices[i].ColorB *= (c * s.B)
		if s.A != nil {
			vertices[i].ColorA = *s.A
		}
	}
}

// TranslateColor

type TranslateColor struct {
	BaseTransform
	R, G, B float32
	A       *float32
}

func (t *TranslateColor) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if t == nil {
		return
	}
	// Try to activate the transformation
	c, ok := t.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	for i := range vertices {
		vertices[i].ColorR = (c * t.R)
		vertices[i].ColorG = (c * t.G)
		vertices[i].ColorB = (c * t.B)
		if t.A != nil {
			vertices[i].ColorA = *t.A
		}
	}
}

// RotateColor

type RotateColor struct {
	BaseTransform
	R, G, B float32
	A       *float32
}

func (r *RotateColor) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if r == nil {
		return
	}
	// Try to activate the transformation
	c, ok := r.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	for i := range vertices {
		vertices[i].ColorR = float32(math.Mod(float64(vertices[i].ColorR+(c*r.R)), 1.))
		vertices[i].ColorG = float32(math.Mod(float64(vertices[i].ColorR+(c*r.G)), 1.))
		vertices[i].ColorB = float32(math.Mod(float64(vertices[i].ColorR+(c*r.B)), 1.))
		if r.A != nil {
			vertices[i].ColorA = *r.A
		}
	}
}

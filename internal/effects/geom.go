package effects

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Translate

type Translate struct {
	BaseTransform
	Dx, Dy float32
}

func (t *Translate) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if t == nil {
		return
	}
	// Try to activate the transformation
	c, ok := t.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	for i := range vertices {
		vertices[i].DstX += c * t.Dx
		vertices[i].DstY += c * t.Dy
	}
}

// Scale

type Scale struct {
	BaseTransform
	Dx, Dy float32
}

func (s *Scale) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if s == nil {
		return
	}
	// Try to activate the transformation
	c, ok := s.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	dx, dy := c*(s.Dx-1), c*(s.Dy-1)
	switch {
	case len(vertices)%4 == 0:
		for i := 0; i < len(vertices); i += 4 {
			minx, miny := vertices[i+0].DstX, vertices[i+0].DstY
			maxx, maxy := vertices[i+3].DstX, vertices[i+3].DstY
			offx, offy := dx*abs(maxx-minx)/2, dy*abs(maxy-miny)/2
			vertices[i+0].DstX -= offx
			vertices[i+0].DstY -= offy
			vertices[i+1].DstX += offx
			vertices[i+1].DstY -= offy
			vertices[i+2].DstX -= offx
			vertices[i+2].DstY += offy
			vertices[i+3].DstX += offx
			vertices[i+3].DstY += offy
		}
	case len(vertices)%3 == 0:
		for i := 0; i < len(vertices); i += 3 {
			// Note: obviously wrong, no height or width for a triangle but probably fun
			minx, miny := vertices[i+0].DstX, vertices[i+0].DstY
			maxx, maxy := vertices[i+2].DstX, vertices[i+2].DstY
			offx, offy := dx*abs(maxx-minx)/2, dy*abs(maxy-miny)/2
			vertices[i+0].DstX -= offx
			vertices[i+0].DstY -= offy
			vertices[i+1].DstX += offx
			vertices[i+1].DstY -= offy
			vertices[i+2].DstX -= offx
			vertices[i+2].DstY += offy
		}
	}
}

// ScaleLeast

type ScaleLeast struct {
	BaseTransform
	Dx, Dy float32
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func (s *ScaleLeast) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if s == nil {
		return
	}
	// Try to activate the transformation
	_, ok := s.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	switch {
	case len(vertices)%4 == 0:
		for i := 0; i < len(vertices); i += 4 {
			minx, miny := vertices[i+0].DstX, vertices[i+0].DstY
			maxx, maxy := vertices[i+3].DstX, vertices[i+3].DstY
			width, height := abs(maxx-minx), abs(maxy-miny)
			dx, dy := float32(0), float32(0)
			if width < s.Dx {
				dx = (s.Dx - width) / 2
			}
			if height < s.Dy {
				dy = (s.Dy - height) / 2
			}
			vertices[i+0].DstX -= dx
			vertices[i+0].DstY -= dy
			vertices[i+1].DstX += dx
			vertices[i+1].DstY -= dy
			vertices[i+2].DstX -= dx
			vertices[i+2].DstY += dy
			vertices[i+3].DstX += dx
			vertices[i+3].DstY += dy
		}
	case len(vertices)%3 == 0:
		// Ignore % 3
		return
		/*for i := 0; i < len(vertices); i += 3 {
			vertices[i+0].DstX -= c * s.Dx
			vertices[i+0].DstY -= c * s.Dy
			vertices[i+1].DstX += c * s.Dx
			vertices[i+1].DstY -= c * s.Dy
			vertices[i+2].DstX -= c * s.Dx
			vertices[i+2].DstY += c * s.Dy
		}*/
	}
}

// Rotate

type Rotate struct {
	BaseTransform
	Sx   float32
	Sy   float32
	CMin float32
	CMax float32
}

func (r *Rotate) Apply(vertices []ebiten.Vertex, rng *rand.Rand, tick uint64, active *uint64) {
	if r == nil {
		return
	}
	// Try to activate the transformation
	c, ok := r.ensureTransform(rng, tick, active)
	if !ok {
		return
	}
	c = r.CMin + c*(r.CMax-r.CMin)
	cos := r.Sx * float32(math.Cos(float64(c*math.Pi)))
	sin := r.Sy * float32(math.Sin(float64(c*math.Pi)))
	switch {
	case len(vertices)%4 == 0:
		for i := 0; i < len(vertices); i += 4 {
			minx, miny := vertices[i+0].DstX, vertices[i+0].DstY
			maxx, maxy := vertices[i+3].DstX, vertices[i+3].DstY
			cx := maxx - abs(maxx-minx)/2
			cy := maxy - abs(maxy-miny)/2
			for j := 0; j < 4; j++ {
				vertices[i+j].DstX -= cx
				vertices[i+j].DstY -= cy
			}
			px, py := vertices[i+0].DstX, vertices[i+0].DstY
			vertices[i+0].DstX = px*cos - py*sin
			vertices[i+0].DstY = px*sin + py*cos
			px, py = vertices[i+1].DstX, vertices[i+1].DstY
			vertices[i+1].DstX = px*cos - py*sin
			vertices[i+1].DstY = px*sin + py*cos
			px, py = vertices[i+2].DstX, vertices[i+2].DstY
			vertices[i+2].DstX = px*cos - py*sin
			vertices[i+2].DstY = px*sin + py*cos
			px, py = vertices[i+3].DstX, vertices[i+3].DstY
			vertices[i+3].DstX = px*cos - py*sin
			vertices[i+3].DstY = px*sin + py*cos
			for j := 0; j < 4; j++ {
				vertices[i+j].DstX += cx
				vertices[i+j].DstY += cy
			}
		}
	case len(vertices)%3 == 0:
		for i := 0; i < len(vertices); i += 3 {
			minx, miny := vertices[i+0].DstX, vertices[i+0].DstY
			maxx, maxy := vertices[i+2].DstX, vertices[i+2].DstY
			cx := maxx - abs(maxx-minx)/2
			cy := maxy - abs(maxy-miny)/2
			for j := 0; j < 3; j++ {
				vertices[i+j].DstX -= cx
				vertices[i+j].DstY -= cy
			}
			px, py := vertices[i+0].DstX, vertices[i+0].DstY
			vertices[i+0].DstX = px*cos - py*sin
			vertices[i+0].DstY = px*sin + py*cos
			px, py = vertices[i+1].DstX, vertices[i+1].DstY
			vertices[i+1].DstX = px*cos - py*sin
			vertices[i+1].DstY = px*sin + py*cos
			px, py = vertices[i+2].DstX, vertices[i+2].DstY
			vertices[i+2].DstX = px*cos - py*sin
			vertices[i+2].DstY = px*sin + py*cos
			for j := 0; j < 3; j++ {
				vertices[i+j].DstX += cx
				vertices[i+j].DstY += cy
			}
		}
	}
}

package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var rectIndices = [6]uint16{0, 1, 2, 1, 2, 3}

func QuadVerticesIndices(dst, src *ebiten.Image, geom *ebiten.GeoM, colorScale *ebiten.ColorScale) ([]ebiten.Vertex, []uint16) {
	var r, g, b, a float32 = 1, 1, 1, 1
	if colorScale != nil {
		r = colorScale.R()
		g = colorScale.G()
		b = colorScale.B()
		a = colorScale.A()
	}
	sx, sy := float32(src.Bounds().Min.X), float32(src.Bounds().Min.Y)
	sw, sh := float32(src.Bounds().Dx()), float32(src.Bounds().Dy())
	dw, dh := float32(src.Bounds().Dx()), float32(src.Bounds().Dy())
	x0, y0 := float32(dst.Bounds().Min.X), float32(dst.Bounds().Min.Y)
	x1, y1 := x0+dw, y0
	x2, y2 := x0, y0+dh
	x3, y3 := x0+dw, y0+dh
	if geom != nil {
		ax, ay := geom.Apply(float64(x0), float64(y0))
		x0, y0 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x1), float64(y1))
		x1, y1 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x2), float64(y2))
		x2, y2 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x3), float64(y3))
		x3, y3 = float32(ax), float32(ay)
	}

	return []ebiten.Vertex{
		{
			DstX:   x0,
			DstY:   y0,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y1,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x2,
			DstY:   y2,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x3,
			DstY:   y3,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, rectIndices[:]
}

func QuadVerticesIndicesWithDims(dst, src *ebiten.Image, width, height int, geom *ebiten.GeoM, colorScale *ebiten.ColorScale) ([]ebiten.Vertex, []uint16) {
	var r, g, b, a float32 = 1, 1, 1, 1
	if colorScale != nil {
		r = colorScale.R()
		g = colorScale.G()
		b = colorScale.B()
		a = colorScale.A()
	}
	var sx, sy, sw, sh float32
	if src != nil {
		sx, sy = float32(src.Bounds().Min.X), float32(src.Bounds().Min.Y)
		sw, sh = float32(src.Bounds().Dx()), float32(src.Bounds().Dy())
	}
	dw, dh := float32(width), float32(height)
	x0, y0 := float32(dst.Bounds().Min.X), float32(dst.Bounds().Min.Y)
	x1, y1 := x0+dw, y0
	x2, y2 := x0, y0+dh
	x3, y3 := x0+dw, y0+dh
	if geom != nil {
		ax, ay := geom.Apply(float64(x0), float64(y0))
		x0, y0 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x1), float64(y1))
		x1, y1 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x2), float64(y2))
		x2, y2 = float32(ax), float32(ay)
		ax, ay = geom.Apply(float64(x3), float64(y3))
		x3, y3 = float32(ax), float32(ay)
	}

	return []ebiten.Vertex{
		{
			DstX:   x0,
			DstY:   y0,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y1,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x2,
			DstY:   y2,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x3,
			DstY:   y3,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, rectIndices[:]
}

func AdaptVerticesToCustomImage(vertices []ebiten.Vertex, src *ebiten.Image) {
	minx, miny := float32(src.Bounds().Min.X), float32(src.Bounds().Min.Y)
	dx, dy := float32(src.Bounds().Dx()), float32(src.Bounds().Dy())

	switch {
	case len(vertices)%4 == 0:
		for i := 0; i < len(vertices); i += 4 {
			vertices[i+0].SrcX = minx
			vertices[i+0].SrcY = miny
			vertices[i+1].SrcX = minx + dx
			vertices[i+1].SrcY = miny
			vertices[i+2].SrcX = minx
			vertices[i+2].SrcY = miny + dy
			vertices[i+3].SrcX = minx + dx
			vertices[i+3].SrcY = miny + dy
		}
	case len(vertices)%3 == 0:
		for i := 0; i < len(vertices); i += 3 {
			vertices[i+0].SrcX = minx
			vertices[i+0].SrcY = miny
			vertices[i+1].SrcX = minx + dx
			vertices[i+1].SrcY = miny
			vertices[i+2].SrcX = minx
			vertices[i+2].SrcY = miny + dy
		}
	}
}

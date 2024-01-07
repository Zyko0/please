package graphics

import (
	"github.com/Zyko0/please/internal/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

var screen *ebiten.Image

func SetScreen(img *ebiten.Image) {
	screen = img
}

func Screen() *ebiten.Image {
	return screen
}

func DrawFullscreenEffect(dst, src *ebiten.Image, geom *ebiten.GeoM, effect any) {
	vertices, indices := QuadVerticesIndices(dst, src, geom, nil)
	dst.DrawTrianglesShader(vertices, indices, assets.ShaderRelief, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: EffectUniforms(),
		Images: [4]*ebiten.Image{
			src,
		},
	})
}

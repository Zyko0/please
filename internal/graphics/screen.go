package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screen *ebiten.Image
	geom   *ebiten.GeoM
)

func SetScreen(img *ebiten.Image) {
	screen = img
}

func Screen() *ebiten.Image {
	return screen
}

func DrawFullscreenEffect(dst, src *ebiten.Image, geom *ebiten.GeoM, shader *ebiten.Shader) {
	vertices, indices := QuadVerticesIndices(dst, src, geom, nil)
	dst.DrawTrianglesShader(vertices, indices, shader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: EffectUniforms(),
		Images: [4]*ebiten.Image{
			src,
		},
		AntiAlias: true,
	})
}

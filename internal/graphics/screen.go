package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	offscreen *ebiten.Image
	screen    *ebiten.Image
)

func SetScreen(img *ebiten.Image) {
	screen = img
	if screen != nil && (offscreen == nil || (screen.Bounds().Dx() != offscreen.Bounds().Dx() && screen.Bounds().Dy() != offscreen.Bounds().Dy())) {
		offscreen = ebiten.NewImageWithOptions(
			screen.Bounds(), &ebiten.NewImageOptions{
				Unmanaged: true,
			},
		)
	}
}

func Screen() *ebiten.Image {
	return screen
}

func Offscreen() *ebiten.Image {
	return offscreen
}

func DrawFullscreenEffect(dst, src *ebiten.Image, geom *ebiten.GeoM, shader *ebiten.Shader) {
	vertices, indices := QuadVerticesIndices(dst, src, geom, nil)
	dst.Clear()
	dst.DrawTrianglesShader(vertices, indices, shader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: EffectUniforms(),
		Images: [4]*ebiten.Image{
			src,
		},
	})
}

package assets

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaders/noop.kage
	srcShaderNoop []byte
	ShaderNoop    *ebiten.Shader

	//go:embed shaders/relief.kage
	srcShaderRelief []byte
	ShaderRelief    *ebiten.Shader

	//go:embed shaders/projection.kage
	srcShaderProjection []byte
	ShaderProjection    *ebiten.Shader

	//go:embed shaders/chromatic_ab.kage
	srcShaderChromaticAb []byte
	ShaderChromaticAb    *ebiten.Shader
)

func init() {
	var err error

	ShaderNoop, err = ebiten.NewShader(srcShaderNoop)
	if err != nil {
		log.Fatal("err:", err)
	}

	ShaderRelief, err = ebiten.NewShader(srcShaderRelief)
	if err != nil {
		log.Fatal("err:", err)
	}

	ShaderProjection, err = ebiten.NewShader(srcShaderProjection)
	if err != nil {
		log.Fatal("err:", err)
	}

	ShaderChromaticAb, err = ebiten.NewShader(srcShaderChromaticAb)
	if err != nil {
		log.Fatal("err:", err)
	}
}

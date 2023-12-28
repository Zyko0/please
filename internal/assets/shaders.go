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
)

func init() {
	var err error

	ShaderNoop, err = ebiten.NewShader(srcShaderNoop)
	if err != nil {
		log.Fatal("err:", err)
	}
}

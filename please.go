package please

import (
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/locker"
	"github.com/Zyko0/please/internal/patch"
	"github.com/Zyko0/please/internal/runtime"
	"github.com/Zyko0/please/metrics"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	patches []*patch.Patch
)

func init() {
	// Ensure frame tick package is initialized
	_ = frame.Current()

	patches = []*patch.Patch{
		patch.NewPatchMethod(&ebiten.Image{}, "DrawImage", patch.DrawImageReplacer),
		patch.NewPatchMethod(&ebiten.Image{}, "DrawTriangles", patch.DrawTrianglesReplacer),
		patch.NewPatchMethod(&ebiten.Image{}, "DrawRectShader", patch.DrawRectShaderReplacer),
		patch.NewPatchFunc(text.DrawWithOptions, patch.TextDrawReplacer),
		patch.NewPatchMethod(&ebiten.Image{}, "DrawTrianglesShader", patch.DrawTrianglesShaderReplacer),
		patch.NewPatchFunc(ebiten.NewImage, patch.NewImageReplacer),
		patch.NewPatchFunc(ebiten.NewImageWithOptions, patch.NewImageWithOptionsReplacer),
	}
	// Set patches reference internally to allow for live unpatching
	patch.SetTextPatch(patches[3])
	patch.SetTrianglesShaderPatch(patches[4])
	patch.SetNewImagePatch(patches[5])
	patch.SetNewImageWithOptionsPatch(patches[6])
}

// Sure!
func GlitchMe() {
	locker.Lock()
	defer locker.Unlock()
	for _, p := range patches {
		p.Enable()
	}
}

// :(
func DontGlitchMe() {
	locker.Lock()
	defer locker.Unlock()
	for _, p := range patches {
		p.Disable()
	}
}

// Fine..
func GiveMeSomethingUsefull() *metrics.Metrics {
	m := &metrics.Metrics{}
	runtime.FillMetrics(m)
	return m
}

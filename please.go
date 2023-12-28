package please

import (
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/locker"
	"github.com/Zyko0/please/internal/patch"
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
	}
	// Give patch address to DrawTrianglesShader replace in order to enable/disable its own patch
	patch.SetTextPatch(patches[3])
	patch.SetTrianglesShaderPatch(patches[4])
}

// Since it's asked so nicely
func GlitchMe() {
	locker.Lock()
	defer locker.Unlock()
	for _, p := range patches {
		p.Enable()
	}
}

// Okay, I'm sorry
func DontGlitchMe() {
	locker.Lock()
	defer locker.Unlock()
	for _, p := range patches {
		p.Disable()
	}
}

// Bored already?
func Shuffle() {

}

type SomethingUsefull struct {
}

// Fine..
func GiveMeSomethingUsefull() {

}

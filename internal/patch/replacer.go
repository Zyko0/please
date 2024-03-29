package patch

import (
	"image"

	"github.com/Zyko0/please/internal/caller"
	"github.com/Zyko0/please/internal/graphics"
	"github.com/Zyko0/please/internal/locker"
	"github.com/Zyko0/please/internal/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Keep patch references in order to call the originals
var (
	patchText                *Patch
	patchTrianglesShader     *Patch
	patchNewImage            *Patch
	patchNewImageWithOptions *Patch
)

func SetTextPatch(p *Patch) {
	patchText = p
}

func SetTrianglesShaderPatch(p *Patch) {
	patchTrianglesShader = p
}

func SetNewImagePatch(p *Patch) {
	patchNewImage = p
}

func SetNewImageWithOptionsPatch(p *Patch) {
	patchNewImageWithOptions = p
}

func NewImageReplacer(width, height int) *ebiten.Image {
	locker.Lock()
	defer locker.Unlock()

	info := caller.ExtractInfo()
	if info.Origin == caller.OriginUser {
		runtime.RegisterNewImage()
	}

	patchNewImage.Disable()
	defer patchNewImage.Enable()
	return ebiten.NewImage(width, height)
}

func NewImageWithOptionsReplacer(bounds image.Rectangle, options *ebiten.NewImageOptions) *ebiten.Image {
	locker.Lock()
	defer locker.Unlock()

	info := caller.ExtractInfo()
	if info.Origin == caller.OriginUser {
		runtime.RegisterNewImage()
	}

	patchNewImageWithOptions.Disable()
	defer patchNewImageWithOptions.Enable()
	return ebiten.NewImageWithOptions(bounds, options)
}

func DrawImageReplacer(dst, src *ebiten.Image, opts *ebiten.DrawImageOptions) {
	locker.Lock()
	defer locker.Unlock()
	info := caller.ExtractInfo()
	// If DrawFinalScreen is called, update the global manager
	// https://github.com/hajimehoshi/ebiten/blob/v2.6.3/gameforui.go#L173-L178
	var geom *ebiten.GeoM
	var colorScale *ebiten.ColorScale
	if opts != nil {
		geom = &opts.GeoM
		colorScale = &opts.ColorScale
		if sc := graphics.ColorMAsScale(opts.ColorM); sc != nil {
			runtime.RecordEbitengine("ColorM is deprecated.")
			colorScale.ScaleWithColorScale(*sc)
		}
	}
	if info.Origin == caller.OriginEbitengineDrawFinalScreen {
		var abort bool
		runtime.Update(dst)
		patchTrianglesShader.Disable()
		// If there's a fullscreen effect display screen with the effect and return
		if evt := runtime.GetScreenEvent(); evt != nil && !evt.Expired() {
			abort = true
			graphics.DrawFullscreenEffect(dst, src, geom, evt.Shader())
		}
		patchTrianglesShader.Enable()
		if abort {
			return
		}
	}
	// Get active effect
	effect := runtime.GetEffect(info)
	// New options
	newopts := &ebiten.DrawTrianglesShaderOptions{
		Uniforms: graphics.EffectUniforms(),
		Images: [4]*ebiten.Image{
			src,
		},
		AntiAlias: false,
	}
	if opts != nil {
		newopts.CompositeMode = opts.CompositeMode
		newopts.Blend = opts.Blend
	}
	// Register call and image information
	runtime.RegisterCall(info, src, geom)
	// Build vertices, indices
	vertices, indices := graphics.QuadVerticesIndices(dst, src, geom, colorScale)
	if img := effect.Image(); img != nil {
		newopts.Images[0] = img
		graphics.AdaptVerticesToCustomImage(vertices, img)
	}
	// Apply effect transformations
	effect.ApplyTransformations(vertices)
	// Draw triangles shader
	patchTrianglesShader.Disable()
	defer patchTrianglesShader.Enable()
	dst.DrawTrianglesShader(vertices, indices, effect.Shader(), newopts)
}

func DrawTrianglesReplacer(dst *ebiten.Image, vertices []ebiten.Vertex, indices []uint16, src *ebiten.Image, opts *ebiten.DrawTrianglesOptions) {
	locker.Lock()
	defer locker.Unlock()
	info := caller.ExtractInfo()
	// Get active effect
	effect := runtime.GetEffect(info)
	// Register call and image information
	runtime.RegisterCall(info, src, nil)
	// Update src image if there's one provided with the active effect
	if img := effect.Image(); img != nil {
		src = img
		graphics.AdaptVerticesToCustomImage(vertices, src)
	}
	// New options
	newopts := &ebiten.DrawTrianglesShaderOptions{
		Uniforms: graphics.EffectUniforms(),
		Images: [4]*ebiten.Image{
			src,
		},
		AntiAlias: false,
	}
	// Apply effect transformations
	effect.ApplyTransformations(vertices)
	// Draw triangles shader
	if opts != nil {
		newopts.CompositeMode = opts.CompositeMode
		newopts.Blend = opts.Blend
		newopts.FillRule = opts.FillRule
		newopts.AntiAlias = opts.AntiAlias
	}
	patchTrianglesShader.Disable()
	defer patchTrianglesShader.Enable()
	dst.DrawTrianglesShader(vertices, indices, effect.Shader(), newopts)
}

func DrawRectShaderReplacer(dst *ebiten.Image, width, height int, shader *ebiten.Shader, opts *ebiten.DrawRectShaderOptions) {
	locker.Lock()
	defer locker.Unlock()
	info := caller.ExtractInfo()
	// If DrawFinalScreen is called, update the global manager
	// https://github.com/hajimehoshi/ebiten/blob/v2.6.3/gameforui.go#L184
	if info.Origin == caller.OriginEbitengineDrawFinalScreen {
		var abort bool
		runtime.Update(opts.Images[0])
		patchTrianglesShader.Disable()
		// If there's a fullscreen effect display screen with the effect and return
		if evt := runtime.GetScreenEvent(); evt != nil && !evt.Expired() {
			abort = true
			graphics.DrawFullscreenEffect(dst, opts.Images[0], &opts.GeoM, evt.Shader())
		}
		patchTrianglesShader.Enable()
		if abort {
			return
		}
	}
	// Get active effect
	effect := runtime.GetEffect(info)
	// Apply effect transformations
	var src *ebiten.Image
	var geom *ebiten.GeoM
	var colorScale *ebiten.ColorScale
	newopts := &ebiten.DrawTrianglesShaderOptions{}
	if opts != nil {
		src = opts.Images[0]
		geom = &opts.GeoM
		colorScale = &opts.ColorScale
		newopts.Blend = opts.Blend
		newopts.CompositeMode = opts.CompositeMode
		newopts.Images = opts.Images
		newopts.Uniforms = opts.Uniforms
	}
	// Register call and image information
	runtime.RegisterCall(info, src, geom)
	// Build vertices, indices
	vertices, indices := graphics.QuadVerticesIndicesWithDims(dst, src, width, height, geom, colorScale)
	// If the effect forces a new image
	if effect.Shader() != nil && effect.Image() != nil {
		shader = effect.Shader()
		newopts.Uniforms = graphics.EffectUniforms()
		newopts.Images[0] = effect.Image()
		graphics.AdaptVerticesToCustomImage(vertices, effect.Image())
	}
	effect.ApplyTransformations(vertices)
	// Draw triangles shader
	patchTrianglesShader.Disable()
	defer patchTrianglesShader.Enable()
	dst.DrawTrianglesShader(vertices, indices, shader, newopts)
}

func TextDrawReplacer(dst *ebiten.Image, str string, face font.Face, opts *ebiten.DrawImageOptions) {
	locker.Lock()
	info := caller.ExtractInfo()
	// Register call information
	runtime.RegisterCall(info, nil, nil)
	// Get active effect
	effect := runtime.GetEffect(info)
	// Apply effect transformations
	// TODO: text char replacer, translate, color break, split by char
	str = effect.ApplyText(str)
	// Draw text with potentially altered colorm/coordinates
	patchText.Disable()
	defer patchText.Enable()
	// Note: Special global unlock because text.Draw is calling the patched DrawImage (which is also locking)
	locker.Unlock()
	text.DrawWithOptions(dst, str, face, opts)
}

func DrawTrianglesShaderReplacer(dst *ebiten.Image, vertices []ebiten.Vertex, indices []uint16, shader *ebiten.Shader, opts *ebiten.DrawTrianglesShaderOptions) {
	locker.Lock()
	defer locker.Unlock()
	info := caller.ExtractInfo()
	var src *ebiten.Image
	if opts == nil {
		opts = &ebiten.DrawTrianglesShaderOptions{}
	} else {
		src = opts.Images[0]
	}
	// Register call and image information
	runtime.RegisterCall(info, src, nil)
	// Get active effect
	effect := runtime.GetEffect(info)
	if effect.Shader() != nil && effect.Image() != nil {
		shader = effect.Shader()
		opts.Uniforms = graphics.EffectUniforms()
		opts.Images[0] = effect.Image()
		graphics.AdaptVerticesToCustomImage(vertices, effect.Image())
	}
	// Apply effect transformations
	effect.ApplyTransformations(vertices)
	// Draw triangles shader
	patchTrianglesShader.Disable()
	defer patchTrianglesShader.Enable()
	dst.DrawTrianglesShader(vertices, indices, shader, opts)
}

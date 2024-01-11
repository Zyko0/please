package metrics

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Caller struct {
	Name   string
	Count  int
	Parent *Caller
	Child  *Caller
}

type Metrics struct {
	Global struct {
		// NewImageCount represents the number of image created by the ebitengine app.
		NewImageCount int
		// Warnings defines a list of "hints" in an attempt to help improve the app.
		Warnings []string
	}
	Frame struct {
		// ImagesInUse records the number of different images referenced last frame.
		// It does not include subimages.
		ImagesInUse int
		// ImageMemory represents an estimation of used GPU memory by the user images.
		// Note that this is a rough estimation calculated by user images and their size
		// It does not account for potential mipmapping, extra texture atlas or padding made
		// by ebitengine internally.
		// For each image: mem += image.Dx()*image.Dy()*4
		ImageMemory uint64
		// LargestImage represents the largest image used in this frame
		LargestImage *ebiten.Image
		// DrawCount records the number of draw calls issued to ebitengine last frame.
		// It does not represent the graphic driver draw calls and doesn't take
		// batching performed by ebitengine either, this is just a library usage info
		DrawCount int
		// DrawCallers returns the Go functions calling ebitengine Draw functions.
		DrawCallers []*Caller
	}
}

func formatMem(size uint64) string {
	const unit = uint64(1024)

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := unit, 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB (%d)", float64(size)/float64(div), "KMGTPE"[exp], size)
}

func (m *Metrics) Print() {
	fmt.Println("Metrics")
	// Images
	fmt.Println("NewImage count:", m.Global.NewImageCount)
	fmt.Println("Images in use:", m.Frame.ImagesInUse)
	fmt.Println("Images memory:", formatMem(m.Frame.ImageMemory))
	if img := m.Frame.LargestImage; img != nil {
		fmt.Printf("Largest image: [%d x %d] => %s\n",
			img.Bounds().Dx(), img.Bounds().Dy(),
			formatMem(uint64(img.Bounds().Dx()*img.Bounds().Dy()*4)),
		)
	}
	// Draw calls
	fmt.Println("Draw count:", m.Frame.DrawCount)
	fmt.Println("Draw callers:")
	var largestName, largestCount int
	for _, c := range m.Frame.DrawCallers {
		largestName = max(largestName, len(c.Name))
		largestCount = max(largestCount, c.Count)
	}
	fmtFn := strconv.FormatInt(int64(largestName), 10)
	cnt := int64(math.Floor(
		max(math.Log10(float64(largestCount)), 0) + 1,
	))
	fmtCount := strconv.FormatInt(cnt, 10)
	fmtFull := "* %-" + fmtFn + "s => %" + fmtCount + "d\n"
	for _, c := range m.Frame.DrawCallers {
		fmt.Printf(fmtFull, c.Name, c.Count)
	}
	// Extra
	if len(m.Global.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, w := range m.Global.Warnings {
			fmt.Println("- " + w)
		}
	}
}

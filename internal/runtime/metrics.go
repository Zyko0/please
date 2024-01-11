package runtime

import (
	"fmt"
	"path"

	"github.com/Zyko0/please/internal/caller"
	"github.com/Zyko0/please/internal/graphics"
	"github.com/Zyko0/please/internal/locker"
	"github.com/Zyko0/please/metrics"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	imagesUses    = map[*ebiten.Image]int{}
	newImageCount = 0

	m = &metrics.Metrics{}
)

func callerName(c *caller.Caller) string {
	_, file := path.Split(c.File)
	return fmt.Sprintf("%s(%s:%d)", c.Func, file, c.Line)
}

func imgMemSize(img *ebiten.Image) uint64 {
	return uint64(img.Bounds().Dx() * img.Bounds().Dy() * 4)
}

func updateMetrics() {
	// Reset
	m.Global.NewImageCount = 0
	m.Global.Warnings = m.Global.Warnings[:0]
	m.Frame.ImagesInUse = 0
	m.Frame.ImageMemory = 0
	m.Frame.DrawCount = 0
	m.Frame.DrawCallers = m.Frame.DrawCallers[:0]
	// Images
	var largest *ebiten.Image
	for img := range imagesUses {
		if graphics.IsDisposed(img) {
			continue
		}
		size := imgMemSize(img)
		m.Frame.ImagesInUse++
		m.Frame.ImageMemory += size
		if largest == nil || size > imgMemSize(largest) {
			largest = img
		}
	}
	m.Frame.LargestImage = largest
	m.Global.NewImageCount = newImageCount
	// Callers
	for hash, info := range callersByHash {
		c := info.Current
		if info.Origin != caller.OriginUser {
			if info.User == nil {
				continue
			}
			c = info.User
		}
		count := int(callsByHash[hash])
		mc := &metrics.Caller{
			Name:   callerName(c),
			Count:  count,
			Parent: nil,
			Child:  nil, // Let's assume no child under the estimated user call
		}
		m.Frame.DrawCallers = append(m.Frame.DrawCallers, mc)
		// Populate parents and children
		parent := c.Prev
		for parent != nil {
			// Abort as soon as the callers seems to come from ebitengine or Go
			if parent.ParseOrigin() != caller.OriginUser {
				break
			}
			mc.Parent = &metrics.Caller{
				Name:   callerName(parent),
				Count:  count,
				Parent: nil,
				Child:  mc,
			}
			mc = mc.Parent
			parent = parent.Prev
		}
		// Increment total draw count
		m.Frame.DrawCount += count
	}
	// Warnings
	for w := range ebitengineWarnings {
		m.Global.Warnings = append(m.Global.Warnings, w)
	}
}

func FillMetrics(dst *metrics.Metrics) {
	locker.Lock()
	defer locker.Unlock()
	*dst = *m
}

package runtime

import (
	"math"
	"math/rand"
	"time"

	"github.com/Zyko0/please/event"
	"github.com/Zyko0/please/internal/caller"
	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/effects"
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/graphics"
	"github.com/Zyko0/please/internal/heuristics"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	noopEffect = effects.NewNoopEffect()
)

var (
	ignoredTick = false

	rng        = rand.New(rand.NewSource(time.Now().UnixNano()))
	lastUpdate uint64

	effectsByHash   = map[caller.Hash]*effects.Effect{}
	callsByHash     = map[caller.Hash]uint{}
	callersByHash   = map[caller.Hash]*caller.Info{}
	srcBoundsByHash = map[caller.Hash][2]uint{}
	/*lastCallsByHash     = map[caller.Hash]uint{}
	lastCallersByHash   = map[caller.Hash]*caller.Info{}
	lastSrcBoundsByHash = map[caller.Hash][2]uint{}*/

	heuristicsMapping = map[caller.Hash]*heuristics.Confidence{}

	activeEvent       *event.Event
	activeScreenEvent *event.ScreenEvent
)

func RNG() *rand.Rand {
	return rng
}

func GetEffect(info *caller.Info) *effects.Effect {
	if frame.Chilling() {
		return noopEffect
	}
	hash := info.Hash()
	e, haveEffect := effectsByHash[hash]
	h, haveHeuristic := heuristicsMapping[hash]
	// If the current tick is ignored, do not update effects
	if !ignoredTick {
		// If no active effect for the hash, instanciate one
		// Or if an heuristic has been missed last frame (UNKNOWN)
		// => try instanciate a more adequate effect
		if !haveEffect ||
			(haveHeuristic && h.ID != heuristics.Unknown && e.IsNoop()) {
			if activeEvent == nil || activeEvent.Expired() {
				return noopEffect
			}
			if !haveEffect && !haveHeuristic {
				e = activeEvent.NewEffect(heuristics.Unknown)
			} else {
				e = activeEvent.NewEffect(h.ID)
			}
			effectsByHash[hash] = e
		}
		defer e.UpdateCounter()
	}
	if !e.Active() {
		return noopEffect
	}

	return e
}

func GetScreenEvent() *event.ScreenEvent {
	if frame.Chilling() || activeScreenEvent.Expired() {
		return nil
	}

	return activeScreenEvent
}

func sign32(n float64) float32 {
	if n < 0 {
		return -1
	}
	return 1
}

func RegisterCall(info *caller.Info, src *ebiten.Image, geom *ebiten.GeoM) {
	// If current tick is ignored, do not register anything
	if ignoredTick {
		return
	}
	callsByHash[info.Hash()]++
	callersByHash[info.Hash()] = info
	if src != nil {
		bounds := src.Bounds()
		x0, y0 := float64(0), float64(0)
		x1, y1 := float64(bounds.Dx()), float64(bounds.Dy())
		if geom != nil {
			x0, y0 = geom.Apply(x0, y0)
			x1, y1 = geom.Apply(x1, y1)
		}
		width := uint(math.Abs(x1 - x0))
		height := uint(math.Abs(y1 - y0))
		srcBoundsByHash[info.Hash()] = [2]uint{width, height}
		// Store the image for stats if not a subImage
		if !graphics.IsSubImage(src) {
			imagesUses[src]++
		}
	}
}

func RegisterNewImage() {
	newImageCount++
	if !frame.Chilling() {
		RecordEbitengine("Calling NewImage() at runtime should be avoided.")
	}
}

func infoString(hash caller.Hash) (string, bool) {
	c, ok := callersByHash[hash]
	if !ok || c == nil {
		return "", false
	}
	if c.User != nil {
		return c.User.String(), true
	}
	if c.Current != nil {
		return c.Current.String(), true
	}
	return "", false
}

func Update(screen *ebiten.Image) {
	// Record the instance of screen image
	graphics.SetScreen(screen)
	// If we're on the same logical tick than last update, skip
	tick := frame.Current()
	if lastUpdate == tick {
		ignoredTick = true
		return
	}
	lastUpdate = tick
	ignoredTick = false
	// If not chilling anymore and no active event, make a new one
	if !frame.Chilling() && activeEvent.Expired() {
		//activeEvent = event.NewEventNoop(rng)
		if config.Noop {
			activeEvent = event.NewEventNoop(rng)
		} else {
			activeEvent = event.EventPool[rng.Intn(len(event.EventPool))](rng)
		}
		// Reset last effects so that they get populated by new event
		clear(effectsByHash)
	}
	activeEvent.Update()
	// If not chilling anymore and no active screen event, and it's time for one
	if !frame.Chilling() && activeScreenEvent.Expired() {
		//activeScreenEvent = event.NewScreenEventRelief(rng)
		if config.Noop {
			activeScreenEvent = event.NewScreenEventNoop(rng)
		} else {
			activeScreenEvent = event.ScreenEventPool[rng.Intn(len(event.ScreenEventPool))](rng)
		}
	}
	activeScreenEvent.Update()
	// Reset all effects' counters for next frame
	for _, e := range effectsByHash {
		e.ResetCounter()
	}
	// Debug
	var biggestFn, biggestCount int
	for hash, count := range callsByHash {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		biggestFn = max(biggestFn, len(str))
		biggestCount = max(biggestCount, int(count))
	}
	/*fmtFn := strconv.FormatInt(int64(biggestFn), 10)
	cnt := int64(math.Floor(
		max(math.Log10(float64(biggestCount)), 0) + 1,
	))
	fmtCount := strconv.FormatInt(cnt, 10)*/

	/*fmt.Println("Statistics")
	fmtFull := "%-" + fmtFn + "s => %" + fmtCount + "d\n"
	for hash, count := range lastCallsByHash {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		fmt.Printf(fmtFull, str, count)
	}*/
	// Update heuristics (who is a player, an enemy, a projectile, etc..)
	heuristicsMapping = heuristics.Compute(callersByHash, srcBoundsByHash)
	/*fmt.Println("Heuristics")
	for hash, score := range heuristicsMapping {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		fmt.Printf("%s => %v\n", str, score)
	}
	fmt.Println("------------")*/
	updateMetrics()
	m.Print()
	// Clean up
	clear(callsByHash)
	clear(callersByHash)
	clear(srcBoundsByHash)
	clear(imagesUses)
}

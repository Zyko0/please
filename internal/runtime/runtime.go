package runtime

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/Zyko0/please/event"
	"github.com/Zyko0/please/internal/caller"
	"github.com/Zyko0/please/internal/effects"
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/graphics"
	"github.com/Zyko0/please/internal/heuristics"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	noopEffect = effects.NewNoopEffect()
	noopEvent  = &event.Event{}
)

var (
	rng        = rand.New(rand.NewSource(time.Now().UnixNano()))
	lastUpdate uint64

	effectsByHash     = map[caller.Hash]*effects.Effect{}
	callsByHash       = map[caller.Hash]uint64{}
	callersByHash     = map[caller.Hash]*caller.Info{}
	lastCallsByHash   = map[caller.Hash]uint64{}
	lastCallersByHash = map[caller.Hash]*caller.Info{}

	heuristicsMapping = map[caller.Hash]*heuristics.Confidence{}

	activeEvent *event.Event
	screenEvent *event.Event // TODO:
	//commandQueue = queue.NewTickQueue[func()]()
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
	if !e.Active() {
		return noopEffect
	}

	return e
}

func sign32(n float64) float32 {
	if n < 0 {
		return -1
	}
	return 1
}

func RegisterCall(info *caller.Info) {
	callsByHash[info.Hash()]++
	callersByHash[info.Hash()] = info
}

func infoString(hash caller.Hash) (string, bool) {
	c, ok := lastCallersByHash[hash]
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
		return
	}
	lastUpdate = tick
	// If not chilling anymore and no active event, make a new one
	if !frame.Chilling() && activeEvent.Expired() {
		activeEvent = event.NewEbitengineEvent(rng) //event.NewDefaultEvent(rng)
		// Reset last effects so that they get populated by new event
		clear(effectsByHash)
	}
	activeEvent.Update()
	// Clear recorded draw entries and backup last frame
	clear(lastCallsByHash)
	clear(lastCallersByHash)
	for k, v := range callsByHash {
		lastCallsByHash[k] = v
	}
	for k, v := range callersByHash {
		lastCallersByHash[k] = v
	}
	clear(callsByHash)
	clear(callersByHash)
	// Reset all effects' counters for next frame
	for _, e := range effectsByHash {
		e.ResetCounter()
	}
	// Debug
	var biggestFn, biggestCount int
	for hash, count := range lastCallsByHash {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		biggestFn = max(biggestFn, len(str))
		biggestCount = max(biggestCount, int(count))
	}
	fmtFn := strconv.FormatInt(int64(biggestFn), 10)
	cnt := int64(math.Floor(
		max(math.Log10(float64(biggestCount)), 0) + 1,
	))
	fmtCount := strconv.FormatInt(cnt, 10)

	fmt.Println("Statistics")
	fmtFull := "%-" + fmtFn + "s => %" + fmtCount + "d\n"
	for hash, count := range lastCallsByHash {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		fmt.Printf(fmtFull, str, count)
	}
	// Set active events if any
	/*activeEvent.Update()
	if activeEvent.Expired() {
		activeEvent = nil
		// TODO: enqueue new event
	}*/
	// Update heuristics (who is a player, an enemy, a projectile, etc..)
	heuristicsMapping = heuristics.Compute(lastCallsByHash, lastCallersByHash)
	fmt.Println("Heuristics")
	for hash, score := range heuristicsMapping {
		str, ok := infoString(hash)
		if !ok {
			continue
		}
		fmt.Printf("%s => %v\n", str, score)
	}
	fmt.Println("------------")
	// TODO: based on m.lastCallsByHash
}

package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/effects"
	"github.com/Zyko0/please/internal/heuristics"
)

type NewEffectFunc func(variation bool) *effects.Effect

type Event struct {
	effectInstances [heuristics.Count]int

	Name             string
	EffectInstancers [heuristics.Count]NewEffectFunc
	Start            uint64
	Duration         uint64
}

func (e *Event) Update() {
	if e == nil || e.Duration == 0 {
		return
	}

	e.Duration--
}

func (e *Event) Expired() bool {
	return e == nil || e.Duration == 0
}

func (e *Event) NewEffect(identifier heuristics.ID) *effects.Effect {
	e.effectInstances[identifier]++
	// Return a variation of the effect if the base one is already active
	return e.EffectInstancers[identifier](e.effectInstances[identifier] > 1)
}

var (
	EventPool = []func(*rand.Rand) *Event{
		NewEventDefault,
		NewEventEbitengine,
	}
)

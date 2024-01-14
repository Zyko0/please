package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/heuristics"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewEventNoop(rng *rand.Rand) *Event {
	return &Event{
		effectInstances: [heuristics.Count]int{},
		Name:            "Noop",
		EffectInstancers: [heuristics.Count]NewEffectFunc{
			heuristics.Player:     NewNoopEffect,
			heuristics.Enemy:      NewNoopEffect,
			heuristics.Resource:   NewNoopEffect,
			heuristics.Projectile: NewNoopEffect,
			heuristics.Block:      NewNoopEffect,
			heuristics.UI:         NewNoopEffect,
			heuristics.Text:       NewNoopEffect,
			heuristics.Unknown:    NewNoopEffect,
		},
		Start:    frame.Current(),
		Duration: uint64(config.EventDuration * float64(ebiten.TPS())),
	}
}

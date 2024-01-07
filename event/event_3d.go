package event

import (
	"math/rand"

	"github.com/Zyko0/please/internal/config"
	"github.com/Zyko0/please/internal/frame"
	"github.com/Zyko0/please/internal/heuristics"
	"github.com/hajimehoshi/ebiten/v2"
)

func New3DEvent(rng *rand.Rand) *Event {
	return &Event{
		effectInstances: [heuristics.Count]int{},
		Name:            "3D",
		EffectInstancers: [heuristics.Count]NewEffectFunc{
			heuristics.Player:     newRandomEffect(rng),
			heuristics.Enemy:      newRandomEffect(rng),
			heuristics.Resource:   newRandomEffect(rng),
			heuristics.Projectile: newRandomEffect(rng),
			heuristics.Block:      newRandomEffect(rng),
			heuristics.UI:         newRandomEffect(rng),
			heuristics.Text:       newRandomEffect(rng),
			heuristics.Unknown:    NewNoopEffect,
		},
		Start:    frame.Current(),
		Duration: uint64(config.EventDuration * float64(ebiten.TPS())),
	}
}

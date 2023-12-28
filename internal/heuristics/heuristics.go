package heuristics

import (
	"fmt"
	"strings"

	"github.com/Zyko0/please/internal/caller"
)

type ID byte

const (
	Player ID = iota
	Enemy
	Resource
	Projectile
	Block
	UI
	Text
	Unknown
	Count
)

var playerCoeffs = map[string]float64{
	"player":    1.,
	"hero":      1.,
	"ally":      1.,
	"allie":     1.,
	"friend":    1.,
	"actor":     0.75,
	"unit":      0.65,
	"character": 0.5,
	"troop":     0.45,
}

func player(hash string, count uint64) float64 {
	value := 0.
	for name, coeff := range playerCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	// Assume the least calls there is, the most likely it is to be the player
	return value / float64(count)
}

var enemyCoeffs = map[string]float64{
	"enemy":      1.,
	"enemie":     1.,
	"opponent":   1.,
	"boss":       1.,
	"monster":    1.,
	"creature":   1.,
	"elite":      1.,
	"hostile":    1.,
	"wave":       0.75,
	"summon":     0.5,
	"invocation": 0.5,
	"bad":        0.25,
	"angry":      0.25,
	"mean":       0.25,
	/*"entity":     0.1,
	"entitie":    0.1,*/
}

func enemy(hash string) float64 {
	value := 0.
	for name, coeff := range enemyCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	return value
}

var projectileCoeffs = map[string]float64{
	"projectile":  1.,
	"bullet":      1.,
	"rocket":      1.,
	"laser":       1.,
	"lazer":       1., // ??
	"projecticle": 1., // Hi kettek
	"ball":        0.75,
	"beam":        0.75,
	"shot":        0.5,
	"shoot":       0.5,
	"thunder":     0.5,
}

func projectile(hash string) float64 {
	value := 0.
	for name, coeff := range projectileCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	// Assume the most calls there is, the most likely it is to be a projectile
	return value
}

var resourceCoeffs = map[string]float64{
	"resource":  1.,
	"powerup":   1.,
	"ammo":      1.,
	"orb":       1.,
	"mana":      1.,
	"stamina":   1.,
	"armor":     1.,
	"item":      1.,
	"shield":    1.,
	"boost":     1.,
	"food":      1.,
	"object":    1.,
	"component": 1.,
	"weapon":    1.,
	"scrap":     1.,
	"health":    1.,
	"bonus":     1.,
	"pickup":    1.,
	"malus":     1.,
	"buff":      0.75,
	"wood":      0.5,
	"metal":     0.5,
}

func resource(hash string) float64 {
	value := 0.
	for name, coeff := range resourceCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	// Assume the most calls there is, the most likely it is to be a resource
	return value
}

var blockCoeffs = map[string]float64{
	"block":       1.,
	"tile":        1.,
	"box":         1.,
	"crate":       1.,
	"obstacle":    1.,
	"environment": 1.,
	"terrain":     1.,
	"map":         1.,
	"spawner":     1.,
	"wall":        1.,
	"tree":        1.,
	"bush":        1.,
	"asteroid":    1.,
	"snow":        1.,
	"grass":       1.,
	"platform":    1.,
	"portal":      1.,
	"door":        1.,
	"house":       1.,
	"ground":      1.,
	"floor":       1.,
	"building":    1.,
	"structure":   1.,
	"turret":      1.,
	"rock":        1.,
	"water":       1.,
	"fire":        1.,
	"mountain":    1.,
	"cloud":       1.,
	"wave":        0.75,
	"element":     0.75,
	"static":      0.75,
	"collider":    0.5,
}

func block(hash string) float64 {
	value := 0.
	for name, coeff := range blockCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	// Assume the most calls there is, the most likely it is to be a block
	return value
}

var uiCoeffs = map[string]float64{
	"button":      1.,
	"slider":      1.,
	"textbox":     1.,
	"label":       1.,
	"hud":         1.,
	"view":        1.,
	"description": 1.,
	"instruction": 1.,
	"hint":        1.,
	"tooltip":     1.,
	"chat":        1.,
	"message":     1.,
	"warning":     1.,
	"dialog":      1.,
	"bar":         1.,
	"progress":    1.,
}

func ui(hash string) float64 {
	value := 0.
	for name, coeff := range uiCoeffs {
		occ := float64(strings.Count(hash, name))
		value += occ * coeff
	}
	// Assume the most calls there is, the most likely it is to be a UI element
	return value
}

type Heuristic [Count]float64

type Confidence struct {
	ID        ID
	Score     float64
	Heuristic Heuristic
}

func (c *Confidence) String() string {
	if c == nil {
		return "<nil>"
	}
	var sb strings.Builder
	switch c.ID {
	case Player:
		sb.WriteString("PLAYER")
	case Enemy:
		sb.WriteString("ENEMY")
	case Resource:
		sb.WriteString("RESOURCE")
	case Projectile:
		sb.WriteString("PROJECTILE")
	case Block:
		sb.WriteString("BLOCK")
	case UI:
		sb.WriteString("UI")
	case Text:
		sb.WriteString("TEXT")
	default:
		sb.WriteString("UNKNOWN")
	}
	sb.WriteString(fmt.Sprintf(": %.2f", c.Score))

	return sb.String()
}

func Compute(calls map[caller.Hash]uint64, infos map[caller.Hash]*caller.Info) map[caller.Hash]*Confidence {
	heuristics := make(map[caller.Hash]*Heuristic, len(calls))
	// Compute guesses for each hash
	for hash, count := range calls {
		heuristics[hash] = nil
		if count == 0 {
			continue
		}
		info, ok := infos[hash]
		if !ok || info == nil {
			// Shouldn't happen
			continue
		}
		// If it comes from a text draw call don't try to run heuristics
		if info.Origin == caller.OriginEbitengineText {
			heuristics[hash] = &Heuristic{
				Text: 1.,
			}
			continue
		}
		// If it comes from ebitengine internal, skip it
		if info.Origin != caller.OriginUser && info.User == nil {
			continue
		}
		// Compute heuristic score based on all previous callers
		heuristics[hash] = &Heuristic{}
		for _, c := range info.AllCallers {
			// Don't populate heuristics based on non-user call
			if c.ParseOrigin() != caller.OriginUser {
				continue
			}
			normalized := strings.ToLower(c.Func)
			heuristics[hash][Player] += player(normalized, count)
			heuristics[hash][Enemy] += enemy(normalized)
			heuristics[hash][Resource] += resource(normalized)
			heuristics[hash][Projectile] += projectile(normalized)
			heuristics[hash][Block] += block(normalized)
			heuristics[hash][UI] += ui(normalized)
			heuristics[hash][Text] = 0
		}
	}
	// Find best matching ID for each hash
	scores := [Count][]caller.Hash{}
	for hash, heuristic := range heuristics {
		// Skip if no heuristic
		if heuristic == nil {
			continue
		}
		foundID := ID(255)
		foundScore := 0.
		for id := ID(0); id < Count; id++ {
			// Skip if null score
			if heuristic[id] == 0 {
				continue
			}
			// Check if there's a better score found
			if heuristic[id] > foundScore {
				foundID = id
				foundScore = heuristic[id]
			}
		}
		// Register hash if relevant ID
		if foundScore > 0. {
			scores[foundID] = append(scores[foundID], hash)
		}
	}
	// Build up confidence heuristics
	var results = map[caller.Hash]*Confidence{}
	for id, hashes := range scores {
		if hashes == nil {
			continue
		}
		for _, hash := range hashes {
			c := &Confidence{
				ID: ID(id),
			}
			if h, ok := heuristics[hash]; ok {
				c.Score = h[id]
				c.Heuristic = *h
			}
			results[hash] = c
		}
	}

	return results
}

package caller

import (
	"fmt"
	"hash/fnv"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type Hash uint64

type Origin byte

const (
	OriginUser Origin = iota
	OriginInternal
	OriginGo
	OriginEbitengine
	OriginEbitengineText
	OriginEbitengineTextV2
	OriginEbitengineDrawFinalScreen
)

func parseFuncOrigin(name string) Origin {
	switch {
	case strings.HasPrefix(name, "github.com/hajimehoshi/ebiten"):
		switch name {
		case "github.com/hajimehoshi/ebiten/v2/text.drawGlyph":
			return OriginEbitengineText
		case "github.com/hajimehoshi/ebiten/v2/text/v2.Draw":
			return OriginEbitengineTextV2
		case "github.com/hajimehoshi/ebiten/v2.(*gameForUI).DrawFinalScreen":
			return OriginEbitengineDrawFinalScreen
		default:
			return OriginEbitengine
		}
	case strings.HasPrefix(name, "github.com/Zyko0/please"):
		return OriginInternal
	case strings.HasPrefix(name, "golang.org/"):
		fallthrough
	case strings.HasPrefix(name, "runtime."):
		return OriginGo
	default:
		return OriginUser
	}
}

func computeUserHash(callers []*Caller) Hash {
	h := fnv.New64a()
	for _, c := range callers {
		if c.ParseOrigin() != OriginUser {
			continue
		}
		h.Write([]byte(c.File + c.Func + strconv.FormatInt(int64(c.Line), 10)))
	}
	return Hash(h.Sum64())
}

type Caller struct {
	File string
	Func string
	Line int

	Prev *Caller
	Next *Caller
}

func (c *Caller) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.Func, c.Line)
}

func (c *Caller) ParseOrigin() Origin {
	return parseFuncOrigin(c.Func)
}

type Info struct {
	hash Hash

	AllCallers []*Caller
	Origin     Origin
	Current    *Caller
	User       *Caller
}

func (i *Info) Hash() Hash {
	return i.hash
}

func ExtractInfo() *Info {
	// Closest caller
	currentFn, currentFile, currentLine, _ := runtime.Caller(2)
	current := &Caller{
		File: currentFile,
		Func: runtime.FuncForPC(currentFn).Name(),
		Line: currentLine,
	}
	origin := current.ParseOrigin()
	// All Callers
	var allCallers []*Caller
	pcs := make([]uintptr, 16)
	count := runtime.Callers(2, pcs)
	pcs = pcs[:count]
	if len(pcs) > 0 {
		allCallers = make([]*Caller, 0, len(pcs))
		frames := runtime.CallersFrames(pcs)
		more := true
		for more {
			var f runtime.Frame
			f, more = frames.Next()
			if f.Func != nil {
				allCallers = append(allCallers, &Caller{
					File: f.File,
					Func: f.Func.Name(),
					Line: f.Line,
				})
			} else {
				break
			}
		}
		// Arrange next, prevs for all callers
		slices.Reverse(allCallers)
		for i := range allCallers {
			if i > 0 {
				allCallers[i].Prev = allCallers[i-1]
			}
			if i < len(allCallers)-1 {
				allCallers[i].Next = allCallers[i+1]
			}
		}
	}
	// User caller if closest one isn't a user one
	var user *Caller
	if origin != OriginUser && origin != OriginEbitengineDrawFinalScreen {
		var file string
		var name string
		var line int
		userOrigin := origin
		skip := 3
		ok := true
		for ok && userOrigin != OriginUser {
			var fn uintptr
			fn, file, line, ok = runtime.Caller(skip)
			name = runtime.FuncForPC(fn).Name()
			userOrigin = parseFuncOrigin(name)
			skip++
		}
		if ok {
			user = &Caller{
				File: file,
				Func: name,
				Line: line,
			}
		}
	}

	return &Info{
		hash: computeUserHash(allCallers),

		AllCallers: allCallers,
		Origin:     origin,
		Current:    current,
		User:       user,
	}
}

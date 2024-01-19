# Please

Please is a fun library to make an ebitengine (https://github.com/hajimehoshi/ebiten) application glitchy.
It is using https://github.com/agiledragon/gomonkey to patch ebitengine functions at runtime and add extra effects to draw calls.

This is my entry for the https://itch.io/jam/ebitengine-holiday-hack-2023 jam and its theme: "bug/glitch".

# Usage

```
go run cmd/main.go <repository_path> [relative_build_folder:optional]
```

Flags:
- `--golog`: Outputs the logs resulting from `go` commands
- `--mode`: Defines the level of glitchness (none,default,medium,unsafe)

## Examples

- `go run cmd/main.go github.com/mharv/scrapyard-charter`
- `go run cmd/main.go github.com/elamre/attractive_defense`
- `go run cmd/main.go --mode=unsafe github.com/tinne26/transition`
- `go run cmd/main.go github.com/ketMix/retromancer`

With a relative folder to run from:
- `go run cmd/main.go --mode=medium github.com/hajimehoshi/ebiten/v2@latest examples/flappy`
- `go run cmd/main.go github.com/hajimehoshi/ebiten/v2@latest examples/snake`
- `go run cmd/main.go github.com/hajimehoshi/ebiten/v2@latest examples/blocks`
- `go run cmd/main.go github.com/kettek/ebijam22 cmd/magnet`

# Library mode

Just ask `please.GlitchMe()` somewhere! (ideally before `ebiten.RunGame`).
To disable the patches at runtime, you can `please.DontGlitchMe()`.

You can also specify a "mode" to define a different level of glitchness:
```go
please.SetMode(please.None) // Keep the features / function hooks but disable all glitches
please.SetMode(please.Unsafe) // Effects frequency and various factors are "maxed"
```

In order for this library not to be totally useless, there is a way to fetch a few metrics by doing:
```go
metrics := please.GiveMeSomethingUseful()
metrics.Print()
```

# In-game captures

* `go run cmd/main.go --mode=unsafe github.com/kettek/ebijam22 cmd/magnet`

![magnet](https://github.com/Zyko0/please/assets/13394516/da390452-d02a-4fb3-836a-4d0f49bad1f1)

* `go run cmd/main.go --mode=unsafe github.com/tinne26/transition`

![transition](https://github.com/Zyko0/please/assets/13394516/8fcd35d5-a660-4f3e-b630-06257e7688de)

# Notes

- Despite the attempt at putting some locks here and there, this library is definitely not thread-safe (if you're calling Draw commands from multiple goroutines at the same time).
- Not many fancy effects unfortunately, but some extra could be plugged easily after the package gets cleaned up a bit (jam code)
- Only tested on Windows!

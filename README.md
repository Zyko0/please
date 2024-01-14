# Please

Please is a troll library to corrupt / make an ebitengine application glitchy.

This is my entry for the https://itch.io/jam/ebitengine-holiday-hack-2023 jam and its theme: "bug/glitch".

# Usage

```
go run cmd/main.go <repository_path> [relative_build_folder:optional]
```

Flags:
- `--golog`: Outputs the logs resulting from `go` commands
- `--mode`: Defines the level of glitchness (none,default,medium,unsafe)

# Library-mode

Just ask `please.GlitchMe()` somewhere! (ideally before `ebiten.RunGame`).
To disable the patches at runtime, you can `please.DontGlitchMe()`.

You can also specify a "mode" to define a different level of glitchness:
```go
please.SetMode(please.None) // Keep the features / function hooks but disable all glitches
please.SetMode(please.Unsafe) // Effects frequency and various factors are "maxed"
```

In order for this library not to be totally useless, there is a way to fetch a few metrics by doing:
```go
metrics := please.GiveMeSomethingUsefull()
metrics.Print()
```

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

### Notes

- Despite the attempt at putting some locks here and there, this library is definitely not thread-safe (if you're calling DrawX commands from multiple goroutines at the same time).
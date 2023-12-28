package graphics

import "github.com/hajimehoshi/ebiten/v2"

var screen *ebiten.Image

func SetScreen(img *ebiten.Image) {
	screen = img
}

func Screen() *ebiten.Image {
	return screen
}

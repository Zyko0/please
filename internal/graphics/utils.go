package graphics

import (
	_ "unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:linkname isSubImage github.com/hajimehoshi/ebiten/v2.(*Image).isSubImage
func isSubImage(img *ebiten.Image) bool

func IsSubImage(img *ebiten.Image) bool {
	return img != nil && isSubImage(img)
}

//go:linkname isDisposed github.com/hajimehoshi/ebiten/v2.(*Image).isDisposed
func isDisposed(img *ebiten.Image) bool

func IsDisposed(img *ebiten.Image) bool {
	return img != nil && isDisposed(img)
}

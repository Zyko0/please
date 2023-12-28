package graphics

import "github.com/hajimehoshi/ebiten/v2"

func ColorMAsScale(colorm ebiten.ColorM) *ebiten.ColorScale {
	r := float32(colorm.Element(0, 0))
	g := float32(colorm.Element(1, 1))
	b := float32(colorm.Element(2, 2))
	a := float32(colorm.Element(3, 3))
	if r == 0 && g == 0 && b == 0 && a == 0 {
		return nil
	}

	sc := &ebiten.ColorScale{}
	if r < 0 || g < 0 || b < 0 || a < 0 || r > 1 || g > 1 || b > 1 {
		sc.Scale(1, 1, 1, 1)
		return sc
	}

	sc.Scale(r, g, b, a)
	sc.ScaleAlpha(a)

	return sc
}

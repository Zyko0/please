package graphics

import "github.com/hajimehoshi/ebiten/v2"

const (
	warningBubbleDuration = 120
)

type WarningBubble struct {
	tick    uint64
	avatar  *ebiten.Image
	message string
}

func NewWarningBubble(message string, avatar *ebiten.Image) *WarningBubble {
	return &WarningBubble{
		tick:    0,
		avatar:  avatar,
		message: message,
	}
}

func (wb *WarningBubble) Expired() bool {
	return wb == nil || wb.tick > warningBubbleDuration
}

func (wb *WarningBubble) Update() {
	wb.tick++
}

func (wb *WarningBubble) DrawNative(screen *ebiten.Image) {

}
